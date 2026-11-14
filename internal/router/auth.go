package router

import "github.com/gin-gonic/gin"

func (r Router) GenerateAuthRoutes(engine *gin.Engine) {
	auth := engine.Group("/auth")
	{
		auth.GET("/login", r.AuthHandler.LoginPage)
		auth.POST("/login", r.AuthHandler.Login)

		auth.GET("/logout", r.AuthHandler.Logout)
	}
}
