package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/MninaTB/vacadm/api/v1/team"
	"github.com/MninaTB/vacadm/api/v1/user"
	"github.com/MninaTB/vacadm/api/v1/vaccation"
	vaccationrequest "github.com/MninaTB/vacadm/api/v1/vaccation_request"
	vaccationressources "github.com/MninaTB/vacadm/api/v1/vaccation_ressource"
	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/database/inmemory"
	"github.com/MninaTB/vacadm/pkg/jwt"
	"github.com/MninaTB/vacadm/pkg/middleware"
	"github.com/MninaTB/vacadm/pkg/model"
)

func main() {
	router := mux.NewRouter()

	logger := logrus.New().WithField("api", "v1")

	var db database.Database = inmemory.NewInmemoryDB()

	usrSvc := user.NewUserService(db, logger)
	router.Path("/user").Methods(http.MethodPut).HandlerFunc(usrSvc.Create)
	router.Path("/user/{userID}").Methods(http.MethodGet).HandlerFunc(usrSvc.GetByID)
	router.Path("/user").Methods(http.MethodGet).HandlerFunc(usrSvc.List)
	router.Path("/user").Methods(http.MethodPatch).HandlerFunc(usrSvc.Update)
	router.Path("/user/{userID}").Methods(http.MethodDelete).HandlerFunc(usrSvc.Delete)

	teamSvc := team.NewTeamService(db, logger)
	router.Path("/team").Methods(http.MethodPut).HandlerFunc(teamSvc.Create)
	router.Path("/team/{teamID}").Methods(http.MethodGet).HandlerFunc(teamSvc.GetByID)
	router.Path("/team/{teamID}/list-users").Methods(http.MethodGet).HandlerFunc(teamSvc.ListTeamUsers)
	// team/{teamID}/vaccation/ressources
	// team/{teamID}/vaccation
	router.Path("/team/{teamID}/capacity").Methods(http.MethodGet).HandlerFunc(teamSvc.ListCapacity)
	router.Path("/team").Methods(http.MethodGet).HandlerFunc(teamSvc.List)
	router.Path("/team").Methods(http.MethodPatch).HandlerFunc(teamSvc.Update)
	router.Path("/team/{teamID}").Methods(http.MethodDelete).HandlerFunc(teamSvc.Delete)

	vacSvc := vaccation.NewVaccation(db, logger)
	router.Path("/user/{userID}/vaccation/{vaccationID}").Methods(http.MethodGet).HandlerFunc(vacSvc.GetByID)
	router.Path("/user/{userID}/vaccation").Methods(http.MethodGet).HandlerFunc(vacSvc.List)
	router.Path("/user/{userID}/vaccation/{vaccationID}").Methods(http.MethodDelete).HandlerFunc(vacSvc.Delete)

	vacReqSvc := vaccationrequest.NewVaccationRequest(db, logger)
	router.Path("/user/{userID}/vaccation/request").Methods(http.MethodPut).HandlerFunc(vacReqSvc.Create)
	router.Path("/user/{userID}/vaccation/request/{vaccation-requestID}").Methods(http.MethodGet).HandlerFunc(vacReqSvc.GetByID)
	router.Path("/user/{userID}/vaccation/request").Methods(http.MethodGet).HandlerFunc(vacReqSvc.List)
	router.Path("/user/{userID}/vaccation/request").Methods(http.MethodPatch).HandlerFunc(vacReqSvc.Update)
	router.Path("/user/{userID}/vaccation/request/{vaccation-requestID}").Methods(http.MethodDelete).HandlerFunc(vacReqSvc.Delete)

	vacResSvc := vaccationressources.NewVaccationRessource(db, logger)
	router.Path("/user/{userID}/vaccation/ressource").Methods(http.MethodPut).HandlerFunc(vacResSvc.Create)
	router.Path("/user/{userID}/vaccation/ressource/{vaccation-ressourceID}").Methods(http.MethodGet).HandlerFunc(vacResSvc.GetByID)
	router.Path("/user/{userID}/vaccation/ressource").Methods(http.MethodGet).HandlerFunc(vacResSvc.List)
	router.Path("/user/{userID}/vaccation/ressource").Methods(http.MethodPatch).HandlerFunc(vacResSvc.Update)
	router.Path("/user/{userID}/vaccation/ressource/{vaccation-ressourceID}").Methods(http.MethodDelete).HandlerFunc(vacResSvc.Delete)

	jwtKey := []byte("my-secret")
	t := jwt.NewTokenizer(jwtKey, 365*24*time.Hour)
	router.Use(middleware.Logging())
	router.Use(middleware.Auth(t))
	u := model.User{
		FirstName: "nina",
		LastName:  "olear",
		Email:     "admin@inform.de",
	}

	newUser, err := db.CreateUser(&u)
	if err != nil {
		logger.Fatalln(err)
	}
	token, err := t.Generate(newUser)
	if err != nil {
		logger.Fatalln(err)
	}

	logger.Info("admin token is: ", token)

	const addr = ":8080"
	log.Println("Starte Server auf Port", addr)
	logger.Info("Starte Server auf Port", addr)
	server := http.Server{
		Addr:    addr,
		Handler: router,
	}
	server.ListenAndServe()
}
