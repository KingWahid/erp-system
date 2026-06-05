package modules

import (
	"erp-system/pkg/config"

	"go.uber.org/fx"
)

// ConfigModule provides application configuration
var ConfigModule = fx.Module("config",
	fx.Provide(config.Load),
)
