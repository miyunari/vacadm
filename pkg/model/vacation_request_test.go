package model

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestVacationRequest_Copy(t *testing.T) {
	now := time.Now()
	tt := []struct {
		name     string
		original *VacationRequest
	}{
		{
			name: "expected",
			original: &VacationRequest{
				ID:        "test-vacation-resource-id",
				UserID:    "test-user-id",
				From:      now.Add(time.Minute),
				To:        now.Add(time.Hour),
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
			got.ID += "vacation-request-id"
			got.UserID = "user-id-request"
			got.From = time.Now()
			got.To = time.Now()
			got.CreatedAt = nil
			got.UpdatedAt = nil
			got.DeletedAt = nil
			if cmp.Equal(tc.original, got) {
				t.Fatal("copy should not be equal")
			}
		})
	}
}
