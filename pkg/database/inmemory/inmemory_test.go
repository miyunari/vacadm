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

func TestInmemoryDB_CreateTeam(t *testing.T) {
	tt := []struct {
		name      string
		team      *model.Team
		userStore []*model.User
		teamStore []*model.Team
		teamCount int
		wantErr   bool
	}{
		{
			name: "owner not found",
			team: &model.Team{
				OwnerID: "owner-not-found",
				Name:    "A-Team",
			},
			teamCount: 0,
			wantErr:   true,
		},
		{
			name: "owner not found",
			userStore: []*model.User{
				{
					ID: "owner-found",
				},
			},
			team: &model.Team{
				OwnerID: "owner-found",
				Name:    "A-Team",
			},
			teamCount: 1,
			wantErr:   false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.userStore != nil {
				db.userStore = tc.userStore
			}
			if tc.teamStore != nil {
				db.teamStore = tc.teamStore
			}
			newTeam, err := db.CreateTeam(context.Background(), tc.team)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if tc.teamCount != len(db.teamStore) {
				t.Fatalf("invalid number of teams in store, want: %d, got: %d",
					tc.teamCount, len(db.teamStore),
				)
			}

			// NOTE: check if uuid is set
			_, err = uuid.Parse(newTeam.ID)
			if err != nil {
				t.Error(err)
			}

			// NOTE: check if createdAt timestamp is set
			if newTeam.CreatedAt == nil {
				t.Error("missing timestamp created_at")
			}

			ignoreFields := cmp.FilterPath(func(p cmp.Path) bool {
				return strings.Contains(p.String(), "ID") ||
					strings.Contains(p.String(), "CreatedAt")
			}, cmp.Ignore())

			// NOTE: compare original struct, ignore ID and CreatedAt (should be different)
			if !cmp.Equal(tc.team, newTeam, ignoreFields) {
				t.Fatal(cmp.Diff(tc.team, newTeam, ignoreFields))
			}
		})
	}
}

func TestInmemoryDB_GetTeamByID(t *testing.T) {
	tt := []struct {
		name      string
		teamID    string
		teamStore []*model.Team
		expect    *model.Team
		wantErr   bool
	}{
		{
			name:    "team does not exist",
			teamID:  "does-not-exist",
			wantErr: true,
		},
		{
			name: "get team by id as expected",
			teamStore: []*model.Team{
				{
					ID:   "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
					Name: "A-Team",
				},
			},
			teamID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
			expect: &model.Team{
				ID:   "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				Name: "A-Team",
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.teamStore != nil {
				db.teamStore = tc.teamStore
			}
			newTeam, err := db.GetTeamByID(context.Background(), tc.teamID)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if !cmp.Equal(tc.expect, newTeam) {
				t.Fatal(cmp.Diff(tc.expect, newTeam))
			}
		})
	}
}

