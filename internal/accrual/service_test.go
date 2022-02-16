package accrual

import (
	"bytes"
	"context"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: RoundTripFunc(fn),
	}
}

func TestExternalService_GetOrder(t *testing.T) {
	type fields struct {
		Client  *http.Client
		baseURL string
	}
	type args struct {
		ctx         context.Context
		orderNumber string
	}

	orderNumber := "1"
	baseURL := "http://localhost:8081"
	client := NewTestClient(func(req *http.Request) *http.Response {
		assert.Equal(t, req.URL.String(), baseURL+"/api/orders/"+orderNumber)

		return &http.Response{
			StatusCode: 200,
			Body:       ioutil.NopCloser(bytes.NewBufferString(`{"status":"PROCESSING"}`)),
			Header:     make(http.Header),
		}
	})

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *OrderResponse
		wantErr bool
	}{
		{
			name:    "should call proper endpoint and properly parse response with status",
			fields:  fields{Client: client, baseURL: baseURL},
			args:    args{ctx: context.TODO(), orderNumber: orderNumber},
			want:    &OrderResponse{Status: "PROCESSING"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &ExternalService{
				Client:  tt.fields.Client,
				baseURL: tt.fields.baseURL,
			}

			got, err := service.GetOrder(tt.args.ctx, tt.args.orderNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOrder() got = %v, want %v", got, tt.want)
			}
		})
	}
}
