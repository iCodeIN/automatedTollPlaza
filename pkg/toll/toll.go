package toll

import (
	// tollErr "automatedTollPlaze/pkg/errors"
	"automatedTollPlaze/pkg/platform/appcontext"
	"context"
)

//Service ..
type Service interface {
	IssueToll(ctx context.Context, ticket *TicketToll) error
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
func (s *handler) IssueToll(ctx context.Context, ticket *TicketToll) error {
	// return tollErr.NewErrorWithCode("ERR.PANIC.NOTFOUND", "not found")
	return nil
}
