package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joeljunstrom/go-luhn"
	"github.com/magmel48/go-musthave-diploma/internal/auth"
	"github.com/magmel48/go-musthave-diploma/internal/balances"
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

	if !luhn.Valid(withdrawalRequest.Order) {
		context.Status(http.StatusUnprocessableEntity)
		return
	}

	err = app.balances.Change(context, userID, withdrawalRequest.Sum)
	if err != nil {
		if errors.Is(err, balances.ErrInsufficientFunds) {
			context.JSON(http.StatusPaymentRequired, gin.H{"error": "insufficient funds"})
		} else {
			context.Status(http.StatusInternalServerError)
		}

		return
	}

	// if balance check is ok - add the request into withdrawals list
	// TODO implement

	context.Status(http.StatusOK)
}
