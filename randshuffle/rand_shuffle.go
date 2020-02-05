package randshuffle

import "math/rand"

// RandShuffle implements the Shuffle interface and implements the shuffle
// from the rand package
type RandShuffle struct{}

// Shuffle takes a list of emails and shuffles them using the Fisher-Yates
// algorithm
func (*RandShuffle) Shuffle(emails []string) []string {
	// rand.Shuffle mutates the same slice so make a copy to maintain the order
	// of the original email slice
	assigned := make([]string, len(emails))
	copy(assigned, emails)

	rand.Shuffle(len(assigned), func(i, j int) {
		assigned[i], assigned[j] = assigned[j], assigned[i]
	})
	return assigned
}

// Type returns the type of shuffle
func Type() string {
	return "Rand.Shuffle"
}
