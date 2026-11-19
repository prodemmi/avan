package handlers

import (
	"avan/internal/services"
	"avan/views/api_keys"
	"avan/views/components"
	"avan/views/layout"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiKeyHandler struct {
	ApiKeyService    services.ApiKeyService
	WorkspaceService services.WorkspaceService
}

func (w ApiKeyHandler) workspaceID(c *gin.Context) (int64, error) {
	requested := c.Query("workspace-id")
	if requested == "" && c.Request.Method != http.MethodGet {
		requested = c.PostForm("workspace-id")
	}
	if requested == "" {
		return GetWorkspaceIDFromContext(c), nil
	}
	id, err := ConvertToInt64(requested)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid workspace")
	}
	workspaces, err := w.WorkspaceService.GetWorkspaces(c.Request.Context(), GetUserIDFromContext(c))
	if err != nil {
		return 0, err
	}
	for _, workspace := range workspaces {
		if workspace.ID == id {
			return id, nil
		}
	}
	return 0, fmt.Errorf("workspace not found")
}

func (w ApiKeyHandler) Index(c *gin.Context) {
	workspaceID, err := w.workspaceID(c)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	_apiKeys, err := w.ApiKeyService.GetApiKeys(workspaceID)
	if err != nil {
		c = ContextWithError(c, err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	page := components.PageData{}

	defaultContent := api_keys.Index(_apiKeys, page, workspaceID)
	var items []layout.MenuItem
	base := layout.Layout(defaultContent, "", items)
	err = base.Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (w ApiKeyHandler) Create(c *gin.Context) {
	workspaceID, err := w.workspaceID(c)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	defaultContent := api_keys.Create(workspaceID)

	base := layout.Layout(defaultContent, "", nil)
	err = base.Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (w ApiKeyHandler) Store(c *gin.Context) {
	var body struct {
		Name     string `form:"name"`
		ExpireAt int    `form:"expire_at"`
	}
	if err := c.Bind(&body); err != nil {
		c = ContextWithError(c, err.Error())
		w.Create(c)
		return
	}

	var expireAt *time.Time
	if body.ExpireAt != 6 {
		var duration time.Duration
		switch body.ExpireAt {
		case 0:
			duration = time.Hour * 1
		case 1:
			duration = time.Hour * 4
		case 2:
			duration = time.Hour * 24
		case 3:
			duration = time.Hour * 24 * 7
		case 4:
			duration = time.Hour * 24 * 30
		case 5:
			duration = time.Hour * 24 * 30 * 3
		}
		_expireAt := time.Now().Add(duration)
		expireAt = &_expireAt
	} else {
		expireAt = nil
	}

	apiKey, err := w.ApiKeyService.CreateToken()
	if err != nil {
		c = ContextWithError(c, err.Error())
		w.Create(c)
		return
	}

	workspaceID, err := w.workspaceID(c)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	ownerID := GetUserIDFromContext(c)

	_, err = w.ApiKeyService.CreateApiKey(workspaceID, ownerID, body.Name, apiKey, expireAt)
	if err != nil {
		c = ContextWithError(c, err.Error())
		w.Create(c)
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/dashboard/api-keys?workspace-id=%d", workspaceID))
}

func (w ApiKeyHandler) Delete(c *gin.Context) {
	apiKeyIDParam := c.Param("id")
	if apiKeyIDParam == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	apiKeyID, err := ConvertToInt64(apiKeyIDParam)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	workspaceID, err := w.workspaceID(c)
	if err != nil {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	err = w.ApiKeyService.DeleteApiKey(apiKeyID, workspaceID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/dashboard/api-keys?workspace-id=%d", workspaceID))
}

func NewApiKeyHandler(apiKeyService services.ApiKeyService, workspaceService services.WorkspaceService) ApiKeyHandler {
	return ApiKeyHandler{ApiKeyService: apiKeyService, WorkspaceService: workspaceService}
}
