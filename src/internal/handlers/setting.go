package handlers

import (
	"avan/internal/services"
	"avan/views/layout"
	"avan/views/settings"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type SettingHandler struct {
	SettingService services.SettingService
}

func (w SettingHandler) Index(c *gin.Context) {
	workspaceID := GetWorkspaceIDFromContext(c)
	list, err := w.SettingService.GetSettings(workspaceID)

	if err != nil {
		c = ContextWithError(c, err.Error())
	}

	defaultContent := settings.Index(list)
	var items []layout.MenuItem
	base := layout.Layout(defaultContent, "", items)
	err = base.Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (w SettingHandler) Update(c *gin.Context) {
	var body struct {
		Key   string  `form:"key"`
		Value *string `form:"value"`
	}

	sendError := func(err error) {
		if err != nil {
			c = ContextWithError(c, err.Error())
			w.Index(c)
		}
	}

	if err := c.Bind(&body); err != nil {
		sendError(err)
		return
	}
	if body.Key == services.APISecurityPolicySettingKey {
		value := "{}"
		if body.Value != nil {
			value = *body.Value
		}
		if _, err := services.ParseSecurityPolicy(value); err != nil {
			c = ContextWithError(c, err.Error())
			w.Index(c)
			return
		}
	}

	workspaceID := GetWorkspaceIDFromContext(c)
	err := w.SettingService.UpdateSetting(workspaceID, body.Key, body.Value)
	if err != nil {
		sendError(err)
		return
	}

	time.Sleep(1 * time.Second)

	message := "setting updated"
	c = ContextWithMessage(c, message)

	c.Redirect(http.StatusFound, c.Request.URL.String())
}

func NewSettingHandler(templateService services.SettingService) SettingHandler {
	return SettingHandler{SettingService: templateService}
}
