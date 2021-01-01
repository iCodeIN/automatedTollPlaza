package db

import (
	"context"
)

// Service is the contract for all of the cache backends that are supported by
// this package
type Service interface {
	Health(ctx context.Context) error
	Find(ctx context.Context, database, collection string, filter map[string]interface{}, projection map[string]interface{}, result interface{}) error
	UpsertMany(ctx context.Context, database, collection string, filter map[string]interface{}, updateData map[string]interface{}, upsert bool) error
}
