package auth

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Middleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		session := sessions.Default(context)
		sessionUserID := session.Get(UserIDKey)

		userID := int64(0)

		switch sessionUserID.(type) {
		case int64:
			userID = sessionUserID.(int64)
		default:
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		context.Set(UserIDKey, userID)
		context.Next()
	}
}
