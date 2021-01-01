package mongo

import (
	"context"
	"testing"

	mgo "go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type mockMongoHandler struct {
	PingFn     func(ctx context.Context, rp *readpref.ReadPref) error
	DatabaseFn func(name string, opts ...*options.DatabaseOptions) *mgo.Database
}

func (m mockMongoHandler) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	return m.PingFn(ctx, rp)
}

func (m mockMongoHandler) Database(name string, opts ...*options.DatabaseOptions) *mgo.Database {
	return m.DatabaseFn(name, opts...)
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
					PingFn: func(ctx context.Context, rp *readpref.ReadPref) error {
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
