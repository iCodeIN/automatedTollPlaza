package errors

import (
	"errors"
	"reflect"
	"testing"
)

func TestTollErr_Error(t *testing.T) {
	type fields struct {
		Message string
		Code    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "New Error String",
			fields: fields{
				Message: "Message",
				Code:    "ERR",
			},
			want: "ERR : Message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewErrorWithCode(tt.fields.Code, tt.fields.Message)
			if got := e.Error(); got != tt.want {
				t.Errorf("TollErr.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToTollError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want *TollErr
	}{
		{
			name: "no error",
			args: args{
				err: nil,
			},
			want: nil,
		},
		{
			name: "toll error type",
			args: args{
				err: ErrMissingFields,
			},
			want: &ErrMissingFields,
		},
		{
			name: "not a toll error type",
			args: args{
				err: errors.New("Test"),
			},
			want: &TollErr{
				Code:    "ERR",
				Message: errors.New("Test").Error(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToTollError(tt.args.err); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToTollError() = %v, want %v", got, tt.want)
			}
		})
	}
}
