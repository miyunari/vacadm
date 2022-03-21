package team

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/model"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func NewTeamService(store database.Database, logger logrus.FieldLogger) *teamService {
	return &teamService{
		store:  store,
		logger: logger.WithField("component", "team-service"),
	}
}

type teamService struct {
	store  database.Database
	logger logrus.FieldLogger
}

func (t *teamService) Create(w http.ResponseWriter, r *http.Request) {
	logger := t.logger.WithField("component", "create")
	var team model.Team
	err := json.NewDecoder(r.Body).Decode(&team)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tm, err := t.store.CreateTeam(&team)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(tm)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t.logger.Info("create team with id: ", team.ID)
	w.WriteHeader(http.StatusCreated)
}

func (t *teamService) GetByID(w http.ResponseWriter, r *http.Request) {
	logger := t.logger.WithField("component", "read")
	teamID, err := extractTeamID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	team, err := t.store.GetTeamByID(teamID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(team)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	t.logger.Info("get team with id: ", teamID)
}

func (t *teamService) List(w http.ResponseWriter, r *http.Request) {
	logger := t.logger.WithField("component", "list")
	list, err := t.store.ListTeams()
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(&list)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	t.logger.Info("get list of teams")
}

func (t *teamService) Update(w http.ResponseWriter, r *http.Request) {
	logger := t.logger.WithField("component", "update")
	var team model.Team
	err := json.NewDecoder(r.Body).Decode(&team)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	uTeam, err := t.store.UpdateTeam(&team)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.NewEncoder(w).Encode(&uTeam)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	t.logger.Info("update team with id: ", team.ID)
}

func (t *teamService) Delete(w http.ResponseWriter, r *http.Request) {
	logger := t.logger.WithField("component", "delete")
	teamID, err := extractTeamID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = t.store.DeleteTeam(teamID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t.logger.Info("delete team with id: ", teamID)
	w.WriteHeader(http.StatusAccepted)
}

func extractTeamID(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	teamID, ok := vars["teamID"]
	if !ok {
		return "", errors.New("could not extract teamID")
	}
	return teamID, nil
}
