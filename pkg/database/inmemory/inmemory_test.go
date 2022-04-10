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

func TestInmemoryDB_GetUserByID(t *testing.T) {
	tt := []struct {
		name      string
		userID    string
		userStore []*model.User
		expect    *model.User
		wantErr   bool
	}{
		{
			name:    "user does not exist",
			userID:  "does-not-exist",
			wantErr: true,
		},
		{
			name: "get user by id as expected",
			userStore: []*model.User{
				{
					ID:        "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
					FirstName: "firstname-found",
					LastName:  "lastname-found",
					Email:     "found@inform.de",
				},
			},
			userID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
			expect: &model.User{
				ID:        "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				FirstName: "firstname-found",
				LastName:  "lastname-found",
				Email:     "found@inform.de",
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.userStore != nil {
				db.userStore = tc.userStore
			}
			newUser, err := db.GetUserByID(context.Background(), tc.userID)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if !cmp.Equal(tc.expect, newUser) {
				t.Fatal(cmp.Diff(tc.expect, newUser))
			}
		})
	}
}

func TestInmemoryDB_ListUsers(t *testing.T) {
	tt := []struct {
		name      string
		userStore []*model.User
		wantErr   bool
	}{
		{
			name:    "empty store",
			wantErr: false,
		},
		{
			name: "list users as expected",
			userStore: []*model.User{
				{
					ID:        "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
					FirstName: "firstname-one",
					LastName:  "lastname-one",
					Email:     "one@inform.de",
				},
				{
					ID:        "fed75474-29df-4d99-a792-09f0bf7ae848",
					FirstName: "firstname-two",
					LastName:  "lastname-two",
					Email:     "two@inform.de",
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.userStore != nil {
				db.userStore = tc.userStore
			}
			users, err := db.ListUsers(context.Background())
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if !cmp.Equal(db.userStore, users) {
				t.Fatal(cmp.Diff(db.userStore, users))
			}
		})
	}
}

func TestInmemoryDB_UpdateUser(t *testing.T) {
	tt := []struct {
		name      string
		user      *model.User
		userStore []*model.User
		wantErr   bool
	}{
		{
			name: "user does not exist",
			user: &model.User{
				FirstName: "firstname-update",
				LastName:  "lastname-update",
				Email:     "admin@inform.de",
			},
			wantErr: true,
		},
		{
			name: "update expected",
			userStore: []*model.User{
				{
					ID:        "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
					FirstName: "firstname-existing",
					LastName:  "lastname-existing",
					Email:     "admin@inform.de",
				},
			},
			user: &model.User{
				ID:        "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				FirstName: "firstname-new",
				LastName:  "lastname-new",
				Email:     "new-123@inform.de",
			},
			wantErr: false,
		},
		{
			name: "update user but parent does not exist",
			user: &model.User{
				ParentID: func() *string {
					tmp := "f95128f7-733d-48b3-9306-cc5fe27cf6a5"
					return &tmp
				}(),
				FirstName: "firstname-missing",
				LastName:  "lastname-missing",
				Email:     "missing@inform.de",
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.userStore != nil {
				db.userStore = tc.userStore
			}
			newUser, err := db.UpdateUser(context.Background(), tc.user)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			// NOTE: check if uuid is set
			_, err = uuid.Parse(newUser.ID)
			if err != nil {
				t.Error(err)
			}

			// NOTE: check if createdAt timestamp is set
			if newUser.UpdatedAt == nil {
				t.Error("missing timestamp updated_at")
			}

			ignoreFields := cmp.FilterPath(func(p cmp.Path) bool {
				return strings.Contains(p.String(), "ID") ||
					strings.Contains(p.String(), "UpdatedAt")
			}, cmp.Ignore())

			// NOTE: compare original struct, ignore ID and CreatedAt (should be different)
			if !cmp.Equal(tc.user, newUser, ignoreFields) {
				t.Fatal(cmp.Diff(tc.user, newUser, ignoreFields))
			}
		})
	}
}

func TestInmemoryDB_DeleteUser(t *testing.T) {
	tt := []struct {
		name      string
		userID    string
		userStore []*model.User
		wantErr   bool
	}{
		{
			name:    "user does not exist",
			userID:  "does-not-exist",
			wantErr: true,
		},
		{
			name: "get user by id as expected",
			userStore: []*model.User{
				{
					ID:        "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
					FirstName: "firstname-found",
					LastName:  "lastname-found",
					Email:     "found@inform.de",
				},
			},
			userID:  "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.userStore != nil {
				db.userStore = tc.userStore
			}
			expectCount := len(tc.userStore) - 1
			err := db.DeleteUser(context.Background(), tc.userID)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if expectCount != len(db.userStore) {
				t.Fatalf("invalid count, want: %d, got: %d", expectCount, len(db.userStore))
			}
		})
	}
}
