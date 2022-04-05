package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	v1 "github.com/MninaTB/vacadm/api/v1"
	"github.com/MninaTB/vacadm/api/v1/team"
	"github.com/MninaTB/vacadm/api/v1/user"
	"github.com/MninaTB/vacadm/api/v1/vaccation"
	vaccationrequest "github.com/MninaTB/vacadm/api/v1/vaccation_request"
	vaccationressources "github.com/MninaTB/vacadm/api/v1/vaccation_ressource"
	"github.com/MninaTB/vacadm/assets/swagger"
	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/database/inmemory"
	"github.com/MninaTB/vacadm/pkg/jwt"
	"github.com/MninaTB/vacadm/pkg/middleware"
	"github.com/MninaTB/vacadm/pkg/model"
)

func main() {
	var (
		// TODO: set default value to false
		swaggerEnabled = flag.Bool("swagger.enable", true, "enables /swagger endpoint")
	)
	flag.Parse()
	router := mux.NewRouter()

	logger := logrus.New().WithField("api", "v1")

	var db database.Database = inmemory.NewInmemoryDB()

	if *swaggerEnabled {
		logger.Info("swagger endpoint \"/swagger\" enabled")
		router.Path("/swagger").Methods(http.MethodGet).HandlerFunc(swagger.Index)
		router.Path("/swagger/api.yaml").Methods(http.MethodGet).HandlerFunc(v1.API)
		router.Handle("/swagger/{rest}",
			http.StripPrefix("/swagger/", http.FileServer(http.FS(swagger.Get()))))
	}

	usrSvc := user.NewUserService(db, logger)
	router.Path("/user").Methods(http.MethodPut).HandlerFunc(usrSvc.Create)
	router.Path("/user/{userID}").Methods(http.MethodGet).HandlerFunc(usrSvc.GetByID)
	router.Path("/user").Methods(http.MethodGet).HandlerFunc(usrSvc.List)
	router.Path("/user/{userID}").Methods(http.MethodPatch).HandlerFunc(usrSvc.Update)
	router.Path("/user/{userID}").Methods(http.MethodDelete).HandlerFunc(usrSvc.Delete)

	teamSvc := team.NewTeamService(db, logger)
	router.Path("/team").Methods(http.MethodPut).HandlerFunc(teamSvc.Create)
	router.Path("/team/{teamID}").Methods(http.MethodGet).HandlerFunc(teamSvc.GetByID)
	router.Path("/team/{teamID}/list-users").Methods(http.MethodGet).HandlerFunc(teamSvc.ListTeamUsers)
	// team/{teamID}/vaccation/ressources
	// team/{teamID}/vaccation
	// router.Path("/team/{teamID}/capacity").Methods(http.MethodGet).HandlerFunc(teamSvc.ListCapacity)
	router.Path("/team").Methods(http.MethodGet).HandlerFunc(teamSvc.List)
	router.Path("/team/{teamID}").Methods(http.MethodPatch).HandlerFunc(teamSvc.Update)
	router.Path("/team/{teamID}").Methods(http.MethodDelete).HandlerFunc(teamSvc.Delete)

	vacSvc := vaccation.NewVaccation(db, logger)
	router.Path("/user/{userID}/vaccation/{vaccationID}").Methods(http.MethodGet).HandlerFunc(vacSvc.GetByID)
	router.Path("/user/{userID}/vaccation").Methods(http.MethodGet).HandlerFunc(vacSvc.List)
	router.Path("/user/{userID}/vaccation/{vaccationID}").Methods(http.MethodDelete).HandlerFunc(vacSvc.Delete)

	vacReqSvc := vaccationrequest.NewVaccationRequest(db, logger)
	router.Path("/user/{userID}/vaccation/request").Methods(http.MethodPut).HandlerFunc(vacReqSvc.Create)
	router.Path("/user/{userID}/vaccation/request/{vaccation-requestID}").Methods(http.MethodGet).HandlerFunc(vacReqSvc.GetByID)
	router.Path("/user/{userID}/vaccation/request").Methods(http.MethodGet).HandlerFunc(vacReqSvc.List)
	router.Path("/user/{userID}/vaccation/request/{vaccation-requestID}").Methods(http.MethodPatch).HandlerFunc(vacReqSvc.Update)
	router.Path("/user/{userID}/vaccation/request/{vaccation-requestID}").Methods(http.MethodDelete).HandlerFunc(vacReqSvc.Delete)

	vacResSvc := vaccationressources.NewVaccationRessource(db, logger)
	router.Path("/user/{userID}/vaccation/ressource").Methods(http.MethodPut).HandlerFunc(vacResSvc.Create)
	router.Path("/user/{userID}/vaccation/ressource/{vaccation-ressourceID}").Methods(http.MethodGet).HandlerFunc(vacResSvc.GetByID)
	router.Path("/user/{userID}/vaccation/ressource").Methods(http.MethodGet).HandlerFunc(vacResSvc.List)
	router.Path("/user/{userID}/vaccation/ressource/{vaccation-ressourceID}").Methods(http.MethodPatch).HandlerFunc(vacResSvc.Update)
	router.Path("/user/{userID}/vaccation/ressource/{vaccation-ressourceID}").Methods(http.MethodDelete).HandlerFunc(vacResSvc.Delete)

	jwtKey := []byte("my-secret")
	t := jwt.NewTokenizer(jwtKey, 365*24*time.Hour)
	router.Use(middleware.Logging())
	//router.Use(middleware.Auth(t, database.NewRelationDB(db)))
	u := model.User{
		FirstName: "nina",
		LastName:  "olear",
		Email:     "admin@inform.de",
	}
	newUser, err := db.CreateUser(context.Background(), &u)
	if err != nil {
		logger.Fatalln(err)
	}
	token, err := t.Generate(newUser)
	if err != nil {
		logger.Fatalln(err)
	}
	team := model.Team{
		OwnerID: u.ID,
		Name:    "root-Team",
	}
	_, err = db.CreateTeam(context.Background(), &team)
	if err != nil {
		logger.Fatalln(err)
	}
	logger.Info("admin token is: ", token)

	const addr = ":8080"
	log.Println("start server on port", addr)
	logger.Info("start server on port", addr)
	server := http.Server{
		Addr:    addr,
		Handler: router,
	}
	server.ListenAndServe()
}
