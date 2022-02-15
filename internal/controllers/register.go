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
		logger.Error("POST /register: unmarshal json payload error", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if user.Login == "" || user.Password == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}

	u, err := app.auth.CreateNew(context, user)
	if err != nil {
		logger.Error("POST /register: user create error", zap.Error(err))

		if errors.Is(err, users.ErrConflict) {
			context.JSON(http.StatusConflict, gin.H{"error": "login already exists"})
		} else {
			context.Status(http.StatusInternalServerError)
		}

		return
	}

	session := sessions.Default(context)
	session.Set(auth.UserIDKey, u.ID)
	err = session.Save()

	if err != nil {
		logger.Error("POST /register: saving session error")
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{"id": u.ID})
}
