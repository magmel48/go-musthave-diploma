package config

import "github.com/alexflint/go-arg"

type config struct {
	baseServiceAddress    string `arg:"-a,env:RUN_ADDRESS" default:"localhost:8080"`
	databaseDSN           string `arg:"-d,env:DATABASE_URI"`
	accrualServiceAddress string `arg:"-r,env:ACCRUAL_SYSTEM_ADDRESS" default:"localhost:8081"`
}

var (
	BaseServiceAddress    string
	DatabaseDSN           string
	AccrualServiceAddress string
)

func init() {
	cfg := config{}
	arg.MustParse(&cfg)
}
