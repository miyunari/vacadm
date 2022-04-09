package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestUser_Copy(t *testing.T) {
	now := time.Now()
	tt := []struct {
		name     string
		original *User
	}{
		{
			name: "expected",
			original: &User{
				ID:        "test-user-id",
				ParentID:  func() *string { str := "test-parent-id"; return &str }(),
				TeamID:    func() *string { str := "test-team-id"; return &str }(),
				FirstName: "test-firstname",
				LastName:  "test-lastname",
				Email:     "test-email",
				CreatedAt: func() *time.Time { tmp := now.Add(10 * time.Minute); return &tmp }(),
				UpdatedAt: func() *time.Time { tmp := now.Add(15 * time.Minute); return &tmp }(),
				DeletedAt: func() *time.Time { tmp := now.Add(30 * time.Minute); return &tmp }(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.original.Copy()
			if !cmp.Equal(tc.original, got) {
				t.Fatal(cmp.Diff(tc.original, got))
			}
			got.ID += "user-id"
			got.ParentID = nil
			got.TeamID = nil
			got.FirstName = "firstname"
			got.LastName = "lastname"
			got.Email = "email"
			got.CreatedAt = nil
			got.UpdatedAt = nil
			got.DeletedAt = nil
			if cmp.Equal(tc.original, got) {
				t.Fatal("copy should not be equal")
			}
		})
	}
}
