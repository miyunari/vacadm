package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestTeam_Copy(t *testing.T) {
	now := time.Now()
	tt := []struct {
		name     string
		original *Team
	}{
		{
			name: "expected",
			original: &Team{
				ID:        "test-team-id",
				OwnerID:   "test-owner-id",
				Name:      "test-team-name",
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
			got.ID += "team-id"
			got.OwnerID = "owner-id"
			got.Name = "team-name"
			got.CreatedAt = nil
			got.UpdatedAt = nil
			got.DeletedAt = nil
			if cmp.Equal(tc.original, got) {
				t.Fatal("copy should not be equal")
			}
		})
	}
}
