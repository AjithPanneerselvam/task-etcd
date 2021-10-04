package main

import (
	"net/http"

	"github.com/AjithPanneerselvam/todo/config"
	"github.com/AjithPanneerselvam/todo/db"
	"github.com/AjithPanneerselvam/todo/router"
	"github.com/AjithPanneerselvam/todo/store/user"
	"github.com/AjithPanneerselvam/todo/util"

	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := config.Load()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	log.Info("loaded environment configs")

	log.Infof("log level: %v", config.LogLevel)
	util.SetupLog(config.LogLevel)

	db, err := db.NewEtcdClient(config.EtcdURLS)
	if err != nil {
		log.Fatal("error creating a etcd client", err)
	}
	log.Info("etcd client instantiated")

	userStore := user.New(db)

	router := router.NewRouter()
	router.AddRoutes(config, userStore)

	log.Infof("starting server at port %v", config.ListenPort)
	http.ListenAndServe(":"+config.ListenPort, router)
}
