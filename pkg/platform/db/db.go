package db

import (
	"context"
)

// Service is the contract for all of the cache backends that are supported by
// this package
type Service interface {
	Health(ctx context.Context) error
	InsertOne(ctx context.Context, database, collection string, data interface{}) error
	Find(ctx context.Context, database, collection string, filter map[string]interface{}, result interface{}) error
	Upsert(ctx context.Context, database, collection string, filter map[string]interface{}, updateData interface{}) error
}
