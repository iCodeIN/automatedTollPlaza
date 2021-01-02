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

// FindOne ..
func (h *handler) FindOne(ctx context.Context, param db.Params) error {
	bsonFilter := bson.M{}
	for key, val := range param.Filter {
		bsonFilter[key] = val
	}
	if err := h.getDatabase(param.Database).C(param.Collection).Find(bsonFilter).One(param.Result); err != nil && err != mgo.ErrNotFound {
		return err
	}
	return nil
}

// FindAll ..
func (h *handler) FindAll(ctx context.Context, param db.Params) error {
	bsonFilter := bson.M{}
	for key, val := range param.Filter {
		bsonFilter[key] = val
	}
	query := h.getDatabase(param.Database).C(param.Collection).Find(bsonFilter)
	if param.Pagination != nil {
		query.Skip(param.Pagination.Start)
		query.Limit(param.Pagination.Limit)
	}
	if err := query.All(param.Result); err != nil {
		return err
	}
	return nil
}

// Count ..
func (h *handler) Count(ctx context.Context, params db.Params) int {
	bsonFilter := bson.M{}
	for key, val := range params.Filter {
		bsonFilter[key] = val
	}
	count, _ := h.getDatabase(params.Database).C(params.Collection).Find(bsonFilter).Count()
	return count
}

// Upsert ..
func (h *handler) Upsert(ctx context.Context, params db.Params) error {
	bsonFilter := bson.M{}
	for key, val := range params.Filter {
		bsonFilter[key] = val
	}
	params.UpsertData = map[string]interface{}{
		"$set": params.UpsertData,
	}
	_, err := h.getDatabase(params.Database).C(params.Collection).Upsert(bsonFilter, params.UpsertData)
	return err
}

// InsertOne ..
func (h *handler) InsertOne(ctx context.Context, params db.Params) error {
	return h.getDatabase(params.Database).C(params.Collection).Insert(params.UpsertData)
}
