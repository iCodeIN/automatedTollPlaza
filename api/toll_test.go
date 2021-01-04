package api

import (
	"automatedTollPlaze/pkg"
	"automatedTollPlaze/pkg/platform/appcontext"
	"automatedTollPlaze/pkg/platform/db"
	"automatedTollPlaze/pkg/toll"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bnkamalesh/webgo/v4"
)

type mockDBHandler struct {
	HealthFn    func(ctx context.Context) error
	InsertOneFn func(ctx context.Context, params db.Params) error
	FindOneFn   func(ctx context.Context, params db.Params) error
	FindAllFn   func(ctx context.Context, params db.Params) error
	CountFn     func(ctx context.Context, params db.Params) int
	UpsertFn    func(ctx context.Context, params db.Params) error
}

func (m mockDBHandler) Health(ctx context.Context) error {
	return m.HealthFn(ctx)
}

func (m mockDBHandler) InsertOne(ctx context.Context, params db.Params) error {
	return m.InsertOneFn(ctx, params)
}

func (m mockDBHandler) FindOne(ctx context.Context, params db.Params) error {
	return m.FindOneFn(ctx, params)
}

func (m mockDBHandler) FindAll(ctx context.Context, params db.Params) error {
	return m.FindAllFn(ctx, params)
}

func (m mockDBHandler) Count(ctx context.Context, params db.Params) int {
	return m.CountFn(ctx, params)
}

func (m mockDBHandler) Upsert(ctx context.Context, params db.Params) error {
	return m.UpsertFn(ctx, params)
}

type mockHandler struct {
	GetTollTicketDetailsFn func(ctx context.Context, ticket *toll.TicketToll) (*toll.TicketToll, error)
	GetTicketIssueListFn   func(ctx context.Context, params *toll.TicketListRequest) toll.TicketListResponse
	IssueTollTicketFn      func(ctx context.Context, ticket *toll.TicketToll) error
	RedeemTollTicketFn     func(ctx context.Context, ticket *toll.TicketToll) (*toll.TicketToll, error)
}

func (m mockHandler) GetTollTicketDetails(ctx context.Context, ticket *toll.TicketToll) (*toll.TicketToll, error) {
	return m.GetTollTicketDetailsFn(ctx, ticket)
}

func (m mockHandler) GetTicketIssueList(ctx context.Context, params *toll.TicketListRequest) toll.TicketListResponse {
	return m.GetTicketIssueListFn(ctx, params)
}

func (m mockHandler) IssueTollTicket(ctx context.Context, ticket *toll.TicketToll) error {
	return m.IssueTollTicketFn(ctx, ticket)
}

func (m mockHandler) RedeemTollTicket(ctx context.Context, ticket *toll.TicketToll) (*toll.TicketToll, error) {
	return m.RedeemTollTicketFn(ctx, ticket)
}

func mockAppContext() *appcontext.AppContext {
	return &appcontext.AppContext{
		DbClient: &mockDBHandler{
			HealthFn: func(ctx context.Context) error {
				return nil
			},
		},
		StartTime: time.Time{},
	}
}

func New(appCtx *appcontext.AppContext, initCfg *webgo.Config) *HTTP {
	cfg := &webgo.Config{
		Host:         initCfg.Host,
		Port:         initCfg.Port,
		ReadTimeout:  initCfg.ReadTimeout * time.Second,
		WriteTimeout: initCfg.WriteTimeout * time.Second,
	}
	h := &HTTP{
		AppContext: appCtx,
		APIHandler: nil,
	}
	h.Server = webgo.NewRouter(cfg, h.Routes())
	return h
}

func mockHTTPServer() *HTTP {
	httpCfg := &webgo.Config{
		Host:         "127.0.0.1",
		Port:         "9090",
		ReadTimeout:  15,
		WriteTimeout: 60,
	}
	http := New(mockAppContext(), httpCfg)
	return http
}

