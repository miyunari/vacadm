package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestVacation_Copy(t *testing.T) {
	now := time.Now()
	tt := []struct {
		name     string
		original *Vacation
	}{
		{
			name: "expected",
			original: &Vacation{
				ID:         "test-vacation-id",
				UserID:     "test-user-id",
				ApprovedBy: func() *string { str := "test-approvedBy-id"; return &str }(),
				From:       now.Add(time.Minute),
				To:         now.Add(time.Hour),
				CreatedAt:  func() *time.Time { tmp := now.Add(10 * time.Minute); return &tmp }(),
				DeletedAt:  func() *time.Time { tmp := now.Add(30 * time.Minute); return &tmp }(),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.original.Copy()
			if !cmp.Equal(tc.original, got) {
				t.Fatal(cmp.Diff(tc.original, got))
			}
			got.ID += "vacation-id"
			got.UserID = "user-id-vacation"
			got.ApprovedBy = nil
			got.From = time.Now()
			got.To = time.Now()
			got.CreatedAt = nil
			got.DeletedAt = nil
			if cmp.Equal(tc.original, got) {
				t.Fatal("copy should not be equal")
			}
		})
	}
}
