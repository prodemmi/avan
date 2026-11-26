package router

import (
	"avan/internal/db/sqlc"
	"avan/internal/handlers"
	"avan/internal/middlewares"
	"avan/views"
	"avan/views/developer"
	"avan/views/layout"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r Router) GenerateDashboardRoutes() {
	authMiddleware := middlewares.AuthMiddleware(r.AuthHandler.AuthService, r.WorkspaceHandler.WorkspaceService)
	dashboard := r.Server.Engine.Group("/dashboard", authMiddleware)
	{
		dashboard.GET("/", func(c *gin.Context) {
			// مقدار پیش‌فرض برای محتوای صفحه‌ی داشبورد
			nodes, healthy, unhealthy := r.AnalysisHandler.AnalysisService.GetNodeAnalysis()
			consumers := r.AnalysisHandler.AnalysisService.AnalysisRepository.Base.Config.RedisConsumerCount
			defaultContent := views.Home(views.HomeStatus{Nodes: nodes, Healthy: healthy, Unhealthy: unhealthy, Consumers: consumers})
			var items []layout.MenuItem
			base := layout.Layout(defaultContent, "", items)
			err := base.Render(c.Request.Context(), c.Writer)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		})

	}
	developerRoutes := dashboard.Group("/developer", middlewares.PermissionMiddleware(r.UserHandler.UserService, "developer.read"))
	developerRoutes.GET("", r.developerIndex)
	developerRoutes.GET("/openapi.json", r.developerOpenAPI)
	dashboard.GET("/switch-workspace", func(c *gin.Context) {
		id, err := handlers.ConvertToInt64(c.Query("workspace-id"))
		workspaceValue, ok := c.Get("workspaces")
		workspaces, valid := workspaceValue.([]sqlc.Workspace)
		if err != nil || id <= 0 || !ok || !valid {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if !workspaceAllowed(workspaces, id) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("workspace_id", fmt.Sprint(id), 365*24*60*60, "/dashboard", "", c.Request.TLS != nil, true)
		c.Redirect(http.StatusSeeOther, "/dashboard/")
	})

	r.GenerateAuthRoutes(r.Server.Engine)
	r.GenerateApiRoutes(r.Server.Engine)
	r.GenerateApiKeyRoutes(dashboard)
	r.GenerateUserRoutes(dashboard)
	r.GenerateProviderRoutes(dashboard)
	r.GenerateWorkspaceRoutes(dashboard)
	r.GenerateTemplateRoutes(dashboard)
	r.GenerateMessageRoutes(dashboard)
	r.GenerateAnalysisRoutes(dashboard)
	r.GenerateSettingsRoutes(dashboard)
}

func workspaceAllowed(workspaces []sqlc.Workspace, id int64) bool {
	for _, workspace := range workspaces {
		if workspace.ID == id {
			return true
		}
	}
	return false
}

func (r Router) selectedDeveloperWorkspace(c *gin.Context) ([]sqlc.Workspace, sqlc.Workspace, error) {
	userID := handlers.GetUserIDFromContext(c)
	workspaces, err := r.WorkspaceHandler.WorkspaceService.GetWorkspaces(c.Request.Context(), userID)
	if err != nil {
		return nil, sqlc.Workspace{}, err
	}
	selectedID := handlers.GetWorkspaceIDFromContext(c)
	if requested := c.Query("workspace-id"); requested != "" {
		selectedID, err = handlers.ConvertToInt64(requested)
		if err != nil {
			return nil, sqlc.Workspace{}, err
		}
	}
	for _, workspace := range workspaces {
		if workspace.ID == selectedID {
			return workspaces, workspace, nil
		}
	}
	return nil, sqlc.Workspace{}, errors.New("workspace not found")
}

func (r Router) selectedDeveloperProvider(c *gin.Context, workspace sqlc.Workspace) ([]sqlc.Provider, sqlc.Provider, error) {
	providers, err := r.ProviderHandler.ProviderService.GetProviders(workspace.ID, handlers.GetUserIDFromContext(c))
	if err != nil {
		return nil, sqlc.Provider{}, err
	}
	if len(providers) == 0 {
		return providers, sqlc.Provider{}, nil
	}
	selected := providers[0]
	if requested := c.Query("provider-id"); requested != "" {
		id, err := handlers.ConvertToInt64(requested)
		if err != nil {
			return nil, sqlc.Provider{}, err
		}
		for _, provider := range providers {
			if provider.ID == id {
				return providers, provider, nil
			}
		}
		return nil, sqlc.Provider{}, errors.New("provider not found")
	}
	return providers, selected, nil
}

func (r Router) developerIndex(c *gin.Context) {
	workspaces, selected, err := r.selectedDeveloperWorkspace(c)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	providers, selectedProvider, err := r.selectedDeveloperProvider(c, selected)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	base := layout.Layout(developer.Index(workspaces, selected, providers, selectedProvider), "", nil)
	if err := base.Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (r Router) developerOpenAPI(c *gin.Context) {
	_, selected, err := r.selectedDeveloperWorkspace(c)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	_, selectedProvider, err := r.selectedDeveloperProvider(c, selected)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	providerID := selectedProvider.ID
	requestSchema := gin.H{
		"type":        "object",
		"description": "Provide either message, or template_id with template data.",
		"required":    []string{"provider_id", "recipient"},
		"properties": gin.H{
			"provider_id": gin.H{"type": "integer", "format": "int64", "example": providerID, "description": "Use the selected provider ID from this workspace."},
			"template_id": gin.H{"type": "integer", "format": "int64"},
			"recipient":   gin.H{"type": "string", "example": "+989120000000"},
			"message":     gin.H{"type": "string", "example": "Hello"},
			"retry":       gin.H{"type": "boolean", "default": false, "description": "Retry provider delivery up to three times when true."},
			"schedule_at": gin.H{"type": "string", "format": "date-time", "description": "Optional RFC 3339 send time; scheduled messages are queued when due."},
			"data":        gin.H{"type": "array", "items": gin.H{"type": "object", "required": []string{"key", "value"}, "properties": gin.H{"key": gin.H{"type": "string"}, "value": gin.H{"type": "string"}}}},
			"payload":     gin.H{"type": "object", "additionalProperties": true},
		},
	}
	response := func(description string, schema gin.H, example any) gin.H {
		return gin.H{"description": description, "content": gin.H{"application/json": gin.H{"schema": schema, "example": example}}}
	}
	success := gin.H{"$ref": "#/components/schemas/Success"}
	failure := gin.H{"$ref": "#/components/schemas/Error"}
	messageOperation := func(quick bool) gin.H {
		statusExample := "pending"
		if quick {
			statusExample = "sent"
		}
		return gin.H{
			"security":    []gin.H{{"ApiKey": []string{}}},
			"parameters":  []gin.H{{"name": "Idempotency-Key", "in": "header", "required": true, "schema": gin.H{"type": "string", "maxLength": 255}, "description": "Unique per logical send within this workspace. Reuse it only for retries of the same request."}},
			"requestBody": gin.H{"required": true, "content": gin.H{"application/json": gin.H{"schema": requestSchema, "example": gin.H{"provider_id": providerID, "recipient": "+989120000000", "message": "Hello", "retry": true, "schedule_at": "2026-12-01T09:30:00Z", "payload": gin.H{"campaign": "welcome"}}}}},
			"responses": gin.H{
				"200": response("Send completed. The data includes message_id, status, and duplicate.", success, gin.H{"data": gin.H{"message_id": 123, "status": statusExample, "duplicate": false}}),
				"202": response("Message is pending or scheduled. Duplicate requests return the existing state and message ID without sending again.", success, gin.H{"data": gin.H{"message_id": 123, "status": "pending", "duplicate": true}}),
				"400": response("Malformed JSON, missing fields/message, invalid template data, or missing Idempotency-Key.", failure, gin.H{"error": "a valid Idempotency-Key header is required"}),
				"401": response("Missing, invalid, or expired API key.", failure, gin.H{"error": "invalid or expired API key"}),
				"403": response("Workspace or provider IP/API-key policy denied the request.", failure, gin.H{"error": "request rejected by security policy"}),
				"409": response("The Idempotency-Key was already used with a different request.", failure, gin.H{"error": "idempotency key was already used with a different request"}),
				"502": response("The existing idempotent message has failed; it will not be sent again by replaying this key.", success, gin.H{"data": gin.H{"message_id": 123, "status": "failed", "duplicate": true}}),
				"500": response("Provider, database, or queue failure.", failure, gin.H{"error": "description"}),
			},
		}
	}
	inquiryOperation := gin.H{
		"security":   []gin.H{{"ApiKey": []string{}}},
		"parameters": []gin.H{{"name": "id", "in": "path", "required": true, "schema": gin.H{"type": "integer", "format": "int64"}}},
		"responses": gin.H{
			"200": response("Workspace-scoped message lifecycle status, including provider delivery/read receipts.", success, gin.H{"data": gin.H{"id": 123, "status": "delivered", "created_at": "2026-12-01T09:30:00Z", "sent_at": "2026-12-01T09:30:02Z", "updated_at": "2026-12-01T09:30:04Z"}}),
			"401": response("Missing, invalid, or expired API key.", failure, gin.H{"error": "invalid or expired API key"}),
			"403": response("Workspace or provider IP/API-key policy denied the request.", failure, gin.H{"error": "request rejected by security policy"}),
			"404": response("Message does not exist in the API key's workspace.", failure, gin.H{"error": "message not found"}),
		},
	}
	statusOperation := gin.H{
		"security":   []gin.H{{"ApiKey": []string{}}},
		"parameters": []gin.H{{"name": "id", "in": "path", "required": true, "schema": gin.H{"type": "integer", "format": "int64"}}},
		"requestBody": gin.H{"required": true, "content": gin.H{"application/json": gin.H{
			"schema":  gin.H{"type": "object", "required": []string{"status"}, "properties": gin.H{"status": gin.H{"type": "string", "enum": []string{"sent", "delivered", "failed", "read"}}}},
			"example": gin.H{"status": "delivered"},
		}}},
		"responses": gin.H{
			"200": response("Message status recorded.", success, gin.H{"data": gin.H{"id": 123, "status": "delivered"}}),
			"400": response("Invalid message ID or status.", failure, gin.H{"error": "status must be sent, delivered, failed, or read"}),
			"401": response("Missing, invalid, or expired API key.", failure, gin.H{"error": "invalid or expired API key"}),
			"403": response("Workspace or provider security policy denied the request.", failure, gin.H{"error": "request rejected by security policy"}),
			"404": response("Message does not exist in the API key's workspace.", failure, gin.H{"error": "message not found"}),
			"409": response("The reported event is not a valid lifecycle transition.", failure, gin.H{"error": "message status transition is not allowed"}),
		},
	}
	c.JSON(http.StatusOK, gin.H{
		"openapi": "3.0.3",
		"info":    gin.H{"title": "Avan API", "version": "1.0.0", "description": fmt.Sprintf("API reference for workspace %s using provider %s (ID %d).", selected.Name, selectedProvider.Name, selectedProvider.ID)},
		"servers": []gin.H{{"url": "/"}},
		"paths": gin.H{
			"/api/send":                 gin.H{"post": messageOperation(false)},
			"/api/quick-send":           gin.H{"post": messageOperation(true)},
			"/api/messages/{id}":        gin.H{"get": inquiryOperation},
			"/api/messages/{id}/status": gin.H{"post": statusOperation},
			"/health":                   gin.H{"get": gin.H{"summary": "Service health", "responses": gin.H{"200": response("Service is running.", gin.H{"type": "object", "properties": gin.H{"status": gin.H{"type": "string"}}}, gin.H{"status": "ok"})}}},
		},
		"components": gin.H{
			"securitySchemes": gin.H{"ApiKey": gin.H{"type": "apiKey", "in": "header", "name": "X-API-KEY"}},
			"schemas": gin.H{
				"Success": gin.H{"type": "object", "required": []string{"data"}, "properties": gin.H{"data": gin.H{"type": "object", "additionalProperties": true}}},
				"Error":   gin.H{"type": "object", "required": []string{"error"}, "properties": gin.H{"error": gin.H{"type": "string"}}},
			},
		},
	})
}
