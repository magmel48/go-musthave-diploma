package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/magmel48/go-musthave-diploma/internal/auth"
	"github.com/magmel48/go-musthave-diploma/internal/balances"
	"github.com/magmel48/go-musthave-diploma/internal/orders"
	"github.com/magmel48/go-musthave-diploma/internal/users"
	"github.com/magmel48/go-musthave-diploma/internal/withdrawals"
	"github.com/magmel48/go-musthave-diploma/internal/withdrawals/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApp_getWithdrawals(t *testing.T) {
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
		statusCode int
	}

	recorder := httptest.NewRecorder()

	withdrawalsMock := &mocks.Repository{}
	withdrawalsMock.On("ListByUser", mock.Anything, mock.Anything).Return(nil, nil)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "should return 204 if no withdrawals were created before",
			fields: fields{ctx: context.TODO(), withdrawals: withdrawalsMock},
			args:   args{context: createGinContext(recorder, "")},
			want:   want{statusCode: http.StatusNoContent},
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

			app.getWithdrawals(tt.args.context)
			assert.Equal(t, tt.want.statusCode, tt.args.context.Writer.Status())
		})
	}
}
