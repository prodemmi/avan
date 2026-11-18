package repositories

import (
	"avan/internal/db/sqlc"
	"context"
)

type SettingRepository struct {
	Base Repository
}

func (w SettingRepository) GetSettings(ctx context.Context, workspaceID int64) ([]sqlc.Setting, error) {
	return w.Base.DB.Queries.GetSettingsPostgresql(ctx, workspaceID)
}

func (w SettingRepository) GetSetting(ctx context.Context, workspaceID int64, key string) (sqlc.Setting, error) {
	return w.Base.DB.Queries.GetSettingPostgresql(ctx, sqlc.GetSettingPostgresqlParams{
		WorkspaceID: workspaceID,
		Key:         key,
	})
}

func (w SettingRepository) CreateSettings(ctx context.Context, workspaceID int64, key string, value *string) error {
	return w.Base.DB.Queries.CreateSettingPostgresql(ctx, sqlc.CreateSettingPostgresqlParams{
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       value,
	})
}

func (w SettingRepository) UpdateSetting(ctx context.Context, workspaceID int64, key string, value *string) error {
	return w.Base.DB.Queries.UpdateSettingPostgresql(ctx, sqlc.UpdateSettingPostgresqlParams{
		WorkspaceID: workspaceID,
		Key:         key,
		Value:       value,
	})
}

func NewSettingRepository(base Repository) SettingRepository {
	return SettingRepository{Base: base}
}
