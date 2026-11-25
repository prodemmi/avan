package services

import (
	"avan/internal/db/sqlc"
	"go.uber.org/zap"
	"time"
)

type ApiService struct {
	Logger         zap.Logger
	MessageService MessageService
}

func (a ApiService) Send(workspaceID, providerID int64, templateID *int64, recipient, content string, retry bool, scheduleAt *time.Time, payload *string, idempotencyKey, requestHash string) (sqlc.Message, bool, error) {
	return a.MessageService.ProduceMessage(workspaceID, providerID, templateID, recipient, content, retry, scheduleAt, payload, idempotencyKey, requestHash)
}

func (a ApiService) QuickSend() error {
	return nil
}

func NewApiService(logger zap.Logger, messageService MessageService) ApiService {
	return ApiService{Logger: logger, MessageService: messageService}
}
