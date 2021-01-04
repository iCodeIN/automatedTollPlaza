package toll

import "time"

// TicketToll ..
type TicketToll struct {
	TicketID         string    `json:"ticketId,omitempty" bson:"ticketId"`
	TollID           string    `json:"tollId" validate:"required" bson:"tollId"`
	VehicleType      string    `json:"vehicleType" validate:"required,oneof=light moderate heavy" bson:"vehicleType"`
	RegistrationNo   string    `json:"vehicleRegistrationNo" validate:"required" bson:"vehicleRegistrationNo"`
	ReturnTollTicket bool      `json:"twoWayToll" bson:"twoWayToll"`
	Status           string    `json:"status" validate:"required,oneof=ISSUED REDEEMED" bson:"status"`
	Price            float64   `json:"TotalAmount,omitempty" bson:"TotalAmount"`
	IssuedTimeStamp  time.Time `json:"issuedTimeStamp" bson:"issuedTimeStamp"`
	RedeemBy         time.Time `json:"-" bson:"redeemBy"`
	UpdatedTimeStamp time.Time `json:"updatedTimeStamp" bson:"updatedTimeStamp"`
}

// TicketListRequest ..
type TicketListRequest struct {
	TollID         string `validate:"required_without=RegistrationNo"`
	RegistrationNo string `validate:"required_without=TollID"`
	Status         string `validate:"oneof=ISSUED REDEEMED"`
	Start          int
	Limit          int
}

// TicketListResponse ..
type TicketListResponse struct {
	List  []TicketToll `json:"list"`
	Count int          `json:"count"`
	Start int          `json:"start"`
	Limit int          `json:"limit"`
}