func TestInmemoryDB_ListTeams(t *testing.T) {
	tt := []struct {
		name      string
		teamStore []*model.Team
		wantErr   bool
	}{
		{
			name:    "empty store",
			wantErr: false,
		},
		{
			name: "list users as expected",
			teamStore: []*model.Team{
				{
					ID:   "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
					Name: "A-Team",
				},
				{
					ID:   "fed75474-29df-4d99-a792-09f0bf7ae848",
					Name: "B-Team",
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.teamStore != nil {
				db.teamStore = tc.teamStore
			}
			teams, err := db.ListTeams(context.Background())
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if !cmp.Equal(db.teamStore, teams) {
				t.Fatal(cmp.Diff(db.teamStore, teams))
			}
		})
	}
}

func TestInmemoryDB_ListTeamUsers(t *testing.T) {
	tt := []struct {
		name        string
		teamID      string
		userStore   []*model.User
		teamStore   []*model.Team
		expectUsers int
		wantErr     bool
	}{
		{
			name:        "empty store",
			expectUsers: 0,
			wantErr:     false,
		},
		{
			name:        "list users as expected",
			expectUsers: 2,
			teamStore: []*model.Team{
				{
					ID: "d4ee305c-18cc-4f1e-a752-764b6913ab67",
				},
			},
			teamID: "d4ee305c-18cc-4f1e-a752-764b6913ab67",
			userStore: []*model.User{
				{
					ID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
					TeamID: func() *string {
						tmp := "d4ee305c-18cc-4f1e-a752-764b6913ab67"
						return &tmp
					}(),
					FirstName: "firstname-one",
					LastName:  "lastname-one",
					Email:     "one@inform.de",
				},
				{
					ID: "fed75474-29df-4d99-a792-09f0bf7ae848",
					TeamID: func() *string {
						tmp := "d4ee305c-18cc-4f1e-a752-764b6913ab67"
						return &tmp
					}(),
					FirstName: "firstname-two",
					LastName:  "lastname-two",
					Email:     "two@inform.de",
				},
				{
					ID:        "23ebe54c-2d91-4a00-8e13-6b12a1a47000",
					FirstName: "firstname-other-team",
					LastName:  "lastname-other-team",
					Email:     "other@inform.de",
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
			if tc.teamStore != nil {
				db.teamStore = tc.teamStore
			}
			users, err := db.ListTeamUsers(context.Background(), tc.teamID)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if tc.expectUsers != len(users) {
				t.Fatalf("invalid number of users, want: %d, got: %d", tc.expectUsers, len(users))
			}
		})
	}
}

func TestInmemoryDB_UpdateTeam(t *testing.T) {
	tt := []struct {
		name      string
		team      *model.Team
		teamStore []*model.Team
		wantErr   bool
	}{
		{
			name: "team does not exist",
			team: &model.Team{
				Name: "team-update",
			},
			wantErr: true,
		},
		{
			name: "update expected",
			teamStore: []*model.Team{
				{
					ID:   "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
					Name: "team-existing",
				},
			},
			team: &model.Team{
				ID:   "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				Name: "team-new",
			},
			wantErr: false,
		},
		{
			name: "update team but owner does not exist",
			team: &model.Team{
				OwnerID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				Name:    "owner-missing",
			},
			wantErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.teamStore != nil {
				db.teamStore = tc.teamStore
			}
			newTeam, err := db.UpdateTeam(context.Background(), tc.team)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			// NOTE: check if createdAt timestamp is set
			if newTeam.UpdatedAt == nil {
				t.Error("missing timestamp updated_at")
			}

			ignoreFields := cmp.FilterPath(func(p cmp.Path) bool {
				return strings.Contains(p.String(), "UpdatedAt")
			}, cmp.Ignore())

			// NOTE: compare original struct, ignore ID and CreatedAt (should be different)
			if !cmp.Equal(tc.team, newTeam, ignoreFields) {
				t.Fatal(cmp.Diff(tc.team, newTeam, ignoreFields))
			}
		})
	}
}

func TestInmemoryDB_DeleteTeam(t *testing.T) {
	tt := []struct {
		name      string
		teamID    string
		teamStore []*model.Team
		wantErr   bool
	}{
		{
			name:    "team does not exist",
			teamID:  "does-not-exist",
			wantErr: true,
		},
		{
			name:   "get user by id as expected",
			teamID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
			teamStore: []*model.Team{
				{
					ID:   "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
					Name: "A-Team",
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.teamStore != nil {
				db.teamStore = tc.teamStore
			}
			expectCount := len(tc.teamStore) - 1
			err := db.DeleteTeam(context.Background(), tc.teamID)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if expectCount != len(db.userStore) {
				t.Fatalf("invalid count, want: %d, got: %d", expectCount, len(db.teamStore))
			}
		})
	}
}

func TestInmemoryDB_CreateVacation(t *testing.T) {
	tt := []struct {
		name          string
		vacation      *model.Vacation
		vacationStore []*model.Vacation
		vacationCount int
		wantErr       bool
	}{
		{
			name:          "missing userID",
			vacation:      &model.Vacation{},
			vacationCount: 1,
			wantErr:       true,
		},
		{
			name: "missing approver",
			vacation: &model.Vacation{
				UserID: "abc",
			},
			vacationCount: 1,
			wantErr:       true,
		},
		{
			name: "creation expected",
			vacationStore: []*model.Vacation{
				{
					ID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				},
			},
			vacation: &model.Vacation{
				UserID:     "abc",
				ApprovedBy: func() *string { tmp := "xyz"; return &tmp }(),
			},
			vacationCount: 2,
			wantErr:       false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.vacationStore != nil {
				db.vacationStore = tc.vacationStore
			}
			newVacation, err := db.CreateVacation(context.Background(), tc.vacation)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if tc.vacationCount != len(db.vacationStore) {
				t.Fatalf("invalid number of vacations in store, want: %d, got: %d",
					tc.vacationCount, len(db.vacationStore),
				)
			}

			// NOTE: check if uuid is set
			_, err = uuid.Parse(newVacation.ID)
			if err != nil {
				t.Error(err)
			}

			// NOTE: check if createdAt timestamp is set
			if newVacation.CreatedAt == nil {
				t.Error("missing timestamp created_at")
			}

			ignoreFields := cmp.FilterPath(func(p cmp.Path) bool {
				return strings.Contains(p.String(), "ID") ||
					strings.Contains(p.String(), "CreatedAt")
			}, cmp.Ignore())

			// NOTE: compare original struct, ignore ID and CreatedAt (should be different)
			if !cmp.Equal(tc.vacation, newVacation, ignoreFields) {
				t.Fatal(cmp.Diff(tc.vacation, newVacation, ignoreFields))
			}
		})
	}
}

func TestInmemoryDB_GetVacationByID(t *testing.T) {
	tt := []struct {
		name          string
		vacationID    string
		vacationStore []*model.Vacation
		expect        *model.Vacation
		wantErr       bool
	}{
		{
			name:       "vacation does not exist",
			vacationID: "does-not-exist",
			wantErr:    true,
		},
		{
			name: "get vacation by id as expected",
			vacationStore: []*model.Vacation{
				{
					ID:     "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
					UserID: "some-user-id",
				},
			},
			vacationID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
			expect: &model.Vacation{
				ID:     "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				UserID: "some-user-id",
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.vacationStore != nil {
				db.vacationStore = tc.vacationStore
			}
			newVacation, err := db.GetVacationByID(context.Background(), tc.vacationID)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if !cmp.Equal(tc.expect, newVacation) {
				t.Fatal(cmp.Diff(tc.expect, newVacation))
			}
		})
	}
}

func TestInmemoryDB_ListVacations(t *testing.T) {
	tt := []struct {
		name          string
		vacationStore []*model.Vacation
		wantErr       bool
	}{
		{
			name:    "empty store",
			wantErr: false,
		},
		{
			name: "list vacations as expected",
			vacationStore: []*model.Vacation{
				{
					ID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				},
				{
					ID: "fed75474-29df-4d99-a792-09f0bf7ae848",
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.vacationStore != nil {
				db.vacationStore = tc.vacationStore
			}
			vacations, err := db.ListVacations(context.Background())
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if !cmp.Equal(db.vacationStore, vacations) {
				t.Fatal(cmp.Diff(db.vacationStore, vacations))
			}
		})
	}
}

func TestInmemoryDB_DeleteVacation(t *testing.T) {
	tt := []struct {
		name          string
		vacationID    string
		vacationStore []*model.Vacation
		wantErr       bool
	}{
		{
			name:       "vacation does not exist",
			vacationID: "does-not-exist",
			wantErr:    true,
		},
		{
			name: "get vacation by id as expected",
			vacationStore: []*model.Vacation{
				{
					ID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				},
			},
			vacationID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
			wantErr:    false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.vacationStore != nil {
				db.vacationStore = tc.vacationStore
			}
			expectCount := len(tc.vacationStore) - 1
			err := db.DeleteVacation(context.Background(), tc.vacationID)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if expectCount != len(db.vacationStore) {
				t.Fatalf("invalid count, want: %d, got: %d", expectCount, len(db.vacationStore))
			}
		})
	}
}

func TestInmemoryDB_CreateVacationRequest(t *testing.T) {
	tt := []struct {
		name                 string
		vacationRequest      *model.VacationRequest
		vacationRequestStore []*model.VacationRequest
		vacationRequestCount int
		wantErr              bool
	}{
		{
			name:                 "missing userID",
			vacationRequest:      &model.VacationRequest{},
			vacationRequestCount: 1,
			wantErr:              true,
		},
		{
			name: "creation expected",
			vacationRequestStore: []*model.VacationRequest{
				{
					ID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				},
			},
			vacationRequest: &model.VacationRequest{
				UserID: "abc",
			},
			vacationRequestCount: 2,
			wantErr:              false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.vacationRequestStore != nil {
				db.vacationRequestStore = tc.vacationRequestStore
			}
			newVacationRequest, err := db.CreateVacationRequest(context.Background(), tc.vacationRequest)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if tc.vacationRequestCount != len(db.vacationRequestStore) {
				t.Fatalf("invalid number of vacationRequests in store, want: %d, got: %d",
					tc.vacationRequestCount, len(db.vacationRequestStore),
				)
			}

			// NOTE: check if uuid is set
			_, err = uuid.Parse(newVacationRequest.ID)
			if err != nil {
				t.Error(err)
			}

			// NOTE: check if createdAt timestamp is set
			if newVacationRequest.CreatedAt == nil {
				t.Error("missing timestamp created_at")
			}

			ignoreFields := cmp.FilterPath(func(p cmp.Path) bool {
				return strings.Contains(p.String(), "ID") ||
					strings.Contains(p.String(), "CreatedAt")
			}, cmp.Ignore())

			// NOTE: compare original struct, ignore ID and CreatedAt (should be different)
			if !cmp.Equal(tc.vacationRequest, newVacationRequest, ignoreFields) {
				t.Fatal(cmp.Diff(tc.vacationRequest, newVacationRequest, ignoreFields))
			}
		})
	}
}

