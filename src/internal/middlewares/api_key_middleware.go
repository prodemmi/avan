package middlewares

import (
	"avan/internal/services"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ApiKeyMiddleware is the middleware to check api key validity
func ApiKeyMiddleware(apiKeyService services.ApiKeyService, workspaceService services.WorkspaceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("X-API-KEY")
		if tokenString == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		apiKey, err := apiKeyService.ValidateApiKey(tokenString)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		workspace, err := workspaceService.GetWorkspace(apiKey.WorkspaceID)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(c.Request.Context(), "workspace", workspace)
		c.Set("workspace", workspace)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
