package router

import (
	"avan/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func (r Router) GenerateApiRoutes(engine *gin.Engine) {
	apiKeyMiddleware := middlewares.ApiKeyMiddleware(r.ApiKeyHandler.ApiKeyService, r.WorkspaceHandler.WorkspaceService)

	api := engine.Group("/api", apiKeyMiddleware)
	{
		api.POST("/send", r.ApiHandler.Send)
		api.POST("/quick-send", r.ApiHandler.QuickSend)
		api.GET("/messages/:id", r.ApiHandler.Inquiry)
		api.POST("/messages/:id/status", r.ApiHandler.ReportStatus)
	}
}
