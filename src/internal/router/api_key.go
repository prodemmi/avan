package router

import (
	"avan/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func (r Router) GenerateApiKeyRoutes(engine *gin.RouterGroup) {
	apiKeys := engine.Group("/api-keys", middlewares.PermissionMiddleware(r.UserHandler.UserService, "api_keys.manage"))
	{
		apiKeys.GET("/", r.ApiKeyHandler.Index)
		apiKeys.GET("/create", r.ApiKeyHandler.Create)

		apiKeys.POST("/", r.ApiKeyHandler.Store)
		apiKeys.POST("/:id/delete", r.ApiKeyHandler.Delete)
	}
}
