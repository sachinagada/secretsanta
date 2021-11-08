package shuffle

import "math/rand"

// Rand implements the Shuffle interface and uses the rand package to
// pseudo randomize the order in the array
type Rand struct{}

// Shuffle takes a list of emails and shuffles them using the Fisher-Yates
// algorithm
func (*Rand) Shuffle(participants []string) {
	// rand.Shuffle mutates the same slice
	rand.Shuffle(len(participants), func(i, j int) {
		participants[i], participants[j] = participants[j], participants[i]
	})
}
