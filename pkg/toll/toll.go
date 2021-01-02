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
	log "github.com/sirupsen/logrus"
)

//Service ..
type Service interface {
	GetTollTicketDetails(ctx context.Context, ticket *TicketToll) (*TicketToll, error)
	GetTicketIssueList(ctx context.Context, params *TicketListRequest) TicketListResponse
	IssueTollTicket(ctx context.Context, ticket *TicketToll) error
	RedeemTollTicket(ctx context.Context, ticket *TicketToll) (*TicketToll, error)
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

// IssueTollTicket ..
func (s *handler) IssueTollTicket(ctx context.Context, ticket *TicketToll) error {
	filter := newTicketFilter(ticket).setRegistrationNoFilter().getFilter()
	filter["status"] = "ISSUED"
	dbParams := db.Params{
		Database:   "toll",
		Collection: "tickets",
		Filter:     filter,
		Result:     &TicketToll{},
	}
	err := s.AppCtx.DbClient.FindOne(ctx, dbParams)
	switch {
	case err != nil && err != db.ErrNotFound:
		log.Error(err)
		return errors.ToTollError(err)
	case err == nil:
		log.Error(ErrPendingTollTickets)
		return ErrPendingTollTickets
	}

	ticket.TicketID = func() string {
		id := rand.NewSource(time.Now().UnixNano())
		return ticket.TollID + strconv.Itoa(int(id.Int63()))
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
		log.Error(err)
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

// GetTollTicketDetails ..
func (s *handler) GetTollTicketDetails(ctx context.Context, ticket *TicketToll) (*TicketToll, error) {
	filter := newTicketFilter(ticket).setTicketIDFilter().filter
	dbTicket := &TicketToll{}
	dbParams := db.Params{
		Database:   "toll",
		Collection: "tickets",
		Filter:     filter,
		Result:     dbTicket,
	}
	if err := s.AppCtx.DbClient.FindOne(ctx, dbParams); err != nil {
		log.Error(err)
		return nil, errors.ToTollError(err)
	}
	return dbTicket, nil
}

// RedeemTollTicket ..
func (s *handler) RedeemTollTicket(ctx context.Context, ticket *TicketToll) (*TicketToll, error) {
	filter := newTicketFilter(ticket).setTicketIDFilter().getFilter()
	dbParams := db.Params{
		Database:   "toll",
		Collection: "tickets",
		Filter:     filter,
		Result:     ticket,
	}
	err := s.AppCtx.DbClient.FindOne(ctx, dbParams)
	if err != nil {
		if err == db.ErrNotFound {
			err = ErrInvalidTollTicket
		}
		log.Error(err)
		return nil, errors.ToTollError(err)
	}
	if ticket.Status == "REDEEMED" {
		log.Error(ErrAlreadyRedeemed)
		return nil, ErrAlreadyRedeemed
	}
	ticket.Status = "REDEEMED"
	ticket.UpdatedTimeStamp = time.Now().UTC()
	dbParams = db.Params{
		Database:   "toll",
		Collection: "tickets",
		Filter:     filter,
		UpsertData: ticket,
	}
	if err := s.AppCtx.DbClient.Upsert(ctx, dbParams); err != nil {
		log.Error(err)
		return nil, errors.ToTollError(err)
	}
	return ticket, nil
}
