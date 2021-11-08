package secretsanta

import (
	"fmt"
)

// Shuffler shuffles the slice of participants which will be the receivers
// for secret santa
type Shuffler interface {
	Shuffle([]string)
}

// PickSecretSanta takes a list of emails and returns a map where the key is
// the secret santa and the value is the person receiving the gift
func PickSecretSanta(emails []string, s Shuffler) (map[string]string, error) {
	if len(emails) < 3 {
		return nil, fmt.Errorf("received %d participants; cannot have fewer than 3 participants", len(emails))
	}

	receivers := make([]string, len(emails))
	copy(receivers, emails)

	s.Shuffle(receivers)

	assigned := assign(emails, receivers)
	return assigned, nil
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
