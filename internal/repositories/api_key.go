package repositories

import (
	"avan/internal/db/sqlc"
	"context"
	"database/sql"
	"time"
)

type ApiKeyRepository struct {
	Base Repository
}

func (w ApiKeyRepository) FindAllApiKeys(ctx context.Context, workspaceID, ownerID int64) ([]sqlc.ApiKey, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.GetApiKeysMysql(ctx, sqlc.GetApiKeysMysqlParams{
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
		})
	case "postgres":
		return w.Base.DB.Queries.GetApiKeysPostgresql(ctx, sqlc.GetApiKeysPostgresqlParams{
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
		})
	case "sqlite":
		return w.Base.DB.Queries.GetApiKeysSqlite(ctx, sqlc.GetApiKeysSqliteParams{
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
		})
	}
	return nil, nil
}

func (w ApiKeyRepository) CreateApiKey(ctx context.Context, workspaceID, ownerID int64, name, apiKey string, expireAt *time.Time) (bool, error) {
	var expireAtTimeValue *time.Time
	validateExpireAtTime := expireAt != nil

	if validateExpireAtTime {
		expireAtTimeValue = expireAt
	} else {
		expireAtTimeValue = new(time.Time)
	}

	var err error
	switch w.Base.GetDBDialector() {
	case "mysql":
		err = w.Base.DB.Queries.CreateAPIKeyMysql(ctx, sqlc.CreateAPIKeyMysqlParams{
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
			Name:        name,
			ApiKey:      apiKey,
			ExpireAt: sql.NullTime{
				Time:  *expireAtTimeValue,
				Valid: validateExpireAtTime,
			},
		})
		break
	case "postgres":
		err = w.Base.DB.Queries.CreateAPIKeyPostgresql(ctx, sqlc.CreateAPIKeyPostgresqlParams{
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
			Name:        name,
			ApiKey:      apiKey,
			ExpireAt: sql.NullTime{
				Time:  *expireAtTimeValue,
				Valid: validateExpireAtTime,
			},
		})
		break
	case "sqlite":
		err = w.Base.DB.Queries.CreateAPIKeySqlite(ctx, sqlc.CreateAPIKeySqliteParams{
			WorkspaceID: workspaceID,
			OwnerID:     ownerID,
			Name:        name,
			ApiKey:      apiKey,
			ExpireAt:    expireAtTimeValue,
		})
		break
	}

	if err != nil {
		return false, err
	}
	return true, nil
}

func (w ApiKeyRepository) FindApiKeyByToken(ctx context.Context, apiKey string) (sqlc.ApiKey, error) {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.FindAPIKeyMysql(ctx, apiKey)
	case "postgres":
		return w.Base.DB.Queries.FindAPIKeyPostgresql(ctx, apiKey)
	case "sqlite":
		return w.Base.DB.Queries.FindAPIKeySqlite(ctx, apiKey)
	}
	return sqlc.ApiKey{}, nil
}

func (w ApiKeyRepository) DeleteApiKey(ctx context.Context, id int64) error {
	switch w.Base.GetDBDialector() {
	case "mysql":
		return w.Base.DB.Queries.DeleteAPIKeyMysql(ctx, id)
	case "postgres":
		return w.Base.DB.Queries.DeleteAPIKeyPostgresql(ctx, int32(id))
	case "sqlite":
		return w.Base.DB.Queries.DeleteAPIKeySqlite(ctx, id)
	}
	return nil
}

func NewApiKeyRepository(base Repository) ApiKeyRepository {
	return ApiKeyRepository{Base: base}
}
