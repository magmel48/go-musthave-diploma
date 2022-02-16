package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/magmel48/go-musthave-diploma/internal/auth"
	"github.com/magmel48/go-musthave-diploma/internal/logger"
	"go.uber.org/zap"
	"net/http"
)

func (app *App) getWithdrawals(context *gin.Context) {
	userID := context.MustGet(auth.UserIDKey).(int64)

	withdrawals, err := app.withdrawals.FindByUser(context, userID)
	if err != nil {
		logger.Error("GET /withdrawals error", zap.Error(err))
		context.Status(http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		context.Status(http.StatusNoContent)
		return
	}

	context.JSON(http.StatusOK, withdrawals)
}
