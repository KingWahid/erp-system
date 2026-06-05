package modules

import (
	"erp-system/pkg/config"
	"erp-system/pkg/logger"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

// LoggerModule provides zap logger
var LoggerModule = fx.Module("logger",
	fx.Provide(newLogger),
)

func newLogger(cfg *config.Config) (*zap.Logger, error) {
	return logger.New(cfg.App.Env)
}
