package inmemory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/MninaTB/vacadm/pkg/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// NewInmemoryDB returns initialized inmemoryDB that fulfills
// the database interface.
func NewInmemoryDB() *InmemoryDB {
	return &InmemoryDB{
		userStore:              make([]*model.User, 0),
		teamStore:              make([]*model.Team, 0),
		vacationStore:          make([]*model.Vacation, 0),
		vacationRequestStore:   make([]*model.VacationRequest, 0),
		vacationRessourceStore: make([]*model.VacationRessource, 0),
		logger:                 logrus.New().WithField("component", "inmemoryDB"),
	}
}

// InmemoryDB is a threadsafe inmemory database implementation.
type InmemoryDB struct {
	muUserStore sync.Mutex
	userStore   []*model.User

	muTeamStore sync.Mutex
	teamStore   []*model.Team

	muVacationStore sync.Mutex
	vacationStore   []*model.Vacation

	muVacationRequestStore sync.Mutex
	vacationRequestStore   []*model.VacationRequest

	muVacationRessourceStore sync.Mutex
	vacationRessourceStore   []*model.VacationRessource

	logger logrus.FieldLogger
}

// CreateUser stores an internal copy of the given user, if email address is
// not already in use, given parentID and/or teamID exists.
// Returns copy with assigned userID.
func (i *InmemoryDB) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	i.muUserStore.Lock()
	defer i.muUserStore.Unlock()
	if user.Email == "" {
		return nil, fmt.Errorf("missing email address")
	}
	for _, e := range i.userStore {
		if e.Email == user.Email {
			return nil, fmt.Errorf("email address must be unique")
		}
	}

	if user.ParentID != nil {
		var foundParent bool
		for _, p := range i.userStore {
			if p.ID == *user.ParentID {
				foundParent = true
				break
			}
		}
		if !foundParent {
			return nil, fmt.Errorf("missing parent with id %s", *user.ParentID)
		}
	}

	if user.TeamID != nil {
		_, err := i.GetTeamByID(ctx, *user.TeamID)
		if err != nil {
			return nil, err
		}
	}

	createdAt := time.Now()
	user.CreatedAt = &createdAt
	user.ID = uuid.NewString()
	usrCopy := user.Copy()

	i.logger.Info("create user with id: ", user.ID)
	i.userStore = append(i.userStore, usrCopy)
	return user, nil
}

