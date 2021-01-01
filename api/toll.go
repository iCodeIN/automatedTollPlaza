package api

import (
	"net/http"

	"github.com/bnkamalesh/webgo/v4"
)

// Health ..
func (api *API) Health(w http.ResponseWriter, r *http.Request) {
	healthResponse := struct {
		Database bool `json:"database"`
	}{
		Database: api.AppContext.DbClient.Health(r.Context()) == nil,
	}
	webgo.R200(w, healthResponse)
}

// issueTollTicket ..
func (api *API) issueTollTicket(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Issused bool `json:"issued"`
	}
	err := api.Handler.TollHandler.IssueToll(r.Context())
	if err != nil {
		webgo.R400(w, err)
		return
	}
	data.Issused = err == nil
	webgo.R200(w, data)
}
