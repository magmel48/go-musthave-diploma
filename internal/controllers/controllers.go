package controllers

import (
	"context"
	"database/sql"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/magmel48/go-musthave-diploma/internal/auth"
	"github.com/magmel48/go-musthave-diploma/internal/balances"
	"github.com/magmel48/go-musthave-diploma/internal/config"
	"github.com/magmel48/go-musthave-diploma/internal/daemons"
	"github.com/magmel48/go-musthave-diploma/internal/orders"
	"github.com/magmel48/go-musthave-diploma/internal/users"
	"github.com/magmel48/go-musthave-diploma/internal/withdrawals"
	"time"
)

type App struct {
	ctx         context.Context
	auth        auth.Auth
	users       users.Repository
	orders      orders.Repository
	balances    balances.Repository
	withdrawals withdrawals.Repository
}

func (app *App) Init(ctx context.Context) error {
	if config.DatabaseDSN == "" {
		return errors.New("no db connection details provided")
	}

	db, err := sql.Open("pgx", config.DatabaseDSN)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(30)
	db.SetMaxIdleConns(30)
	db.SetConnMaxIdleTime(10 * time.Second)

	app.users = users.NewPostgreSQLRepository(db)
	app.auth = auth.NewService(app.users)

	app.orders = orders.NewPostgreSQLRepository(db)
	app.balances = balances.NewPostgreSQLRepository(db)
	app.withdrawals = withdrawals.NewPostgreSQLRepository(db)

	// run daemon for order updates checking
	daemon := daemons.NewExternalAccrualJob(ctx, db)
	s := gocron.NewScheduler(time.UTC)
	_, err = s.Every(5).Seconds().Do(daemon.Start)
	if err != nil {
		return err
	}

	s.StartAsync()

	return nil
}

func (app *App) Handler() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())

	store := cookie.NewStore([]byte(config.SessionsSecret))
	r.Use(sessions.Sessions("session", store))

	r.POST("/api/user/register", app.register)
	r.POST("/api/user/login", app.login)

	authorized := r.Group("/api/user")
	authorized.Use(auth.Middleware())
	{
		authorized.POST("/orders", app.calculateOrder)
		authorized.GET("/orders", app.orderList)
		authorized.GET("/balance", app.getBalance)
		authorized.POST("/balance/withdraw", app.withdraw)
		authorized.GET("/balance/withdrawals", app.getWithdrawals)
	}

	return r
}
