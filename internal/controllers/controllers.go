package controllers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/magmel48/go-musthave-diploma/internal/auth"
	"github.com/magmel48/go-musthave-diploma/internal/config"
	"github.com/magmel48/go-musthave-diploma/internal/orders"
	"github.com/magmel48/go-musthave-diploma/internal/users"
	"time"
)

type App struct {
	Context context.Context
	db      *sql.DB
	users   users.Repository
	orders  orders.Repository
	auth    auth.Auth
}

const UserIDKey = "user_id"

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

	app.users = users.NewRepository(app.db)
	app.auth = auth.NewService(app.users)

	app.orders = orders.NewRepository(app.db)

	return nil
}

func (app *App) Handler() *gin.Engine {
	r := gin.Default()

	store := cookie.NewStore([]byte(config.SessionsSecret))
	r.Use(sessions.Sessions("session", store))

	r.POST("/api/user/register", app.register)
	r.POST("/api/user/login", app.login)
	r.POST("/api/user/orders", app.calculateOrder)
	r.GET("/api/user/orders", orderList)
	r.GET("/api/user/balance", getBalance)
	r.POST("/api/user/balance/withdraw", withdraw)
	r.GET("/api/user/balance/withdrawals", getWithdrawals)

	return r
}
