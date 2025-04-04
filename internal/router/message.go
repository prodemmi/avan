package router

import (
	"github.com/gin-gonic/gin"
)

func (r Router) GenerateMessageRoutes(group *gin.RouterGroup) {
	messages := group.Group("/messages")
	{
		messages.GET("/", r.MessageHandler.Index)
	}
}
