package services

import (
	"avan/internal/db/sqlc"
	"avan/internal/redis"
	"avan/internal/repositories"
	"context"
	"fmt"
)

type SettingService struct {
	Repo  repositories.SettingRepository
	Cache redis.Cache[sqlc.Setting]

	UserService UserService
}

func (s SettingService) GetSettings(workspaceID int64) ([]sqlc.Setting, error) {
	ctx := context.Background()
	return s.Repo.GetSettings(ctx, workspaceID)
}

func (s SettingService) GetSetting(workspaceID int64, key string) (sqlc.Setting, error) {
	ctx := context.Background()
	return s.Repo.GetSetting(ctx, workspaceID, key)
}

func (s SettingService) CreateSettings(workspaceID int64, key string, value *string) error {
	ctx := context.Background()
	return s.Repo.CreateSettings(ctx, workspaceID, key, value)
}

func (s SettingService) UpdateSetting(workspaceID int64, key string, value *string) error {
	ctx := context.Background()
	return s.Repo.UpdateSetting(ctx, workspaceID, key, value)
}

func (s SettingService) CreateDefaultSettings(workspaceID int64) {
	ctx := context.Background()

	admin, _ := s.UserService.GetDefaultUser()

	theme := "auto"
	primaryColor := "#780C28"
	timezone := "UTC"
	adminEmail := admin.Email
	logDir := s.Repo.Base.Config.LogsDirectory
	apiSecurityPolicy := "{}"

	settings := []sqlc.Setting{
		sqlc.Setting{
			Key:   "theme",
			Value: &theme,
		},
		sqlc.Setting{
			Key:   "primary-color",
			Value: &primaryColor,
		},
		sqlc.Setting{
			Key:   "timezone",
			Value: &timezone,
		},
		sqlc.Setting{
			Key:   "admin-email",
			Value: &adminEmail,
		},
		sqlc.Setting{
			Key:   "log-dir",
			Value: &logDir,
		},
		sqlc.Setting{
			Key:   APISecurityPolicySettingKey,
			Value: &apiSecurityPolicy,
		},
	}
	fmt.Println("settings", settings)

	for _, setting := range settings {
		if _, err := s.GetSetting(workspaceID, setting.Key); err == nil {
			continue
		}
		_ = s.Repo.CreateSettings(ctx, workspaceID, setting.Key, setting.Value)
	}
}

func NewSettingService(repo repositories.SettingRepository, userService UserService) SettingService {
	cache := redis.NewCache[sqlc.Setting](repo.Base.Redis)
	return SettingService{Repo: repo, Cache: cache, UserService: userService}
}