func TestInmemoryDB_GetVacationRequestByID(t *testing.T) {
	tt := []struct {
		name                 string
		vacationRequestID    string
		vacationRequestStore []*model.VacationRequest
		expect               *model.VacationRequest
		wantErr              bool
	}{
		{
			name:              "vacationRequest does not exist",
			vacationRequestID: "does-not-exist",
			wantErr:           true,
		},
		{
			name: "get vacationRequest by id as expected",
			vacationRequestStore: []*model.VacationRequest{
				{
					ID:     "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
					UserID: "valid-user-id",
				},
			},
			vacationRequestID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
			expect: &model.VacationRequest{
				ID:     "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				UserID: "valid-user-id",
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.vacationRequestStore != nil {
				db.vacationRequestStore = tc.vacationRequestStore
			}
			newVacationRequest, err := db.GetVacationRequestByID(context.Background(), tc.vacationRequestID)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if !cmp.Equal(tc.expect, newVacationRequest) {
				t.Fatal(cmp.Diff(tc.expect, newVacationRequest))
			}
		})
	}
}

func TestInmemoryDB_ListVacationRequests(t *testing.T) {
	tt := []struct {
		name                 string
		vacationRequestStore []*model.VacationRequest
		wantErr              bool
	}{
		{
			name:    "empty store",
			wantErr: false,
		},
		{
			name: "list vacationRequests as expected",
			vacationRequestStore: []*model.VacationRequest{
				{
					ID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				},
				{
					ID: "fed75474-29df-4d99-a792-09f0bf7ae848",
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.vacationRequestStore != nil {
				db.vacationRequestStore = tc.vacationRequestStore
			}
			vacationRequests, err := db.ListVacationRequests(context.Background())
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if !cmp.Equal(db.vacationRequestStore, vacationRequests) {
				t.Fatal(cmp.Diff(db.vacationRequestStore, vacationRequests))
			}
		})
	}
}

func TestInmemoryDB_DeleteVacationRequest(t *testing.T) {
	tt := []struct {
		name                 string
		vacationRequestID    string
		vacationRequestStore []*model.VacationRequest
		wantErr              bool
	}{
		{
			name:              "vacationRequest does not exist",
			vacationRequestID: "does-not-exist",
			wantErr:           true,
		},
		{
			name: "get vacationRequest by id as expected",
			vacationRequestStore: []*model.VacationRequest{
				{
					ID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				},
			},
			vacationRequestID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
			wantErr:           false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.vacationRequestStore != nil {
				db.vacationRequestStore = tc.vacationRequestStore
			}
			expectCount := len(tc.vacationRequestStore) - 1
			err := db.DeleteVacationRequest(context.Background(), tc.vacationRequestID)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if expectCount != len(db.vacationRequestStore) {
				t.Fatalf("invalid count, want: %d, got: %d", expectCount, len(db.vacationRequestStore))
			}
		})
	}
}

