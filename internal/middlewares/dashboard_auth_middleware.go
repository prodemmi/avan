package middlewares

import (
	"avan/internal/services"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is the middleware to check JWT validity
func AuthMiddleware(authService services.AuthService, workspaceService services.WorkspaceService) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, _ := c.Cookie("auth_token")
		if tokenString == "" {
			c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, claims, err := authService.ParseJWT(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Store claims in the context
		userID := claims["user_id"].(float64)
		c.Set("user_id", userID)

		sendErr := func(err error) {
			if err != nil {
				c.JSON(http.StatusInternalServerError, nil)
				c.Abort()
				return
			}
		}

		ctx := c.Request.Context()
		workspace, err := workspaceService.GetDefaultWorkspace(int64(userID))
		sendErr(err)
		ctx = context.WithValue(c.Request.Context(), "workspace", workspace)
		c.Set("workspace", workspace)

		primaryColor, err := workspaceService.GetPrimaryColor(int64(userID))
		sendErr(err)

		ctx = context.WithValue(ctx, "primaryColor", primaryColor)
		c.Set("primaryColor", primaryColor)

		theme, err := workspaceService.GetTheme(int64(userID))
		sendErr(err)

		ctx = context.WithValue(ctx, "theme", theme)
		c.Set("theme", theme)

		timezone, err := workspaceService.GetTimezone(int64(userID))
		sendErr(err)

		ctx = context.WithValue(ctx, "timezone", timezone)
		c.Set("timezone", timezone)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
