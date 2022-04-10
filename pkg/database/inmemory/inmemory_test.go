package inmemory

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"github.com/MninaTB/vacadm/pkg/model"
)

func TestInmemoryDB_CreateUser(t *testing.T) {
	tt := []struct {
		name      string
		user      *model.User
		userStore []*model.User
		userCount int
		wantErr   bool
	}{
		{
			name: "normal creation",
			user: &model.User{
				FirstName: "firstname",
				LastName:  "lastname",
				Email:     "admin@inform.de",
			},
			userCount: 1,
			wantErr:   false,
		},
		{
			name: "creation matches parent",
			userStore: []*model.User{
				{
					ID:        "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
					FirstName: "firstname-parent",
					LastName:  "lastname-parent",
					Email:     "admin@inform.de",
				},
			},
			user: &model.User{
				ParentID: func() *string {
					tmp := "f95128f7-733d-48b3-9306-cc5fe27cf6a5"
					return &tmp
				}(),
				FirstName: "firstname-child",
				LastName:  "lastname-child",
				Email:     "child-123@inform.de",
			},
			userCount: 2,
			wantErr:   false,
		},
		{
			name: "create user but parent does not exist",
			user: &model.User{
				ParentID: func() *string {
					tmp := "f95128f7-733d-48b3-9306-cc5fe27cf6a5"
					return &tmp
				}(),
				FirstName: "firstname-someone",
				LastName:  "lastname-someone",
				Email:     "someone@inform.de",
			},
			userCount: 0,
			wantErr:   true,
		},
		{
			name: "create user but email address is empty",
			user: &model.User{
				FirstName: "firstname-email",
				LastName:  "lastname-email",
				Email:     "",
			},
			userCount: 0,
			wantErr:   true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.userStore != nil {
				db.userStore = tc.userStore
			}
			newUser, err := db.CreateUser(context.Background(), tc.user)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if tc.userCount != len(db.userStore) {
				t.Fatalf("invalid number of users in store, want: %d, got: %d",
					tc.userCount, len(db.userStore),
				)
			}

			// NOTE: check if uuid is set
			_, err = uuid.Parse(newUser.ID)
			if err != nil {
				t.Error(err)
			}

			// NOTE: check if createdAt timestamp is set
			if newUser.CreatedAt == nil {
				t.Error("missing timestamp created_at")
			}

			ignoreFields := cmp.FilterPath(func(p cmp.Path) bool {
				return strings.Contains(p.String(), "ID") ||
					strings.Contains(p.String(), "CreatedAt")
			}, cmp.Ignore())

			// NOTE: compare original struct, ignore ID and CreatedAt (should be different)
			if !cmp.Equal(tc.user, newUser, ignoreFields) {
				t.Fatal(cmp.Diff(tc.user, newUser, ignoreFields))
			}
		})
	}
}
