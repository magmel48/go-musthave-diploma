package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/magmel48/go-musthave-diploma/internal/auth"
	"github.com/magmel48/go-musthave-diploma/internal/balances"
	"github.com/magmel48/go-musthave-diploma/internal/balances/mocks"
	"github.com/magmel48/go-musthave-diploma/internal/orders"
	"github.com/magmel48/go-musthave-diploma/internal/users"
	"github.com/magmel48/go-musthave-diploma/internal/withdrawals"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"testing"
)

func TestApp_getBalance(t *testing.T) {
	type fields struct {
		ctx         context.Context
		auth        auth.Auth
		users       users.Repository
		orders      orders.Repository
		balances    balances.Repository
		withdrawals withdrawals.Repository
	}
	type args struct {
		context *gin.Context
	}
	type want struct {
		body string
	}

	balancesMock := &mocks.Repository{}
	balancesMock.On("FindByUser", mock.Anything, mock.Anything).Return(&balances.Balance{Current: 99}, nil)

	recorder := httptest.NewRecorder()

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "should return current and withdrawn balances of the user",
			fields: fields{ctx: context.TODO(), balances: balancesMock},
			args:   args{context: createGinContext(recorder, "")},
			want:   want{body: `{"current": 99, "withdrawn": 0}`},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &App{
				ctx:         tt.fields.ctx,
				auth:        tt.fields.auth,
				users:       tt.fields.users,
				orders:      tt.fields.orders,
				balances:    tt.fields.balances,
				withdrawals: tt.fields.withdrawals,
			}

			app.getBalance(tt.args.context)
			assert.JSONEq(t, tt.want.body, recorder.Body.String())
		})
	}
}
