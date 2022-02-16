package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joeljunstrom/go-luhn"
	"github.com/magmel48/go-musthave-diploma/internal/auth"
	"github.com/magmel48/go-musthave-diploma/internal/balances"
	"github.com/magmel48/go-musthave-diploma/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

type WithdrawalRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func (app *App) withdraw(context *gin.Context) {
	userID := context.MustGet(auth.UserIDKey).(int64)

	var withdrawalRequest WithdrawalRequest
	err := json.NewDecoder(context.Request.Body).Decode(&withdrawalRequest)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	if !luhn.Valid(withdrawalRequest.Order) || withdrawalRequest.Sum <= 0 {
		context.Status(http.StatusUnprocessableEntity)
		return
	}

	err = app.balances.Change(context, userID, -withdrawalRequest.Sum)
	if err != nil {
		if errors.Is(err, balances.ErrInsufficientFunds) {
			context.JSON(http.StatusPaymentRequired, gin.H{"error": "insufficient funds"})
		} else {
			context.Status(http.StatusInternalServerError)
		}

		return
	}

	// if balance check is ok - add the request into withdrawals list
	err = app.withdrawals.Create(context, userID, withdrawalRequest.Order, withdrawalRequest.Sum)
	if err != nil {
		logger.Error("POST /withdrawals error", zap.Error(err))
		context.Status(http.StatusInternalServerError)
		return
	}

	context.Status(http.StatusOK)
}
