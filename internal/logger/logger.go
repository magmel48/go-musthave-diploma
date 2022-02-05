package logger

import "go.uber.org/zap"

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

func Sync() error {
	return logger.Sync()
}
