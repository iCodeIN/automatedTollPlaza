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

// Routes ...
func (api *API) Routes() []*webgo.Route {
	return []*webgo.Route{
		{
			Name:    "health",
			Method:  http.MethodGet,
			Pattern: "/",
			Handlers: []http.HandlerFunc{
				api.Health,
			},
		},
		{
			Name:    "Issue New Toll Ticket",
			Method:  http.MethodGet,
			Pattern: "/issue",
			Handlers: []http.HandlerFunc{
				api.issueTollTicket,
			},
		},
	}
}
