package controllers

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/magmel48/go-musthave-diploma/internal/logger"
	"github.com/magmel48/go-musthave-diploma/internal/orders"
	"github.com/theplant/luhn"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strconv"
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

	// TODO change luhn to own implementation
	id, err := strconv.Atoi(string(payload))
	if err != nil || !luhn.Valid(id) {
		context.Status(http.StatusUnprocessableEntity)
		return
	}

	order, err := app.orders.Create(context, orders.Order{UserID: userID, Number: strconv.Itoa(id)})
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
