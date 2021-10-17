package main

import (
	"net/http"

	"github.com/AjithPanneerselvam/task-etcd/config"
	"github.com/AjithPanneerselvam/task-etcd/db"
	"github.com/AjithPanneerselvam/task-etcd/router"
	"github.com/AjithPanneerselvam/task-etcd/store/task"
	"github.com/AjithPanneerselvam/task-etcd/util"

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

	taskStore := task.New(db)

	router := router.NewRouter()
	router.AddRoutes(config, taskStore)

	log.Infof("starting server at port %v", config.ListenPort)
	http.ListenAndServe(":"+config.ListenPort, router)
}
