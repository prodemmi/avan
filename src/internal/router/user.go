package router

import "github.com/gin-gonic/gin"

import "avan/internal/middlewares"

func (r Router) GenerateUserRoutes(group *gin.RouterGroup) {
	users := group.Group("/users", middlewares.PermissionMiddleware(r.UserHandler.UserService, "users.manage"))
	users.GET("", r.UserHandler.Index)
	users.POST("/roles", r.UserHandler.SaveRole)
	users.POST("/roles/:id", r.UserHandler.SaveRole)
	users.POST("/roles/:id/delete", r.UserHandler.DeleteRole)
	users.POST("/:id/role", r.UserHandler.AssignRole)
	users.POST("/:id/workspaces", r.UserHandler.AssignWorkspaces)
	users.POST("/:id/edit", r.UserHandler.Update)
	users.POST("/:id/delete", r.UserHandler.Delete)
	users.POST("/", r.UserHandler.Store)
}
