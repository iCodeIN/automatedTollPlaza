package appcontext

import (
	"automatedTollPlaze/pkg/platform/db"
	"time"
)

// AppContext ..
type AppContext struct {
	DbClient  db.Service
	StartTime time.Time
}
