package handlers

import (
	"avan/internal/db/sqlc"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func GetUserIDFromContext(c *gin.Context) int64 {
	return int64(c.MustGet("user_id").(float64))
}

func GetWorkspaceFromContext(c *gin.Context) int64 {
	workspaceIDQuery, _ := c.GetQuery("workspace-id")
	workspaceID, _ := strconv.ParseInt(workspaceIDQuery, 10, 64)
	return workspaceID
}

func GetWorkspaceIDFromContext(c *gin.Context) int64 {
	_workspace, _ := c.Get("workspace")
	workspace := _workspace.(sqlc.Workspace)
	return workspace.ID
}

func ContextWithError(c *gin.Context, errors ...string) *gin.Context {
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "errors", errors))
	return c
}

func ContextWithMessage(c *gin.Context, messages ...string) *gin.Context {
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), "messages", messages))
	return c
}

func ToTitle(s string) string {
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func GetConfigTypeLabel(configType int64) string {
	switch configType {
	case 0:
		return "sms"
	case 1:
		return "email"
	case 2:
		return "telegram"
	case 3:
		return "slack"
	case 4:
		return "webhook"
	}
	return ""
}

func ConvertToInt64(input string) (int64, error) {
	i, err := strconv.ParseInt(input, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func ConvertToInt32(input string) (int32, error) {
	i, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}

func SendUnprocessable(c *gin.Context, err error) {
	//c.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"error": err.(string)})
}

func SendError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
}

func SendSuccess(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"data": data,
	})
}
