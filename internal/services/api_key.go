package services

import (
	"avan/internal/db/sqlc"
	"avan/internal/repositories"
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type ApiKeyService struct {
	ApiKeyRepository repositories.ApiKeyRepository
}

func (a ApiKeyService) CreateToken() (string, error) {
	input := uuid.New().String()
	hash, err := bcrypt.GenerateFromPassword([]byte(input), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (a ApiKeyService) GetApiKeys(workspaceID, ownerID int64) ([]sqlc.ApiKey, error) {
	ctx := context.Background()
	return a.ApiKeyRepository.FindAllApiKeys(ctx, workspaceID, ownerID)
}

func (a ApiKeyService) CreateApiKey(workspaceID, ownerID int64, name, apiKey string, expireAt *time.Time) (bool, error) {
	ctx := context.Background()
	return a.ApiKeyRepository.CreateApiKey(ctx, workspaceID, ownerID, name, apiKey, expireAt)
}

func (a ApiKeyService) FindApiKeyByToken(token string) (sqlc.ApiKey, error) {
	ctx := context.Background()
	return a.ApiKeyRepository.FindApiKeyByToken(ctx, token)
}

func (a ApiKeyService) DeleteApiKey(id int64) error {
	ctx := context.Background()
	return a.ApiKeyRepository.DeleteApiKey(ctx, id)
}

func (a ApiKeyService) ValidateApiKey(apiKey string) bool {
	entity, err := a.FindApiKeyByToken(apiKey)
	if err != nil {
		return false
	}
	if entity.ExpireAt == nil {
		return true
	}

	return entity.ExpireAt.Before(time.Now())
}

func NewApiKeyService(apiKeyRepository repositories.ApiKeyRepository) ApiKeyService {
	return ApiKeyService{ApiKeyRepository: apiKeyRepository}
}
