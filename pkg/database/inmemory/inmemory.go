package inmemory

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MninaTB/vacadm/pkg/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func NewInmemoryDB() *InmemoryDB {
	return &InmemoryDB{
		userStore:               make([]*model.User, 0),
		teamStore:               make([]*model.Team, 0),
		vacationStore:          make([]*model.Vacation, 0),
		vacationRequestStore:   make([]*model.VacationRequest, 0),
		vacationRessourceStore: make([]*model.VacationRessource, 0),
		logger:                  logrus.New().WithField("component", "inmemoryDB"),
	}
}

type InmemoryDB struct {
	userStore               []*model.User
	teamStore               []*model.Team
	vacationStore          []*model.Vacation
	vacationRequestStore   []*model.VacationRequest
	vacationRessourceStore []*model.VacationRessource
	logger                  logrus.FieldLogger
}

func (i *InmemoryDB) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
	for _, e := range i.userStore {
		if e.Email == user.Email {
			return nil, fmt.Errorf("email address must be unique")
		}
	}

	if user.ParentID != nil {
		var foundParent bool
		for _, p := range i.userStore {
			if *p.ParentID == *user.ParentID {
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

func (i *InmemoryDB) GetUserByID(_ context.Context, id string) (*model.User, error) {
	for _, s := range i.userStore {
		if s.ID == id {
			i.logger.Info("get user with id: ", s.ID)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no user found")
	return nil, errors.New("no user found")
}

func (i *InmemoryDB) ListUsers(_ context.Context) ([]*model.User, error) {
	i.logger.Info("get list of users")
	return i.userStore, nil
}

func (i *InmemoryDB) UpdateUser(_ context.Context, user *model.User) (*model.User, error) {
	updatededAt := time.Now()
	for x := 0; x < len(i.userStore); x++ {
		if i.userStore[x].ID == user.ID {
			if user.Email != "" {
				i.userStore[x].Email = user.Email
			}
			i.userStore[x].FirstName = user.FirstName
			i.userStore[x].LastName = user.LastName
			i.userStore[x].UpdatedAt = &updatededAt
			i.logger.Info("update user with id: ", user.ID)
			return i.userStore[x], nil
		}
	}
	i.logger.Error("update failed: no user found")
	return nil, errors.New("update failed: no user found")
}

func (i *InmemoryDB) DeleteUser(_ context.Context, id string) error {
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

func (i *InmemoryDB) CreateTeam(_ context.Context, team *model.Team) (*model.Team, error) {
	createdAt := time.Now()
	team.CreatedAt = &createdAt
	team.ID = uuid.NewString()
	teamCopy := team.Copy()

	i.logger.Info("create team with id: ", team.ID)
	i.teamStore = append(i.teamStore, teamCopy)
	return team, nil
}

func (i *InmemoryDB) GetTeamByID(_ context.Context, id string) (*model.Team, error) {
	for _, s := range i.teamStore {
		if s.ID == id {
			i.logger.Info("get team with id: ", s.ID)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no team found")
	return nil, errors.New("no team found")
}

func (i *InmemoryDB) ListTeams(_ context.Context) ([]*model.Team, error) {
	i.logger.Info("get list of teams")
	return i.teamStore, nil
}

func (i *InmemoryDB) ListTeamUsers(ctx context.Context, teamID string) ([]*model.User, error) {
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

func (i *InmemoryDB) UpdateTeam(_ context.Context, team *model.Team) (*model.Team, error) {
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

func (i *InmemoryDB) DeleteTeam(_ context.Context, id string) error {
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

func (i *InmemoryDB) GetVacationByID(_ context.Context, id string) (*model.Vacation, error) {
	for _, s := range i.vacationStore {
		if s.ID == id {
			i.logger.Info("get vacation with id: ", id)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no vacation found")
	return nil, errors.New("no vacation found")
}

func (i *InmemoryDB) ListVacations(_ context.Context) ([]*model.Vacation, error) {
	i.logger.Info("get list of vacations")
	return i.vacationStore, nil
}

func (i *InmemoryDB) DeleteVacation(_ context.Context, id string) error {
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

func (i *InmemoryDB) CreateVacationRequest(_ context.Context, v *model.VacationRequest) (*model.VacationRequest, error) {
	createdAt := time.Now()
	v.CreatedAt = &createdAt
	v.ID = uuid.NewString()
	vCopy := v.Copy()

	i.logger.Info("create vacation-request with id: ", v.ID)
	i.vacationRequestStore = append(i.vacationRequestStore, vCopy)
	return v, nil
}

func (i *InmemoryDB) GetVacationRequestByID(_ context.Context, id string) (*model.VacationRequest, error) {
	for _, s := range i.vacationRequestStore {
		if s.ID == id {
			i.logger.Info("get vacation-request with id: ", id)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no vacation-request found")
	return nil, errors.New("no vacation-request found")
}

func (i *InmemoryDB) ListVacationRequests(_ context.Context) ([]*model.VacationRequest, error) {
	i.logger.Info("get list of vacation-requests")
	return i.vacationRequestStore, nil
}

func (i *InmemoryDB) UpdateVacationRequest(_ context.Context, v *model.VacationRequest) (*model.VacationRequest, error) {
	i.logger.Error("update failed: no update on vacation-request possible")
	return nil, errors.New("update failed: no update on vacation-request possible")
}

func (i *InmemoryDB) DeleteVacationRequest(_ context.Context, id string) error {
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

func (i *InmemoryDB) CreateVacationRessource(_ context.Context, v *model.VacationRessource) (*model.VacationRessource, error) {
	createdAt := time.Now()
	v.CreatedAt = &createdAt
	v.ID = uuid.NewString()
	vCopy := v.Copy()

	i.logger.Info("create vacation-ressource with id: ", v.ID)
	i.vacationRessourceStore = append(i.vacationRessourceStore, vCopy)
	return v, nil
}

func (i *InmemoryDB) GetVacationRessourceByID(_ context.Context, id string) (*model.VacationRessource, error) {
	for _, s := range i.vacationRessourceStore {
		if s.ID == id {
			i.logger.Info("get vacation-ressource with id: ", id)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no vacation-ressource found")
	return nil, errors.New("no vacation-ressource found")
}

func (i *InmemoryDB) ListVacationRessource(_ context.Context) ([]*model.VacationRessource, error) {
	i.logger.Info("get list of vacation-ressource")
	return i.vacationRessourceStore, nil
}

func (i *InmemoryDB) UpdateVacationRessource(_ context.Context, v *model.VacationRessource) (*model.VacationRessource, error) {
	i.logger.Error("update failed: no update on vacation-ressource possible")
	return nil, errors.New("update failed: no update on vacation-ressource possible")
}

func (i *InmemoryDB) DeleteVacationRessource(_ context.Context, id string) error {
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
