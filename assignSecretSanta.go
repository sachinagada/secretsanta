package secretsanta

import (
	"fmt"

	"github.com/sachinagada/secretsanta/randshuffle"
)

// PickSecretSanta takes a list of emails and returns a map[string]string with
// the assigned secret santa
func PickSecretSanta(emails []string) (map[string]string, error) {
	if len(emails) < 2 {
		return nil, fmt.Errorf("Cannot have less than 2 participants")
	}

	r := &randshuffle.RandShuffle{}
	assigned := r.Shuffle(emails)

	secretSantas := assign(emails, assigned)
	return secretSantas, nil
}

// assign takes a list of the input emails and the shuffled emails and
// assigns the secret santa to individuals and ensures that no one has themselves
// as their own secret santa. The returned map's key will be the person receiving
// the email and the value is their assigned secret santa
func assign(emails, assigned []string) map[string]string {
	assignedMap := make(map[string]string, len(emails))
	for i := 0; i < len(emails); i++ {
		// if the person is their own secret santa, switch with the
		// next participant
		if emails[i] == assigned[i] {
			nextIndex := (i + 1) % len(assigned)
			// ensure you assign the next index to the current because i could
			// be the last index
			assignedMap[emails[nextIndex]] = assigned[i]
			assigned[i], assigned[nextIndex] = assigned[nextIndex], assigned[i]
		}
		assignedMap[emails[i]] = assigned[i]
	}

	return assignedMap
}
