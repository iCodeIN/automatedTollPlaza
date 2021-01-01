package toll

import (
	// tollErr "automatedTollPlaze/pkg/errors"
	"automatedTollPlaze/pkg/platform/appcontext"
	"context"
)

//Service ..
type Service interface {
	IssueToll(context.Context) error
}

// service ..
type handler struct {
	AppCtx *appcontext.AppContext
}

// NewTollService ..
func NewTollService(ctx context.Context, appCtx *appcontext.AppContext) Service {
	return &handler{
		AppCtx: appCtx,
	}
}

// IssueToll ..
func (s *handler) IssueToll(ctx context.Context) error {
	// return tollErr.NewErrorWithCode("ERR.PANIC.NOTFOUND", "not found")
	return nil
}
