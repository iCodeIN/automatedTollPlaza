package toll

import "time"

// TicketToll ..
type TicketToll struct {
	TicketID         string    `json:"ticketId,omitempty" bson:"ticketId"`
	TollID           string    `json:"tollId" validate:"required" bson:"tollId"`
	RegistrationNo   string    `json:"vehicleRegistrationNo" validate:"required" bson:"vehicleRegistrationNo"`
	ReturnTollTicket bool      `json:"twoWayToll" bson:"twoWayToll"`
	Status           string    `json:"status" validate:"required,oneof=ISSED REDEEMED" bson:"status"`
	Price            float64   `json:"TotalAmount,omitempty" bson:"TotalAmount"`
	IssuedTimeStamp  time.Time `json:"issuedTimeStamp" bson:"issuedTimeStamp"`
	UpdatedTimeStamp time.Time `json:"updatedTimeStamp" bson:"updatedTimeStamp"`
}
