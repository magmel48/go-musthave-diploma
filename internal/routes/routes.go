package routes

import (
	"database/sql"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/magmel48/go-musthave-diploma/internal/config"
	"time"
)

type App struct {
	db *sql.DB
}

func (app *App) Init() error {
	if config.DatabaseDSN == "" {
		return errors.New("no db connection details provided")
	}

	var err error
	app.db, err = sql.Open("pgx", config.DatabaseDSN)
	if err != nil {
		return err
	}

	app.db.SetMaxOpenConns(30)
	app.db.SetMaxIdleConns(30)
	app.db.SetConnMaxIdleTime(10 * time.Second)

	// TODO create repositories with passing connection into all of them
	return nil
}

func (app *App) Handler() *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte(config.SessionsSecret))
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
