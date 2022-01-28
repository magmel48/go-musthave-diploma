package main

import (
	"context"
	"github.com/alexflint/go-arg"
	"golang.org/x/sync/errgroup"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	BaseServiceAddress    string `arg:"-a,env:RUN_ADDRESS" default:"localhost:8080"`
	DatabaseDSN           string `arg:"-d,env:DATABASE_URI"`
	AccrualServiceAddress string `arg:"-r,env:ACCRUAL_SYSTEM_ADDRESS" default:"localhost:8081"`
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	cfg := config{}
	arg.MustParse(&cfg)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	server := http.Server{
		Addr:        cfg.BaseServiceAddress,
		IdleTimeout: time.Second,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	eg, ctx := errgroup.WithContext(ctx)

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
