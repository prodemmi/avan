package db

import (
	"avan/internal/config"
	"avan/internal/db/sqlc"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
	"io"
	"log"
)

type Database struct {
	Config  config.Config
	DB      *sql.DB
	Queries *sqlc.Queries
	Logger  zap.Logger
}

func getDSN(config config.Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		config.DatabaseUser, config.DatabasePassword,
		config.DatabaseHost, config.DatabasePort, config.DatabaseName,
		config.DatabaseSSLMode)
}

func NewDatabase(config config.Config, logger zap.Logger) Database {
	dsn := getDSN(config)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Fatal("failed to open the database", zap.Error(err))
	}

	if err := db.Ping(); err != nil {
		logger.Fatal("database connection failed", zap.Error(err))
	}

	goose.SetLogger(log.New(io.Discard, "", 0))
	MigratePostgres(logger, db)

	return Database{
		Config:  config,
		DB:      db,
		Queries: sqlc.New(db),
		Logger:  logger,
	}
}
