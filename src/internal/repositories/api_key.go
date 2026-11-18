package repositories

import (
	"avan/internal/db/sqlc"
	"context"
	"time"
)

type ApiKeyRepository struct {
	Base Repository
}

func (w ApiKeyRepository) FindAllApiKeys(ctx context.Context, workspaceID, ownerID int64) ([]sqlc.ApiKey, error) {
	return w.Base.DB.Queries.GetApiKeysPostgresql(ctx, sqlc.GetApiKeysPostgresqlParams{
		WorkspaceID: workspaceID,
		OwnerID:     ownerID,
	})
}

func (w ApiKeyRepository) CreateApiKey(ctx context.Context, workspaceID, ownerID int64, name, apiKey string, expireAt *time.Time) (bool, error) {
	if err := w.Base.DB.Queries.CreateAPIKeyPostgresql(ctx, sqlc.CreateAPIKeyPostgresqlParams{
		WorkspaceID: workspaceID,
		OwnerID:     ownerID,
		Name:        name,
		ApiKey:      apiKey,
		ExpireAt:    expireAt,
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (w ApiKeyRepository) FindApiKeyByToken(ctx context.Context, apiKey string) (sqlc.ApiKey, error) {
	return w.Base.DB.Queries.FindAPIKeyPostgresql(ctx, apiKey)
}

func (w ApiKeyRepository) DeleteApiKey(ctx context.Context, id int64) error {
	return w.Base.DB.Queries.DeleteAPIKeyPostgresql(ctx, id)
}

func NewApiKeyRepository(base Repository) ApiKeyRepository {
	return ApiKeyRepository{Base: base}
}
