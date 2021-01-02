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
	if f.ticket != nil && len(f.ticket.RegistrationNo) > 0 {
		f.filter["vehicleRegistrationNo"] = f.ticket.RegistrationNo
	}
	return f
}

// setTicketIDFilter ..
func (f filterMap) setTicketIDFilter() filterMap {
	if f.ticket != nil && len(f.ticket.TicketID) > 0 {
		f.filter["ticketId"] = f.ticket.TicketID
	}
	return f
}

// setTollIDFilter ..
func (f filterMap) setTollIDFilter() filterMap {
	if f.ticket != nil && len(f.ticket.TollID) > 0 {
		f.filter["tollId"] = f.ticket.TollID
	}
	return f
}

// setStatusFilter ..
func (f filterMap) setStatusFilter() filterMap {
	if f.ticket != nil && len(f.ticket.Status) > 0 {
		f.filter["status"] = f.ticket.Status
	}
	return f
}

// getFilter ..
func (f filterMap) getFilter() map[string]interface{} {
	return f.filter
}
