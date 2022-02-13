package controllers

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/magmel48/go-musthave-diploma/internal/users"
	"net/http"
)

func (app *App) register(context *gin.Context) {
	var user users.User
	err := json.NewDecoder(context.Request.Body).Decode(&user)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	if user.Login == "" || user.Password == "" {
		context.Status(http.StatusBadRequest)
		return
	}

	id, err := app.users.Create(context, user)
	if err != nil {
		context.Status(http.StatusBadRequest)
		return
	}

	session := sessions.Default(context)
	session.Set("user_id", id)
	err = session.Save()

	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{"id": id})
}
