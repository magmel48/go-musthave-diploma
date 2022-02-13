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

func (app *App) register(context *gin.Context) {
	var user users.User
	err := json.NewDecoder(context.Request.Body).Decode(&user)
	if err != nil {
		logger.Error("/register: unmarshal json payload error", zap.Error(err))

		context.Status(http.StatusBadRequest)
		return
	}

	if user.Login == "" || user.Password == "" {
		context.Status(http.StatusBadRequest)
		return
	}

	id, err := app.auth.CreateNew(context, user)
	if err != nil {
		logger.Error("/register: user create error", zap.Error(err))

		if errors.Is(err, auth.ErrConflict) {
			context.Status(http.StatusConflict)
		} else {
			context.Status(http.StatusBadRequest)
		}

		return
	}

	session := sessions.Default(context)
	session.Set("user_id", id)
	err = session.Save()

	if err != nil {
		logger.Error("/register: saving session error")

		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{"id": id})
}
