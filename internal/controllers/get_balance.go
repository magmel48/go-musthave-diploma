package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/magmel48/go-musthave-diploma/internal/auth"
	"net/http"
)

func (app *App) getBalance(context *gin.Context) {
	userID := context.MustGet(auth.UserIDKey).(int64)

	balance, err := app.balances.FindByUser(context, userID)
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, balance)
}
