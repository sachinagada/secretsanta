package secretsanta

import (
	"fmt"
	"math/rand"
)

// SimpleAssign assigns secret santas by using the Fisher–Yates shuffle algorithm.
// It will shuffle the participants around and the participants in the two arrays
// at the same index are the assigned secret santas. It also validates that the
// same person hasn't been assigned to themselves as the secret
// santa. If they have, then swap with the participant next to them
func SimpleAssign(emails []string) (map[string]string, error) {
	if len(emails) < 2 {
		return nil, fmt.Errorf("Cannot have less than 2 participants")
	}

	// copy the array so we can compare with the original and know which
	// participant is mapped to whom
	assigned := make([]string, len(emails))
	assignedMap := make(map[string]string, len(emails))
	for i, email := range emails {
		assigned[i] = email
	}

	// shuffle the participants and the participants in the same index will
	// be the assigned receiver
	rand.Shuffle(len(assigned), func(i, j int) {
		assigned[i], assigned[j] = assigned[j], assigned[i]
	})

	for i := 0; i < len(emails); i++ {
		// if the person is their own secret santa, switch with the
		// next participant
		if emails[i] == assigned[i] {
			nextIndex := (i + 1) % len(assigned)
			assigned[i], assigned[nextIndex] = assigned[nextIndex], assigned[i]
		}

		assignedMap[emails[i]] = assigned[i]
	}

	return assignedMap, nil
}
