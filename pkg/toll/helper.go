package toll

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
