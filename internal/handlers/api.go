package handlers

import (
	"avan/internal/db/sqlc"
	"avan/internal/dto"
	"avan/internal/notifiers"
	"avan/internal/services"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"time"
)

type ApiHandler struct {
	ApiService      services.ApiService
	MessageService  services.MessageService
	ProviderService services.ProviderService
	Notifier        notifiers.Notifier
}

type SendBody struct {
	ProviderID int64               `json:"provider_id" binding:"required"`
	TemplateID *int64              `json:"template_id"`
	Recipient  string              `json:"recipient" binding:"required"`
	Message    string              `json:"message"`
	ScheduleAt *time.Time          `json:"schedule_at"`
	Data       *[]dto.TemplateData `json:"data"`
	Payload    any                 `json:"payload"`
}

func (a ApiHandler) Send(c *gin.Context) {
	var body SendBody
	err := c.Bind(&body)
	if err != nil {
		SendError(c, err)
		return
	}
	hasTemplate := body.TemplateID != nil
	hasMessage := body.Message != ""

	if !hasTemplate && !hasMessage {
		err = errors.New("message is not set")
		SendError(c, err)
		return
	}

	message := body.Message
	workspaceID := GetWorkspaceIDFromContext(c)
	if hasTemplate {
		message, err = a.MessageService.BindMessageData(*body.TemplateID, body.Data)
		if err != nil {
			SendError(c, err)
			return
		}
	}

	var payload *string
	if body.Payload != nil {
		pay, err := json.Marshal(body.Payload)
		if err != nil {
			SendError(c, err)
			return
		}
		p := string(pay)
		payload = &p
	}

	if err := a.ApiService.Send(workspaceID, body.ProviderID, body.TemplateID, body.Recipient, message, body.ScheduleAt, payload); err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, "message added to the queue")
}

func (a ApiHandler) QuickSend(c *gin.Context) {
	var body SendBody
	err := c.Bind(&body)
	if err != nil {
		SendError(c, err)
		return
	}
	hasTemplate := body.TemplateID != nil
	hasMessage := body.Message != ""

	if !hasTemplate && !hasMessage {
		err = errors.New("message is not set")
		SendError(c, err)
		return
	}

	message := body.Message
	workspaceID := GetWorkspaceIDFromContext(c)
	if hasTemplate {
		message, err = a.MessageService.BindMessageData(*body.TemplateID, body.Data)
		if err != nil {
			SendError(c, err)
			return
		}
	}

	var payload *string
	if body.Payload != nil {
		pay, err := json.Marshal(body.Payload)
		if err != nil {
			SendError(c, err)
			return
		}
		p := string(pay)
		payload = &p
	}

	sentAt := time.Now()
	id, err := a.MessageService.CreateMessage(workspaceID, body.ProviderID, body.TemplateID, body.Recipient, message, 1, nil, body.ScheduleAt, &sentAt, payload)
	if err != nil {
		SendError(c, err)
		return
	}

	now := time.Now()
	msg := sqlc.Message{
		ID:          id,
		WorkspaceID: workspaceID,
		ProviderID:  body.ProviderID,
		TemplateID:  body.TemplateID,
		Recipient:   body.Recipient,
		Content:     message,
		Status:      1,
		ScheduleAt:  nil,
		SentAt:      now,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	provider, err := a.ProviderService.GetProvider(msg.ProviderID)
	if err != nil {
		SendError(c, err)
		return
	}

	commitMessage := func(err error) {
		if err == nil {
			_ = a.Notifier.MessageService.MarkMessageAsSuccess(msg.ID)
		} else {
			_ = a.Notifier.MessageService.MarkMessageAsError(msg.ID)
		}
	}

	retry := 3
	for retry > 0 {
		err = a.Notifier.SendMessage(provider, msg)
		commitMessage(err)
		if err == nil {
			SendSuccess(c, "message sent quickly")
			return
		}
		retry--
		time.Sleep(500 * time.Millisecond)
	}

	if err := a.ApiService.Send(workspaceID, body.ProviderID, body.TemplateID, body.Recipient, message, body.ScheduleAt, payload); err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, "message added to the queue")
}

func NewApiHandler(authService services.ApiService, messageService services.MessageService, notifier notifiers.Notifier, providerService services.ProviderService) ApiHandler {
	return ApiHandler{ApiService: authService, MessageService: messageService, Notifier: notifier, ProviderService: providerService}
}
