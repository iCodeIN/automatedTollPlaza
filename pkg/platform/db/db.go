package db

import (
	"automatedTollPlaze/pkg/errors"
	"context"
)

// Params ..
type Params struct {
	Database   string
	Collection string
	Filter     map[string]interface{}
	UpsertData interface{}
	Result     interface{}
	Pagination *Paginate
}

// Paginate ..
type Paginate struct {
	Start, Limit int
}

// Service is the contract for all of the cache backends that are supported by
// this package
type Service interface {
	Health(ctx context.Context) error
	InsertOne(ctx context.Context, params Params) error
	FindOne(ctx context.Context, params Params) error
	FindAll(ctx context.Context, params Params) error
	Count(ctx context.Context, params Params) int
	Upsert(ctx context.Context, params Params) error
}

var (
	// ErrNotFound ..
	ErrNotFound = errors.NewErrorWithCode(
		"ERR.DB.NOT_FOUND",
		"Content doesn't exists",
	)
)
