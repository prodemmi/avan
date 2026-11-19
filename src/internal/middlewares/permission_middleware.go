package middlewares

import (
	"avan/internal/db/sqlc"
	"avan/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PermissionMiddleware(userService services.UserService, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}
		workspaceValue, ok := c.Get("workspace")
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}
		workspace, ok := workspaceValue.(sqlc.Workspace)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}
		var id int64
		switch value := userID.(type) {
		case float64:
			id = int64(value)
		case int64:
			id = value
		default:
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}
		allowed, err := userService.HasPermission(c.Request.Context(), workspace.ID, id, permission)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to check permissions"})
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permission denied"})
			return
		}
		c.Next()
	}
}
