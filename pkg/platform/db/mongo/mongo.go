package mongo

import (
	"automatedTollPlaze/pkg/platform/db"
	"context"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// handler ..
type handler struct {
	dbClient client
}

type client interface {
	Ping() error
	DB(name string) *mgo.Database
}

// Cfg ..
type Cfg struct {
	Host string `json:"host"`
}

// NewMongoClient ..
func NewMongoClient(ctx context.Context, cfg Cfg) db.Service {
	mgoClient, err := mgo.DialWithTimeout(cfg.Host, time.Duration(5)*time.Second)
	if err != nil {
		panic(err)
	}
	return &handler{
		dbClient: mgoClient,
	}
}

// getDatabase ..
func (h *handler) getDatabase(dbName string) *mgo.Database {
	return h.dbClient.DB(dbName)
}

// Health ..
func (h *handler) Health(ctx context.Context) error {
	return h.dbClient.Ping()
}

// Find ..
func (h *handler) Find(ctx context.Context, database, collection string, filter map[string]interface{}, result interface{}) error {
	bsonFilter := bson.M{}
	for key, val := range filter {
		bsonFilter[key] = val
	}
	if err := h.getDatabase(database).C(collection).Find(bsonFilter).One(result); err != nil && err != mgo.ErrNotFound {
		return err
	}
	return nil
}

// Upsert ..
func (h *handler) Upsert(ctx context.Context, database, collection string, filter map[string]interface{}, updateData interface{}) error {
	bsonFilter := bson.M{}
	for key, val := range filter {
		bsonFilter[key] = val
	}
	updateData = map[string]interface{}{
		"$set": updateData,
	}
	_, err := h.getDatabase(database).C(collection).Upsert(bsonFilter, updateData)
	return err
}

// InsertOne ..
func (h *handler) InsertOne(ctx context.Context, database, collection string, data interface{}) error {
	return h.getDatabase(database).C(collection).Insert(data)
}
