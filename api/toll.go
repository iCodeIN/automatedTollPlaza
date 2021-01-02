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
)

// Home ..
func (api *API) Home(w http.ResponseWriter, r *http.Request) {
	home := map[string]string{
		"startTime": api.AppContext.StartTime.String(),
		"message":   "Welcome to Automated Toll Plaza",
	}
	webgo.R200(w, home)
}

// Health ..
func (api *API) Health(w http.ResponseWriter, r *http.Request) {
	healthResponse := struct {
		StartTime  string                 `json:"startTime"`
		Dependency map[string]interface{} `json:"dependency"`
	}{
		StartTime: api.AppContext.StartTime.String(),
		Dependency: map[string]interface{}{
			"database": api.AppContext.DbClient.Health(r.Context()) == nil,
		},
	}
	webgo.R200(w, healthResponse)
}

// issueTollTicket ..
func (api *API) issueTollTicket(w http.ResponseWriter, r *http.Request) {
	requestData := toll.TicketToll{}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		webgo.SendError(
			w,
			errors.ErrUnprocessableEntity,
			http.StatusUnprocessableEntity,
		)
		return
	}
	if err := validator.New().Struct(requestData); err != nil {
		webgo.R400(w, errors.ErrMissingFields)
		return
	}
	err := api.Handler.TollHandler.IssueToll(r.Context(), &requestData)
	if err != nil {
		webgo.R400(w, err)
		return
	}
	webgo.R200(w, requestData)
}

// getTicketIssueList ..
func (api *API) getTicketIssueList(w http.ResponseWriter, r *http.Request) {
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
		webgo.R400(w, errors.ErrMissingFields)
		return
	}
	params.Start, params.Limit = getPagination(r)
	list := api.Handler.TollHandler.GetTicketIssueList(r.Context(), &params)
	webgo.R200(w, list)
}

// getTicketIssueList ..
func (api *API) getTicketDetails(w http.ResponseWriter, r *http.Request) {
	params := webgo.Context(r).Params()
	ticket := &toll.TicketToll{
		TicketID: params["ticketId"],
	}
	tollTicketData, err := api.Handler.TollHandler.GetTicketDetails(r.Context(), ticket)
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
