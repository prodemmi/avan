package router

import (
	"avan/internal/handlers"
	"avan/internal/server"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Router struct {
	Server           server.Server
	Logger           zap.Logger
	AuthHandler      handlers.AuthHandler
	ApiHandler       handlers.ApiHandler
	ApiKeyHandler    handlers.ApiKeyHandler
	UserHandler      handlers.UserHandler
	ProviderHandler  handlers.ProviderHandler
	WorkspaceHandler handlers.WorkspaceHandler
	TemplateHandler  handlers.TemplateHandler
	MessageHandler   handlers.MessageHandler
	AnalysisHandler  handlers.AnalysisHandler
	SettingHandler   handlers.SettingHandler
}

func (r Router) GenerateRoutes() {
	r.GenerateDashboardRoutes()

	// redirect to dashboard
	r.Server.Engine.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/dashboard")
	})

	// Health check route
	r.Server.Engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Server.Engine.Static("/public", "./public")
}

func NewRouter(server server.Server, logger zap.Logger, authHandler handlers.AuthHandler, apiHandler handlers.ApiHandler, apiKeyHandler handlers.ApiKeyHandler, userHandler handlers.UserHandler, providerHandler handlers.ProviderHandler, workspaceHandler handlers.WorkspaceHandler, templateHandler handlers.TemplateHandler, messageHandler handlers.MessageHandler, analysisHandler handlers.AnalysisHandler, settingHandler handlers.SettingHandler) Router {
	return Router{Server: server, Logger: logger, AuthHandler: authHandler, ApiHandler: apiHandler, ApiKeyHandler: apiKeyHandler, UserHandler: userHandler, ProviderHandler: providerHandler, WorkspaceHandler: workspaceHandler, TemplateHandler: templateHandler, MessageHandler: messageHandler, AnalysisHandler: analysisHandler, SettingHandler: settingHandler}
}
