package services

import (
	"go.uber.org/zap"
	"time"
)

type ApiService struct {
	Logger         zap.Logger
	MessageService MessageService
}

func (a ApiService) Send(workspaceID, providerID int64, templateID *int64, recipient, content string, scheduleAt *time.Time, payload *string) error {
	return a.MessageService.ProduceMessage(workspaceID, providerID, templateID, recipient, content, scheduleAt, payload)
}

func (a ApiService) QuickSend() error {
	return nil
}

func NewApiService(logger zap.Logger, messageService MessageService) ApiService {
	return ApiService{Logger: logger, MessageService: messageService}
}
