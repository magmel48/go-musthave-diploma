package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func register(context *gin.Context) {
	session := sessions.Default(context)

	session.Set("user_id", "0")
	err := session.Save()
	if err != nil {
		context.Status(http.StatusInternalServerError)
		return
	}

	context.JSON(http.StatusOK, gin.H{"id": "0"})
}
