package toll

import "automatedTollPlaze/pkg/errors"

var (
	// ErrPendingTollTickets ..
	ErrPendingTollTickets = errors.NewErrorWithCode(
		"ERR.TOLL.PENDING_TOLL_TICKETS",
		"This are pending toll ticket for redemption",
	)

	// ErrInvalidTollTicket ..
	ErrInvalidTollTicket = errors.NewErrorWithCode(
		"ERR.TOLL.INVALID_TOLL_TICKET",
		"This is an invalid toll ticket",
	)

	// ErrAlreadyRedeemed ..
	ErrAlreadyRedeemed = errors.NewErrorWithCode(
		"ERR.TOLL.ALREADY_TOLL_TICKET_REDEEMED",
		"This toll ticket is already redeemed",
	)
)
