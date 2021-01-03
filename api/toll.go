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

// issueTollTicket ..
func (h *HTTP) issueTollTicket(w http.ResponseWriter, r *http.Request) {
	requestData := toll.TicketToll{}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		log.Error(err)
		webgo.SendError(
			w,
			errors.ErrUnprocessableEntity,
			http.StatusUnprocessableEntity,
		)
		return
	}
	if err := validator.New().Struct(requestData); err != nil {
		log.Error(err)
		webgo.R400(w, errors.ErrMissingFields)
		return
	}
	err := h.APIHandler.TollHandler.IssueTollTicket(r.Context(), &requestData)
	if err != nil {
		webgo.R400(w, err)
		return
	}
	webgo.R200(w, requestData)
}

// getTicketIssueList ..
func (h *HTTP) getTicketIssueList(w http.ResponseWriter, r *http.Request) {
	params := toll.TicketListRequest{
		TollID:         r.URL.Query().Get("tollId"),
		RegistrationNo: r.URL.Query().Get("registrationNo"),
		Status: func() string {
			status := strings.TrimSpace(r.URL.Query().Get("status"))
			if len(status) == 0 {
				status = "ISSUED"
			}
			return status
		}(),
	}
	if err := validator.New().Struct(params); err != nil {
		log.Error(err)
		webgo.R400(w, errors.ErrMissingFields)
		return
	}
	params.Start, params.Limit = getPagination(r)
	list := h.APIHandler.TollHandler.GetTicketIssueList(r.Context(), &params)
	webgo.R200(w, list)
}

// getTicketIssueList ..
func (h *HTTP) getTicketDetails(w http.ResponseWriter, r *http.Request) {
	params := webgo.Context(r).Params()
	ticket := &toll.TicketToll{
		TicketID: params["ticketId"],
	}
	tollTicketData, err := h.APIHandler.TollHandler.GetTollTicketDetails(r.Context(), ticket)
	if err != nil {
		webgo.R400(w, err)
	}
	webgo.R200(w, tollTicketData)
}

// redeemTollTicket ..
func (h *HTTP) redeemTollTicket(w http.ResponseWriter, r *http.Request) {
	params := webgo.Context(r).Params()
	ticket := &toll.TicketToll{
		TicketID: params["ticketId"],
	}
	tollTicketData, err := h.APIHandler.TollHandler.RedeemTollTicket(r.Context(), ticket)
	if err != nil {
		webgo.R400(w, err)
	}
	webgo.R200(w, tollTicketData)
}

// getPagination ..
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
