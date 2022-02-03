package secretsanta

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sachinagada/secretsanta/send"
)

type sender interface {
	Send(santas []send.Santa) error
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

func NewServer(shuf Shuffler, sender sender) *server {
	return &server{
		shuf:   shuf,
		sender: sender,
	}
}

func addresses(ps participantMap) []string {
	addresses := make([]string, 0, len(ps))
	for address := range ps {
		addresses = append(addresses, address)
	}

	return addresses
}

// ServeHTTP handles the requests from the form and assigns the secret santa to
// each participant. Each participant will receive communication regarding who
// they are assigned.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	addrs := addresses(partMap)

	// assigned key is the secret santa and the value is person receiving the
	// gift
	assigned := pickSecretSanta(addrs, s.shuf)

	if informErr := s.inform(partMap, assigned); informErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(informErr.Error()))
		return
	}

	w.Write([]byte("Santas have been chosen and informed. Happy Holidays!"))
}

func (s *server) inform(ps participantMap, assigned map[string]string) error {
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

	return s.sender.Send(santas)
}

// pickSecretSanta takes a list of emails and returns a map where the key is
// the secret santa and the value is the person receiving the gift
func pickSecretSanta(emails []string, s Shuffler) map[string]string {
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
