package api

import (
	"automatedTollPlaze/pkg"
	"automatedTollPlaze/pkg/platform/appcontext"
	"net/http"
	"time"

	"github.com/bnkamalesh/webgo/v4"
)

// API model
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

// Home is a handler to serve home page request
func (h *HTTP) Home(w http.ResponseWriter, r *http.Request) {
	home := map[string]interface{}{
		"startTime": h.AppContext.StartTime,
		"message":   "Welcome to Automated Toll Plaza",
	}
	webgo.R200(w, home)
}

// Health is a handler to check the health status of the application and dependencies
func (h *HTTP) Health(w http.ResponseWriter, r *http.Request) {
	healthResponse := struct {
		StartTime  time.Time              `json:"startTime"`
		Dependency map[string]interface{} `json:"dependency"`
	}{
		StartTime: h.AppContext.StartTime,
		Dependency: map[string]interface{}{
			"database": h.AppContext.DbClient.Health(r.Context()) == nil,
		},
	}
	webgo.R200(w, healthResponse)
}

// Routes initializes and returns the list of routes for our application
func (h *HTTP) Routes() []*webgo.Route {
	return []*webgo.Route{
		{
			Name:    "Health",
			Method:  http.MethodGet,
			Pattern: "/health",
			Handlers: []http.HandlerFunc{
				h.Health,
			},
		},
		{
			Name:    "Issue New Toll Ticket",
			Method:  http.MethodPost,
			Pattern: "/tickets/issue",
			Handlers: []http.HandlerFunc{
				h.issueTollTicket,
			},
		},
		{
			Name:    "List Ticket based on tollId, registration number and/or status",
			Method:  http.MethodGet,
			Pattern: "/tickets",
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
