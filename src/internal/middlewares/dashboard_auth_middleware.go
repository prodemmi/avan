package middlewares

import (
	"avan/internal/services"
	"context"
	"net/http"
	"strconv"
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

		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}
		c.Set("user_id", userID)

		ctx := c.Request.Context()
		workspaces, err := workspaceService.GetWorkspaces(ctx, int64(userID))
		if err != nil || len(workspaces) == 0 {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		workspace, err := workspaceService.GetDefaultWorkspace(int64(userID))
		if err != nil {
			workspace = workspaces[0]
		}
		if cookie, cookieErr := c.Cookie("workspace_id"); cookieErr == nil {
			if selectedID, parseErr := strconv.ParseInt(cookie, 10, 64); parseErr == nil {
				for _, accessible := range workspaces {
					if accessible.ID == selectedID {
						workspace = accessible
						break
					}
				}
			}
		}
		ctx = context.WithValue(c.Request.Context(), "workspace", workspace)
		ctx = context.WithValue(ctx, "workspaces", workspaces)
		c.Set("workspace", workspace)
		c.Set("workspaces", workspaces)

		setting := func(key, fallback string) (string, error) {
			value, err := workspaceService.SettingService.GetSetting(workspace.ID, key)
			if err != nil {
				return "", err
			}
			if value.Value == nil || *value.Value == "" {
				return fallback, nil
			}
			return *value.Value, nil
		}
		primaryColor, err := setting("primary-color", "#780C28")
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx = context.WithValue(ctx, "primaryColor", primaryColor)
		c.Set("primaryColor", primaryColor)

		theme, err := setting("theme", "auto")
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx = context.WithValue(ctx, "theme", theme)
		c.Set("theme", theme)

		timezone, err := setting("timezone", "UTC")
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx = context.WithValue(ctx, "timezone", timezone)
		c.Set("timezone", timezone)

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
