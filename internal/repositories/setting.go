package repositories

import (
	"avan/internal/db/sqlc"
	"context"
	"database/sql"
)

type SettingRepository struct {
	Base Repository
}

func (w SettingRepository) GetSettings(ctx context.Context, workspaceID int64) ([]sqlc.Setting, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.GetSettingsMysql(ctx, workspaceID)
	case "postgres":
		return w.Base.DB.Queries.GetSettingsPostgresql(ctx, workspaceID)
	case "sqlite":
		return w.Base.DB.Queries.GetSettingsSqlite(ctx, workspaceID)
	}
	return nil, nil
}

func (w SettingRepository) GetSetting(ctx context.Context, workspaceID int64, key string) (sqlc.Setting, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.GetSettingMysql(ctx, sqlc.GetSettingMysqlParams{
			WorkspaceID: workspaceID,
			Key:         key,
		})
	case "postgres":
		return w.Base.DB.Queries.GetSettingPostgresql(ctx, sqlc.GetSettingPostgresqlParams{
			WorkspaceID: workspaceID,
			Key:         key,
		})
	case "sqlite":
		return w.Base.DB.Queries.GetSettingSqlite(ctx, sqlc.GetSettingSqliteParams{
			WorkspaceID: workspaceID,
			Key:         key,
		})
	}
	return sqlc.Setting{}, nil
}

func (w SettingRepository) CreateSettings(ctx context.Context, workspaceID int64, key string, value *string) error {
	var nullValue sql.NullString

	if value != nil {
		nullValue = sql.NullString{String: *value, Valid: true}
	} else {
		nullValue = sql.NullString{Valid: false}
	}

	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.CreateSettingMysql(ctx, sqlc.CreateSettingMysqlParams{
			WorkspaceID: workspaceID,
			Key:         key,
			Value:       nullValue,
		})
	case "postgres":
		return w.Base.DB.Queries.CreateSettingPostgresql(ctx, sqlc.CreateSettingPostgresqlParams{
			WorkspaceID: workspaceID,
			Key:         key,
			Value:       nullValue,
		})
	case "sqlite":
		return w.Base.DB.Queries.CreateSettingSqlite(ctx, sqlc.CreateSettingSqliteParams{
			WorkspaceID: workspaceID,
			Key:         key,
			Value:       value,
		})
	}
	return nil
}

func (w SettingRepository) UpdateSetting(ctx context.Context, workspaceID int64, key string, value *string) error {
	var nullValue sql.NullString
	if value != nil {
		nullValue = sql.NullString{String: *value, Valid: true}
	} else {
		nullValue = sql.NullString{Valid: false}
	}

	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.UpdateSettingMysql(ctx, sqlc.UpdateSettingMysqlParams{
			WorkspaceID: workspaceID,
			Value:       nullValue,
			Key:         key,
		})
	case "postgres":
		return w.Base.DB.Queries.UpdateSettingPostgresql(ctx, sqlc.UpdateSettingPostgresqlParams{
			WorkspaceID: workspaceID,
			Key:         key,
			Value:       nullValue,
		})
	case "sqlite":
		return w.Base.DB.Queries.UpdateSettingSqlite(ctx, sqlc.UpdateSettingSqliteParams{
			WorkspaceID: workspaceID,
			Key:         key,
			Value:       value,
		})
	}
	return nil
}

func NewSettingRepository(base Repository) SettingRepository {
	return SettingRepository{Base: base}
}
