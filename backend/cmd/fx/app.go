package fxapp

import (
	"erp-system/cmd/fx/modules"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

// New builds and returns the fx application.
// Dependency graph is resolved automatically by fx.
func New() *fx.App {
	return fx.New(
		// Suppress default fx logging, use zap instead
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),

		// Modules (order does not matter — fx resolves the graph)
		modules.ConfigModule,   // config.Config
		modules.LoggerModule,   // *zap.Logger
		modules.DatabaseModule, // *gorm.DB  + auto-migrate
		modules.CacheModule,    // *cache.RedisCache, *cache.CacheManager
		modules.JWTModule,      // *jwt.Manager
		modules.AuthModule,     // UserRepository, RoleRepository, AuthUsecase, AuthAdapter
		modules.SupplierModule, // SupplierRepository (cached), SupplierUsecase, SupplierAdapter
		modules.ServerModule,   // Echo, Router, AuthMiddleware, lifecycle hooks
	)
}
