package router

import "github.com/gin-gonic/gin"

func (r Router) GenerateUserRoutes(group *gin.RouterGroup) {
	group.GET("/users", r.UserHandler.Index)
}
