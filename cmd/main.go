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

	// reading the current directly path
	mydir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
		return
	}

	// reading the configuration file
	byt, err := ioutil.ReadFile(mydir + "/config/config.json")
	if err != nil {
		log.Panic(err)
		return
	}

	// unmarshalling/mapping to configuration variable
	initCfg := config.Cfg{}
	if err := json.Unmarshal(byt, &initCfg); err != nil {
		log.Panic(err)
		return
	}

	// creating an application level context with necessary dependencies
	appCtx := &appcontext.AppContext{
		DbClient: mongo.NewMongoClient(ctx, mongo.Cfg{
			Host: initCfg.MongoConfig.Host,
		}),
		Config:    initCfg,
		StartTime: time.Now().UTC(),
	}

	// creating a service handler with toll service
	apiHandler := pkg.ServiceHandler{
		TollHandler: toll.NewTollService(ctx, appCtx),
	}

	// reading the configurations to start the web server
	cfg := &webgo.Config{
		Host:         initCfg.ServerConfig.Host,
		Port:         initCfg.ServerConfig.Port,
		ReadTimeout:  initCfg.ServerConfig.ReadTimeout * time.Second,
		WriteTimeout: initCfg.ServerConfig.WriteTimeout * time.Second,
	}

	// injecting the service handler into http server
	httpServer := &api.HTTP{
		AppContext: appCtx,
		APIHandler: &apiHandler,
	}

	// initializing the server with the routes
	httpServer.Server = webgo.NewRouter(cfg, httpServer.Routes())

	// initializing the NotFound http handler
	httpServer.Server.NotFound = api.NotFound()

	// check if webserver required request logs or not.
	if initCfg.HTTPLog {
		httpServer.Server.Use(middleware.AccessLog)
	}

	// starting the web server.
	log.Info("Server has started on port ", initCfg.ServerConfig.Port)
	httpServer.Server.Start()
}
