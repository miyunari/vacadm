package inmemory

import (
	"context"
	"errors"
	"time"

	"github.com/MninaTB/vacadm/pkg/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func NewInmemoryDB() *inmemoryDB {
	return &inmemoryDB{
		userStore:               make([]*model.User, 0),
		teamStore:               make([]*model.Team, 0),
		vaccationStore:          make([]*model.Vaccation, 0),
		vaccationRequestStore:   make([]*model.VaccationRequest, 0),
		vaccationRessourceStore: make([]*model.VaccationRessource, 0),
		logger:                  logrus.New().WithField("component", "inmemoryDB"),
	}
}

type inmemoryDB struct {
	userStore               []*model.User
	teamStore               []*model.Team
	vaccationStore          []*model.Vaccation
	vaccationRequestStore   []*model.VaccationRequest
	vaccationRessourceStore []*model.VaccationRessource
	logger                  logrus.FieldLogger
}

func (i *inmemoryDB) CreateUser(_ context.Context, user *model.User) (*model.User, error) {
	createdAt := time.Now()
	user.CreatedAt = &createdAt
	user.ID = uuid.NewString()
	usrCopy := user.Copy()

	i.logger.Info("create user with id: ", user.ID)
	i.userStore = append(i.userStore, usrCopy)
	return user, nil
}