func TestInmemoryDB_CreateVacationResource(t *testing.T) {
	tt := []struct {
		name                   string
		vacationResource      *model.VacationResource
		vacationResourceStore []*model.VacationResource
		vacationResourceCount int
		wantErr                bool
	}{
		{
			name:                   "missing userID",
			vacationResource:      &model.VacationResource{},
			vacationResourceCount: 1,
			wantErr:                true,
		},
		{
			name: "creation matches parent",
			vacationResourceStore: []*model.VacationResource{
				{
					ID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				},
			},
			vacationResource: &model.VacationResource{
				UserID: "some-user-id",
			},
			vacationResourceCount: 2,
			wantErr:                false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.vacationResourceStore != nil {
				db.vacationResourceStore = tc.vacationResourceStore
			}
			newVacationResource, err := db.CreateVacationResource(context.Background(), tc.vacationResource)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if tc.vacationResourceCount != len(db.vacationResourceStore) {
				t.Fatalf("invalid number of vacationResource in store, want: %d, got: %d",
					tc.vacationResourceCount, len(db.vacationResourceStore),
				)
			}

			// NOTE: check if uuid is set
			_, err = uuid.Parse(newVacationResource.ID)
			if err != nil {
				t.Error(err)
			}

			// NOTE: check if createdAt timestamp is set
			if newVacationResource.CreatedAt == nil {
				t.Error("missing timestamp created_at")
			}

			ignoreFields := cmp.FilterPath(func(p cmp.Path) bool {
				return strings.Contains(p.String(), "ID") ||
					strings.Contains(p.String(), "CreatedAt")
			}, cmp.Ignore())

			// NOTE: compare original struct, ignore ID and CreatedAt (should be different)
			if !cmp.Equal(tc.vacationResource, newVacationResource, ignoreFields) {
				t.Fatal(cmp.Diff(tc.vacationResource, newVacationResource, ignoreFields))
			}
		})
	}
}

