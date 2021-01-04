package toll

import (
	"automatedTollPlaze/config"
	"automatedTollPlaze/pkg/platform/appcontext"
	"automatedTollPlaze/pkg/platform/db"
	"context"
	"reflect"
	"testing"
	"time"
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

func mockAppContext() *appcontext.AppContext {
	return &appcontext.AppContext{
		DbClient:  &mockDBHandler{},
		StartTime: time.Now().UTC(),
		Config: config.Cfg{
			Pricing: config.Pricing{
				Default: config.PriceValue{
					OneWay: 100,
					TwoWay: 200,
				},
				VehicleType: map[string]config.PriceValue{
					"light": {
						OneWay: 100,
						TwoWay: 200,
					},
					"moderate": {
						OneWay: 100,
						TwoWay: 200,
					},
					"heavy": {
						OneWay: 100,
						TwoWay: 200,
					},
				},
			},
		},
	}
}

func Test_handler_IssueTollTicket(t *testing.T) {
	type fields struct {
		AppCtx *appcontext.AppContext
	}
	type args struct {
		ctx       context.Context
		ticket    *TicketToll
		dbHandler db.Service
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error",
			fields: fields{
				AppCtx: mockAppContext(),
			},
			args: args{
				ctx: context.Background(),
				ticket: &TicketToll{
					RegistrationNo: "KA 01 EA 1234",
				},
				dbHandler: mockDBHandler{
					FindOneFn: func(ctx context.Context, params db.Params) error {
						return ErrInvalidTollTicket
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no error - pending toll tickets",
			fields: fields{
				AppCtx: mockAppContext(),
			},
			args: args{
				ctx: context.Background(),
				ticket: &TicketToll{
					RegistrationNo: "KA 01 EA 1234",
				},
				dbHandler: mockDBHandler{
					FindOneFn: func(ctx context.Context, params db.Params) error {
						return nil
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Error in generating toll tickets",
			fields: fields{
				AppCtx: mockAppContext(),
			},
			args: args{
				ctx: context.Background(),
				ticket: &TicketToll{
					TollID:           "1",
					ReturnTollTicket: false,
					RegistrationNo:   "KA 01 EA 1234",
				},
				dbHandler: mockDBHandler{
					FindOneFn: func(ctx context.Context, params db.Params) error {
						return db.ErrNotFound
					},
					InsertOneFn: func(ctx context.Context, params db.Params) error {
						return ErrInvalidTollTicket
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Successfully created a new toll ticket",
			fields: fields{
				AppCtx: mockAppContext(),
			},
			args: args{
				ctx: context.Background(),
				ticket: &TicketToll{
					TollID:           "1",
					VehicleType:      "light",
					ReturnTollTicket: false,
					RegistrationNo:   "KA 01 EA 1234",
				},
				dbHandler: mockDBHandler{
					FindOneFn: func(ctx context.Context, params db.Params) error {
						return db.ErrNotFound
					},
					InsertOneFn: func(ctx context.Context, params db.Params) error {
						return nil
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.AppCtx.DbClient = tt.args.dbHandler
			s := NewTollService(tt.args.ctx, tt.fields.AppCtx)
			if err := s.IssueTollTicket(tt.args.ctx, tt.args.ticket); (err != nil) != tt.wantErr {
				t.Errorf("handler.IssueTollTicket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_handler_GetTicketIssueList(t *testing.T) {
	type fields struct {
		AppCtx *appcontext.AppContext
	}
	type args struct {
		ctx       context.Context
		params    *TicketListRequest
		dbHandler db.Service
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   TicketListResponse
	}{
		{
			name: "List",
			fields: fields{
				AppCtx: mockAppContext(),
			},
			args: args{
				ctx: context.Background(),
				params: &TicketListRequest{
					TollID:         "1",
					RegistrationNo: "KA 01 EA 1234",
					Status:         "REDEEMED",
				},
				dbHandler: mockDBHandler{
					FindAllFn: func(ctx context.Context, params db.Params) error {
						return nil
					},
					CountFn: func(ctx context.Context, params db.Params) int {
						return 0
					},
				},
			},
			want: TicketListResponse{
				List: make([]TicketToll, 0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.AppCtx.DbClient = tt.args.dbHandler
			s := NewTollService(tt.args.ctx, tt.fields.AppCtx)
			if got := s.GetTicketIssueList(tt.args.ctx, tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handler.GetTicketIssueList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_GetTollTicketDetails(t *testing.T) {
	type fields struct {
		AppCtx *appcontext.AppContext
	}
	type args struct {
		ctx       context.Context
		ticket    *TicketToll
		dbHandler db.Service
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TicketToll
		wantErr bool
	}{
		{
			name: "Not found",
			fields: fields{
				AppCtx: mockAppContext(),
			},
			args: args{
				ctx: context.Background(),
				ticket: &TicketToll{
					TicketID: "1",
				},
				dbHandler: mockDBHandler{
					FindOneFn: func(ctx context.Context, params db.Params) error {
						return db.ErrNotFound
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "found",
			fields: fields{
				AppCtx: mockAppContext(),
			},
			args: args{
				ctx: context.Background(),
				ticket: &TicketToll{
					TicketID: "1",
				},
				dbHandler: mockDBHandler{
					FindOneFn: func(ctx context.Context, params db.Params) error {
						if val, ok := params.Result.(*TicketToll); ok {
							val.TicketID = "1"
							val.TollID = "1"
							val.RegistrationNo = "KA 01 EA 1234"
							val.ReturnTollTicket = true
							val.Status = "ISSUED"
							val.Price = 200
						}
						return nil
					},
				},
			},
			want: &TicketToll{
				TicketID:         "1",
				TollID:           "1",
				RegistrationNo:   "KA 01 EA 1234",
				ReturnTollTicket: true,
				Status:           "ISSUED",
				Price:            200,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.AppCtx.DbClient = tt.args.dbHandler
			s := NewTollService(tt.args.ctx, tt.fields.AppCtx)
			got, err := s.GetTollTicketDetails(tt.args.ctx, tt.args.ticket)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.GetTollTicketDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handler.GetTollTicketDetails() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handler_RedeemTollTicket(t *testing.T) {
	type fields struct {
		AppCtx *appcontext.AppContext
	}
	type args struct {
		ctx       context.Context
		ticket    *TicketToll
		dbHandler db.Service
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *TicketToll
		wantErr bool
	}{
		{
			name: "Not found",
			fields: fields{
				AppCtx: mockAppContext(),
			},
			args: args{
				ctx: context.Background(),
				ticket: &TicketToll{
					TicketID: "1",
				},
				dbHandler: mockDBHandler{
					FindOneFn: func(ctx context.Context, params db.Params) error {
						return db.ErrNotFound
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Already Redeemed",
			fields: fields{
				AppCtx: mockAppContext(),
			},
			args: args{
				ctx: context.Background(),
				ticket: &TicketToll{
					TicketID: "1",
				},
				dbHandler: mockDBHandler{
					FindOneFn: func(ctx context.Context, params db.Params) error {
						if val, ok := params.Result.(*TicketToll); ok {
							val.TicketID = "1"
							val.Status = "REDEEMED"
						}
						return nil
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Updated Failed",
			fields: fields{
				AppCtx: mockAppContext(),
			},
			args: args{
				ctx: context.Background(),
				ticket: &TicketToll{
					TicketID: "1",
				},
				dbHandler: mockDBHandler{
					FindOneFn: func(ctx context.Context, params db.Params) error {
						if val, ok := params.Result.(*TicketToll); ok {
							val.TicketID = "1"
							val.Status = "ISSUED"
						}
						return nil
					},
					UpsertFn: func(ctx context.Context, params db.Params) error {
						return ErrInvalidTollTicket
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Successfully Redeemed",
			fields: fields{
				AppCtx: mockAppContext(),
			},
			args: args{
				ctx: context.Background(),
				ticket: &TicketToll{
					TicketID: "1",
				},
				dbHandler: mockDBHandler{
					FindOneFn: func(ctx context.Context, params db.Params) error {
						if val, ok := params.Result.(*TicketToll); ok {
							t := time.Now().Add(time.Hour * 24)
							val.TicketID = "1"
							val.RedeemBy = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
							val.Status = "ISSUED"
						}
						return nil
					},
					UpsertFn: func(ctx context.Context, params db.Params) error {
						return nil
					},
				},
			},
			want: &TicketToll{
				TicketID: "1",
				Status:   "REDEEMED",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fields.AppCtx.DbClient = tt.args.dbHandler
			s := NewTollService(tt.args.ctx, tt.fields.AppCtx)
			got, err := s.RedeemTollTicket(tt.args.ctx, tt.args.ticket)
			if (err != nil) != tt.wantErr {
				t.Errorf("handler.RedeemTollTicket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				got.IssuedTimeStamp = time.Time{}
				got.RedeemBy = time.Time{}
				got.UpdatedTimeStamp = time.Time{}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("handler.RedeemTollTicket() = %v, want %v", got, tt.want)
			}
		})
	}
}
