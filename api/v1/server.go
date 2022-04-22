package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/MninaTB/vacadm/api/v1/team"
	"github.com/MninaTB/vacadm/api/v1/user"
	"github.com/MninaTB/vacadm/api/v1/vacation"
	vacationrequest "github.com/MninaTB/vacadm/api/v1/vacation_request"
	vacationresources "github.com/MninaTB/vacadm/api/v1/vacation_resource"
	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/notify"
)

// Tokenizer implements methods to verify auth tokens.
type TokenValidator interface {
	// Valid if a token is valid, userID and teamID are returned.
	// if a token is invalid, an error is returned.
	Valid(token string) (userID string, teamID string, err error)
}

type server struct {
	logger   logrus.FieldLogger
	db       database.Database
	mw       []mux.MiddlewareFunc
	tv       TokenValidator
	notifier notify.Notifier
}

// NewServer returns a new http.Handler.
func NewServer(
	db database.Database,
	notifier notify.Notifier,
	tokenValidator TokenValidator,
	middleware ...mux.MiddlewareFunc,
) http.Handler {
	return &server{
		logger:   logrus.New().WithField("api", "v1"),
		mw:       middleware,
		db:       db,
		notifier: notifier,
		tv:       tokenValidator,
	}
}

// ServeHTTP dispatches the handler registered in the matched route.
//
// When there is a match, the route variables can be retrieved calling
// mux.Vars(request).
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	usrSvc := user.NewUserService(s.db, s.logger)

	teamSvc := team.NewTeamService(s.db, s.logger, s.tv)

	vacSvc := vacation.NewVacationService(s.db, s.logger)

	vacReqSvc := vacationrequest.NewVacationRequestService(s.db, s.notifier, s.logger)

	vacResSvc := vacationresources.NewVacationResourceService(s.db, s.logger)

	router := mux.NewRouter()
	router.Path("/user").Methods(http.MethodPut).HandlerFunc(usrSvc.Create)
	router.Path("/user/{userID}").Methods(http.MethodGet).HandlerFunc(usrSvc.GetByID)
	router.Path("/user").Methods(http.MethodGet).HandlerFunc(usrSvc.List)
	router.Path("/user/{userID}").Methods(http.MethodPatch).HandlerFunc(usrSvc.Update)
	router.Path("/user/{userID}").Methods(http.MethodDelete).HandlerFunc(usrSvc.Delete)

	router.Path("/team").Methods(http.MethodPut).HandlerFunc(teamSvc.Create)
	router.Path("/team/{teamID}").Methods(http.MethodGet).HandlerFunc(teamSvc.GetByID)
	router.Path("/team/{teamID}/list-users").Methods(http.MethodGet).HandlerFunc(teamSvc.ListTeamUsers)
	router.Path("/team/list-vacation").Methods(http.MethodPost).HandlerFunc(teamSvc.ListCapacity)

	router.Path("/team").Methods(http.MethodGet).HandlerFunc(teamSvc.List)
	router.Path("/team/{teamID}").Methods(http.MethodPatch).HandlerFunc(teamSvc.Update)
	router.Path("/team/{teamID}").Methods(http.MethodDelete).HandlerFunc(teamSvc.Delete)

	router.Path("/user/{userID}/vacation/{vacationID}").Methods(http.MethodGet).HandlerFunc(vacSvc.GetByID)
	router.Path("/user/{userID}/vacation").Methods(http.MethodGet).HandlerFunc(vacSvc.List)
	router.Path("/user/{userID}/vacation/{vacationID}").Methods(http.MethodDelete).HandlerFunc(vacSvc.Delete)

	router.Path("/user/{userID}/vacation/request").Methods(http.MethodPut).HandlerFunc(vacReqSvc.Create)
	router.Path("/user/{userID}/vacation/request/{vacation-requestID}").Methods(http.MethodGet).HandlerFunc(vacReqSvc.GetByID)
	router.Path("/user/{userID}/vacation/request/{vacation-requestID}/approve/{parentID}").Methods(http.MethodPut).HandlerFunc(vacReqSvc.Approve)
	router.Path("/user/{userID}/vacation/request").Methods(http.MethodGet).HandlerFunc(vacReqSvc.List)
	router.Path("/user/{userID}/vacation/request/{vacation-requestID}").Methods(http.MethodPatch).HandlerFunc(vacReqSvc.Update)
	router.Path("/user/{userID}/vacation/request/{vacation-requestID}").Methods(http.MethodDelete).HandlerFunc(vacReqSvc.Delete)

	router.Path("/user/{userID}/vacation/resource").Methods(http.MethodPut).HandlerFunc(vacResSvc.Create)
	router.Path("/user/{userID}/vacation/resource/{vacation-resourceID}").Methods(http.MethodGet).HandlerFunc(vacResSvc.GetByID)
	router.Path("/user/{userID}/vacation/resource").Methods(http.MethodGet).HandlerFunc(vacResSvc.List)
	router.Path("/user/{userID}/vacation/resource/{vacation-resourceID}").Methods(http.MethodPatch).HandlerFunc(vacResSvc.Update)
	router.Path("/user/{userID}/vacation/resource/{vacation-resourceID}").Methods(http.MethodDelete).HandlerFunc(vacResSvc.Delete)
	if s.mw != nil {
		router.Use(s.mw...)
	}

	router.ServeHTTP(w, r)
}
