package db

import (
	"avan/internal/config"
	"avan/internal/db/sqlc"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
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

func getDSN(config config.Config, dbType string) string {
	switch dbType {
	case "sqlite3":
		return fmt.Sprintf("file:%s?_fk=1", config.DatabaseName)
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
			config.DatabaseUser, config.DatabasePassword,
			config.DatabaseHost, config.DatabasePort, config.DatabaseName)
	case "postgres":
		return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			config.DatabaseUser, config.DatabasePassword,
			config.DatabaseHost, config.DatabasePort, config.DatabaseName,
			config.DatabaseSSLMode)
	default:
		return ""
	}
}

func NewDatabase(config config.Config, logger zap.Logger) Database {
	dbType := config.DatabaseType
	if dbType == "" {
		logger.Fatal("DB_TYPE environment variable is not set")
	}

	dsn := getDSN(config, dbType)
	if dsn == "" {
		logger.Fatal("unsupported database type", zap.String("db_type", dbType))
	}

	db, err := sql.Open(dbType, dsn)
	if err != nil {
		logger.Fatal("failed to open the database", zap.Error(err))
	}

	if err := db.Ping(); err != nil {
		logger.Fatal("database connection failed", zap.Error(err))
	}

	var queries *sqlc.Queries

	goose.SetLogger(log.New(io.Discard, "", 0))

	switch dbType {
	case "mysql":
		queries = sqlc.New(db)
		MigrateMysql(logger, db)
		break
	case "postgres":
		queries = sqlc.New(db)
		MigratePostgres(logger, db)
		break
	case "sqlite3":
		queries = sqlc.New(db)
		MigrateSqlite(logger, db)
		break
	default:
		logger.Fatal("unsupported database type")
	}

	return Database{
		Config:  config,
		DB:      db,
		Queries: queries,
		Logger:  logger,
	}
}
