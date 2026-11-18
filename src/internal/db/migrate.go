package db

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"os"
	"path/filepath"
)

func postgresMigrationsPath() string {
	path, _ := os.Getwd()
	return filepath.Join(path, "internal", "db", "migrations", "postgresql")
}

func MigratePostgres(logger zap.Logger, db *sql.DB) {
	if err := goose.SetDialect("postgres"); err != nil {
		logger.Fatal("failed to set migration dialect", zap.Error(err))
	}

	if err := goose.Up(db, postgresMigrationsPath()); err != nil {
		logger.Fatal("failed to migrate", zap.Error(err))
	}
}
