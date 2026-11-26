package handlers

import (
	"avan/internal/db/sqlc"
	"avan/internal/dto"
	"avan/internal/notifiers"
	"avan/internal/services"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
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
	idempotencyKey, ok := requestIdempotencyKey(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "a valid Idempotency-Key header is required"})
		return
	}
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
	if !a.providerPolicyAllows(c, provider.ID) {
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

	requestHash := sendRequestHash(body, message, payload)
	result, created, err := a.ApiService.Send(workspaceID, body.ProviderID, body.TemplateID, body.Recipient, message, body.Retry, body.ScheduleAt, payload, idempotencyKey, requestHash)
	if errors.Is(err, services.ErrIdempotencyKeyConflict) {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		SendError(c, err)
		return
	}
	writeSendResult(c, result, !created)
}

func (a ApiHandler) QuickSend(c *gin.Context) {
	idempotencyKey, ok := requestIdempotencyKey(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "a valid Idempotency-Key header is required"})
		return
	}
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
	if !a.providerPolicyAllows(c, provider.ID) {
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
	requestHash := sendRequestHash(body, message, payload)
	if body.ScheduleAt != nil && body.ScheduleAt.After(time.Now()) {
		result, created, err := a.ApiService.Send(workspaceID, body.ProviderID, body.TemplateID, body.Recipient, message, body.Retry, body.ScheduleAt, payload, idempotencyKey, requestHash)
		if errors.Is(err, services.ErrIdempotencyKeyConflict) {
			c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if err != nil {
			SendError(c, err)
			return
		}
		writeSendResult(c, result, !created)
		return
	}

	msg, created, err := a.MessageService.CreateIdempotentMessage(workspaceID, body.ProviderID, body.TemplateID, body.Recipient, message, body.Retry, nil, payload, idempotencyKey, requestHash)
	if errors.Is(err, services.ErrIdempotencyKeyConflict) {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		SendError(c, err)
		return
	}
	if !created {
		writeSendResult(c, msg, true)
		return
	}

	err = a.Notifier.SendWithFailover(provider, msg, body.Retry)
	if err == nil {
		if err := a.MessageService.MarkMessageAsSuccess(msg.ID); err != nil {
			SendError(c, err)
			return
		}
		msg.Status = 0
		sentAt := time.Now()
		msg.SentAt = &sentAt
		writeSendResult(c, msg, false)
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
		_ = a.MessageService.MarkMessageAsError(msg.ID)
		SendError(c, err)
		return
	}

	writeSendResult(c, msg, false)
}

func (a ApiHandler) Inquiry(c *gin.Context) {
	messageID, err := ConvertToInt64(c.Param("id"))
	if err != nil || messageID <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
		return
	}
	message, err := a.MessageService.GetMessageForWorkspace(messageID, GetWorkspaceIDFromContext(c))
	if errors.Is(err, sql.ErrNoRows) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "message not found"})
		return
	}
	if err != nil {
		SendError(c, err)
		return
	}
	if !a.providerPolicyAllows(c, message.ProviderID) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{
		"id":           message.ID,
		"status":       messageStatusName(message.Status, message.ScheduleAt),
		"created_at":   message.CreatedAt,
		"scheduled_at": message.ScheduleAt,
		"sent_at":      message.SentAt,
		"updated_at":   message.UpdatedAt,
	}})
}

func (a ApiHandler) ReportStatus(c *gin.Context) {
	messageID, err := ConvertToInt64(c.Param("id"))
	if err != nil || messageID <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid message ID"})
		return
	}
	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}
	status, ok := reportedMessageStatus(body.Status)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "status must be sent, delivered, failed, or read"})
		return
	}
	workspaceID := GetWorkspaceIDFromContext(c)
	message, err := a.MessageService.GetMessageForWorkspace(messageID, workspaceID)
	if errors.Is(err, sql.ErrNoRows) {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "message not found"})
		return
	}
	if err != nil {
		SendError(c, err)
		return
	}
	if !a.providerPolicyAllows(c, message.ProviderID) {
		return
	}
	if err := a.MessageService.ReportMessageStatus(messageID, workspaceID, status); errors.Is(err, services.ErrInvalidMessageStatusTransition) {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	} else if err != nil {
		SendError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"id": messageID, "status": messageStatusName(status, nil)}})
}

func (a ApiHandler) providerPolicyAllows(c *gin.Context, providerID int64) bool {
	policy, err := a.ProviderService.GetProviderSecurityPolicy(providerID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "provider security policy is invalid or unavailable"})
		return false
	}
	value, exists := c.Get("api_key")
	apiKey, ok := value.(sqlc.ApiKey)
	if !exists || !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired API key"})
		return false
	}
	if !policy.Allows(c.ClientIP(), apiKey.ID) {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "request rejected by provider security policy"})
		return false
	}
	return true
}

func reportedMessageStatus(value string) (int64, bool) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "sent":
		return 0, true
	case "failed":
		return 2, true
	case "delivered":
		return 3, true
	case "read":
		return 4, true
	default:
		return 0, false
	}
}

func requestIdempotencyKey(c *gin.Context) (string, bool) {
	key := strings.TrimSpace(c.GetHeader("Idempotency-Key"))
	return key, key != "" && len(key) <= 255
}

func sendRequestHash(body SendBody, content string, payload *string) string {
	request, _ := json.Marshal(struct {
		ProviderID int64      `json:"provider_id"`
		TemplateID *int64     `json:"template_id"`
		Recipient  string     `json:"recipient"`
		Message    string     `json:"message"`
		Retry      bool       `json:"retry"`
		ScheduleAt *time.Time `json:"schedule_at"`
		Payload    *string    `json:"payload"`
	}{body.ProviderID, body.TemplateID, body.Recipient, content, body.Retry, body.ScheduleAt, payload})
	hash := sha256.Sum256(request)
	return hex.EncodeToString(hash[:])
}

func writeSendResult(c *gin.Context, message sqlc.Message, duplicate bool) {
	status := messageStatusName(message.Status, message.ScheduleAt)
	statusCode := http.StatusOK
	if status == "pending" || status == "scheduled" {
		statusCode = http.StatusAccepted
	} else if status == "failed" {
		statusCode = http.StatusBadGateway
	}
	c.JSON(statusCode, gin.H{"data": gin.H{
		"message_id":   message.ID,
		"status":       status,
		"duplicate":    duplicate,
		"scheduled_at": message.ScheduleAt,
		"sent_at":      message.SentAt,
	}})
}

func messageStatusName(status int64, scheduleAt *time.Time) string {
	switch status {
	case 0:
		return "sent"
	case 1:
		if scheduleAt != nil {
			return "scheduled"
		}
		return "pending"
	case 2:
		return "failed"
	case 3:
		return "delivered"
	case 4:
		return "read"
	default:
		return "unknown"
	}
}

func NewApiHandler(authService services.ApiService, messageService services.MessageService, notifier notifiers.Notifier, providerService services.ProviderService) ApiHandler {
	return ApiHandler{ApiService: authService, MessageService: messageService, Notifier: notifier, ProviderService: providerService}
}
