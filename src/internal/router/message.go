package router

import (
	"avan/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func (r Router) GenerateMessageRoutes(group *gin.RouterGroup) {
	messages := group.Group("/messages", middlewares.PermissionMiddleware(r.UserHandler.UserService, "messages.read"))
	{
		messages.GET("/", r.MessageHandler.Index)
	}
}
