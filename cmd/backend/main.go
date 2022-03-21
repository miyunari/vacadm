package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/MninaTB/vacadm/api/v1/team"
	"github.com/MninaTB/vacadm/api/v1/user"
	"github.com/MninaTB/vacadm/api/v1/vaccation"
	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/database/inmemory"
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
	router.Path("/team").Methods(http.MethodGet).HandlerFunc(teamSvc.List)
	router.Path("/team").Methods(http.MethodPatch).HandlerFunc(teamSvc.Update)
	router.Path("/team/{teamID}").Methods(http.MethodDelete).HandlerFunc(teamSvc.Delete)

	vacSvc := vaccation.NewVaccation(db, logger)
	router.Path("/vaccation").Methods(http.MethodPut).HandlerFunc(vacSvc.Create)
	router.Path("/vaccation/{vaccationID}").Methods(http.MethodGet).HandlerFunc(vacSvc.GetByID)
	router.Path("/vaccation").Methods(http.MethodGet).HandlerFunc(vacSvc.List)
	router.Path("/vaccation").Methods(http.MethodPatch).HandlerFunc(vacSvc.Update)
	router.Path("/vaccation/{vaccationID}").Methods(http.MethodDelete).HandlerFunc(vacSvc.Delete)

	const addr = ":8080"
	log.Println("Starte Server auf Port", addr)
	logger.Info("Starte Server auf Port", addr)
	server := http.Server{
		Addr:    addr,
		Handler: router,
	}
	server.ListenAndServe()
}
