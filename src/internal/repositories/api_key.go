package repositories

import (
	"avan/internal/db/sqlc"
	"context"
	"time"
)

type ApiKeyRepository struct {
	Base Repository
}

func (w ApiKeyRepository) FindAllApiKeys(ctx context.Context, workspaceID int64) ([]sqlc.ApiKey, error) {
	return w.Base.DB.Queries.GetApiKeysPostgresql(ctx, workspaceID)
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

func (w ApiKeyRepository) DeleteApiKey(ctx context.Context, id, workspaceID int64) error {
	return w.Base.DB.Queries.DeleteAPIKeyPostgresql(ctx, sqlc.DeleteAPIKeyPostgresqlParams{
		ID:          id,
		WorkspaceID: workspaceID,
	})
}

func NewApiKeyRepository(base Repository) ApiKeyRepository {
	return ApiKeyRepository{Base: base}
}
