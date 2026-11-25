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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing API key"})
			return
		}

		apiKey, err := apiKeyService.ValidateApiKey(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired API key"})
			return
		}

		workspace, err := workspaceService.GetWorkspace(apiKey.WorkspaceID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired API key"})
			return
		}
		policy, err := workspaceService.GetAPISecurityPolicy(workspace.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{"error": "workspace security policy is invalid or unavailable"})
			return
		}
		if !policy.Allows(c.ClientIP(), apiKey.ID) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "request rejected by workspace security policy"})
			return
		}

		ctx := context.WithValue(c.Request.Context(), "workspace", workspace)
		c.Set("workspace", workspace)
		c.Set("api_key", apiKey)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
