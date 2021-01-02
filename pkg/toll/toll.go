package toll

import (
	// tollErr "automatedTollPlaze/pkg/errors"
	"automatedTollPlaze/pkg/errors"
	"automatedTollPlaze/pkg/platform/appcontext"
	"automatedTollPlaze/pkg/platform/db"
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
)

//Service ..
type Service interface {
	GetTicketDetails(ctx context.Context, ticket *TicketToll) (*TicketToll, error)
	GetTicketIssueList(ctx context.Context, params *TicketListRequest) TicketListResponse
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
	// filter := newTicketFilter(ticket).setRegistrationNoFilter().setTollIDFilter().getFilter()
	// dbParams := db.Params{
	// 	Database:   "toll",
	// 	Collection: "tickets",
	// 	Filter:     filter,
	// 	Result:     &TicketToll{},
	// }
	// if err := s.AppCtx.DbClient.FindOne(ctx, dbParams); err != nil {
	// 	return errors.ToTollError(err)
	// }
	ticket.TicketID = func() string {
		id := rand.NewSource(time.Now().UnixNano())
		return ticket.TollID + "-" + strconv.Itoa(int(id.Int63()))
	}()
	ticket.Price = func() float64 {
		price := 200.0
		if !ticket.ReturnTollTicket {
			price = 100.0
			ticket.Status = "REDEEMED"
		}
		return price
	}()
	ticket.IssuedTimeStamp = time.Now().UTC()
	ticket.UpdatedTimeStamp = time.Now().UTC()
	params := db.Params{
		Database:   "toll",
		Collection: "tickets",
		UpsertData: ticket,
	}
	if err := s.AppCtx.DbClient.InsertOne(ctx, params); err != nil {
		return errors.ToTollError(err)
	}
	return nil
}

// GetTicketIssueList ..
func (s *handler) GetTicketIssueList(ctx context.Context, params *TicketListRequest) TicketListResponse {
	ticket := &TicketToll{}
	copier.Copy(ticket, params)
	filter := newTicketFilter(ticket).setRegistrationNoFilter().setTollIDFilter().setStatusFilter().getFilter()
	ticketList := make([]TicketToll, 0)
	dbParams := db.Params{
		Database:   "toll",
		Collection: "tickets",
		Filter:     filter,
		Result:     &ticketList,
		Pagination: &db.Paginate{
			Start: params.Start,
			Limit: params.Limit,
		},
	}
	s.AppCtx.DbClient.FindAll(ctx, dbParams)
	list := TicketListResponse{
		List:  ticketList,
		Count: s.AppCtx.DbClient.Count(ctx, dbParams),
		Start: params.Start,
		Limit: params.Limit,
	}
	return list
}

// GetTicketIssueList ..
func (s *handler) GetTicketDetails(ctx context.Context, ticket *TicketToll) (*TicketToll, error) {
	filter := newTicketFilter(ticket).setTicketIDFilter().filter
	dbTicket := &TicketToll{}
	dbParams := db.Params{
		Database:   "toll",
		Collection: "tickets",
		Filter:     filter,
		Result:     dbTicket,
	}
	if err := s.AppCtx.DbClient.FindOne(ctx, dbParams); err != nil {
		return nil, errors.ToTollError(err)
	}
	return dbTicket, nil
}
