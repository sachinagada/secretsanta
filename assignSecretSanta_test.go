package secretsanta

import (
	"testing"
)

func TestAssing(t *testing.T) {
	testcases := []struct {
		descripton string
		original   []string
		shuffled   []string
	}{
		{
			descripton: "already_shuffled",
			original:   []string{"a", "b", "c"},
			shuffled:   []string{"b", "c", "a"},
		},
		{
			descripton: "partially_shuffled",
			original:   []string{"a", "b", "c", "d"},
			shuffled:   []string{"b", "a", "c", "d"},
		},
		{
			descripton: "unshuffled_last_index",
			original:   []string{"a", "b", "c", "d"},
			shuffled:   []string{"c", "a", "b", "d"},
		},
		{
			descripton: "completely unshuffled",
			original:   []string{"a", "b", "c", "d"},
			shuffled:   []string{"a", "b", "c", "d"},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.descripton, func(t *testing.T) {
			t.Parallel()
			assigned := assign(tc.original, tc.shuffled)
			if len(tc.original) != len(assigned) {
				t.Fatalf("returned map has length %d and not equal to number of participants: %d", len(assigned), len(tc.original))
			}

			m := make(map[string]struct{}, len(tc.shuffled))
			for k, v := range assigned {
				if k == v {
					t.Errorf("%q is assigned to themselves", k)
				}

				if _, ok := m[v]; ok {
					t.Errorf("value %q already in the map. Value is assigned to multiple keys", v)
				}
				m[v] = struct{}{}
			}
		})
	}
}
