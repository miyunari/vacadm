package middleware

import (
	"context"
	"net/http"

	"github.com/MninaTB/vacadm/api/v1/util"
	jwt "github.com/MninaTB/vacadm/pkg/jwt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Validator interface {
	Valid(token string) (userID string, teamID string, err error)
}

type RelationDB interface {
	IsParentUser(ctx context.Context, userID, parentID string) (bool, error)
	IsTeamMember(ctx context.Context, teamID, userID string) (bool, error)
	IsTeamOwner(ctx context.Context, teamID, userID string) (bool, error)
}

func shallPass(r *http.Request, db RelationDB, rUserID, rTeamID string) (bool, error) {
	userID, errUserID := util.UserIDFromRequest(r)
	teamID, errTeamID := util.TeamIDFromRequest(r)

	if errUserID == util.ErrDoesNotExistUserID && errTeamID == util.ErrDoesNotExistUserID {
		return true, nil
	}

	isUser := userID == rUserID
	isParent, err := db.IsParentUser(r.Context(), userID, rUserID)
	if err != nil {
		return false, err
	}
	if errUserID == nil && errTeamID == util.ErrDoesNotExistTeamID {
		return isUser || isParent, nil
	}
	isMember, err := db.IsTeamMember(r.Context(), teamID, rTeamID)
	if err != nil {
		return false, err
	}
	isOwner, err := db.IsTeamOwner(r.Context(), teamID, rUserID)
	if err != nil {
		return false, err
	}
	if errUserID == util.ErrDoesNotExistUserID && errTeamID == nil {
		return isMember || isOwner, nil
	}

	return (isUser || isParent) && (isMember || isOwner), nil
}

func Auth(v Validator, db RelationDB) mux.MiddlewareFunc {
	logger := logrus.WithField("component", "auth-middleware")
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger = logrus.WithField("path", r.URL.Path)
			token, err := jwt.ExtractToken(r)
			if err != nil {
				logger.Error(err)
				w.WriteHeader(http.StatusForbidden)
				return
			}
			userID, teamID, err := v.Valid(token)
			if err != nil {
				logger.Error(err)
				w.WriteHeader(http.StatusForbidden)
				return
			}
			allowed, err := shallPass(r, db, userID, teamID)
			if err != nil {
				logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !allowed {
				logger.Error("access denied!")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func Logging() mux.MiddlewareFunc {
	logger := logrus.WithField("component", "log-middleware")
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger = logger.WithField("path", r.URL.String())
			logger.Info("new request")
			h.ServeHTTP(w, r)
			logger.Info("end request")
		})
	}
}
