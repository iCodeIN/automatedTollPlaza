package mongo

import (
	"automatedTollPlaze/pkg/platform/db"
	"context"

	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// handler ..
type handler struct {
	dbClient client
}

type client interface {
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	Database(name string, opts ...*options.DatabaseOptions) *mgo.Database
}

// Cfg ..
type Cfg struct {
	Host string `json:"host"`
}

// NewMongoClient ..
func NewMongoClient(ctx context.Context, cfg Cfg) db.Service {
	mgoClient, err := mgo.Connect(ctx, options.Client().ApplyURI(cfg.Host))
	if err != nil {
		panic(err)
	}
	return &handler{
		dbClient: mgoClient,
	}
}

// getDatabase ..
func (h *handler) getDatabase(dbName string) *mgo.Database {
	return h.dbClient.Database(dbName)
}

// Health ..
func (h *handler) Health(ctx context.Context) error {
	return h.dbClient.Ping(ctx, readpref.Primary())
}

// Find ..
func (h *handler) Find(ctx context.Context, database, collection string, filter map[string]interface{}, projection map[string]interface{}, result interface{}) error {
	opts := options.Find()
	opts.Projection = projection
	cursor, err := h.getDatabase(database).Collection(collection).Find(ctx, filter, options.Find(), opts)
	if err != nil {
		return err
	}
	return cursor.All(ctx, &result)
}

// UpsertMany ..
func (h *handler) UpsertMany(ctx context.Context, database, collection string, filter map[string]interface{}, updateData map[string]interface{}, upsert bool) error {
	opts := options.Update()
	opts.SetUpsert(upsert)
	updateData = map[string]interface{}{
		"$set": updateData,
	}
	_, err := h.getDatabase(database).Collection(collection).UpdateMany(ctx, filter, updateData, opts)
	return err
}
