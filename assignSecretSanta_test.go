package secretsanta

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssing(t *testing.T) {
	testcases := []struct {
		descripton string
		original   []string
		shuffled   []string
	}{
		{
			descripton: "Perfectly shuffled",
			original:   []string{"a", "b", "c"},
			shuffled:   []string{"b", "c", "a"},
		},
		{
			descripton: "Partially shuffled",
			original:   []string{"a", "b", "c", "d"},
			shuffled:   []string{"b", "a", "c", "d"},
		},
		{
			descripton: "Unshuffled last index",
			original:   []string{"a", "b", "c", "d"},
			shuffled:   []string{"c", "a", "b", "d"},
		},
	}

	for _, tc := range testcases {
		tc := tc
		assigned := assign(tc.original, tc.shuffled)
		require.Equal(t, len(tc.original), len(assigned))

		for k, v := range assigned {
			assert.NotEqual(t, k, v)
		}
	}
}
