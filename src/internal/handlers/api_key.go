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
	ApiKeyService services.ApiKeyService
}

func (w ApiKeyHandler) Index(c *gin.Context) {
	workspaceID := GetWorkspaceIDFromContext(c)
	_apiKeys, err := w.ApiKeyService.GetApiKeys(workspaceID)
	if err != nil {
		c = ContextWithError(c, err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	page := components.PageData{}

	defaultContent := api_keys.Index(_apiKeys, page)
	var items []layout.MenuItem
	base := layout.Layout(defaultContent, "", items)
	err = base.Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (w ApiKeyHandler) Create(c *gin.Context) {
	defaultContent := api_keys.Create()

	base := layout.Layout(defaultContent, "", nil)
	err := base.Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

var redirectApiKeyURL = "http://%s/dashboard/api-keys/"

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

	workspaceID := GetWorkspaceIDFromContext(c)
	ownerID := GetUserIDFromContext(c)

	_, err = w.ApiKeyService.CreateApiKey(workspaceID, ownerID, body.Name, apiKey, expireAt)
	if err != nil {
		c = ContextWithError(c, err.Error())
		w.Create(c)
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf(redirectApiKeyURL, c.Request.Host))
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
	err = w.ApiKeyService.DeleteApiKey(apiKeyID, GetWorkspaceIDFromContext(c))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Redirect(http.StatusFound, "/dashboard/api-keys")
}

func NewApiKeyHandler(apiKeyService services.ApiKeyService) ApiKeyHandler {
	return ApiKeyHandler{ApiKeyService: apiKeyService}
}
