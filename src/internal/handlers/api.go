package handlers

import (
	"avan/internal/db/sqlc"
	"avan/internal/dto"
	"avan/internal/notifiers"
	"avan/internal/services"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
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
	Retry      bool                `json:"retry"`
	ScheduleAt *time.Time          `json:"schedule_at"`
	Data       *[]dto.TemplateData `json:"data"`
	Payload    any                 `json:"payload"`
}

func (a ApiHandler) Send(c *gin.Context) {
	var body SendBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hasTemplate := body.TemplateID != nil
	hasMessage := body.Message != ""

	if !hasTemplate && !hasMessage {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "message or template_id is required"})
		return
	}

	message := body.Message
	workspaceID := GetWorkspaceIDFromContext(c)
	provider, err := a.ProviderService.GetProvider(body.ProviderID)
	if err != nil || provider.WorkspaceID != workspaceID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "provider not found in this workspace"})
		return
	}
	if hasTemplate {
		template, err := a.MessageService.TemplateService.GetTemplate(*body.TemplateID)
		if err != nil || template.WorkspaceID != workspaceID {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "template not found in this workspace"})
			return
		}
		message, err = a.MessageService.BindMessageData(*body.TemplateID, body.Data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	var payload *string
	if body.Payload != nil {
		pay, err := json.Marshal(body.Payload)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		p := string(pay)
		payload = &p
	}

	if err := a.ApiService.Send(workspaceID, body.ProviderID, body.TemplateID, body.Recipient, message, body.Retry, body.ScheduleAt, payload); err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, "message added to the queue")
}

func (a ApiHandler) QuickSend(c *gin.Context) {
	var body SendBody
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hasTemplate := body.TemplateID != nil
	hasMessage := body.Message != ""

	if !hasTemplate && !hasMessage {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "message or template_id is required"})
		return
	}

	message := body.Message
	workspaceID := GetWorkspaceIDFromContext(c)
	provider, err := a.ProviderService.GetProvider(body.ProviderID)
	if err != nil || provider.WorkspaceID != workspaceID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "provider not found in this workspace"})
		return
	}
	if hasTemplate {
		template, err := a.MessageService.TemplateService.GetTemplate(*body.TemplateID)
		if err != nil || template.WorkspaceID != workspaceID {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "template not found in this workspace"})
			return
		}
		message, err = a.MessageService.BindMessageData(*body.TemplateID, body.Data)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	var payload *string
	if body.Payload != nil {
		pay, err := json.Marshal(body.Payload)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		p := string(pay)
		payload = &p
	}
	if body.ScheduleAt != nil && body.ScheduleAt.After(time.Now()) {
		if err := a.ApiService.Send(workspaceID, body.ProviderID, body.TemplateID, body.Recipient, message, body.Retry, body.ScheduleAt, payload); err != nil {
			SendError(c, err)
			return
		}
		SendSuccess(c, "message scheduled")
		return
	}

	id, err := a.MessageService.CreateMessage(workspaceID, body.ProviderID, body.TemplateID, body.Recipient, message, 1, nil, body.Retry, nil, nil, payload)
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
		Retry:       body.Retry,
		ScheduleAt:  nil,
		Payload:     payload,
		SentAt:      nil,
		CreatedAt:   &now,
		UpdatedAt:   &now,
	}

	err = notifiers.SendWithRetry(func() error { return a.Notifier.SendMessage(provider, msg) }, body.Retry)
	if err == nil {
		if err := a.MessageService.MarkMessageAsSuccess(msg.ID); err != nil {
			SendError(c, err)
			return
		}
		SendSuccess(c, "message sent quickly")
		return
	}

	if !body.Retry {
		_ = a.MessageService.MarkMessageAsError(msg.ID)
		SendError(c, err)
		return
	}
	if err := a.MessageService.MarkMessageAsPending(msg.ID); err != nil {
		SendError(c, err)
		return
	}
	if err := a.MessageService.ProduceMessageDirectly(msg.ProviderID, msg.ID); err != nil {
		SendError(c, err)
		return
	}

	SendSuccess(c, "message added to the queue")
}

func NewApiHandler(authService services.ApiService, messageService services.MessageService, notifier notifiers.Notifier, providerService services.ProviderService) ApiHandler {
	return ApiHandler{ApiService: authService, MessageService: messageService, Notifier: notifier, ProviderService: providerService}
}
