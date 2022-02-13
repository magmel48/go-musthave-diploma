package config

import (
	"github.com/alexflint/go-arg"
)

type config struct {
	BaseServiceAddress    string `arg:"-a,env:RUN_ADDRESS" default:"localhost:8080"`
	DatabaseDSN           string `arg:"-d,env:DATABASE_URI"`
	AccrualServiceAddress string `arg:"-r,env:ACCRUAL_SYSTEM_ADDRESS" default:"localhost:8081"`
	SessionsSecret        string `arg:"-s,env:SESSIONS_SECRET" default:"very_big_secret"`
}

var (
	BaseServiceAddress    string
	DatabaseDSN           string
	AccrualServiceAddress string
	SessionsSecret        string
)

func init() {
	cfg := config{}
	arg.MustParse(&cfg)

	BaseServiceAddress = cfg.BaseServiceAddress
	DatabaseDSN = cfg.DatabaseDSN
	AccrualServiceAddress = cfg.AccrualServiceAddress
	SessionsSecret = cfg.SessionsSecret
}
