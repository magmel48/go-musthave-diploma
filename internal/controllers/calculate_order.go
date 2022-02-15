package controllers

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joeljunstrom/go-luhn"
	"github.com/magmel48/go-musthave-diploma/internal/logger"
	"github.com/magmel48/go-musthave-diploma/internal/orders"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func (app *App) calculateOrder(context *gin.Context) {
	session := sessions.Default(context)
	sessionUserID := session.Get(UserIDKey)

	userID := int64(0)

	switch sessionUserID.(type) {
	case int64:
		userID = sessionUserID.(int64)
	default:
		context.Status(http.StatusUnauthorized)
		return
	}

	payload, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		logger.Error("POST /orders: read request body error", zap.Error(err))
		context.Status(http.StatusInternalServerError)
		return
	}

	id := string(payload)
	if !luhn.Valid(id) {
		context.Status(http.StatusUnprocessableEntity)
		return
	}

	order, err := app.orders.Create(context, orders.Order{UserID: userID, Number: id})
	if err != nil {
		if errors.Is(err, orders.ErrConflict) {
			context.JSON(http.StatusConflict, gin.H{"error": "order already exists"})
		} else {
			context.Status(http.StatusInternalServerError)
		}

		return
	}

	context.JSON(http.StatusAccepted, gin.H{"id": order.ID})
}
