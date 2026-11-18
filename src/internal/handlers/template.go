package handlers

import (
	"avan/internal/services"
	"avan/views/components"
	"avan/views/layout"
	"avan/views/templates"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TemplateHandler struct {
	TemplateService services.TemplateService
}

func (w TemplateHandler) Index(c *gin.Context) {
	userID := GetUserIDFromContext(c)
	workspaceID := GetWorkspaceIDFromContext(c)
	_templates, err := w.TemplateService.GetTemplates(workspaceID, userID)
	if err != nil {
		c = ContextWithError(c, err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	currentPageStr := c.Query("page")
	currentPage, _ := ConvertToInt32(currentPageStr)
	page := components.PageData{
		Route:       "templates",
		CurrentPage: currentPage,
		MaxResults:  int32(len(_templates)),
		Bucket:      10,
	}

	defaultContent := templates.Index(_templates, page)
	var items []layout.MenuItem
	base := layout.Layout(defaultContent, "", items)
	err = base.Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (w TemplateHandler) Create(c *gin.Context) {
	defaultContent := templates.Create()

	base := layout.Layout(defaultContent, "", nil)
	err := base.Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (w TemplateHandler) Edit(c *gin.Context) {
	sendError := func(err error) {
		if err != nil {
			c = ContextWithError(c, err.Error())
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	templateIDParam := c.Param("id")
	if templateIDParam == "" {
		c.Redirect(http.StatusNotFound, "/dashboard/template")
		return
	}

	// Convert the template ID from string to int64
	templateID, err := ConvertToInt64(templateIDParam)
	sendError(err)

	// Fetch the template from the service
	template, err := w.TemplateService.GetTemplate(templateID)
	sendError(err)

	// Render the template edit page
	defaultContent := templates.Edit(template)
	base := layout.Layout(defaultContent, "", nil)

	err = base.Render(c.Request.Context(), c.Writer)
	sendError(err)
}

func (w TemplateHandler) Update(c *gin.Context) {
	var body struct {
		Name    string `form:"name"`
		Content string `form:"content"`
	}

	sendError := func(err error) {
		if err != nil {
			c = ContextWithError(c, err.Error())
			w.Create(c)
		}
	}

	if err := c.Bind(&body); err != nil {
		sendError(err)
		return
	}

	templateIDParam := c.Param("id")
	if templateIDParam == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Convert the template ID from string to int64
	templateID, err := ConvertToInt64(templateIDParam)
	sendError(err)

	err = w.TemplateService.UpdateTemplate(templateID, body.Name, body.Content)
	sendError(err)

	message := "template updated"
	c = ContextWithMessage(c, message)

	w.Edit(c)
}

var redirectTemplateURL = "http://%s/dashboard/templates/"

func (w TemplateHandler) Store(c *gin.Context) {
	var body struct {
		Name    string `form:"name"`
		Content string `form:"content"`
	}
	if err := c.Bind(&body); err != nil {
		c = ContextWithError(c, err.Error())
		w.Create(c)
		return
	}
	ownerID := GetUserIDFromContext(c)
	workspaceID := GetWorkspaceIDFromContext(c)
	_, err := w.TemplateService.CreateTemplate(body.Name, body.Content, ownerID, workspaceID)
	if err != nil {
		c = ContextWithError(c, err.Error())
		w.Create(c)
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf(redirectTemplateURL, c.Request.Host))
}

func (w TemplateHandler) Delete(c *gin.Context) {
	templateIDParam := c.Param("id")
	if templateIDParam == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	templateID, err := ConvertToInt64(templateIDParam)
	err = w.TemplateService.DeleteTemplate(templateID)
	if err != nil {
		c = ContextWithError(c, err.Error())
	}

	c.Redirect(http.StatusFound, "/dashboard/templates")
}

func NewTemplateHandler(templateService services.TemplateService) TemplateHandler {
	return TemplateHandler{TemplateService: templateService}
}
