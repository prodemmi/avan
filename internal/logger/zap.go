package logger

import (
	"go.uber.org/zap"
)

func NewLogger() zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.DisableStacktrace = true
	cfg.DisableCaller = true

	logger, _ := cfg.Build()

	defer logger.Sync()

	return *logger
}
