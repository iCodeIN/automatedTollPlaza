package appcontext

import (
	"automatedTollPlaze/config"
	"automatedTollPlaze/pkg/platform/db"
	"time"
)

// AppContext ..
type AppContext struct {
	DbClient  db.Service
	Config    config.Cfg
	StartTime time.Time
}
