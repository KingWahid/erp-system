package modules

import (
	"erp-system/pkg/config"
	"erp-system/pkg/database"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DatabaseModule provides PostgreSQL connection only.
// Schema is managed via SQL migrations (migrations/001_init.sql),
// NOT via GORM AutoMigrate to avoid constraint naming conflicts.
var DatabaseModule = fx.Module("database",
	fx.Provide(newDatabase),
)

func newDatabase(cfg *config.Config, logger *zap.Logger) (*gorm.DB, error) {
	db, err := database.NewPostgres(cfg)
	if err != nil {
		return nil, err
	}

	logger.Info("database connected",
		zap.String("host", cfg.Database.Host),
		zap.String("name", cfg.Database.Name),
	)
	return db, nil
}
