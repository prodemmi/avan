package handlers

import (
	"avan/internal/services"
	"avan/views/components"
	"avan/views/layout"
	"avan/views/messages"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	MessageService  services.MessageService
	TemplateService services.TemplateService
	ProviderService services.ProviderService
}

type ListMessagesFilters struct {
	Page      int32   `form:"page"`
	Provider  *int64  `form:"provider"`
	Template  *int64  `form:"template"`
	Recipient *string `form:"recipient"`
	Content   *string `form:"content"`
	Status    *int16  `form:"status"`
}

func (w MessageHandler) Index(c *gin.Context) {
	sendError := func(err error) {
		if err != nil {
			c = ContextWithError(c, err.Error())
		}
	}

	var filters ListMessagesFilters
	err := c.ShouldBindQuery(&filters)
	sendError(err)

	workspaceID := GetWorkspaceIDFromContext(c)
	ownerID := GetUserIDFromContext(c)

	_messages, _ := w.MessageService.GetMessages(filters.Page, workspaceID, filters.Provider, filters.Template, filters.Recipient, filters.Content, filters.Status)
	_templates, _ := w.TemplateService.GetTemplates(workspaceID, ownerID)
	_providers, _ := w.ProviderService.GetProviders(workspaceID, ownerID)

	_messageCount, _ := w.MessageService.MessageCount(workspaceID, filters.Provider, filters.Template, filters.Recipient, filters.Content, filters.Status)
	bucket := int32(10)

	page := components.PageData{
		Route:       "messages",
		CurrentPage: filters.Page,
		MaxResults:  int32(_messageCount),
		Bucket:      bucket,
	}

	defaultContent := messages.Index(_messages, _templates, _providers, page)
	var items []layout.MenuItem
	base := layout.Layout(defaultContent, "", items)
	err = base.Render(c.Request.Context(), c.Writer)
	sendError(err)
}

func NewMessageHandler(messageService services.MessageService, templateService services.TemplateService, providerService services.ProviderService) MessageHandler {
	return MessageHandler{MessageService: messageService, TemplateService: templateService, ProviderService: providerService}
}
