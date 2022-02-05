package routes

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func Handler() *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte("very_big_secret")) // could be env var
	r.Use(sessions.Sessions("session", store))

	r.POST("/api/user/register", register)
	r.POST("/api/user/login", login)
	r.POST("/api/user/orders", calculateOrder)
	r.GET("/api/user/orders", orderList)
	r.GET("/api/user/balance", getBalance)
	r.POST("/api/user/balance/withdraw", withdraw)
	r.GET("/api/user/balance/withdrawals", getWithdrawals)

	return r
}
