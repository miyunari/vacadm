package token

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/MninaTB/vacadm/api/v1/util"
	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/jwt"
	"github.com/MninaTB/vacadm/pkg/model"
)

// Tokenizer implements methods to verify and generate new auth tokens.
type Tokenizer interface {
	Valid(token string) (userID string, teamID string, err error)
	Generate(u *model.User) (string, error)
}

// NewTokenService returns a new TokenService
func NewTokenService(db database.Database, t Tokenizer) *TokenService {
	return &TokenService{
		tokenizer:     t,
		store:         db,
		relationStore: database.NewRelationDB(db),
		logger:        logrus.WithField("component", "token-service"),
	}
}

// TokenService is used to refresh/create tokens for existing users.
type TokenService struct {
	store         database.Database
	relationStore database.RelationDB
	tokenizer     Tokenizer
	logger        logrus.FieldLogger
}

// Refresh verifies user permissions based on the given token. If a user is
// autorized a new token will be generated and returned in the response body.
// Payload example:
// {"token":"a6f9f420-0c43-4527-8178-fe53a2a66302"}
func (t *TokenService) Refresh(w http.ResponseWriter, r *http.Request) {
	logger := t.logger.WithField("method", "refresh")

	token, err := jwt.ExtractToken(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	parentID, _, err := t.tokenizer.Valid(token)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userID, err := util.UserIDFromRequest(r)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	isParent, err := t.relationStore.IsParentUser(r.Context(), userID, parentID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if parentID != userID || !isParent {
		logger.Error(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	usr, err := t.store.GetUserByID(r.Context(), userID)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	newToken, err := t.tokenizer.Generate(usr)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: newToken,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
