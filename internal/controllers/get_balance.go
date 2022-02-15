package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (app *App) getBalance(context *gin.Context) {
	// TODO implement
	context.Status(http.StatusInternalServerError)
}
