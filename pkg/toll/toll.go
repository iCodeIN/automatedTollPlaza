package toll

import (
	// tollErr "automatedTollPlaze/pkg/errors"
	"automatedTollPlaze/pkg/platform/appcontext"
	"context"
	"fmt"
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
	filter := getFilters(ticket)
	delete(filter, "ticketId")
	project := make(map[string]interface{}, 0)
	list := make([]map[string]interface{}, 0)
	if err := s.AppCtx.DbClient.Find(ctx, "toll", "ticket", filter, project, list); err != nil {
		return err
	}
	fmt.Println(list)
	return s.AppCtx.DbClient.InsertOne(ctx, "toll", "tickets", ticket)
}

func getFilters(ticket *TicketToll) map[string]interface{} {
	filter := make(map[string]interface{}, 0)

	return filter
}
