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

	v1 "github.com/MninaTB/vacadm/api/v1"
	"github.com/MninaTB/vacadm/assets/swagger"
	"github.com/MninaTB/vacadm/pkg/database"
	"github.com/MninaTB/vacadm/pkg/database/inmemory"
	"github.com/MninaTB/vacadm/pkg/database/mariadb"
	"github.com/MninaTB/vacadm/pkg/jwt"
	"github.com/MninaTB/vacadm/pkg/middleware"
	"github.com/MninaTB/vacadm/pkg/model"
)

func main() {
	var (
		// TODO: set swagger default value to false
		swaggerEnabled = flag.Bool("swagger.enable", true, "enables /swagger endpoint")
		// TODO: set initRoot default value to false
		initRoot = flag.Bool("init.root", true, "create root user on startup")
		address  = flag.String("address", "localhost:8080", "ip:port")
		// TODO: remove weak default secret
		jwtKey     = flag.String("secret", "my-secret", "secret for jwt token")
		sqlConnStr = flag.String("sql.conn", "", "sql connection str. user:password@/dbname example: root:my-secret-pw@(127.0.0.1:3306)/test?parseTime=true")
		srvTimeout = flag.Duration("timeout", 30*time.Second, "server timeout")
	)
	flag.Parse()

	logger := logrus.New()
	var db database.Database = inmemory.NewInmemoryDB()
	if *sqlConnStr != "" {
		sqlDB, err := sql.Open("mysql", *sqlConnStr)
		if err != nil {
			logger.Fatal(err)
		}
		defer sqlDB.Close()
		db = mariadb.NewMariaDB(sqlDB)
	}

	secret := []byte(*jwtKey)
	if len(secret) == 0 {
		logger.Fatal("missing jwt secret")
	}
	t := jwt.NewTokenizer(secret, 365*24*time.Hour)
	apiv1 := v1.NewServer(db, middleware.Logging(), middleware.Auth(t, database.NewRelationDB(db)))
	router := mux.NewRouter()
	const pathPrefixV1 = "/v1"
	router.Handle(fmt.Sprintf("%s/{dummy}", pathPrefixV1), http.StripPrefix(pathPrefixV1, apiv1))

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
		token, err := t.Generate(u)
		if err != nil {
			logger.Fatalln(err)
		}
		logger.Info("admin token is: ", token)
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
