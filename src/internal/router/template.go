package router

import (
	"avan/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func (r Router) GenerateTemplateRoutes(routerGroup *gin.RouterGroup) {
	workspaces := routerGroup.Group("/templates", middlewares.PermissionMiddleware(r.UserHandler.UserService, "templates.manage"))
	{
		workspaces.GET("/", r.TemplateHandler.Index)
		workspaces.GET("/create", r.TemplateHandler.Create)
		workspaces.GET("/:id/edit", r.TemplateHandler.Edit)

		workspaces.POST("/", r.TemplateHandler.Store)
		workspaces.POST("/:id/edit", r.TemplateHandler.Update)
		workspaces.POST("/:id/delete", r.TemplateHandler.Delete)
	}
}