func TestHTTP_home(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		url         string
		apiServices *pkg.ServiceHandler
		code        int
		expected    string
	}{
		{
			name:     "Home Page",
			method:   http.MethodGet,
			url:      "/",
			code:     200,
			expected: `{"data":{"message":"Welcome to Automated Toll Plaza","startTime":"0001-01-01 00:00:00 +0000 UTC"},"status":200}`,
		},
	}
	server := mockHTTPServer()
	router := server.Server
	for _, tt := range tests {
		server.APIHandler = tt.apiServices
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			response := httptest.NewRecorder()
			// handler := http.HandlerFunc(h.content)

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			router.ServeHTTP(response, request)

			// Check the status code is what we expect.
			if response.Code != tt.code {
				t.Errorf("wrong status code: got %v want %v", response.Code, tt.code)
			}

			// Check the response body is what we expect.
			if !reflect.DeepEqual(strings.TrimSpace(response.Body.String()), tt.expected) {
				t.Errorf("unexpected body: got %v want %v", response.Body.String(), tt.expected)
			}
		})
	}
}

func TestHTTP_health(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		url         string
		apiServices *pkg.ServiceHandler
		code        int
		expected    string
	}{
		{
			name:     "Get the health status",
			method:   http.MethodGet,
			url:      "/health",
			code:     200,
			expected: `{"data":{"startTime":"0001-01-01 00:00:00 +0000 UTC","dependency":{"database":true}},"status":200}`,
		},
	}
	server := mockHTTPServer()
	router := server.Server
	for _, tt := range tests {
		server.APIHandler = tt.apiServices
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			response := httptest.NewRecorder()
			// handler := http.HandlerFunc(h.content)

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			router.ServeHTTP(response, request)

			// Check the status code is what we expect.
			if response.Code != tt.code {
				t.Errorf("wrong status code: got %v want %v", response.Code, tt.code)
			}

			// Check the response body is what we expect.
			if !reflect.DeepEqual(strings.TrimSpace(response.Body.String()), tt.expected) {
				t.Errorf("unexpected body: got %v want %v", response.Body.String(), tt.expected)
			}
		})
	}
}

func TestHTTP_issueTollTicket(t *testing.T) {
	basePath := "/tickets/issue"

	tests := []struct {
		name        string
		method      string
		url         string
		body        string
		apiServices *pkg.ServiceHandler
		code        int
		expected    string
	}{
		{
			name:   "Invalid Payload",
			method: http.MethodPost,
			url:    basePath,
			body:   `{""}`,
			apiServices: &pkg.ServiceHandler{
				TollHandler: mockHandler{},
			},
			code:     422,
			expected: `{"errors":{"message":"Invalid Json","code":"ERR.HTTP.UNPROCESSABLEENTITY"},"status":422}`,
		},
		{
			name:   "Missing Fields",
			method: http.MethodPost,
			url:    basePath,
			body:   `{"tollId": "1"}`,
			apiServices: &pkg.ServiceHandler{
				TollHandler: mockHandler{},
			},
			code:     400,
			expected: `{"errors":{"message":"There are some missing fields that are required","code":"ERR.APP.MISSING_FIELDS"},"status":400}`,
		},
		{
			name:   "Pending toll tickets",
			method: http.MethodPost,
			url:    basePath,
			body:   `{"ticketId":"","tollId":"2","vehicleRegistrationNo":"KA 01 AB 1235","vehicleType":"light","twoWayToll":true,"status":"ISSUED"}`,
			apiServices: &pkg.ServiceHandler{
				TollHandler: mockHandler{
					IssueTollTicketFn: func(ctx context.Context, ticket *toll.TicketToll) error {
						return toll.ErrPendingTollTickets
					},
				},
			},
			code:     400,
			expected: `{"errors":{"message":"This are pending toll ticket for redemption","code":"ERR.TOLL.PENDING_TOLL_TICKETS"},"status":400}`,
		},
		{
			name:   "Successfully Issued toll",
			method: http.MethodPost,
			url:    basePath,
			body:   `{"ticketId":"","tollId":"2","vehicleRegistrationNo":"KA 01 AB 1235","vehicleType":"light","twoWayToll":true,"status":"ISSUED"}`,
			apiServices: &pkg.ServiceHandler{
				TollHandler: mockHandler{
					IssueTollTicketFn: func(ctx context.Context, ticket *toll.TicketToll) error {
						return nil
					},
				},
			},
			code:     200,
			expected: `{"data":{"tollId":"2","vehicleType":"light","vehicleRegistrationNo":"KA 01 AB 1235","twoWayToll":true,"status":"ISSUED","issuedTimeStamp":"0001-01-01T00:00:00Z","updatedTimeStamp":"0001-01-01T00:00:00Z"},"status":200}`,
		},
	}
	server := mockHTTPServer()
	router := server.Server
	for _, tt := range tests {
		server.APIHandler = tt.apiServices
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(tt.method, tt.url, strings.NewReader(tt.body))
			if err != nil {
				t.Fatal(err)
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			response := httptest.NewRecorder()
			// handler := http.HandlerFunc(h.content)

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			router.ServeHTTP(response, request)

			// Check the status code is what we expect.
			if response.Code != tt.code {
				t.Errorf("wrong status code: got %v want %v", response.Code, tt.code)
			}

			// Check the response body is what we expect.
			if !reflect.DeepEqual(strings.TrimSpace(response.Body.String()), tt.expected) {
				t.Errorf("unexpected body: got %v want %v", response.Body.String(), tt.expected)
			}
		})
	}
}

