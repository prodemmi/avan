package handlers

import (
	"avan/internal/services"
	"avan/views/components"
	"avan/views/layout"
	"avan/views/providers"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func validateSMSForm(url *string, method *int, bodySchema *string, bodyEncoding, authType string, authUsername, accountSID, sender *string) error {
	if url == nil || strings.TrimSpace(*url) == "" || method == nil || (*method != 0 && *method != 1) {
		return fmt.Errorf("SMS URL and a supported HTTP method are required")
	}
	if bodyEncoding != "json" && bodyEncoding != "form" {
		return fmt.Errorf("SMS POST body encoding must be JSON or form-urlencoded")
	}
	if authType != "token" && authType != "basic" {
		return fmt.Errorf("SMS authentication type must be token or HTTP Basic")
	}
	if bodySchema != nil && strings.TrimSpace(*bodySchema) != "" {
		var fields map[string]json.RawMessage
		if err := json.Unmarshal([]byte(*bodySchema), &fields); err != nil || fields == nil {
			return fmt.Errorf("SMS body schema must be a JSON object")
		}
	}
	if authType == "basic" && (authUsername == nil || strings.TrimSpace(*authUsername) == "") {
		return fmt.Errorf("HTTP Basic username is required")
	}
	if bodySchema != nil {
		if strings.Contains(*bodySchema, "$username") && (authUsername == nil || strings.TrimSpace(*authUsername) == "") {
			return fmt.Errorf("SMS request username is required by the body template")
		}
		if strings.Contains(*bodySchema, "$sender") && (sender == nil || strings.TrimSpace(*sender) == "") {
			return fmt.Errorf("SMS sender number or ID is required by the body template")
		}
	}
	if strings.Contains(*url, "$account_sid") && (accountSID == nil || strings.TrimSpace(*accountSID) == "") {
		return fmt.Errorf("Account SID is required by the provider URL")
	}
	return nil
}

func optionalTrimmed(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

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
	_providers, err := w.ProviderService.GetProviders(GetWorkspaceIDFromContext(c), userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

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
	if typeParam != "sms" && typeParam != "email" && typeParam != "telegram" && typeParam != "slack" && typeParam != "webhook" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	mode := c.Query("mode")
	if mode == "" {
		mode = c.PostForm("mode")
	}
	if mode == "" {
		mode = "blank"
	}
	templateID := c.Query("template")
	if templateID == "" {
		templateID = c.PostForm("template")
	}
	var template providers.ProviderTemplate
	if templateID != "" {
		var ok bool
		template, ok = providers.FindProviderTemplate(templateID)
		if !ok || template.ConfigType != typeParam {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		mode = "template"
	}
	if mode != "blank" && mode != "template" || mode == "template" && templateID == "" && c.Request.Method != http.MethodGet {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	defaultContent := providers.Create(typeParam, mode, template)
	var items []layout.MenuItem
	base := layout.Layout(defaultContent, "", items)
	err := base.Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (w ProviderHandler) Store(c *gin.Context) {
	var body struct {
		Name     string `form:"name"`
		Mode     string `form:"mode"`
		Template string `form:"template"`
		Priority int64  `form:"priority"`
		Weight   int64  `form:"weight"`
		// Sms Input
		SmsURL            *string `form:"sms_url"`
		SmsMethod         *int    `form:"sms_method"`
		SmsAuthTokenKey   *string `form:"sms_auth_token_key"`
		SmsAuthTokenValue *string `form:"sms_auth_token_value"`
		SmsIsAuthInHeader *bool   `form:"sms_is_auth_in_header,default=false"`
		SmsBodySchema     *string `form:"sms_body_schema"`
		SmsBodyEncoding   string  `form:"sms_body_encoding"`
		SmsAuthType       string  `form:"sms_auth_type"`
		SmsAuthUsername   *string `form:"sms_auth_username"`
		SmsAccountSID     *string `form:"sms_account_sid"`
		SmsSender         *string `form:"sms_sender"`

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
	if err != nil {
		return
	}
	if body.SmsBodyEncoding == "" {
		body.SmsBodyEncoding = "json"
	}
	if body.SmsAuthType == "" {
		body.SmsAuthType = "token"
	}
	body.SmsAuthUsername = optionalTrimmed(body.SmsAuthUsername)
	body.SmsAccountSID = optionalTrimmed(body.SmsAccountSID)
	body.SmsSender = optionalTrimmed(body.SmsSender)
	if typeParam != "sms" && typeParam != "email" && typeParam != "telegram" && typeParam != "slack" && typeParam != "webhook" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if strings.TrimSpace(body.Name) == "" {
		sendError(fmt.Errorf("provider name is required"))
		return
	}
	if body.Priority < 0 || body.Priority > 2147483647 || body.Weight < 1 || body.Weight > 1000 {
		sendError(fmt.Errorf("failover priority must be non-negative and weight must be between 1 and 1000"))
		return
	}
	if body.Mode == "template" {
		template, ok := providers.FindProviderTemplate(body.Template)
		if !ok || template.ConfigType != typeParam {
			sendError(fmt.Errorf("choose a supported template for this channel"))
			return
		}
		if template.ID == "kavenegar" && (body.SmsAuthTokenValue == nil || strings.TrimSpace(*body.SmsAuthTokenValue) == "") {
			sendError(fmt.Errorf("Kavenegar API key is required"))
			return
		}
		if template.ID == "sendgrid" && (body.Password == nil || strings.TrimSpace(*body.Password) == "") {
			sendError(fmt.Errorf("SendGrid API key is required"))
			return
		}
		if template.ID == "melipayamak" && (body.SmsMethod == nil || *body.SmsMethod != 1 || body.SmsBodyEncoding != "form" || body.SmsAuthType != "token" || body.SmsAuthTokenKey == nil || strings.TrimSpace(*body.SmsAuthTokenKey) != "password" || body.SmsAuthTokenValue == nil || strings.TrimSpace(*body.SmsAuthTokenValue) == "" || body.SmsIsAuthInHeader != nil && *body.SmsIsAuthInHeader || body.SmsAuthUsername == nil || body.SmsSender == nil) {
			sendError(fmt.Errorf("Melipayamak requires form-urlencoded POST, API username, password, and sender number"))
			return
		}
		if template.ID == "twilio" && (body.SmsMethod == nil || *body.SmsMethod != 1 || body.SmsBodyEncoding != "form" || body.SmsAuthType != "basic" || body.SmsAuthUsername == nil || body.SmsAccountSID == nil || body.SmsSender == nil || body.SmsAuthTokenValue == nil || strings.TrimSpace(*body.SmsAuthTokenValue) == "" || body.SmsURL == nil || !strings.Contains(*body.SmsURL, "$account_sid")) {
			sendError(fmt.Errorf("Twilio requires form-urlencoded POST, Account SID, API Key SID and secret, and a sender number"))
			return
		}
	} else if body.Mode != "blank" || body.Template != "" {
		sendError(fmt.Errorf("invalid provider setup mode"))
		return
	}

	ownerID := GetUserIDFromContext(c)
	workspaceID := GetWorkspaceIDFromContext(c)
	var configErr error
	var configType int32

	switch typeParam {
	case "sms":
		if err := validateSMSForm(body.SmsURL, body.SmsMethod, body.SmsBodySchema, body.SmsBodyEncoding, body.SmsAuthType, body.SmsAuthUsername, body.SmsAccountSID, body.SmsSender); err != nil {
			sendError(err)
			return
		}
		if body.SmsAuthType == "basic" && (body.SmsAuthTokenValue == nil || strings.TrimSpace(*body.SmsAuthTokenValue) == "") {
			sendError(fmt.Errorf("HTTP Basic password is required"))
			return
		}
		configType = 0
	case "email":
		if body.SMTPHost == nil || strings.TrimSpace(*body.SMTPHost) == "" || body.SMTPPort == nil || *body.SMTPPort < 1 || body.FromEmail == nil || strings.TrimSpace(*body.FromEmail) == "" {
			sendError(fmt.Errorf("SMTP host, a valid port, and sender email are required"))
			return
		}
		configType = 1
	case "telegram":
		if body.BotToken == nil || strings.TrimSpace(*body.BotToken) == "" {
			sendError(fmt.Errorf("Telegram bot token is required"))
			return
		}
		configType = 2
	case "slack":
		if body.WebhookUrl == nil || strings.TrimSpace(*body.WebhookUrl) == "" {
			sendError(fmt.Errorf("Slack webhook URL is required"))
			return
		}
		configType = 3
	case "webhook":
		if body.WebhookURL == nil || strings.TrimSpace(*body.WebhookURL) == "" {
			sendError(fmt.Errorf("webhook URL is required"))
			return
		}
		configType = 4
	}

	provider, err := w.ProviderService.CreateProvider(workspaceID, ownerID, configType, strings.TrimSpace(body.Name), body.Priority, body.Weight)
	if err != nil {
		sendError(err)
		return
	}
	switch configType {
	case 0:
		configErr = w.ProviderService.CreateSmsConfig(workspaceID, provider.ID, ownerID, *body.SmsURL, int32(*body.SmsMethod), body.SmsAuthTokenKey, body.SmsAuthTokenValue, body.SmsIsAuthInHeader, body.SmsBodySchema, body.SmsBodyEncoding, body.SmsAuthType, body.SmsAuthUsername, body.SmsAccountSID, body.SmsSender)
	case 1:
		configErr = w.ProviderService.CreateEmailConfig(workspaceID, provider.ID, ownerID, *body.SMTPHost, *body.SMTPPort, *body.FromEmail, body.Username, body.Password)
	case 2:
		configErr = w.ProviderService.CreateTelegramConfig(workspaceID, provider.ID, ownerID, *body.BotToken)
	case 3:
		configErr = w.ProviderService.CreateSlackConfig(workspaceID, provider.ID, ownerID, *body.WebhookUrl)
	case 4:
		configErr = w.ProviderService.CreateWebhookConfig(workspaceID, provider.ID, ownerID, *body.WebhookURL)
	}
	if configErr != nil {
		_ = w.ProviderService.DeleteProvider(provider.ID)
		sendError(configErr)
		return
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
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if provider.WorkspaceID != GetWorkspaceIDFromContext(c) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

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
		Name     string `form:"name"`
		Priority int64  `form:"priority"`
		Weight   int64  `form:"weight"`
		// Sms Input
		SmsURL            *string `form:"sms_url"`
		SmsMethod         *int    `form:"sms_method"`
		SmsAuthTokenKey   *string `form:"sms_auth_token_key"`
		SmsAuthTokenValue *string `form:"sms_auth_token_value"`
		SmsIsAuthInHeader *bool   `form:"sms_is_auth_in_header,default=false"`
		SmsBodySchema     *string `form:"sms_body_schema"`
		SmsBodyEncoding   string  `form:"sms_body_encoding"`
		SmsAuthType       string  `form:"sms_auth_type"`
		SmsAuthUsername   *string `form:"sms_auth_username"`
		SmsAccountSID     *string `form:"sms_account_sid"`
		SmsSender         *string `form:"sms_sender"`

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
			w.Edit(c)
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
	if err != nil {
		return
	}
	if body.SmsBodyEncoding == "" {
		body.SmsBodyEncoding = "json"
	}
	if body.SmsAuthType == "" {
		body.SmsAuthType = "token"
	}
	body.SmsAuthUsername = optionalTrimmed(body.SmsAuthUsername)
	body.SmsAccountSID = optionalTrimmed(body.SmsAccountSID)
	body.SmsSender = optionalTrimmed(body.SmsSender)

	providerID, err := ConvertToInt64(providerIDParam)
	sendError(err)

	provider, err := w.ProviderService.GetProvider(providerID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if provider.WorkspaceID != GetWorkspaceIDFromContext(c) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if strings.TrimSpace(body.Name) == "" || body.Priority < 0 || body.Priority > 2147483647 || body.Weight < 1 || body.Weight > 1000 {
		c = ContextWithError(c, "Provider name is required; priority must be non-negative and weight must be between 1 and 1000")
		w.Edit(c)
		return
	}
	validationError := ""
	switch provider.ConfigType {
	case 0:
		if err := validateSMSForm(body.SmsURL, body.SmsMethod, body.SmsBodySchema, body.SmsBodyEncoding, body.SmsAuthType, body.SmsAuthUsername, body.SmsAccountSID, body.SmsSender); err != nil {
			validationError = err.Error()
		} else {
			currentSMS, configErr := w.ProviderService.GetSmsConfig(provider.ID)
			if configErr != nil {
				sendError(configErr)
				return
			}
			if currentSMS.AuthType != body.SmsAuthType && (body.SmsAuthTokenValue == nil || strings.TrimSpace(*body.SmsAuthTokenValue) == "") {
				validationError = "enter a new credential when changing the SMS authentication type"
			}
		}
	case 1:
		if body.SMTPHost == nil || strings.TrimSpace(*body.SMTPHost) == "" || body.SMTPPort == nil || *body.SMTPPort < 1 || body.FromEmail == nil || strings.TrimSpace(*body.FromEmail) == "" {
			validationError = "SMTP host, a valid port, and sender email are required"
		}
	case 2:
		if body.BotToken == nil || strings.TrimSpace(*body.BotToken) == "" {
			validationError = "Telegram bot token is required"
		}
	case 3:
		if body.SlackWebhookUrl == nil || strings.TrimSpace(*body.SlackWebhookUrl) == "" {
			validationError = "Slack webhook URL is required"
		}
	case 4:
		if body.WebhookURL == nil || strings.TrimSpace(*body.WebhookURL) == "" {
			validationError = "webhook URL is required"
		}
	}
	if validationError != "" {
		c = ContextWithError(c, validationError)
		w.Edit(c)
		return
	}

	err = w.ProviderService.UpdateProvider(providerID, strings.TrimSpace(body.Name), body.Priority, body.Weight)
	sendError(err)
	if err != nil {
		return
	}

	switch provider.ConfigType {
	case 0:
		err = w.ProviderService.UpdateSmsConfig(provider.ID, *body.SmsURL, int32(*body.SmsMethod), body.SmsAuthTokenKey, body.SmsAuthTokenValue, body.SmsIsAuthInHeader, body.SmsBodySchema, body.SmsBodyEncoding, body.SmsAuthType, body.SmsAuthUsername, body.SmsAccountSID, body.SmsSender)
		sendError(err)
		if err != nil {
			return
		}
		break
	case 1:
		err = w.ProviderService.UpdateEmailConfig(provider.ID, *body.SMTPHost, *body.SMTPPort, body.Username, body.Password, *body.FromEmail)
		sendError(err)
		if err != nil {
			return
		}
	case 2:
		err = w.ProviderService.UpdateTelegramConfig(provider.ID, *body.BotToken)
		sendError(err)
		if err != nil {
			return
		}
	case 3:
		err = w.ProviderService.UpdateSlackConfig(provider.ID, *body.SlackWebhookUrl)
		sendError(err)
		if err != nil {
			return
		}
	case 4:
		err = w.ProviderService.UpdateWebhookConfig(provider.ID, *body.WebhookURL)
		sendError(err)
		if err != nil {
			return
		}
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
	if err != nil {
		return
	}
	provider, err := w.ProviderService.GetProvider(id)
	if err != nil || provider.WorkspaceID != GetWorkspaceIDFromContext(c) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	err = w.ProviderService.DeleteProvider(id)
	sendIfError(err)

	c.Redirect(http.StatusFound, "/dashboard/providers")
}

func NewProviderHandler(providerService services.ProviderService, workspaceService services.WorkspaceService) ProviderHandler {
	return ProviderHandler{ProviderService: providerService, WorkspaceService: workspaceService}
}
