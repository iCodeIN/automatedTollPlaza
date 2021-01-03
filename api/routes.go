package api

import (
	"automatedTollPlaze/pkg"
	"automatedTollPlaze/pkg/platform/appcontext"
	"net/http"

	"github.com/bnkamalesh/webgo/v4"
)

// API ..
type API struct {
	AppContext *appcontext.AppContext
	Handler    pkg.ServiceHandler
}

// HTTP represents structure of Http Requests
type HTTP struct {
	AppContext *appcontext.AppContext
	APIHandler *pkg.ServiceHandler
	Server     *webgo.Router
}

// NotFound NotFound is the 404 handler
func NotFound() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		webgo.R404(w, map[string]interface{}{
			"message": "The route your looking for doesn't exists",
			"code":    "ERR.ROUTE.NOTFOUND",
		})
		return
	})
}

// Home ..
func (h *HTTP) Home(w http.ResponseWriter, r *http.Request) {
	home := map[string]string{
		"startTime": h.AppContext.StartTime.String(),
		"message":   "Welcome to Automated Toll Plaza",
	}
	webgo.R200(w, home)
}

// Health ..
func (h *HTTP) Health(w http.ResponseWriter, r *http.Request) {
	healthResponse := struct {
		StartTime  string                 `json:"startTime"`
		Dependency map[string]interface{} `json:"dependency"`
	}{
		StartTime: h.AppContext.StartTime.String(),
		Dependency: map[string]interface{}{
			"database": h.AppContext.DbClient.Health(r.Context()) == nil,
		},
	}
	webgo.R200(w, healthResponse)
}

// Routes ...
func (h *HTTP) Routes() []*webgo.Route {
	return []*webgo.Route{
		{
			Name:    "health",
			Method:  http.MethodGet,
			Pattern: "/health",
			Handlers: []http.HandlerFunc{
				h.Health,
			},
		},
		{
			Name:    "Issue New Toll Ticket",
			Method:  http.MethodPost,
			Pattern: "/issue",
			Handlers: []http.HandlerFunc{
				h.issueTollTicket,
			},
		},
		{
			Name:    "List Ticket based on tollId, registration number and/or status",
			Method:  http.MethodGet,
			Pattern: "/tickets/issued",
			Handlers: []http.HandlerFunc{
				h.getTicketIssueList,
			},
		},
		{
			Name:    "Get Ticket Details",
			Method:  http.MethodGet,
			Pattern: "/tickets/:ticketId",
			Handlers: []http.HandlerFunc{
				h.getTicketDetails,
			},
		},
		{
			Name:    "Redeem Toll Ticket",
			Method:  http.MethodPatch,
			Pattern: "/tickets/:ticketId",
			Handlers: []http.HandlerFunc{
				h.redeemTollTicket,
			},
		},
		{
			Name:    "Home",
			Method:  http.MethodGet,
			Pattern: "/",
			Handlers: []http.HandlerFunc{
				h.Home,
			},
		},
	}
}
