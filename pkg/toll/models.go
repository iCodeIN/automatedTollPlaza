package toll

// TicketToll ..
type TicketToll struct {
	TicketID         string `json:"ticketId,omitempty"`
	TollID           string `json:"tollId" validate:"required"`
	RegistrationNo   string `json:"vehicleRegistrationNo" validate:"required"`
	ReturnTollTicket bool   `json:"twoWayToll"`
	Status           string `json:"status" validate:"required,oneof=ISSED REDEEMED"`
}
