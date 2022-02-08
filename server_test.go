package secretsanta

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sachinagada/secretsanta/send"
	"github.com/sachinagada/secretsanta/shuffle"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeSender implements sender
type fakeSender struct {
	err error
}

// Send will return any errror for tests
func (f *fakeSender) Send(ctx context.Context, santas []send.Santa) error {
	for _, santa := range santas {
		if strings.Contains(santa.Recipient, "@") {
			return fmt.Errorf("error with recepient %q; expected name, got email address", santa.Recipient)
		}

		if strings.Contains(santa.Name, "@") {
			return fmt.Errorf("error with santa name %q; expected name, got email address", santa.Name)
		}

		if !strings.Contains(santa.Addr, "@") {
			return fmt.Errorf("error with santa address %q; invalid email address", santa.Addr)
		}
	}
	return f.err
}

func TestServer(t *testing.T) {
	shuffler := shuffle.Rand{}

	testCases := map[string]struct {
		participants []participant
		sender       sender
		expStatus    int
		expResponse  string
	}{
		"happy_case": {
			participants: []participant{
				{
					FirstName: "harry",
					LastName:  "potter",
					Email:     "harrypotter@gmail.com",
				}, {
					FirstName: "hermione",
					LastName:  "granger",
					Email:     "hermionegranger@gmail.com",
				}, {
					FirstName: "ron",
					LastName:  "weasley",
					Email:     "ronweasley@gmail.com",
				},
			},
			sender:      &fakeSender{},
			expStatus:   http.StatusOK,
			expResponse: "Santas have been chosen and informed. Happy Holidays!",
		}, "duplicate_emails": {
			participants: []participant{
				{
					FirstName: "harry",
					LastName:  "potter",
					Email:     "harrypotter@gmail.com",
				}, {
					FirstName: "hermione",
					LastName:  "granger",
					Email:     "harrypotter@gmail.com",
				}, {
					FirstName: "ron",
					LastName:  "weasley",
					Email:     "ronweasley@gmail.com",
				},
			},
			sender:      &fakeSender{},
			expStatus:   http.StatusBadRequest,
			expResponse: "Duplicate email address was used for some of the participants. Please check all the information and try again.",
		}, "error_sending": {
			participants: []participant{
				{
					FirstName: "harry",
					LastName:  "potter",
					Email:     "harrypotter@gmail.com",
				}, {
					FirstName: "hermione",
					LastName:  "granger",
					Email:     "hermionegranger@gmail.com",
				}, {
					FirstName: "ron",
					LastName:  "weasley",
					Email:     "ronweasley@gmail.com",
				},
			},
			sender: &fakeSender{
				err: fmt.Errorf("error sending to santa"),
			},
			expStatus:   http.StatusInternalServerError,
			expResponse: "error sending to santa",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			tc := testCase
			t.Parallel()

			b, marshErr := json.Marshal(tc.participants)
			require.NoError(t, marshErr)

			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
			w := httptest.NewRecorder()
			server := NewServer(&shuffler, tc.sender)
			server.ServeHTTP(w, req)

			resp := w.Result()
			assert.Equal(t, tc.expStatus, resp.StatusCode)
			body, readErr := io.ReadAll(resp.Body)
			require.NoError(t, readErr)
			assert.Equal(t, tc.expResponse, string(body))
		})
	}
}

func TestAssign(t *testing.T) {
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