func TestHTTP_getTicketIssueList(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		url         string
		apiServices *pkg.ServiceHandler
		code        int
		expected    string
	}{
		{
			name:   "Missing Fields",
			method: http.MethodGet,
			url:    "/tickets?tollId=",
			apiServices: &pkg.ServiceHandler{
				TollHandler: mockHandler{},
			},
			code:     400,
			expected: `{"errors":{"message":"There are some missing fields that are required","code":"ERR.APP.MISSING_FIELDS"},"status":400}`,
		},
		{
			name:   "Pending toll tickets",
			method: http.MethodGet,
			url:    "/tickets?tollId=1&status=ISSUED",
			apiServices: &pkg.ServiceHandler{
				TollHandler: mockHandler{
					GetTicketIssueListFn: func(ctx context.Context, params *toll.TicketListRequest) toll.TicketListResponse {
						return toll.TicketListResponse{
							List: make([]toll.TicketToll, 0),
						}
					},
				},
			},
			code:     200,
			expected: `{"data":{"list":[],"count":0,"start":0,"limit":0},"status":200}`,
		},
	}
	server := mockHTTPServer()
	router := server.Server
	for _, tt := range tests {
		server.APIHandler = tt.apiServices
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			response := httptest.NewRecorder()
			// handler := http.HandlerFunc(h.content)

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			router.ServeHTTP(response, request)

			// Check the status code is what we expect.
			if response.Code != tt.code {
				t.Errorf("wrong status code: got %v want %v", response.Code, tt.code)
			}

			// Check the response body is what we expect.
			if !reflect.DeepEqual(strings.TrimSpace(response.Body.String()), tt.expected) {
				t.Errorf("unexpected body: got %v want %v", response.Body.String(), tt.expected)
			}
		})
	}
}

func TestHTTP_getTicketDetails(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		url         string
		apiServices *pkg.ServiceHandler
		code        int
		expected    string
	}{
		{
			name:   "Route not found",
			method: http.MethodGet,
			url:    "/tickets/",
			apiServices: &pkg.ServiceHandler{
				TollHandler: mockHandler{},
			},
			code:     404,
			expected: `{"errors":{"code":"ERR.ROUTE.NOTFOUND","message":"The route your looking for doesn't exists"},"status":404}`,
		},
		{
			name:   "Invalid Ticket",
			method: http.MethodGet,
			url:    "/tickets/1",
			apiServices: &pkg.ServiceHandler{
				TollHandler: mockHandler{
					GetTollTicketDetailsFn: func(ctx context.Context, ticket *toll.TicketToll) (*toll.TicketToll, error) {
						return nil, toll.ErrInvalidTollTicket
					},
				},
			},
			code:     400,
			expected: `{"errors":{"message":"This is an invalid toll ticket","code":"ERR.TOLL.INVALID_TOLL_TICKET"},"status":400}`,
		},
		{
			name:   "Toll Tickets Details",
			method: http.MethodGet,
			url:    "/tickets/2",
			apiServices: &pkg.ServiceHandler{
				TollHandler: mockHandler{
					GetTollTicketDetailsFn: func(ctx context.Context, ticket *toll.TicketToll) (*toll.TicketToll, error) {
						return &toll.TicketToll{}, nil
					},
				},
			},
			code:     200,
			expected: `{"data":{"tollId":"","vehicleType":"","vehicleRegistrationNo":"","twoWayToll":false,"status":"","issuedTimeStamp":"0001-01-01T00:00:00Z","updatedTimeStamp":"0001-01-01T00:00:00Z"},"status":200}`,
		},
	}
	server := mockHTTPServer()
	router := server.Server
	router.NotFound = NotFound()
	for _, tt := range tests {
		server.APIHandler = tt.apiServices
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			response := httptest.NewRecorder()
			// handler := http.HandlerFunc(h.content)

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			router.ServeHTTP(response, request)

			// Check the status code is what we expect.
			if response.Code != tt.code {
				t.Errorf("wrong status code: got %v want %v", response.Code, tt.code)
			}

			// Check the response body is what we expect.
			if !reflect.DeepEqual(strings.TrimSpace(response.Body.String()), tt.expected) {
				t.Errorf("unexpected body: got %v want %v", response.Body.String(), tt.expected)
			}
		})
	}
}

