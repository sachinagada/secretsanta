package pick

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/sachinagada/secretsanta/send"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats"
	"go.opencensus.io/trace"
)

type sender interface {
	Send(ctx context.Context, santas []send.Santa) error
}

// Shuffler shuffles the slice of participants
type Shuffler interface {
	Shuffle([]string)
}

// server is a handler that will receive the list of participants and assign
// secret santas
type server struct {
	shuf   Shuffler
	sender sender
}

type participant struct {
	FirstName string
	LastName  string
	Email     string
}

// key is the email and the value is the name of the participant
type participantMap map[string]string

func NewServer(shuf Shuffler, sender sender, port string) *http.Server {
	ss := server{
		shuf:   shuf,
		sender: sender,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/santa", ss.pickSanta)
	mux.HandleFunc("/health", ss.healthCheck)
	mux.Handle("/", http.FileServer(http.Dir("./static")))

	httpSvr := http.Server{
		Addr: net.JoinHostPort("", port),
		Handler: &ochttp.Handler{
			Handler:          mux,
			IsPublicEndpoint: true,
			IsHealthEndpoint: func(r *http.Request) bool {
				return strings.Contains(r.URL.Path, "/health")
			},
		},
	}

	return &httpSvr
}

// healthCheck always returns 200 status. Being able to hit the health check
// endpoint indicates that the server is up and running
func (s *server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// ServeHTTP handles the requests from the form and assigns the secret santa to
// each participant. Each participant will receive communication regarding who
// they are assigned.
func (s *server) pickSanta(w http.ResponseWriter, r *http.Request) {
	b, readErr := io.ReadAll(r.Body)
	if readErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Uh oh. Rudolph is lost! Please try again."))
		return
	}

	var participants []participant
	if unmarshalErr := json.Unmarshal(b, &participants); unmarshalErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Please check all the information is entered correctly and try again."))
		return
	}

	partMap := make(participantMap, len(participants))
	for _, p := range participants {
		partMap[p.Email] = fmt.Sprintf("%s %s", strings.Title(p.FirstName), strings.Title(p.LastName))
	}

	if len(partMap) != len(participants) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Duplicate email address was used for some of the participants. Please check all the information and try again."))
		return
	}

	ctx := r.Context()
	stats.Record(ctx, MParticipants.M(int64(len(participants))))
	addrs := addresses(partMap)

	// assigned key is the secret santa and the value is person receiving the
	// gift
	assigned := pickSecretSanta(ctx, addrs, s.shuf)

	startTime := time.Now()
	defer stats.Record(ctx, MCommunicationLatency.M(time.Since(startTime).Milliseconds()))
	if informErr := s.inform(ctx, partMap, assigned); informErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(informErr.Error()))
		return
	}

	w.Write([]byte("Santas have been chosen and informed. Happy Holidays!"))
}

func (s *server) inform(ctx context.Context, ps participantMap, assigned map[string]string) error {
	santas := make([]send.Santa, 0, len(ps))
	for address, name := range ps {
		recepient := assigned[address]
		recepientName := ps[recepient]
		santa := send.Santa{
			Name:      name,
			Addr:      address,
			Recipient: recepientName,
		}
		santas = append(santas, santa)
	}

	return s.sender.Send(ctx, santas)
}

func addresses(ps participantMap) []string {
	addresses := make([]string, 0, len(ps))
	for address := range ps {
		addresses = append(addresses, address)
	}

	return addresses
}

// pickSecretSanta takes a list of emails and returns a map where the key is
// the secret santa and the value is the person receiving the gift
func pickSecretSanta(ctx context.Context, emails []string, s Shuffler) map[string]string {
	_, span := trace.StartSpan(ctx, "pick_secret_santa")
	defer span.End()

	if len(emails) < 3 {
		panic(fmt.Sprintf("received %d participants; cannot have fewer than 3 participants", len(emails)))
	}

	receivers := make([]string, len(emails))
	copy(receivers, emails)

	s.Shuffle(receivers)

	assigned := assign(emails, receivers)
	return assigned
}

// assign takes a list of the input emails and the shuffled emails and
// assigns the secret santa to individuals. The returned map's key will be the
// person receiving the email (Secret Santa) and the value is the email of the
// individual receiving the gift.
func assign(santa, receiver []string) map[string]string {
	assignedMap := make(map[string]string, len(santa))
	for i := 0; i < len(santa); i++ {
		// if the person is their own secret santa, switch with the
		// next participant
		if santa[i] == receiver[i] {
			nextIndex := (i + 1) % len(receiver)

			receiver[i], receiver[nextIndex] = receiver[nextIndex], receiver[i]
			// i could be the last index so update the assigned map for index 0
			assignedMap[santa[nextIndex]] = receiver[nextIndex]

		}
		assignedMap[santa[i]] = receiver[i]
	}

	return assignedMap
}
