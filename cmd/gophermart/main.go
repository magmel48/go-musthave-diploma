package main

import (
	"context"
	"github.com/magmel48/go-musthave-diploma/internal/config"
	"github.com/magmel48/go-musthave-diploma/internal/controllers"
	"github.com/magmel48/go-musthave-diploma/internal/logger"
	"golang.org/x/sync/errgroup"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// no error handling here, see https://github.com/uber-go/zap/issues/880
	defer logger.Sync()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	app := controllers.App{Context: ctx}
	err := app.Init()
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Addr:        config.BaseServiceAddress,
		IdleTimeout: time.Second,
		Handler:     app.Handler(),
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	eg.Go(func() error {
		logger.Info("starting loyal system service...")
		return server.ListenAndServe()
	})

	eg.Go(func() error {
		<-ctx.Done()

		logger.Info("stopping the service...")

		err := server.Shutdown(ctx)
		return err
	})

	eg.Wait()
}
