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

func (i *inmemoryDB) CreateVaccationRequest(vacReq *model.VaccationRequest) (*model.VaccationRequest, error) {
	vacReq.CreatedAt = time.Now()
	vacReq.ID = uuid.NewString()
	vacReqCopy := vacReq.Copy()

	i.logger.Info("create vaccation-request with id: ", vacReq.ID)
	i.vaccationRequestStore = append(i.vaccationRequestStore, vacReqCopy)
	return vacReq, nil
}

func (i *inmemoryDB) GetVaccationRequestByID(id string) (*model.VaccationRequest, error) {
	for _, s := range i.vaccationRequestStore {
		if s.ID == id {
			i.logger.Info("get vaccation-request with id: ", id)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no vaccation-request found")
	return nil, errors.New("no vaccation-request found")
}

func (i *inmemoryDB) ListVaccationRequests() ([]*model.VaccationRequest, error) {
	i.logger.Info("get list of vaccation-requests")
	return i.vaccationRequestStore, nil
}

func (i *inmemoryDB) UpdateVaccationRequest(vaccationRequest *model.VaccationRequest) (*model.VaccationRequest, error) {
	i.logger.Error("update failed: no update on vaccation-request possible")
	return nil, errors.New("update failed: no update on vaccation-request possible")
}

func (i *inmemoryDB) DeleteVaccationRequest(id string) error {
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

func (i *inmemoryDB) CreateVaccationRessource(vacRes *model.VaccationRessource) (*model.VaccationRessource, error) {
	vacRes.CreatedAt = time.Now()
	vacRes.ID = uuid.NewString()
	vacResCopy := vacRes.Copy()

	i.logger.Info("create vaccation-ressource with id: ", vacRes.ID)
	i.vaccationRessourceStore = append(i.vaccationRessourceStore, vacResCopy)
	return vacRes, nil
}

func (i *inmemoryDB) GetVaccationRessourceByID(id string) (*model.VaccationRessource, error) {
	for _, s := range i.vaccationRessourceStore {
		if s.ID == id {
			i.logger.Info("get vaccation-ressource with id: ", id)
			return s.Copy(), nil
		}
	}
	i.logger.Error("no vaccation-ressource found")
	return nil, errors.New("no vaccation-ressource found")
}

func (i *inmemoryDB) ListVaccationRessource() ([]*model.VaccationRessource, error) {
	i.logger.Info("get list of vaccation-ressource")
	return i.vaccationRessourceStore, nil
}

func (i *inmemoryDB) UpdateVaccationRessource(vaccationRessources *model.VaccationRessource) (*model.VaccationRessource, error) {
	i.logger.Error("update failed: no update on vaccation-ressource possible")
	return nil, errors.New("update failed: no update on vaccation-ressource possible")
}

func (i *inmemoryDB) DeleteVaccationRessource(id string) error {
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