func TestHTTP_redeemTollTicket(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		url         string
		apiServices *pkg.ServiceHandler
		code        int
		expected    string
	}{
		{
			name:   "Route not found",
			method: http.MethodPatch,
			url:    "/tickets/",
			apiServices: &pkg.ServiceHandler{
				TollHandler: mockHandler{},
			},
			code:     404,
			expected: `{"errors":{"code":"ERR.ROUTE.NOTFOUND","message":"The route your looking for doesn't exists"},"status":404}`,
		},
		{
			name:   "Invalid Ticket",
			method: http.MethodPatch,
			url:    "/tickets/1",
			apiServices: &pkg.ServiceHandler{
				TollHandler: mockHandler{
					RedeemTollTicketFn: func(ctx context.Context, ticket *toll.TicketToll) (*toll.TicketToll, error) {
						return nil, toll.ErrInvalidTollTicket
					},
				},
			},
			code:     400,
			expected: `{"errors":{"message":"This is an invalid toll ticket","code":"ERR.TOLL.INVALID_TOLL_TICKET"},"status":400}`,
		},
		{
			name:   "Redeemed Successfully",
			method: http.MethodPatch,
			url:    "/tickets/2",
			apiServices: &pkg.ServiceHandler{
				TollHandler: mockHandler{
					RedeemTollTicketFn: func(ctx context.Context, ticket *toll.TicketToll) (*toll.TicketToll, error) {
						return &toll.TicketToll{}, nil
					},
				},
			},
			code:     200,
			expected: `{"data":{"tollId":"","vehicleType":"","vehicleRegistrationNo":"","twoWayToll":false,"status":"","issuedTimeStamp":"0001-01-01T00:00:00Z","updatedTimeStamp":"0001-01-01T00:00:00Z"},"status":200}`,
		},
	}
	server := mockHTTPServer()
	router := server.Server
	router.NotFound = NotFound()
	for _, tt := range tests {
		server.APIHandler = tt.apiServices
		t.Run(tt.name, func(t *testing.T) {
			request, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			response := httptest.NewRecorder()
			// handler := http.HandlerFunc(h.content)

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			router.ServeHTTP(response, request)

			// Check the status code is what we expect.
			if response.Code != tt.code {
				t.Errorf("wrong status code: got %v want %v", response.Code, tt.code)
			}

			// Check the response body is what we expect.
			if !reflect.DeepEqual(strings.TrimSpace(response.Body.String()), tt.expected) {
				t.Errorf("unexpected body: got %v want %v", response.Body.String(), tt.expected)
			}
		})
	}
}

func Test_getPagination(t *testing.T) {
	type args struct {
		start string
		limit string
	}
	tests := []struct {
		name  string
		args  args
		start int
		limit int
	}{
		{
			name: "Start & Limit from Request",
			args: args{
				start: "0",
				limit: "25",
			},
			start: 0,
			limit: 25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			urlValues := "http://locahost"
			url, _ := url.Parse(urlValues)
			urlVal := url.Query()
			urlVal.Add("start", tt.args.start)
			urlVal.Add("limit", tt.args.limit)
			url.RawQuery = urlVal.Encode()
			req, _ := http.NewRequest("GET", url.String(), nil)
			start, limit := getPagination(req)
			if start != tt.start {
				t.Errorf("getPagination() got = %v, want %v", start, tt.start)
			}
			if limit != tt.limit {
				t.Errorf("getPagination() got1 = %v, want %v", limit, tt.limit)
			}
		})
	}
}
