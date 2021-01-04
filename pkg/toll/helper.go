package toll

type filterMap struct {
	filter map[string]interface{}
	ticket *TicketToll
}

// newTicketFilter returns a new filterMap
func newTicketFilter(ticket *TicketToll) filterMap {
	return filterMap{
		ticket: ticket,
		filter: make(map[string]interface{}, 0),
	}
}

// setRegistrationNoFilter sets the registration query filter to the filterMap
func (f filterMap) setRegistrationNoFilter() filterMap {
	if f.ticket != nil && len(f.ticket.RegistrationNo) > 0 {
		f.filter["vehicleRegistrationNo"] = f.ticket.RegistrationNo
	}
	return f
}

// setTicketIDFilter sets the ticketId query filter to the filterMap
func (f filterMap) setTicketIDFilter() filterMap {
	if f.ticket != nil && len(f.ticket.TicketID) > 0 {
		f.filter["ticketId"] = f.ticket.TicketID
	}
	return f
}

// setTollIDFilter sets the tollId query filter to the filterMap
func (f filterMap) setTollIDFilter() filterMap {
	if f.ticket != nil && len(f.ticket.TollID) > 0 {
		f.filter["tollId"] = f.ticket.TollID
	}
	return f
}

// setStatusFilter sets the status query filter to the filterMap
func (f filterMap) setStatusFilter() filterMap {
	if f.ticket != nil && len(f.ticket.Status) > 0 {
		f.filter["status"] = f.ticket.Status
	}
	return f
}

// getFilter returns the filters created
func (f filterMap) getFilter() map[string]interface{} {
	return f.filter
}