// GetUserByID returns the associated user by the given id.
func (i *InmemoryDB) GetUserByID(_ context.Context, id string) (*model.User, error) {
	i.muUserStore.Lock()
	defer i.muUserStore.Unlock()
	for _, s := range i.userStore {
		if s.ID == id {
			i.logger.Info("get user with id: ", s.ID)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no user found")
	return nil, errors.New("no user found")
}

// ListUsers returns a copy of the internal user list.
func (i *InmemoryDB) ListUsers(_ context.Context) ([]*model.User, error) {
	i.muUserStore.Lock()
	defer i.muUserStore.Unlock()
	i.logger.Info("get list of users")

	userStore := make([]*model.User, len(i.userStore))
	for j, u := range i.userStore {
		userStore[j] = u.Copy()
	}
	return userStore, nil
}

// UpdateUser updates user entry by the given user.
func (i *InmemoryDB) UpdateUser(_ context.Context, user *model.User) (*model.User, error) {
	i.muUserStore.Lock()
	defer i.muUserStore.Unlock()
	updatededAt := time.Now()
	for x := 0; x < len(i.userStore); x++ {
		if i.userStore[x].ID != user.ID {
			continue
		}
		if user.Email != "" {
			i.userStore[x].Email = user.Email
		}
		if user.ParentID != nil {
			var found bool
			for _, u := range i.userStore {
				if u.ID == *user.ParentID {
					i.userStore[x].ParentID = user.ParentID
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("parent with id: '%s' not found", *user.ParentID)
			}
		}
		if user.FirstName != "" {
			i.userStore[x].FirstName = user.FirstName
		}
		if user.LastName != "" {
			i.userStore[x].LastName = user.LastName
		}
		i.userStore[x].UpdatedAt = &updatededAt
		i.logger.Info("update user with id: ", user.ID)
		return i.userStore[x].Copy(), nil
	}
	i.logger.Error("update failed: no user found")
	return nil, errors.New("update failed: no user found")
}

// DeleteUser removes user entry by the given id.
func (i *InmemoryDB) DeleteUser(_ context.Context, id string) error {
	i.muUserStore.Lock()
	defer i.muUserStore.Unlock()
	for x, user := range i.userStore {
		if user.ID == id {
			i.logger.Info("deleted user with id: ", id)
			i.userStore = append(i.userStore[:x], i.userStore[x+1:]...)
			return nil
		}
	}
	i.logger.Error("user didn't exist")
	return errors.New("user didn't exist")
}

// CreateTeam stores an internal copy of the given team.
// Returns copy with assigned teamID.
func (i *InmemoryDB) CreateTeam(_ context.Context, team *model.Team) (*model.Team, error) {
	i.muTeamStore.Lock()
	defer i.muTeamStore.Unlock()
	createdAt := time.Now()
	team.CreatedAt = &createdAt
	team.ID = uuid.NewString()
	teamCopy := team.Copy()

	i.logger.Info("create team with id: ", team.ID)
	i.teamStore = append(i.teamStore, teamCopy)
	return team, nil
}

// GetTeamByID returns the associated team by the given id.
func (i *InmemoryDB) GetTeamByID(_ context.Context, id string) (*model.Team, error) {
	i.muTeamStore.Lock()
	defer i.muTeamStore.Unlock()
	for _, s := range i.teamStore {
		if s.ID == id {
			i.logger.Info("get team with id: ", s.ID)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no team found")
	return nil, errors.New("no team found")
}

// ListTeams returns a copy of the internal team list.
func (i *InmemoryDB) ListTeams(_ context.Context) ([]*model.Team, error) {
	i.muTeamStore.Lock()
	defer i.muTeamStore.Unlock()
	i.logger.Info("get list of teams")
	teamStore := make([]*model.Team, len(i.teamStore))
	for j, t := range i.teamStore {
		teamStore[j] = t.Copy()
	}
	return teamStore, nil
}

// ListTeamUsers returns a list of users associated by the given teamID
func (i *InmemoryDB) ListTeamUsers(ctx context.Context, teamID string) ([]*model.User, error) {
	i.muTeamStore.Lock()
	defer i.muTeamStore.Unlock()
	var users []*model.User
	allUsers, err := i.ListUsers(ctx)
	if err != nil {
		return nil, err
	}
	for _, u := range allUsers {
		if u.TeamID != nil && *u.TeamID == teamID {
			users = append(users, u.Copy())
		}
	}
	return users, nil
}

// UpdateTeam updates team entry by the given team.
func (i *InmemoryDB) UpdateTeam(_ context.Context, team *model.Team) (*model.Team, error) {
	i.muTeamStore.Lock()
	defer i.muTeamStore.Unlock()
	updatededAt := time.Now()
	for x := 0; x < len(i.teamStore); x++ {
		if i.teamStore[x].ID == team.ID {
			i.teamStore[x].Name = team.Name
			i.teamStore[x].UpdatedAt = &updatededAt
			i.logger.Info("update team with id: ", team.ID)
			return i.teamStore[x], nil
		}
	}
	i.logger.Info("update failed: no team found")
	return nil, errors.New("update failed: no team found")
}

// DeleteTeam removes team entry by the given id.
func (i *InmemoryDB) DeleteTeam(_ context.Context, id string) error {
	i.muTeamStore.Lock()
	defer i.muTeamStore.Unlock()
	for x, team := range i.teamStore {
		if team.ID == id {
			i.logger.Info("delete team with id: ", id)
			i.teamStore = append(i.teamStore[:x], i.teamStore[x+1:]...)
			return nil
		}
	}
	i.logger.Error("team didn't exist")
	return errors.New("team didn't exist")
}

// GetVacationByID returns the associated vacation by the given id.
func (i *InmemoryDB) GetVacationByID(_ context.Context, id string) (*model.Vacation, error) {
	i.muVacationStore.Lock()
	defer i.muVacationStore.Unlock()
	for _, s := range i.vacationStore {
		if s.ID == id {
			i.logger.Info("get vacation with id: ", id)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no vacation found")
	return nil, errors.New("no vacation found")
}

// ListVacations returns a copy of the internal vacation list.
func (i *InmemoryDB) ListVacations(_ context.Context) ([]*model.Vacation, error) {
	i.muVacationStore.Lock()
	defer i.muVacationStore.Unlock()
	i.logger.Info("get list of vacations")
	vacationStore := make([]*model.Vacation, len(i.vacationStore))
	for j, v := range i.vacationStore {
		vacationStore[j] = v.Copy()
	}
	return vacationStore, nil
}

// DeleteVacation removes vacation entry by the given id.
func (i *InmemoryDB) DeleteVacation(_ context.Context, id string) error {
	i.muVacationStore.Lock()
	defer i.muVacationStore.Unlock()
	for x, vacation := range i.vacationStore {
		if vacation.ID == id {
			i.logger.Info("delete vacation with id: ", vacation.ID)
			i.vacationStore = append(i.vacationStore[:x], i.vacationStore[x+1:]...)
			return nil
		}
	}
	i.logger.Error("vacation didn't exist")
	return errors.New("vacation didn't exist")
}

// CreateVacationRequest stores an internal copy of the given vacationRequest.
// Returns copy with assigned vacationRequestID.
func (i *InmemoryDB) CreateVacationRequest(_ context.Context, v *model.VacationRequest) (*model.VacationRequest, error) {
	i.muVacationRequestStore.Lock()
	defer i.muVacationRequestStore.Unlock()
	createdAt := time.Now()
	v.CreatedAt = &createdAt
	v.ID = uuid.NewString()
	vCopy := v.Copy()

	i.logger.Info("create vacation-request with id: ", v.ID)
	i.vacationRequestStore = append(i.vacationRequestStore, vCopy)
	return v, nil
}

// GetVacationRequestByID returns the associated vacationRequest by the given id.
func (i *InmemoryDB) GetVacationRequestByID(_ context.Context, id string) (*model.VacationRequest, error) {
	i.muVacationRequestStore.Lock()
	defer i.muVacationRequestStore.Unlock()
	for _, s := range i.vacationRequestStore {
		if s.ID == id {
			i.logger.Info("get vacation-request with id: ", id)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no vacation-request found")
	return nil, errors.New("no vacation-request found")
}

// ListVacationRequests returns a copy of the internal vacationRequest list.
func (i *InmemoryDB) ListVacationRequests(_ context.Context) ([]*model.VacationRequest, error) {
	i.muVacationRequestStore.Lock()
	defer i.muVacationRequestStore.Unlock()
	i.logger.Info("get list of vacation-requests")
	vacationRequestStore := make([]*model.VacationRequest, len(i.vacationRequestStore))
	for j, v := range i.vacationRequestStore {
		vacationRequestStore[j] = v.Copy()
	}
	return vacationRequestStore, nil
}

// UpdateVacationRequest updates vacationRequest entry by the given vacationRequest.
func (i *InmemoryDB) UpdateVacationRequest(_ context.Context, v *model.VacationRequest) (*model.VacationRequest, error) {
	i.logger.Error("update failed: no update on vacation-request possible")
	return nil, errors.New("update failed: no update on vacation-request possible")
}

// DeleteVacationRequest removes vacationRequest entry by the given id.
func (i *InmemoryDB) DeleteVacationRequest(_ context.Context, id string) error {
	i.muVacationRequestStore.Lock()
	defer i.muVacationRequestStore.Unlock()
	for x, vacationRequest := range i.vacationRequestStore {
		if vacationRequest.ID == id {
			i.logger.Info("delete vacation-request with id: ", vacationRequest.ID)
			i.vacationRequestStore = append(i.vacationRequestStore[:x], i.vacationRequestStore[x+1:]...)
			return nil
		}
	}
	i.logger.Error("vacation-request didn't exist")
	return errors.New("vacation-request didn't exist")
}

// CreateVacationRessource stores an internal copy of the given vacationRessource.
// Returns copy with assigned vacationRessourceID.
func (i *InmemoryDB) CreateVacationRessource(_ context.Context, v *model.VacationRessource) (*model.VacationRessource, error) {
	i.muVacationRessourceStore.Lock()
	defer i.muVacationRessourceStore.Unlock()
	createdAt := time.Now()
	v.CreatedAt = &createdAt
	v.ID = uuid.NewString()
	vCopy := v.Copy()

	i.logger.Info("create vacation-ressource with id: ", v.ID)
	i.vacationRessourceStore = append(i.vacationRessourceStore, vCopy)
	return v, nil
}

// GetVacationRessourceByID returns the associated vacationRessource by the given id.
func (i *InmemoryDB) GetVacationRessourceByID(_ context.Context, id string) (*model.VacationRessource, error) {
	i.muVacationRessourceStore.Lock()
	defer i.muVacationRessourceStore.Unlock()
	for _, s := range i.vacationRessourceStore {
		if s.ID == id {
			i.logger.Info("get vacation-ressource with id: ", id)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no vacation-ressource found")
	return nil, errors.New("no vacation-ressource found")
}

// ListVacationRessource returns a copy of the internal vacationRessource list.
func (i *InmemoryDB) ListVacationRessource(_ context.Context) ([]*model.VacationRessource, error) {
	i.muVacationRessourceStore.Lock()
	defer i.muVacationRessourceStore.Unlock()
	i.logger.Info("get list of vacation-ressource")
	vacationRessourceStore := make([]*model.VacationRessource, len(i.vacationRessourceStore))
	for j, v := range i.vacationRessourceStore {
		vacationRessourceStore[j] = v.Copy()
	}
	return vacationRessourceStore, nil
}

// UpdateVacationRessource updates vacationRessource entry by the given vacationRessource.
func (i *InmemoryDB) UpdateVacationRessource(_ context.Context, v *model.VacationRessource) (*model.VacationRessource, error) {
	i.logger.Error("update failed: no update on vacation-ressource possible")
	return nil, errors.New("update failed: no update on vacation-ressource possible")
}

// DeleteVacationRessource removes vacationRessource entry by the given id.
func (i *InmemoryDB) DeleteVacationRessource(_ context.Context, id string) error {
	i.muVacationRessourceStore.Lock()
	defer i.muVacationRessourceStore.Unlock()
	for x, vacationRessource := range i.vacationRessourceStore {
		if vacationRessource.ID == id {
			i.logger.Info("delete vacation-ressource with id: ", vacationRessource.ID)
			i.vacationRessourceStore = append(i.vacationRessourceStore[:x], i.vacationRessourceStore[x+1:]...)
			return nil
		}
	}
	i.logger.Error("vacation-ressource didn't exist")
	return errors.New("vacation-ressource didn't exist")
}
