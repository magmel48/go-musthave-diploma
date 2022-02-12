package config

import "github.com/alexflint/go-arg"

type config struct {
	baseServiceAddress    string `arg:"-a,env:RUN_ADDRESS" default:"localhost:8080"`
	databaseDSN           string `arg:"-d,env:DATABASE_URI"`
	accrualServiceAddress string `arg:"-r,env:ACCRUAL_SYSTEM_ADDRESS" default:"localhost:8081"`
	sessionsSecret        string `arg:"-s,env:SESSIONS_SECRET" default:"very_big_secret"`
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

	BaseServiceAddress = cfg.baseServiceAddress
	DatabaseDSN = cfg.databaseDSN
	AccrualServiceAddress = cfg.accrualServiceAddress
	SessionsSecret = cfg.sessionsSecret
}
