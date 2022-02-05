package main

import (
	"context"
	"github.com/magmel48/go-musthave-diploma/internal/config"
	"github.com/magmel48/go-musthave-diploma/internal/logger"
	"github.com/magmel48/go-musthave-diploma/internal/routes"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	defer func() {
		// each application should flush logs before actual closing
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	eg, ctx := errgroup.WithContext(ctx)
	server := http.Server{
		Addr:        config.BaseServiceAddress,
		IdleTimeout: time.Second,
		Handler:     routes.Handler(),
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	eg.Go(func() error {
		log.Println("starting loyal system service...")
		return server.ListenAndServe()
	})

	eg.Go(func() error {
		<-ctx.Done()

		log.Println("stopping the service...")

		err := server.Shutdown(ctx)
		return err
	})

	log.Println(eg.Wait())
}
