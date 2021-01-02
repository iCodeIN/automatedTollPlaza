package api

import (
	"automatedTollPlaze/pkg/errors"
	"automatedTollPlaze/pkg/toll"
	"encoding/json"
	"net/http"

	"github.com/bnkamalesh/webgo/v4"
	"github.com/go-playground/validator/v10"
)

// Home ..
func (api *API) Home(w http.ResponseWriter, r *http.Request) {
	home := map[string]string{
		"message": "Welcome to Automated Toll Plaza",
	}
	webgo.R200(w, home)
}

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
		webgo.R400(w, err)
		return
	}
	err := api.Handler.TollHandler.IssueToll(r.Context())
	if err != nil {
		webgo.R400(w, err)
		return
	}
	// data.Issused = err == nil
	// webgo.R200(w, data)
}
