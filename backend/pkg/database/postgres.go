package database

import (
	"fmt"
	"net/http"

	"erp-system/pkg/config"
	apperrors "erp-system/pkg/errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgres(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	if cfg.App.Env == "production" {
		gormCfg.Logger = logger.Default.LogMode(logger.Error)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormCfg)
	if err != nil {
		return nil, apperrors.NewCustomError("failed to connect to database").
			WithErrorCode(apperrors.ErrCodeDatabaseConnection).
			WithMessageID("error_database_connection").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, apperrors.NewCustomError("failed to get sql.DB").
			WithErrorCode(apperrors.ErrCodeDatabaseConnection).
			WithMessageID("error_database_connection").
			WithHTTPCode(http.StatusInternalServerError).
			WithCause(err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)

	return db, nil
}
