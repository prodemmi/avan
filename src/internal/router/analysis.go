package router

import "github.com/gin-gonic/gin"

func (r Router) GenerateAnalysisRoutes(routerGroup *gin.RouterGroup) {
	analysis := routerGroup.Group("/analysis")
	{
		analysis.GET("/", r.AnalysisHandler.Index)
	}
}
