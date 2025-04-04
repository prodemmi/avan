package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/golang-migrate/migrate/database/sqlite3"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"os"
)

func getMigrationFilename(filename string) string {
	path, _ := os.Getwd()
	return fmt.Sprintf("%s/internal/db/migrations/%s/", path, filename)
}

func MigrateMysql(logger zap.Logger, db *sql.DB) {
	if err := goose.SetDialect("mysql"); err != nil {
		logger.Fatal("failed to set migration dialect", zap.Error(err))
	}

	if err := goose.Up(db, getMigrationFilename("mysql")); err != nil {
		logger.Fatal("failed to migrate", zap.Error(err))
	}
}

func MigratePostgres(logger zap.Logger, db *sql.DB) {
	driver, _ := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		getMigrationFilename("postgresql"),
		"postgres", driver)
	if err != nil {
		logger.Fatal("preparing migration failed", zap.Error(err))
	}

	err = m.Up()
	if err != nil {
		logger.Fatal("migration failed", zap.Error(err))
	}
}

func MigrateSqlite(logger zap.Logger, db *sql.DB) {
	driver, _ := sqlite3.WithInstance(db, &sqlite3.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		getMigrationFilename("sqlite"),
		"sqlite", driver)
	if err != nil {
		logger.Fatal("preparing migration failed", zap.Error(err))
	}

	err = m.Up()
	if err != nil {
		logger.Fatal("migration failed", zap.Error(err))
	}
}
