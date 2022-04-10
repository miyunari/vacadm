package team

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/MninaTB/vacadm/api/v1/util"
	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/jwt"
	"github.com/MninaTB/vacadm/pkg/model"
	"github.com/sirupsen/logrus"
)

type Tokenizer interface {
	Valid(token string) (userID string, teamID string, err error)
}

func NewTeamService(store database.Database, logger logrus.FieldLogger, t Tokenizer) *teamService {
	return &teamService{
		store:         store,
		relationStore: database.NewRelationDB(store),
		tokenizer:     t,
		logger:        logger.WithField("component", "team-service"),
	}
}

type teamService struct {
	store         database.Database
	relationStore database.RelationDB
	logger        logrus.FieldLogger
	tokenizer     Tokenizer
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
	teamID, err := util.TeamIDFromRequest(r)
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
	teamID, err := util.TeamIDFromRequest(r)
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
	teamID, err := util.TeamIDFromRequest(r)
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

func (t *teamService) ListCapacity(w http.ResponseWriter, r *http.Request) {
	logger := t.logger.WithFields(
		logrus.Fields{
			"method": "list-capacity",
		},
	)

	request := &capacityRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger = logger.WithFields(
		logrus.Fields{
			"from": request.From,
			"to":   request.To,
		},
	)
	logger.Info("list capacity")

	token, err := jwt.ExtractToken(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// NOTE: ignore teamID, team members have only access to a anonymized list
	userID, _, err := t.tokenizer.Valid(token)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var resp []*capacityResponse
	var teamsBundle []*teamBundle

	if request.TeamID != "" {
		team, err := t.store.GetTeamByID(r.Context(), request.TeamID)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		teamsBundle = append(teamsBundle, &teamBundle{
			teamOwnerID: team.OwnerID,
			teamID:      request.TeamID,
		})
	}

	if len(teamsBundle) == 0 {
		teams, err := t.store.ListTeams(r.Context())
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for _, team := range teams {
			teamsBundle = append(teamsBundle, &teamBundle{
				teamOwnerID: team.OwnerID,
				teamID:      team.ID,
			})
		}
	}

	for _, tb := range teamsBundle {
		vacs, err := t.vaccationByTeam(r.Context(), tb.teamID, request.From, request.To)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		window := &capacityResponse{
			TeamID: request.TeamID,
			From:   request.From,
			To:     request.To,
		}

		isOwner, err := t.relationStore.IsTeamOwner(r.Context(), tb.teamID, userID)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		isParentOfOwner, err := t.relationStore.IsParentUser(r.Context(), tb.teamOwnerID, userID)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if isOwner || isParentOfOwner {
			window.vacation = vacs
		}

		users, err := t.store.ListTeamUsers(r.Context(), tb.teamID)
		if err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// NOTE: we dont need to consider pure workdays, since vacation resources
		// also do not consider them.
		workDays := daysBetween(request.From, request.To) * float64(len(users))
		var daysOfVacation float64
		for _, vac := range vacs {
			daysOfVacation += daysBetween(vac.From, vac.To)
		}
		ratio := workDays / daysOfVacation
		if ratio > 0.8 {
			window.Availability = "LOW"
		} else if ratio <= 0.8 && ratio > 0.25 {
			window.Availability = "MEDIUM"
		} else {
			window.Availability = "HIGH"
		}
		resp = append(resp, window)
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (t *teamService) vaccationByTeam(ctx context.Context, teamID string, from, to time.Time) ([]*model.Vacation, error) {
	vacations, err := t.store.GetVacationsByTeamID(ctx, teamID)
	if err != nil {
		return nil, err
	}

	var response []*model.Vacation
	for _, v := range vacations {
		if !(v.From.After(from) && from.Before(to)) {
			continue
		}
		// NOTE: trim start time
		if v.From.Before(from) {
			v.From = from
		}
		// NOTE: trim end time
		if v.To.After(to) {
			v.To = to
		}
		response = append(response, v)
	}

	return response, nil
}

func daysBetween(from, to time.Time) float64 {
	return to.Sub(from).Hours() / 24
}

type capacityRequest struct {
	From   time.Time
	To     time.Time
	TeamID string
}

type capacityResponse struct {
	From   time.Time
	To     time.Time
	TeamID string

	Availability string // NOTE: HIGH, MEDIUM, LOW

	// NOTE: Only with sufficient authorization
	vacation []*model.Vacation
}

type teamBundle struct {
	teamOwnerID string
	teamID      string
}
