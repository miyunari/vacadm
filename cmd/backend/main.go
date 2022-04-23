package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/MninaTB/vacadm/api/token"
	v1 "github.com/MninaTB/vacadm/api/v1"
	"github.com/MninaTB/vacadm/assets/swagger"
	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/database/inmemory"
	"github.com/MninaTB/vacadm/pkg/database/mariadb"
	"github.com/MninaTB/vacadm/pkg/jwt"
	"github.com/MninaTB/vacadm/pkg/middleware"
	"github.com/MninaTB/vacadm/pkg/model"
	"github.com/MninaTB/vacadm/pkg/notify"
	"github.com/MninaTB/vacadm/pkg/version"
)

func main() {
	var (
		swaggerEnabled = flag.Bool("swagger.enable", false, "enables /swagger endpoint")
		initRoot       = flag.Bool("init.root", false, "create root user on startup")
		address        = flag.String("address", "localhost:8080", "ip:port")
		jwtKey         = flag.String("secret", "", "secret for jwt token")
		sqlConnStr     = flag.String("sql.conn", "", `sql connection str. user:password@/dbname
		example: root:my-secret-pw@(127.0.0.1:3306)/test?parseTime=true`)
		srvTimeout   = flag.Duration("timeout", time.Minute, "server timeout")
		smtpHost     = flag.String("smtp.host", "", "address of smtp server")
		smtpPort     = flag.String("smtp.port", "", "port of smtp server")
		smtpUser     = flag.String("smtp.user", "", "smtp user mail address")
		smtpPassword = flag.String("smtp.password", "", "smtp user password")
	)
	flag.Parse()

	logger := logrus.New()
	logger.Info(version.Version())

	var db database.Database = inmemory.NewInmemoryDB()
	if *sqlConnStr != "" {
		sqlDB, err := sql.Open("mysql", *sqlConnStr)
		if err != nil {
			logger.Fatal(err)
		}
		defer sqlDB.Close()
		db = mariadb.NewMariaDB(sqlDB)
	}

	var notifier notify.Notifier = notify.NewNoopNotifier()
	if *smtpHost != "" && *smtpPort != "" && *smtpUser != "" {
		logger.WithFields(logrus.Fields{
			"host": *smtpHost,
			"port": *smtpPort,
		}).Infof("enabled smtp notifier, address: %s", *smtpUser)
		notifier = notify.NewMailer(*smtpHost, *smtpPort, *smtpUser, *smtpPassword, db)
	}

	router := mux.NewRouter()
	secret := []byte(*jwtKey)
	if len(secret) == 0 {
		logger.Fatal("missing jwt secret")
	}
	t := jwt.NewTokenizer(secret, 365*24*time.Hour)
	router.Path("/token/new/{userID}").Methods(http.MethodGet).HandlerFunc(token.NewTokenService(db, t).Refresh)
	apiv1 := v1.NewServer(db, notifier, t, middleware.Logging(), middleware.Auth(t, database.NewRelationDB(db)))
	const pathPrefixV1 = "/v1"
	// HACK: allow sub routes on v1 router.
	router.Handle(fmt.Sprintf("%s/{dummy1}", pathPrefixV1), http.StripPrefix(pathPrefixV1, apiv1))
	router.Handle(fmt.Sprintf("%s/{dummy1}/{dummy2}", pathPrefixV1), http.StripPrefix(pathPrefixV1, apiv1))
	router.Handle(fmt.Sprintf("%s/{dummy1}/{dummy2}/{dummy3}", pathPrefixV1), http.StripPrefix(pathPrefixV1, apiv1))
	router.Handle(fmt.Sprintf("%s/{dummy1}/{dummy2}/{dummy3}/{dummy4}", pathPrefixV1), http.StripPrefix(pathPrefixV1, apiv1))

	if *swaggerEnabled {
		logger.Info("swagger endpoint \"/swagger\" enabled")
		router.Path("/swagger").Methods(http.MethodGet).HandlerFunc(swagger.Index)
		router.Path("/swagger/api.yaml").Methods(http.MethodGet).HandlerFunc(v1.API)
		router.Handle("/swagger/{dummy}",
			http.StripPrefix("/swagger/", http.FileServer(http.FS(swagger.Get()))))
	}

	if *initRoot {
		u := &model.User{
			FirstName: "nina",
			LastName:  "olear",
			Email:     "admin@inform.de",
		}
		u, err := db.CreateUser(context.Background(), u)
		if err != nil {
			logger.Fatalln(err)
		}
		rootToken, err := t.Generate(u)
		if err != nil {
			logger.Fatalln(err)
		}
		logger.Info("admin token is: ", rootToken)
	}

	server := &http.Server{
		Addr:              *address,
		Handler:           router,
		ReadTimeout:       *srvTimeout,
		WriteTimeout:      *srvTimeout,
		IdleTimeout:       *srvTimeout,
		ReadHeaderTimeout: *srvTimeout,
	}
	logger.Info("Listen and serve on address: ", *address)
	if err := server.ListenAndServe(); err != nil {
		logger.Panic(err)
	}
}