func TestInmemoryDB_GetVacationResourceByID(t *testing.T) {
	tt := []struct {
		name                   string
		vacationResourceID    string
		vacationResourceStore []*model.VacationResource
		expect                 *model.VacationResource
		wantErr                bool
	}{
		{
			name:                "vacationResource does not exist",
			vacationResourceID: "does-not-exist",
			wantErr:             true,
		},
		{
			name: "get vacationResource by id as expected",
			vacationResourceStore: []*model.VacationResource{
				{
					ID:     "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
					UserID: "valid-user-id",
				},
			},
			vacationResourceID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
			expect: &model.VacationResource{
				ID:     "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				UserID: "valid-user-id",
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.vacationResourceStore != nil {
				db.vacationResourceStore = tc.vacationResourceStore
			}
			newVacationResource, err := db.GetVacationResourceByID(context.Background(), tc.vacationResourceID)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if !cmp.Equal(tc.expect, newVacationResource) {
				t.Fatal(cmp.Diff(tc.expect, newVacationResource))
			}
		})
	}
}

func TestInmemoryDB_ListVacationResource(t *testing.T) {
	tt := []struct {
		name                   string
		vacationResourceStore []*model.VacationResource
		wantErr                bool
	}{
		{
			name:    "empty store",
			wantErr: false,
		},
		{
			name: "list vacationResources as expected",
			vacationResourceStore: []*model.VacationResource{
				{
					ID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				},
				{
					ID: "fed75474-29df-4d99-a792-09f0bf7ae848",
				},
			},
			wantErr: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.vacationResourceStore != nil {
				db.vacationResourceStore = tc.vacationResourceStore
			}
			vacationResources, err := db.ListVacationResource(context.Background())
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if !cmp.Equal(db.vacationResourceStore, vacationResources) {
				t.Fatal(cmp.Diff(db.vacationResourceStore, vacationResources))
			}
		})
	}
}

func TestInmemoryDB_DeleteVacationResource(t *testing.T) {
	tt := []struct {
		name                   string
		vacationResourceID    string
		vacationResourceStore []*model.VacationResource
		wantErr                bool
	}{
		{
			name:                "vacationResource does not exist",
			vacationResourceID: "does-not-exist",
			wantErr:             true,
		},
		{
			name: "get vacationResource by id as expected",
			vacationResourceStore: []*model.VacationResource{
				{
					ID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
				},
			},
			vacationResourceID: "f95128f7-733d-48b3-9306-cc5fe27cf6a5",
			wantErr:             false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			db := NewInmemoryDB()
			if tc.vacationResourceStore != nil {
				db.vacationResourceStore = tc.vacationResourceStore
			}
			expectCount := len(tc.vacationResourceStore) - 1
			err := db.DeleteVacationResource(context.Background(), tc.vacationResourceID)
			if err != nil && !tc.wantErr {
				t.Fatal(err)
			} else if err != nil && tc.wantErr {
				return
			}

			if expectCount != len(db.vacationResourceStore) {
				t.Fatalf("invalid count, want: %d, got: %d", expectCount, len(db.vacationResourceStore))
			}
		})
	}
}
