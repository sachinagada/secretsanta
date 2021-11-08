package shuffle

import (
	"testing"
)

func TestRand(t *testing.T) {
	r := &Rand{}
	s := []string{"a", "b", "c", "d"}

	c := make([]string, len(s))
	copy(c, s)

	r.Shuffle(c)

	for j, v := range s {
		if v != c[j] && c[j] != "" {
			// at least one of the values doesn't match
			return
		}
	}

	t.Fatalf("all values in the shuffled slice matched. slice wasn't shuffled")
}
