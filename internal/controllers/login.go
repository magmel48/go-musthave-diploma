package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/magmel48/go-musthave-diploma/internal/auth"
	"github.com/magmel48/go-musthave-diploma/internal/logger"
	"github.com/magmel48/go-musthave-diploma/internal/users"
	"go.uber.org/zap"
	"net/http"
)

func (app *App) login(context *gin.Context) {
	var user users.User
	err := json.NewDecoder(context.Request.Body).Decode(&user)
	if err != nil {
		logger.Error("/login: unmarshal json payload error", zap.Error(err))

		context.Status(http.StatusBadRequest)
		return
	}

	if id, err := app.auth.CheckUser(context, user); err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			context.Status(http.StatusUnauthorized)
		} else {
			logger.Error("/login: check user error", zap.Error(err))
			context.Status(http.StatusInternalServerError)
		}

		return
	} else {
		user.ID = id
	}

	context.JSON(http.StatusOK, gin.H{"id": user.ID})
}
