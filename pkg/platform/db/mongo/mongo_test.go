package mongo

import (
	"context"
	"testing"

	"github.com/globalsign/mgo"
)

type mockMongoHandler struct {
	PingFn func() error
	DBFn   func(name string) *mgo.Database
}

func (m mockMongoHandler) Ping() error {
	return m.PingFn()
}

func (m mockMongoHandler) DB(name string) *mgo.Database {
	return m.DBFn(name)
}

func Test_handler_Health(t *testing.T) {
	type fields struct {
		dbClient mockMongoHandler
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				dbClient: mockMongoHandler{
					PingFn: func() error {
						return nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &handler{
				dbClient: tt.fields.dbClient,
			}
			if err := h.Health(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("handler.Health() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
