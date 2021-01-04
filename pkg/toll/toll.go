package toll

import (
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

// Service has list of behavor being provided
type Service interface {
	IssueTollTicket(ctx context.Context, ticket *TicketToll) error
	GetTicketIssueList(ctx context.Context, params *TicketListRequest) TicketListResponse
	GetTollTicketDetails(ctx context.Context, ticket *TicketToll) (*TicketToll, error)
	RedeemTollTicket(ctx context.Context, ticket *TicketToll) (*TicketToll, error)
}

// service handler
type handler struct {
	AppCtx *appcontext.AppContext
}

// NewTollService returns a new toll service of supported methods
func NewTollService(ctx context.Context, appCtx *appcontext.AppContext) Service {
	return &handler{
		AppCtx: appCtx,
	}
}

// IssueTollTicket issues a new toll ticket.
// checking if there are any pending toll ticket corresponding to the provided filter or not.
// generating a unique ticket id
// calculate the price to be charged based on the type of vehicle & type of ticket to be issued - single or return
// added the issue & last updated timestamp as current time and redeemBy timestamp
func (s *handler) IssueTollTicket(ctx context.Context, ticket *TicketToll) error {

	// builds the filter based on vehicle registration number & status as "issued"
	filter := newTicketFilter(ticket).setRegistrationNoFilter().getFilter()
	filter["status"] = IssuedStatus

	// necessary input for querying mongoDB
	dbParams := db.Params{
		Database:   TollDatabaseName,
		Collection: TicketCollectionName,
		Filter:     filter,
		Result:     &TicketToll{},
	}

	// checking if there are any pending toll ticket corresponding to the provided filter or not.
	// if the error return db.ErrNotFound - meaning no data corresponding to provided filter,
	// a new toll ticket can be issued
	err := s.AppCtx.DbClient.FindOne(ctx, dbParams)
	switch {
	case err != nil && err != db.ErrNotFound:
		log.Error(err)
		return errors.ToTollError(err)
	case err == nil:
		log.Error(ErrPendingTollTickets)
		return ErrPendingTollTickets
	}

	// generating a unique ticket id
	ticket.TicketID = func() string {
		id := rand.NewSource(time.Now().UnixNano())
		return ticket.TollID + strconv.Itoa(int(id.Int63()))
	}()

	// calculate the price to be charged based on the type of vehicle
	// price calculation also includes if it is a return type of toll or not.
	ticket.Price = func() float64 {
		priceVal := s.AppCtx.Config.Pricing.Default
		if val, ok := s.AppCtx.Config.Pricing.VehicleType[ticket.VehicleType]; ok {
			priceVal = val
		}
		price := priceVal.TwoWay
		if !ticket.ReturnTollTicket {
			price = priceVal.OneWay
			ticket.Status = RedeemedStatus
		}
		return price
	}()

	// added the issue & last updated timestamp as current time
	issueTime := time.Now().UTC()
	ticket.IssuedTimeStamp = issueTime
	ticket.UpdatedTimeStamp = issueTime

	// adding the redeemeption date before penality is charged. It should be redeemed before midnight
	redeemBy := issueTime.Add(24 * time.Hour)
	ticket.RedeemBy = time.Date(redeemBy.Year(), redeemBy.Month(), redeemBy.Day(), 0, 0, 0, 0, redeemBy.Location())

	// necessary input for create a new toll ticket in mongoDB
	params := db.Params{
		Database:   TollDatabaseName,
		Collection: TicketCollectionName,
		UpsertData: ticket,
	}

	// execute or creating a new toll ticket for the given tollId and registration number
	if err := s.AppCtx.DbClient.InsertOne(ctx, params); err != nil {
		log.Error(err)
		return errors.ToTollError(err)
	}
	return nil
}

// GetTicketIssueList get the list of the toll tickets based on the following parameters
//		1. TollId
// 		2. Vehicle Registration Number
// 		3. Status of toll ticket - optional but defaulted to ISSUED
// Applies start & limit for pagination
func (s *handler) GetTicketIssueList(ctx context.Context, params *TicketListRequest) TicketListResponse {

	// Copy the data to make it flexible to using setter functions to build the filter
	ticket := &TicketToll{}
	copier.Copy(ticket, params)

	// building the filter using ticketId provided for the findOne query
	filter := newTicketFilter(ticket).setRegistrationNoFilter().setTollIDFilter().setStatusFilter().getFilter()
	ticketList := make([]TicketToll, 0)

	// necessary input for querying mongoDB
	dbParams := db.Params{
		Database:   TollDatabaseName,
		Collection: TicketCollectionName,
		Filter:     filter,
		Result:     &ticketList,
		Pagination: &db.Paginate{
			Start: params.Start,
			Limit: params.Limit,
		},
	}

	// get the the list of toll tickets based on the provided filter
	s.AppCtx.DbClient.FindAll(ctx, dbParams)
	list := TicketListResponse{
		List:  ticketList,
		Count: s.AppCtx.DbClient.Count(ctx, dbParams),
		Start: params.Start,
		Limit: params.Limit,
	}
	return list
}

// GetTollTicketDetails get the toll ticket details for the provided ticketId
func (s *handler) GetTollTicketDetails(ctx context.Context, ticket *TicketToll) (*TicketToll, error) {

	// building the filter using ticketId provided for the findOne query
	filter := newTicketFilter(ticket).setTicketIDFilter().getFilter()
	dbTicket := &TicketToll{}

	// necessary input for querying mongoDB
	dbParams := db.Params{
		Database:   TollDatabaseName,
		Collection: TicketCollectionName,
		Filter:     filter,
		Result:     dbTicket,
	}

	// checking if the provided ticketId exists or not.
	if err := s.AppCtx.DbClient.FindOne(ctx, dbParams); err != nil {
		log.Error(err)
		return nil, errors.ToTollError(err)
	}
	return dbTicket, nil
}

// RedeemTollTicket performs redemption of the toll ticket.
// It check if the provided toll ticket is valid or exists in the database or not.
// Check the current status of the toll ticket. If it is already redeemed, it gives an error.
// Check the redeemBy field value to the current time. According it will charge a penality.
// It marks the status of the toll ticket as redeemed along with time stamp.
func (s *handler) RedeemTollTicket(ctx context.Context, ticket *TicketToll) (*TicketToll, error) {

	// building the filter using ticketId provided for the findOne query
	filter := newTicketFilter(ticket).setTicketIDFilter().getFilter()

	// necessary input for querying mongoDB
	dbParams := db.Params{
		Database:   TollDatabaseName,
		Collection: TicketCollectionName,
		Filter:     filter,
		Result:     ticket,
	}

	// checking if the provided ticketId exists or not.
	err := s.AppCtx.DbClient.FindOne(ctx, dbParams)
	if err != nil {
		if err == db.ErrNotFound {
			err = ErrInvalidTollTicket
		}
		log.Error(err)
		return nil, errors.ToTollError(err)
	}

	// if the ticket is already redeemed or not.
	if ticket.Status == RedeemedStatus {
		log.Error(ErrAlreadyRedeemed)
		return nil, ErrAlreadyRedeemed
	}

	// checking the redemption time is valid or not.
	if time.Now().UTC().After(ticket.RedeemBy) {
		log.Warn("Redemption time has past. Charging double penality.")
		ticket.Price = ticket.Price * 2
	}

	// Update the status of the ticket & the timeStamp of redemption
	ticket.Status = RedeemedStatus
	ticket.UpdatedTimeStamp = time.Now().UTC()

	// necessary input for querying mongoDB
	dbParams = db.Params{
		Database:   TollDatabaseName,
		Collection: TicketCollectionName,
		Filter:     filter,
		UpsertData: ticket,
	}

	// updating the data in mongoDB corresponding to the provided ticketId
	if err := s.AppCtx.DbClient.Upsert(ctx, dbParams); err != nil {
		log.Error(err)
		return nil, errors.ToTollError(err)
	}
	return ticket, nil
}
