package main

import (
	"automatedTollPlaze/api"
	"automatedTollPlaze/config"
	"automatedTollPlaze/pkg"
	"automatedTollPlaze/pkg/platform/appcontext"
	"automatedTollPlaze/pkg/platform/db/mongo"
	"automatedTollPlaze/pkg/toll"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/bnkamalesh/webgo/v4"
	"github.com/bnkamalesh/webgo/v4/middleware"

	log "github.com/sirupsen/logrus"
)

func main() {
	// setting logger
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)

	ctx := context.Background()
	mydir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
		return
	}
	byt, err := ioutil.ReadFile(mydir + "/config/config.json")
	if err != nil {
		log.Panic(err)
		return
	}
	initCfg := config.Cfg{}
	if err := json.Unmarshal(byt, &initCfg); err != nil {
		log.Panic(err)
		return
	}
	appCtx := &appcontext.AppContext{
		DbClient: mongo.NewMongoClient(ctx, mongo.Cfg{
			Host: initCfg.MongoConfig.Host,
		}),
		Config:    initCfg,
		StartTime: time.Now().Local(),
	}
	apiHandler := pkg.ServiceHandler{
		TollHandler: toll.NewTollService(ctx, appCtx),
	}
	cfg := &webgo.Config{
		Host:         initCfg.ServerConfig.Host,
		Port:         initCfg.ServerConfig.Port,
		ReadTimeout:  initCfg.ServerConfig.ReadTimeout * time.Second,
		WriteTimeout: initCfg.ServerConfig.WriteTimeout * time.Second,
	}
	httpServer := &api.HTTP{
		AppContext: appCtx,
		APIHandler: &apiHandler,
	}
	httpServer.Server = webgo.NewRouter(cfg, httpServer.Routes())
	httpServer.Server.NotFound = api.NotFound()

	if initCfg.HTTPLog {
		httpServer.Server.Use(middleware.AccessLog)
	}
	log.Info("Server has started on port ", initCfg.ServerConfig.Port)
	httpServer.Server.Start()
}
