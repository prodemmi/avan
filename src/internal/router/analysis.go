package router

import (
	"avan/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func (r Router) GenerateAnalysisRoutes(routerGroup *gin.RouterGroup) {
	analysis := routerGroup.Group("/analysis", middlewares.PermissionMiddleware(r.UserHandler.UserService, "analysis.read"))
	{
		analysis.GET("/", r.AnalysisHandler.Index)
	}
}
