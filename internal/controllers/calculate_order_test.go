package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/magmel48/go-musthave-diploma/internal/auth"
	"github.com/magmel48/go-musthave-diploma/internal/balances"
	"github.com/magmel48/go-musthave-diploma/internal/orders"
	"github.com/magmel48/go-musthave-diploma/internal/orders/mocks"
	"github.com/magmel48/go-musthave-diploma/internal/users"
	"github.com/magmel48/go-musthave-diploma/internal/withdrawals"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func createContext(body string) *gin.Context {
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Set(auth.UserIDKey, int64(1))

	req, _ := http.NewRequest("", "", strings.NewReader(body))
	ctx.Request = req

	return ctx
}

func TestApp_calculateOrder(t *testing.T) {
	type fields struct {
		ctx         context.Context
		auth        auth.Auth
		users       users.Repository
		orders      orders.Repository
		balances    balances.Repository
		withdrawals withdrawals.Repository
		req         *http.Request
	}
	type args struct {
		context     *gin.Context
		orderNumber string
	}
	type want struct {
		statusCode int
	}

	gin.SetMode(gin.TestMode)

	wrongOrderNumber := "123"
	correctOrderNumber := "12345678903"

	ordersMock := &mocks.Repository{}
	ordersMock.On("FindByUser", mock.Anything, mock.Anything, mock.Anything).Return(&orders.Order{}, nil)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "returns 422 if order number is invalid",
			fields: fields{ctx: context.TODO()},
			args:   args{context: createContext(wrongOrderNumber)},
			want:   want{statusCode: http.StatusUnprocessableEntity},
		},
		{
			name:   "returns 200 if the order number was already registered by the user",
			fields: fields{ctx: context.TODO(), orders: ordersMock},
			args:   args{context: createContext(correctOrderNumber)},
			want:   want{statusCode: http.StatusOK},
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

			app.calculateOrder(tt.args.context)
			assert.Equal(t, tt.want.statusCode, tt.args.context.Writer.Status())
		})
	}
}
