package router

import (
	"github.com/gin-gonic/gin"
)

func (r Router) GenerateSettingsRoutes(routerGroup *gin.RouterGroup) {
	settings := routerGroup.Group("/settings")
	{
		settings.GET("/", r.SettingHandler.Index)
		settings.POST("/", r.SettingHandler.Update)
	}
}
