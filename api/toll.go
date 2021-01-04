package api

import (
	"automatedTollPlaze/pkg/errors"
	"automatedTollPlaze/pkg/toll"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/bnkamalesh/webgo/v4"
	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"
)

// issueTollTicket issues a new toll ticket
func (h *HTTP) issueTollTicket(w http.ResponseWriter, r *http.Request) {
	requestData := toll.TicketToll{}

	// decode the data from request body
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Error(err)
		webgo.SendError(
			w,
			errors.ErrUnprocessableEntity,
			http.StatusUnprocessableEntity,
		)
		return
	}

	// validates the mandatory fields for issuing a toll
	if err := validator.New().Struct(requestData); err != nil {
		log.Error(err)
		webgo.R400(w, errors.ErrMissingFields)
		return
	}

	// issues a new toll
	err := h.APIHandler.TollHandler.IssueTollTicket(r.Context(), &requestData)
	if err != nil {
		webgo.R400(w, err)
		return
	}
	webgo.R200(w, requestData)
}

// getTicketIssueList gets the lists of toll tickets based on the filter
func (h *HTTP) getTicketIssueList(w http.ResponseWriter, r *http.Request) {
	params := toll.TicketListRequest{
		TollID:         r.URL.Query().Get("tollId"),
		RegistrationNo: r.URL.Query().Get("registrationNo"),
		Status: func() string {
			status := strings.TrimSpace(r.URL.Query().Get("status"))
			if len(status) == 0 {
				status = toll.IssuedStatus
			}
			return status
		}(),
	}

	// validates the mandatory fields for listing all the toll tickets
	if err := validator.New().Struct(params); err != nil {
		log.Error(err)
		webgo.R400(w, errors.ErrMissingFields)
		return
	}

	// checks for the pagination. if the values are not provide, it will set to default
	params.Start, params.Limit = getPagination(r)

	// gets the list of toll tickets
	list := h.APIHandler.TollHandler.GetTicketIssueList(r.Context(), &params)
	webgo.R200(w, list)
}

// getTicketIssueList gets the toll ticket details
func (h *HTTP) getTicketDetails(w http.ResponseWriter, r *http.Request) {
	params := webgo.Context(r).Params()

	// reads the ticketId from the request path.
	ticket := &toll.TicketToll{
		TicketID: params["ticketId"],
	}

	// gets the ticket details corresponding to the provided ticketId
	tollTicketData, err := h.APIHandler.TollHandler.GetTollTicketDetails(r.Context(), ticket)
	if err != nil {
		webgo.R400(w, err)
	}
	webgo.R200(w, tollTicketData)
}

// redeemTollTicket redeems the toll ticket
func (h *HTTP) redeemTollTicket(w http.ResponseWriter, r *http.Request) {
	params := webgo.Context(r).Params()

	// reads the ticketId from the request path.
	ticket := &toll.TicketToll{
		TicketID: params["ticketId"],
	}

	// redeems the ticket details corresponding to the provided ticketId
	tollTicketData, err := h.APIHandler.TollHandler.RedeemTollTicket(r.Context(), ticket)
	if err != nil {
		webgo.R400(w, err)
	}
	webgo.R200(w, tollTicketData)
}

// getPagination read the pagination keys - start & limit from query parameters
// if the key is not passed, it will return a default value of 0 for start & 25 for limit
func getPagination(r *http.Request) (int, int) {
	start := 0
	limit := 25
	if val := strings.TrimSpace(r.URL.Query().Get("start")); len(val) > 0 {
		if intVal, err := strconv.Atoi(val); err == nil || intVal > 0 {
			start = intVal
		}
	}
	if val := strings.TrimSpace(r.URL.Query().Get("limit")); len(val) > 0 {
		if intVal, err := strconv.Atoi(val); err == nil || intVal > 0 {
			limit = intVal
		}
	}
	return start, limit
}
