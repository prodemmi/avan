package router

import (
	"avan/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func (r Router) GenerateProviderRoutes(group *gin.RouterGroup) {
	providers := group.Group("/providers", middlewares.PermissionMiddleware(r.UserHandler.UserService, "providers.manage"))
	{
		providers.GET("/", r.ProviderHandler.Index)
		providers.GET("/create/:type", r.ProviderHandler.Create)
		providers.GET("/:id/edit", r.ProviderHandler.Edit)

		providers.POST("/:id/edit", r.ProviderHandler.Update)
		providers.POST("/create/:type", r.ProviderHandler.Store)
		providers.POST("/:id/delete", r.ProviderHandler.Delete)
	}
}
