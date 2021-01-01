package appcontext

import "automatedTollPlaze/pkg/platform/db"

// AppContext ..
type AppContext struct {
	DbClient db.Service
}
