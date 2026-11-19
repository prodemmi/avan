package router

import (
	"avan/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func (r Router) GenerateSettingsRoutes(routerGroup *gin.RouterGroup) {
	settings := routerGroup.Group("/settings", middlewares.PermissionMiddleware(r.UserHandler.UserService, "settings.manage"))
	{
		settings.GET("/", r.SettingHandler.Index)
		settings.POST("/", r.SettingHandler.Update)
	}
}
