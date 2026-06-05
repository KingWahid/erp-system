// migrate is a standalone CLI tool for running SQL migrations and seeds.
//
// Usage:
//
//	go run ./cmd/migrate -action=migrate    → run 001_init.sql
//	go run ./cmd/migrate -action=seed-mock  → run mock_data.sql
//	go run ./cmd/migrate -action=drop       → drop all tables (DANGER)
package main

import (
	"flag"
	"fmt"
	"os"

	"erp-system/pkg/config"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func main() {
	action := flag.String("action", "migrate", "Action: migrate | seed-mock | drop")
	flag.Parse()

	// Load .env
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("warning: .env not found, using environment variables")
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Init logger
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() //nolint:errcheck

	// Connect DB
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		logger.Fatal("failed to connect database", zap.Error(err))
	}

	logger.Info("connected to database",
		zap.String("host", cfg.Database.Host),
		zap.String("db", cfg.Database.Name),
	)

	switch *action {
	case "migrate":
		runMigration(db, logger)
	case "seed-mock":
		runSeedMock(db, logger)
	case "drop":
		runDrop(db, logger)
	default:
		logger.Fatal("unknown action", zap.String("action", *action))
	}
}

// runMigration runs the SQL schema migration
func runMigration(db *gorm.DB, logger *zap.Logger) {
	logger.Info("running migration: 001_init.sql")

	sql, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		logger.Fatal("failed to read migration file", zap.Error(err))
	}

	if err := db.Exec(string(sql)).Error; err != nil {
		logger.Fatal("migration failed", zap.Error(err))
	}

	logger.Info("migration completed successfully")
}

// runSeedMock inserts mock data for development/testing
func runSeedMock(db *gorm.DB, logger *zap.Logger) {
	logger.Info("running seed: mock_data.sql")

	sql, err := os.ReadFile("migrations/seed/mock_data.sql")
	if err != nil {
		logger.Fatal("failed to read seed file", zap.Error(err))
	}

	if err := db.Exec(string(sql)).Error; err != nil {
		logger.Fatal("seed failed", zap.Error(err))
	}

	logger.Info("mock data seeded successfully")
	printSeedSummary(db, logger)
}

// runDrop drops all application tables (DANGEROUS — dev only)
func runDrop(db *gorm.DB, logger *zap.Logger) {
	logger.Warn("DROPPING ALL TABLES — this cannot be undone")

	tables := []string{
		"supplier_invoices",
		"supplier_stage_histories",
		"supplier_performance_ratings",
		"supplier_materials",
		"supplier_groups",
		"supplier_addresses",
		"supplier_contacts",
		"suppliers",
		"user_roles",
		"role_permissions",
		"permissions",
		"roles",
		"users",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			logger.Error("failed to drop table", zap.String("table", table), zap.Error(err))
		} else {
			logger.Info("dropped table", zap.String("table", table))
		}
	}

	logger.Info("all tables dropped")
}

// printSeedSummary logs row counts after seeding
func printSeedSummary(db *gorm.DB, logger *zap.Logger) {
	type tableCount struct {
		table string
		count int64
	}

	tables := []string{
		"users", "roles", "permissions", "user_roles", "role_permissions",
		"suppliers", "supplier_addresses", "supplier_contacts", "supplier_groups",
		"supplier_materials", "supplier_performance_ratings", "supplier_stage_histories",
		"supplier_invoices",
	}

	logger.Info("── seed summary ──────────────────")
	for _, table := range tables {
		var count int64
		db.Raw(fmt.Sprintf("SELECT COUNT(*) FROM %s", table)).Scan(&count)
		logger.Info("table count",
			zap.String("table", table),
			zap.Int64("rows", count),
		)
	}
	logger.Info("──────────────────────────────────")
}
