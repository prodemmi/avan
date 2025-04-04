package handlers

import (
	"avan/internal/services"
	"avan/views/components"
	"avan/views/layout"
	"avan/views/providers"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProviderHandler struct {
	ProviderService  services.ProviderService
	WorkspaceService services.WorkspaceService
}

func (w ProviderHandler) Index(c *gin.Context) {
	sendError := func(err error) {
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	userID := GetUserIDFromContext(c)
	workspace, err := w.WorkspaceService.GetDefaultWorkspace(userID)
	sendError(err)

	_providers, _ := w.ProviderService.GetProviders(workspace.ID, userID)

	currentPageStr := c.Query("page")
	currentPage, _ := ConvertToInt32(currentPageStr)
	page := components.PageData{
		Route:       "providers",
		CurrentPage: currentPage,
		MaxResults:  int32(len(_providers)),
		Bucket:      10,
	}

	defaultContent := providers.Index(_providers, page)
	var items []layout.MenuItem
	base := layout.Layout(defaultContent, "", items)
	err = base.Render(c.Request.Context(), c.Writer)
	sendError(err)
}

func (w ProviderHandler) Create(c *gin.Context) {
	typeParam := c.Param("type")
	if typeParam == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	defaultContent := providers.Create(typeParam)
	var items []layout.MenuItem
	base := layout.Layout(defaultContent, "", items)
	err := base.Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (w ProviderHandler) Store(c *gin.Context) {
	var body struct {
		Name string `form:"name"`
		// Sms Input
		SmsURL            *string `form:"sms_url"`
		SmsMethod         *int    `form:"sms_method"`
		SmsAuthTokenKey   *string `form:"sms_auth_token_key"`
		SmsAuthTokenValue *string `form:"sms_auth_token_value"`
		SmsIsAuthInHeader *bool   `form:"sms_is_auth_in_header,default=false"`

		// Email Input
		SMTPHost  *string `form:"smtp_host"`
		SMTPPort  *int64  `form:"smtp_port"`
		Username  *string `form:"username"`
		Password  *string `form:"password"`
		FromEmail *string `form:"from_email"`

		// Telegram Input
		BotToken *string `form:"bot_token"`

		// Slack Input
		WebhookUrl *string `form:"webhook_url"`

		// Webhook Input
		WebhookURL *string `form:"webhook_url"`
	}

	sendError := func(err error) {
		if err != nil {
			c = ContextWithError(c, err.Error())
			w.Create(c)
			c.Abort()
			return
		}
	}

	typeParam := c.Param("type")
	if typeParam == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err := c.Bind(&body)
	sendError(err)

	ownerID := GetUserIDFromContext(c)
	workspace, err := w.WorkspaceService.GetDefaultWorkspace(ownerID)
	sendError(err)

	switch typeParam {
	case "sms":
		provider, err := w.ProviderService.CreateProvider(workspace.ID, ownerID, 0, body.Name)
		sendError(err)
		err = w.ProviderService.CreateSmsConfig(workspace.ID, provider.ID, ownerID, *body.SmsURL, int32(*body.SmsMethod), body.SmsAuthTokenKey, body.SmsAuthTokenValue, body.SmsIsAuthInHeader)
		sendError(err)
		break
	case "email":
		provider, err := w.ProviderService.CreateProvider(workspace.ID, ownerID, 1, body.Name)
		sendError(err)
		err = w.ProviderService.CreateEmailConfig(workspace.ID, provider.ID, ownerID, *body.SMTPHost, *body.SMTPPort, *body.FromEmail, body.Username, body.Password)
		sendError(err)
	case "telegram":
		provider, err := w.ProviderService.CreateProvider(workspace.ID, ownerID, 2, body.Name)
		sendError(err)
		err = w.ProviderService.CreateTelegramConfig(workspace.ID, provider.ID, ownerID, *body.BotToken)
		sendError(err)
	case "slack":
		provider, err := w.ProviderService.CreateProvider(workspace.ID, ownerID, 3, body.Name)
		sendError(err)
		err = w.ProviderService.CreateSlackConfig(workspace.ID, provider.ID, ownerID, *body.WebhookUrl)
		sendError(err)
	case "webhook":
		provider, err := w.ProviderService.CreateProvider(workspace.ID, ownerID, 4, body.Name)
		sendError(err)
		err = w.ProviderService.CreateWebhookConfig(workspace.ID, provider.ID, ownerID, *body.WebhookURL)
		sendError(err)
	}

	c.Redirect(http.StatusFound, "/dashboard/providers")
}

func (w ProviderHandler) Edit(c *gin.Context) {
	sendIfError := func(err error) {
		if err != nil {
			c = ContextWithError(c)
			w.Index(c)
			c.Abort()
			return
		}
	}

	providerIDParam := c.Param("id")
	if providerIDParam == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	providerID, err := ConvertToInt64(providerIDParam)
	sendIfError(err)

	provider, err := w.ProviderService.GetProvider(providerID)
	sendIfError(err)

	var config interface{}
	var configErr error
	var configType string
	switch provider.ConfigType {
	case 0:
		configType = "sms"
		config, configErr = w.ProviderService.GetSmsConfig(provider.ID)
		break
	case 1:
		configType = "email"
		config, configErr = w.ProviderService.GetEmailConfig(provider.ID)
		break
	case 2:
		configType = "telegram"
		config, configErr = w.ProviderService.GetTelegramConfig(provider.ID)
		break
	case 3:
		configType = "slack"
		config, configErr = w.ProviderService.GetSlackConfig(provider.ID)
		break
	case 4:
		configType = "webhook"
		config, configErr = w.ProviderService.GetWebhookConfig(provider.ID)
		break
	}
	sendIfError(configErr)

	defaultContent := providers.Edit(provider, configType, config)
	var items []layout.MenuItem
	base := layout.Layout(defaultContent, "", items)
	err = base.Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (w ProviderHandler) Update(c *gin.Context) {
	var body struct {
		Name string `form:"name"`
		// Sms Input
		SmsURL            *string `form:"sms_url"`
		SmsMethod         *int    `form:"sms_method"`
		SmsAuthTokenKey   *string `form:"sms_auth_token_key"`
		SmsAuthTokenValue *string `form:"sms_auth_token_value"`
		SmsIsAuthInHeader *bool   `form:"sms_is_auth_in_header,default=false"`

		// Email Input
		SMTPHost  *string `form:"smtp_host"`
		SMTPPort  *int64  `form:"smtp_port"`
		Username  *string `form:"username"`
		Password  *string `form:"password"`
		FromEmail *string `form:"from_email"`

		// Telegram Input
		BotToken *string `form:"bot_token"`

		// Slack Input
		SlackWebhookUrl *string `form:"slack_webhook_url"`

		// Webhook Input
		WebhookURL    *string `form:"webhook_url"`
		WebhookMethod *int    `form:"webhook_method"`
	}

	sendError := func(err error) {
		if err != nil {
			c = ContextWithError(c, err.Error())
			w.Create(c)
			c.Abort()
			return
		}
	}

	providerIDParam := c.Param("id")
	if providerIDParam == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err := c.Bind(&body)
	sendError(err)

	providerID, err := ConvertToInt64(providerIDParam)
	sendError(err)

	provider, err := w.ProviderService.GetProvider(providerID)
	sendError(err)

	err = w.ProviderService.UpdateProvider(providerID, body.Name)
	sendError(err)

	switch provider.ConfigType {
	case 0:
		err = w.ProviderService.UpdateSmsConfig(provider.ID, *body.SmsURL, int32(*body.SmsMethod), body.SmsAuthTokenKey, body.SmsAuthTokenValue, body.SmsIsAuthInHeader)
		sendError(err)
		break
	case 1:
		err = w.ProviderService.UpdateEmailConfig(provider.ID, *body.SMTPHost, *body.SMTPPort, body.Username, body.Password, *body.FromEmail)
		sendError(err)
	case 2:
		err = w.ProviderService.UpdateTelegramConfig(provider.ID, *body.BotToken)
		sendError(err)
	case 3:
		err = w.ProviderService.UpdateSlackConfig(provider.ID, *body.SlackWebhookUrl)
		sendError(err)
	case 4:
		err = w.ProviderService.UpdateWebhookConfig(provider.ID, *body.WebhookURL)
		sendError(err)
	}
	message := fmt.Sprintf("%s provider updated", ToTitle(GetConfigTypeLabel(provider.ConfigType)))
	c = ContextWithMessage(c, message)

	w.Edit(c)
}

func (w ProviderHandler) Delete(c *gin.Context) {
	sendIfError := func(err error) {
		if err != nil {
			c = ContextWithError(c)
			w.Index(c)
			c.Abort()
			return
		}
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	id, err := ConvertToInt64(idParam)
	sendIfError(err)

	err = w.ProviderService.DeleteProvider(id)
	sendIfError(err)

	c.Redirect(http.StatusFound, "/dashboard/providers")
}

func NewProviderHandler(providerService services.ProviderService, workspaceService services.WorkspaceService) ProviderHandler {
	return ProviderHandler{ProviderService: providerService, WorkspaceService: workspaceService}
}
