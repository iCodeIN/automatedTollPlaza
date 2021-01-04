package toll

import "automatedTollPlaze/pkg/errors"

var (
	// ErrPendingTollTickets is error when toll tickets are pending for redemption
	ErrPendingTollTickets = errors.NewErrorWithCode(
		"ERR.TOLL.PENDING_TOLL_TICKETS",
		"This are pending toll ticket for redemption",
	)

	// ErrInvalidTollTicket is error when toll ticket is invalid
	ErrInvalidTollTicket = errors.NewErrorWithCode(
		"ERR.TOLL.INVALID_TOLL_TICKET",
		"This is an invalid toll ticket",
	)

	// ErrAlreadyRedeemed is error when the toll ticket is already redeemed
	ErrAlreadyRedeemed = errors.NewErrorWithCode(
		"ERR.TOLL.ALREADY_TOLL_TICKET_REDEEMED",
		"This toll ticket is already redeemed",
	)
)
