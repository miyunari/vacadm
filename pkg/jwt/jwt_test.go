package jwt

import (
	"testing"
	"time"

	"github.com/MninaTB/vacadm/pkg/model"
)

func TestTokenizer(t *testing.T) {
	const (
		testUUID0 = "d2446dd8-a360-404e-93e0-b559a19736ac"
		testUUID1 = "a05bbe22-36fd-40b0-b1fc-18d4b425f730"
	)

	tt := []struct {
		name            string
		tokenizer       *Tokenizer
		user            *model.User
		overwriteSecret []byte
		expectUserID    string
		expectTeamID    string
		expectErr       bool
	}{
		{
			name:      "expect",
			tokenizer: NewTokenizer([]byte("123"), 24*time.Hour),
			user: &model.User{
				ID:     testUUID0,
				TeamID: func() *string { tmp := testUUID1; return &tmp }(),
			},
			expectUserID: testUUID0,
			expectTeamID: testUUID1,
		},
		{
			name:      "missing secret",
			tokenizer: NewTokenizer([]byte(""), 24*time.Hour),
			user: &model.User{
				ID: testUUID0,
			},
			expectUserID: testUUID0,
			expectErr:    true,
		},
		{
			name:      "missing secret",
			tokenizer: NewTokenizer([]byte("123"), 24*time.Hour),
			user: &model.User{
				ID: testUUID0,
			},
			overwriteSecret: []byte("abc"),
			expectUserID:    testUUID0,
			expectErr:       true,
		},
		{
			name:      "token invalid",
			tokenizer: NewTokenizer([]byte("123"), -24*time.Hour),
			user: &model.User{
				ID: testUUID0,
			},
			expectUserID: testUUID0,
			expectErr:    true,
		},
		{
			name:      "userID is not a UUID",
			tokenizer: NewTokenizer([]byte("123"), 24*time.Hour),
			user: &model.User{
				ID: "abc",
			},
			expectErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			token, err := tc.tokenizer.Generate(tc.user)
			if err != nil && tc.expectErr {
				return
			} else if err != nil {
				t.Fatal(err)
			}
			if len(tc.overwriteSecret) != 0 {
				tc.tokenizer.hmacSecret = tc.overwriteSecret
			}
			userID, teamID, err := tc.tokenizer.Valid(token)
			if err != nil && tc.expectErr {
				return
			} else if err != nil {
				t.Fatal(err)
			}
			if tc.expectTeamID != teamID {
				t.Errorf("invalid teamID, want: %s, got: %s", tc.expectTeamID, teamID)
			}
			if tc.expectUserID != userID {
				t.Errorf("invalid userID, want: %s, got: %s", tc.expectUserID, userID)
			}
		})
	}
}
