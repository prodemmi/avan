package repositories

import (
	"avan/internal/config"
	"avan/internal/db"
	"avan/internal/redis"
	"context"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	Config config.Config
	DB     db.Database
	Redis  redis.Redis
	Logger zap.Logger
}

func (r Repository) GetServiceURL() string {
	domain := r.Config.ServerDomain
	port := r.Config.ServerPort
	return fmt.Sprintf("http://%s:%d", domain, port)
}

func (r Repository) HashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		r.Logger.Error("failed to hash password", zap.String("event", "hash_password"))
	}
	return string(hashedPassword)
}

func (r Repository) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (r Repository) IsConfiguredAdmin(ctx context.Context, userID int64) (bool, error) {
	if r.Config.AdminEmail == "" {
		return false, nil
	}
	var isAdmin bool
	err := r.DB.DB.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM users WHERE id = $1 AND LOWER(email) = LOWER($2))`, userID, r.Config.AdminEmail).Scan(&isAdmin)
	return isAdmin, err
}

func NewRepository(config config.Config, DB db.Database, redis redis.Redis, logger zap.Logger) Repository {
	return Repository{Config: config, DB: DB, Redis: redis, Logger: logger}
}
