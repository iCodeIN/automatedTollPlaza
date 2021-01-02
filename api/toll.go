package api

import (
	"automatedTollPlaze/pkg/errors"
	"automatedTollPlaze/pkg/toll"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

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
	requestData.TicketID = func() string {
		id := rand.NewSource(time.Now().UnixNano())
		return requestData.TollID + "-" + strconv.Itoa(int(id.Int63()))
	}()
	requestData.Price = func() float64 {
		price := 100.0
		if requestData.ReturnTollTicket {
			price = 200.0
		}
		return price
	}()
	err := api.Handler.TollHandler.IssueToll(r.Context(), &requestData)
	if err != nil {
		webgo.R400(w, err)
		return
	}
	webgo.R200(w, requestData)
}
