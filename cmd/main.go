package main

import (
	"automatedTollPlaze/api"
	"automatedTollPlaze/config"
	"automatedTollPlaze/pkg"
	"automatedTollPlaze/pkg/platform/appcontext"
	"automatedTollPlaze/pkg/platform/db/mongo"
	"automatedTollPlaze/pkg/toll"
	"automatedTollPlaze/utils"
	"context"
	"os"
	"time"

	"github.com/bnkamalesh/webgo/v4"
	"github.com/bnkamalesh/webgo/v4/middleware"

	log "github.com/sirupsen/logrus"
)

// HTTP represents structure of Http Requests
type HTTP struct {
	AppContext *appcontext.AppContext
	server     *webgo.Router
}

func main() {
	// setting logger
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)

	ctx := context.Background()
	mydir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	initCfg := &config.Cfg{}
	if err := utils.ReadFile(mydir+"/config/config.json", utils.FileData{Data: initCfg}); err != nil {
		log.Panic(err)
	}
	appCtx := &appcontext.AppContext{
		DbClient: mongo.NewMongoClient(ctx, mongo.Cfg{
			Host: initCfg.MongoConfig.Host,
		}),
		StartTime: time.Now().Local(),
	}
	apiHandler := api.API{
		AppContext: appCtx,
		Handler: pkg.ServiceHandler{
			TollHandler: toll.NewTollService(ctx, appCtx),
		},
	}
	cfg := &webgo.Config{
		Host:         initCfg.ServerConfig.Host,
		Port:         initCfg.ServerConfig.Port,
		ReadTimeout:  initCfg.ServerConfig.ReadTimeout * time.Second,
		WriteTimeout: initCfg.ServerConfig.WriteTimeout * time.Second,
	}
	router := webgo.NewRouter(cfg, apiHandler.Routes())
	router.NotFound = api.NotFound()
	httpServer := &HTTP{
		AppContext: appCtx,
		server:     router,
	}
	if initCfg.HTTPLog {
		router.Use(middleware.AccessLog)
	}
	log.Info("Server has started on port ", initCfg.ServerConfig.Port)
	httpServer.server.Start()
}
