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
)

// HTTP represents structure of Http Requests
type HTTP struct {
	AppContext *appcontext.AppContext
	server     *webgo.Router
}

func main() {
	ctx := context.Background()
	mydir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	initCfg := &config.Cfg{}
	if err := utils.ReadFile(mydir+"/config/config.json", utils.FileData{Data: initCfg}); err != nil {
		panic(err)
	}
	appCtx := &appcontext.AppContext{
		DbClient: mongo.NewMongoClient(ctx, mongo.Cfg{
			Host: initCfg.MongoConfig.Host,
		}),
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
	httpServer.server.Start()
}
