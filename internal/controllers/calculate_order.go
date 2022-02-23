package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/joeljunstrom/go-luhn"
	"github.com/magmel48/go-musthave-diploma/internal/auth"
	"github.com/magmel48/go-musthave-diploma/internal/logger"
	"github.com/magmel48/go-musthave-diploma/internal/orders"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func (app *App) calculateOrder(context *gin.Context) {
	userID := context.MustGet(auth.UserIDKey).(int64)

	payload, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		logger.Error("POST /orders: read request body error", zap.Error(err))
		context.Status(http.StatusInternalServerError)
		return
	}

	id := string(payload)
	if id == "" || !luhn.Valid(id) {
		// 422 - spec: invalid order number
		context.Status(http.StatusUnprocessableEntity)
		return
	}

	order, err := app.orders.FindByUser(context, id, userID)
	if err != nil {
		logger.Error("POST /orders: find order error", zap.Error(err))
		context.Status(http.StatusInternalServerError)
		return
	}

	// 200 - spec: order number already loaded by this user
	if order != nil {
		context.JSON(http.StatusOK, gin.H{"id": order.ID})
		return
	}

	order, err = app.orders.Create(context, id, userID)
	if err != nil {
		if errors.Is(err, orders.ErrConflict) {
			// 409 - spec: order number already loaded by another user
			context.JSON(http.StatusConflict, gin.H{"error": "order already exists"})
		} else {
			logger.Error("POST /orders: create order error", zap.Error(err))
			context.Status(http.StatusInternalServerError)
		}

		return
	}

	// 202 - spec: new order number accepted
	context.JSON(http.StatusAccepted, gin.H{"id": order.ID})
}
