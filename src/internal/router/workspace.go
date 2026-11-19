package router

import (
	"avan/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func (r Router) GenerateWorkspaceRoutes(routerGroup *gin.RouterGroup) {
	workspaces := routerGroup.Group("/workspaces", middlewares.PermissionMiddleware(r.UserHandler.UserService, "workspaces.manage"))
	{
		workspaces.GET("/", r.WorkspaceHandler.Index)
		workspaces.GET("/create", r.WorkspaceHandler.Create)
		workspaces.GET("/:id/edit", r.WorkspaceHandler.Edit)

		workspaces.POST("/", r.WorkspaceHandler.Store)
		workspaces.POST("/:id/edit", r.WorkspaceHandler.Update)
		workspaces.POST("/:id/setDefault", r.WorkspaceHandler.SetDefault)
		workspaces.POST("/:id/delete", r.WorkspaceHandler.Delete)
	}
}
