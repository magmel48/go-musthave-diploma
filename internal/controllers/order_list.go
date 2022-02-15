package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/magmel48/go-musthave-diploma/internal/auth"
	"net/http"
)

func (app *App) orderList(context *gin.Context) {
	userID := context.MustGet(auth.UserIDKey).(int64)

	orders, err := app.orders.FindUserOrders(context, userID)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		context.Status(http.StatusNoContent)
		return
	}

	context.JSON(http.StatusOK, orders)
}
