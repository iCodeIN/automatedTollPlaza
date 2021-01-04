package appcontext

import (
	"automatedTollPlaze/config"
	"automatedTollPlaze/pkg/platform/db"
	"time"
)

// AppContext has the application level dependencies
type AppContext struct {
	DbClient  db.Service
	Config    config.Cfg
	StartTime time.Time
}
