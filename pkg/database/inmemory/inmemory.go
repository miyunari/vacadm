package inmemory

import (
	"errors"
	"time"

	"github.com/MninaTB/vacadm/pkg/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func NewInmemoryDB() *inmemoryDB {
	return &inmemoryDB{
		userStore:      make([]*model.User, 0),
		teamStore:      make([]*model.Team, 0),
		vaccationStore: make([]*model.Vaccation, 0),
		logger:         logrus.New().WithField("component", "inmemoryDB"),
	}
}

type inmemoryDB struct {
	userStore      []*model.User
	teamStore      []*model.Team
	vaccationStore []*model.Vaccation
	logger         logrus.FieldLogger
}

func (i *inmemoryDB) CreateUser(user *model.User) (*model.User, error) {
	user.CreatedAt = time.Now()
	user.ID = uuid.NewString()
	usrCopy := user.Copy()

	i.logger.Info("create user with id: ", user.ID)
	i.userStore = append(i.userStore, usrCopy)
	return user, nil
}

func (i *inmemoryDB) GetUserByID(id string) (*model.User, error) {
	for _, s := range i.userStore {
		if s.ID == id {
			i.logger.Info("get user with id: ", s.ID)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no user found")
	return nil, errors.New("no user found")
}

func (i *inmemoryDB) ListUsers() ([]*model.User, error) {
	i.logger.Info("get list of users")
	return i.userStore, nil
}

func (i *inmemoryDB) UpdateUser(user *model.User) (*model.User, error) {
	for x := 0; x < len(i.userStore); x++ {
		if i.userStore[x].ID == user.ID {
			if user.Email != "" {
				i.userStore[x].Email = user.Email
			}
			i.userStore[x].First_name = user.First_name
			i.userStore[x].Last_name = user.Last_name
			i.userStore[x].UpdatedAt = time.Now()
			i.logger.Info("update user with id: ", user.ID)
			return i.userStore[x], nil
		}
	}
	i.logger.Error("update failed: no user found")
	return nil, errors.New("update failed: no user found")
}

func (i *inmemoryDB) DeleteUser(id string) error {
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

func (i *inmemoryDB) CreateTeam(team *model.Team) (*model.Team, error) {
	team.CreatedAt = time.Now()
	team.ID = uuid.NewString()
	teamCopy := team.Copy()

	i.logger.Info("create team with id: ", team.ID)
	i.teamStore = append(i.teamStore, teamCopy)
	return team, nil
}

func (i *inmemoryDB) GetTeamByID(id string) (*model.Team, error) {
	for _, s := range i.teamStore {
		if s.ID == id {
			i.logger.Info("get team with id: ", s.ID)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no team found")
	return nil, errors.New("no team found")
}

func (i *inmemoryDB) ListTeams() ([]*model.Team, error) {
	i.logger.Info("get list of teams")
	return i.teamStore, nil
}

func (i *inmemoryDB) UpdateTeam(team *model.Team) (*model.Team, error) {
	for x := 0; x < len(i.teamStore); x++ {
		if i.teamStore[x].ID == team.ID {
			i.teamStore[x].Name = team.Name
			i.teamStore[x].UpdatedAt = time.Now()
			i.logger.Info("update team with id: ", team.ID)
			return i.teamStore[x], nil
		}
	}
	i.logger.Info("update failed: no team found")
	return nil, errors.New("update failed: no team found")
}

func (i *inmemoryDB) DeleteTeam(id string) error {
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

func (i *inmemoryDB) CreateVaccation(vaccation *model.Vaccation) (*model.Vaccation, error) {
	vaccation.CreatedAt = time.Now()
	vaccation.ID = uuid.NewString()
	vaccationCopy := vaccation.Copy()

	i.logger.Info("create vaccation with id: ", vaccation.ID)
	i.vaccationStore = append(i.vaccationStore, vaccationCopy)
	return vaccation, nil
}

func (i *inmemoryDB) GetVaccationByID(id string) (*model.Vaccation, error) {
	for _, s := range i.vaccationStore {
		if s.ID == id {
			i.logger.Info("get vaccation with id: ", id)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no vaccation found")
	return nil, errors.New("no vaccation found")
}

func (i *inmemoryDB) ListVaccations() ([]*model.Vaccation, error) {
	i.logger.Info("get list of vaccations")
	return i.vaccationStore, nil
}

func (i *inmemoryDB) UpdateVaccation(vaccation *model.Vaccation) (*model.Vaccation, error) {
	i.logger.Error("update failed: no update on vaccation possible")
	return nil, errors.New("update failed: no update on vaccation possible")
}

func (i *inmemoryDB) DeleteVaccation(id string) error {
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
