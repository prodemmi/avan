package router

import (
	"github.com/gin-gonic/gin"
)

func (r Router) GenerateApiKeyRoutes(engine *gin.RouterGroup) {
	apiKeys := engine.Group("/api-keys")
	{
		apiKeys.GET("/", r.ApiKeyHandler.Index)
		apiKeys.GET("/create", r.ApiKeyHandler.Create)

		apiKeys.POST("/", r.ApiKeyHandler.Store)
		apiKeys.POST("/:id/delete", r.ApiKeyHandler.Delete)
	}
}
