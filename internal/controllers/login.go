package controllers

import (
	"encoding/json"
	"errors"
	"github.com/gin-contrib/sessions"
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
		logger.Error("POST /login: unmarshal json payload error", zap.Error(err))
		context.Status(http.StatusBadRequest)
		return
	}

	if id, err := app.auth.CheckUser(context, user); err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			context.Status(http.StatusUnauthorized)
		} else {
			logger.Error("POST /login: check user error", zap.Error(err))
			context.Status(http.StatusInternalServerError)
		}

		return
	} else {
		user.ID = id
	}

	session := sessions.Default(context)
	session.Set(UserIDKey, user.ID)
	err = session.Save()

	if err != nil {
		logger.Error("POST /login: saving session error")
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{"id": user.ID})
}