func (i *inmemoryDB) GetUserByID(_ context.Context, id string) (*model.User, error) {
	for _, s := range i.userStore {
		if s.ID == id {
			i.logger.Info("get user with id: ", s.ID)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no user found")
	return nil, errors.New("no user found")
}

func (i *inmemoryDB) ListUsers(_ context.Context) ([]*model.User, error) {
	i.logger.Info("get list of users")
	return i.userStore, nil
}

func (i *inmemoryDB) UpdateUser(_ context.Context, user *model.User) (*model.User, error) {
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

func (i *inmemoryDB) DeleteUser(_ context.Context, id string) error {
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

func (i *inmemoryDB) CreateTeam(_ context.Context, team *model.Team) (*model.Team, error) {
	createdAt := time.Now()
	team.CreatedAt = &createdAt
	team.ID = uuid.NewString()
	teamCopy := team.Copy()

	i.logger.Info("create team with id: ", team.ID)
	i.teamStore = append(i.teamStore, teamCopy)
	return team, nil
}

func (i *inmemoryDB) GetTeamByID(_ context.Context, id string) (*model.Team, error) {
	for _, s := range i.teamStore {
		if s.ID == id {
			i.logger.Info("get team with id: ", s.ID)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no team found")
	return nil, errors.New("no team found")
}

func (i *inmemoryDB) ListTeams(_ context.Context) ([]*model.Team, error) {
	i.logger.Info("get list of teams")
	return i.teamStore, nil
}

func (i *inmemoryDB) ListTeamUsers(ctx context.Context, teamID string) ([]*model.User, error) {
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

func (i *inmemoryDB) UpdateTeam(_ context.Context, team *model.Team) (*model.Team, error) {
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

func (i *inmemoryDB) DeleteTeam(_ context.Context, id string) error {
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

func (i *inmemoryDB) GetVaccationByID(_ context.Context, id string) (*model.Vaccation, error) {
	for _, s := range i.vaccationStore {
		if s.ID == id {
			i.logger.Info("get vaccation with id: ", id)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no vaccation found")
	return nil, errors.New("no vaccation found")
}

func (i *inmemoryDB) ListVaccations(_ context.Context) ([]*model.Vaccation, error) {
	i.logger.Info("get list of vaccations")
	return i.vaccationStore, nil
}

func (i *inmemoryDB) DeleteVaccation(_ context.Context, id string) error {
	for x, vaccation := range i.vaccationStore {
		if vaccation.ID == id {
			i.logger.Info("delete vaccation with id: ", vaccation.ID)
			i.vaccationStore = append(i.vaccationStore[:x], i.vaccationStore[x+1:]...)
			return nil
		}
	}
	i.logger.Error("vaccation didn't exist")
	return errors.New("vaccation didn't exist")
}

func (i *inmemoryDB) CreateVaccationRequest(_ context.Context, v *model.VaccationRequest) (*model.VaccationRequest, error) {
	createdAt := time.Now()
	v.CreatedAt = &createdAt
	v.ID = uuid.NewString()
	vCopy := v.Copy()

	i.logger.Info("create vaccation-request with id: ", v.ID)
	i.vaccationRequestStore = append(i.vaccationRequestStore, vCopy)
	return v, nil
}

func (i *inmemoryDB) GetVaccationRequestByID(_ context.Context, id string) (*model.VaccationRequest, error) {
	for _, s := range i.vaccationRequestStore {
		if s.ID == id {
			i.logger.Info("get vaccation-request with id: ", id)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no vaccation-request found")
	return nil, errors.New("no vaccation-request found")
}

func (i *inmemoryDB) ListVaccationRequests(_ context.Context) ([]*model.VaccationRequest, error) {
	i.logger.Info("get list of vaccation-requests")
	return i.vaccationRequestStore, nil
}

func (i *inmemoryDB) UpdateVaccationRequest(_ context.Context, v *model.VaccationRequest) (*model.VaccationRequest, error) {
	i.logger.Error("update failed: no update on vaccation-request possible")
	return nil, errors.New("update failed: no update on vaccation-request possible")
}

func (i *inmemoryDB) DeleteVaccationRequest(_ context.Context, id string) error {
	for x, vaccationRequest := range i.vaccationRequestStore {
		if vaccationRequest.ID == id {
			i.logger.Info("delete vaccation-request with id: ", vaccationRequest.ID)
			i.vaccationRequestStore = append(i.vaccationRequestStore[:x], i.vaccationRequestStore[x+1:]...)
			return nil
		}
	}
	i.logger.Error("vaccation-request didn't exist")
	return errors.New("vaccation-request didn't exist")
}

func (i *inmemoryDB) CreateVaccationRessource(_ context.Context, v *model.VaccationRessource) (*model.VaccationRessource, error) {
	createdAt := time.Now()
	v.CreatedAt = &createdAt
	v.ID = uuid.NewString()
	vCopy := v.Copy()

	i.logger.Info("create vaccation-ressource with id: ", v.ID)
	i.vaccationRessourceStore = append(i.vaccationRessourceStore, vCopy)
	return v, nil
}

func (i *inmemoryDB) GetVaccationRessourceByID(_ context.Context, id string) (*model.VaccationRessource, error) {
	for _, s := range i.vaccationRessourceStore {
		if s.ID == id {
			i.logger.Info("get vaccation-ressource with id: ", id)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no vaccation-ressource found")
	return nil, errors.New("no vaccation-ressource found")
}

func (i *inmemoryDB) ListVaccationRessource(_ context.Context) ([]*model.VaccationRessource, error) {
	i.logger.Info("get list of vaccation-ressource")
	return i.vaccationRessourceStore, nil
}

func (i *inmemoryDB) UpdateVaccationRessource(_ context.Context, v *model.VaccationRessource) (*model.VaccationRessource, error) {
	i.logger.Error("update failed: no update on vaccation-ressource possible")
	return nil, errors.New("update failed: no update on vaccation-ressource possible")
}

func (i *inmemoryDB) DeleteVaccationRessource(_ context.Context, id string) error {
	for x, vaccationRessource := range i.vaccationRessourceStore {
		if vaccationRessource.ID == id {
			i.logger.Info("delete vaccation-ressource with id: ", vaccationRessource.ID)
			i.vaccationRessourceStore = append(i.vaccationRessourceStore[:x], i.vaccationRessourceStore[x+1:]...)
			return nil
		}
	}
	i.logger.Error("vaccation-ressource didn't exist")
	return errors.New("vaccation-ressource didn't exist")
}
