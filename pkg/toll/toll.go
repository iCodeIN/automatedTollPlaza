package toll

import (
	// tollErr "automatedTollPlaze/pkg/errors"
	"automatedTollPlaze/pkg/errors"
	"automatedTollPlaze/pkg/platform/appcontext"
	"context"
	"math/rand"
	"strconv"
	"time"
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
	filter := newTicketFilter(ticket).setRegistrationNoFilter().setTollIDFilter().filter
	findOneData := &TicketToll{}
	if err := s.AppCtx.DbClient.Find(ctx, "toll", "tickets", filter, findOneData); err != nil {
		return errors.ToTollError(err)
	}
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
	if err := s.AppCtx.DbClient.InsertOne(ctx, "toll", "tickets", ticket); err != nil {
		return errors.ToTollError(err)
	}
	return nil
}

type filterMap struct {
	filter map[string]interface{}
	ticket *TicketToll
}

// newTicketFilter ..
func newTicketFilter(ticket *TicketToll) filterMap {
	return filterMap{
		ticket: ticket,
		filter: make(map[string]interface{}, 0),
	}
}

// setRegistrationNoFilter ..
func (f filterMap) setRegistrationNoFilter() filterMap {
	if f.ticket != nil {
		f.filter["vehicleRegistrationNo"] = f.ticket.RegistrationNo
	}
	return f
}

// setTicketIDFilter ..
func (f filterMap) setTicketIDFilter() filterMap {
	if f.ticket != nil {
		f.filter["ticketId"] = f.ticket.TicketID
	}
	return f
}

// setTollIDFilter ..
func (f filterMap) setTollIDFilter() filterMap {
	if f.ticket != nil {
		f.filter["tollId"] = f.ticket.TollID
	}
	return f
}

// getFilter ..
func (f filterMap) getFilter() map[string]interface{} {
	return f.filter
}
