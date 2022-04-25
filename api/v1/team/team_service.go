package team

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/MninaTB/vacadm/api/v1/util"
	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/jwt"
	"github.com/MninaTB/vacadm/pkg/model"
)

// Tokenizer implements methods to verify auth tokens.
type Tokenizer interface {
	// Valid if a token is valid, userID and teamID are returned.
	// if a token is invalid, an error is returned.
	Valid(token string) (userID string, teamID string, err error)
}

// NewTeamService returns a new TeamService.
func NewTeamService(store database.Database, logger logrus.FieldLogger, t Tokenizer) *TeamService {
	return &TeamService{
		store:         store,
		relationStore: database.NewRelationDB(store),
		tokenizer:     t,
		logger:        logger.WithField("component", "team-service"),
	}
}

// TeamService implements http.HandlerFunc's to operate on team resources.
type TeamService struct {
	store         database.Database
	relationStore database.RelationDB
	logger        logrus.FieldLogger
	tokenizer     Tokenizer
}

// Create reads the given payload and creates a store representation accordingly.
func (t *TeamService) Create(w http.ResponseWriter, r *http.Request) {
	logger := t.logger.WithField("method", "create")
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

// GetByID extracts a TeamID from URL and writes all team information into the
// given response writer.
func (t *TeamService) GetByID(w http.ResponseWriter, r *http.Request) {
	logger := t.logger.WithField("method", "getByID")
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

// List retuns a list of all teams available on the internal store.
func (t *TeamService) List(w http.ResponseWriter, r *http.Request) {
	logger := t.logger.WithField("method", "list")
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

// ListTeamUsers returns a list of users associated to teamID transmitted in URL.
func (t *TeamService) ListTeamUsers(w http.ResponseWriter, r *http.Request) {
	logger := t.logger.WithField("method", "ListTeamUsers")
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

// Update reads new team settings from the request body and updates the store
// representation accordingly.
func (t *TeamService) Update(w http.ResponseWriter, r *http.Request) {
	logger := t.logger.WithField("method", "update")
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

// Delete a team associated to the given teamID in the URL.
func (t *TeamService) Delete(w http.ResponseWriter, r *http.Request) {
	logger := t.logger.WithField("method", "delete")
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

// ListCapacity lists teams and their availability for the requested period.
// Example request:
// {
//   "from":"2009-11-10T23:00:00Z",
//   "to":"2009-11-10T23:00:00Z",
//   # NOTE: team_id is optional. if no teamID is provided, all teams are taken
//   # into account.
//   "team_id":"1234-456"
// }
//
// Example response requesting Content-Type json/application:
// {
//   "to":"0001-01-01T00:00:00Z",
//   "from":"0001-01-01T00:00:00Z",
//   "team_id":"",
//   # NOTE: Depending on the availability ratio this value can be:
//   # "HIGH", "MEDIUM" or "LOW".
//   "availability":"HIGH",
//   # NOTE: vacations are only displayed if the requesting user is a team owner.
//   # Or the requesting user is the parent of a team owner.
//   # Parent is recursive in this case. This means that the parent of the
//   # parent is also valid.
//   "vacations":[
//     {
//       "id":"",
//       "user_id":"",
//       "approved_by":null,
//       "from":"0001-01-01T00:00:00Z",
//       "to":"0001-01-01T00:00:00Z",
//       "created_at":null
//     }
//   ]
// }
//
// Example response requesting Content-Type csv/application:
// from,to,teamID,availability,vacation-id,vacation-user_id,vacation-approved_by,vacation-from,vacation-to,vacation-created_at,vacation-deleted_at
// 2022-04-19 22:23:40.886412444 +0200 CEST m=-258901.921920057,2022-04-25 22:23:40.886412586 +0200 CEST m=+259498.078080085,a7da8eb8-410f-4f6a-8324-1db65a289a13,HIGH,,,,,,,
//	2022-04-19 22:23:40.886412677 +0200 CEST m=-258901.921919824,2022-04-25 22:23:40.886412747 +0200 CEST m=+259498.078080246,e22b2a12-cf42-44c6-a2ed-c3630ba9583a,HIGH,,,,,,,
//	2022-04-19 22:23:40.886412822 +0200 CEST m=-258901.921919683,2022-04-25 22:23:40.886412887 +0200 CEST m=+259498.078080386,e22b2a12-cf42-44c6-a2ed-c3630ba9583a,HIGH,,,,,,,
func (t *TeamService) ListCapacity(w http.ResponseWriter, r *http.Request) {
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
	} else {
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
			window.Vacation = vacs
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
			window.Availability = "HIGH"
		} else if ratio <= 0.8 && ratio > 0.25 {
			window.Availability = "MEDIUM"
		} else {
			window.Availability = "LOW"
		}
		resp = append(resp, window)
	}

	if r.Header.Get("Content-Type") == "application/csv" {
		if err := new(capacityResponse).WriteCSVHeader(w); err != nil {
			logger.Error(err)
		}
		for _, set := range resp {
			if err := set.WriteCSV(w); err != nil {
				logger.Error(err)
			}
		}
		return
	}

	err = json.NewEncoder(w).Encode(&resp)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (t *TeamService) vaccationByTeam(ctx context.Context, teamID string, from, to time.Time) ([]*model.Vacation, error) {
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
	From time.Time `json:"to"`
	To   time.Time `json:"from"`
	// NOTE: optional
	TeamID string `json:"team_id"`
}

type capacityResponse struct {
	From   time.Time `json:"to"`
	To     time.Time `json:"from"`
	TeamID string    `json:"team_id"`

	// NOTE: HIGH, MEDIUM, LOW
	Availability string `json:"availability"`

	// NOTE: Only with sufficient authorization
	Vacation []*model.Vacation `json:"vacations"`
}

func (c *capacityResponse) WriteCSVHeader(w io.Writer) error {
	wr := csv.NewWriter(w)
	defer wr.Flush()
	wr.Write([]string{
		"from", "to", "teamID", "availability",
		"vacation-id", "vacation-user_id", "vacation-approved_by",
		"vacation-from", "vacation-to", "vacation-created_at", "vacation-deleted_at",
	})
	return nil
}

func (c *capacityResponse) WriteCSV(w io.Writer) error {
	wr := csv.NewWriter(w)
	defer wr.Flush()
	if len(c.Vacation) == 0 {
		wr.Write([]string{
			c.From.String(), c.To.String(), c.TeamID, c.Availability,
			"", "", "", "", "", "", "",
		})
		return nil
	}
	for _, vac := range c.Vacation {
		var approvedBy string
		if vac.ApprovedBy != nil {
			approvedBy = *vac.ApprovedBy
		}
		var createdAt, deletedAt string
		if vac.CreatedAt != nil {
			createdAt = vac.CreatedAt.String()
		}
		if vac.DeletedAt != nil {
			deletedAt = vac.DeletedAt.String()
		}
		wr.Write([]string{
			c.From.String(), c.To.String(), c.TeamID, c.Availability,
			vac.ID, vac.UserID, approvedBy,
			vac.From.String(), vac.To.String(), createdAt, deletedAt,
		})
	}
	return nil
}

type teamBundle struct {
	teamOwnerID string
	teamID      string
}
