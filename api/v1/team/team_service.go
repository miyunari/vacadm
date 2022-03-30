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
	logger.Info("create new team")
	var team model.Team
	err := json.NewDecoder(r.Body).Decode(&team)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tm, err := t.store.CreateTeam(r.Context(), &team)
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
	logger.Info("get team by id")
	teamID, err := extractTeamID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	team, err := t.store.GetTeamByID(r.Context(), teamID)
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
	logger.Info("retrieve team list")
	list, err := t.store.ListTeams(r.Context())
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

func (t *teamService) ListTeamUsers(w http.ResponseWriter, r *http.Request) {
	logger := t.logger.WithField("component", "list-users")
	logger.Info("retrieve list users from team")
	teamID, err := extractTeamID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	teamUser, err := t.store.ListTeamUsers(r.Context(), teamID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(teamUser)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t *teamService) Update(w http.ResponseWriter, r *http.Request) {
	logger := t.logger.WithField("component", "update")
	logger.Info("update team")
	var team model.Team
	err := json.NewDecoder(r.Body).Decode(&team)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	uTeam, err := t.store.UpdateTeam(r.Context(), &team)
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
	logger.Info("delete team")
	teamID, err := extractTeamID(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = t.store.DeleteTeam(r.Context(), teamID)
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
