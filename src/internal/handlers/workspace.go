package handlers

import (
	"avan/internal/services"
	"avan/views/components"
	"avan/views/layout"
	"avan/views/workspaces"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WorkspaceHandler struct {
	WorkspaceService services.WorkspaceService
}

func (w WorkspaceHandler) Index(c *gin.Context) {
	userID := GetUserIDFromContext(c)
	_workspaces, err := w.WorkspaceService.GetWorkspaces(c.Request.Context(), userID)
	if err != nil {
		c = ContextWithError(c, err.Error())
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	currentPageStr := c.Query("page")
	currentPage, _ := ConvertToInt32(currentPageStr)
	page := components.PageData{
		Route:       "workspaces",
		CurrentPage: currentPage,
		MaxResults:  int32(len(_workspaces)),
		Bucket:      10,
	}

	defaultContent := workspaces.Index(_workspaces, page)
	var items []layout.MenuItem
	base := layout.Layout(defaultContent, "", items)
	err = base.Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (w WorkspaceHandler) Create(c *gin.Context) {
	defaultContent := workspaces.WorkspacesCreate()

	base := layout.Layout(defaultContent, "", nil)
	err := base.Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (w WorkspaceHandler) Edit(c *gin.Context) {
	sendError := func(err error) {
		if err != nil {
			c = ContextWithError(c, err.Error())
			_ = c.AbortWithError(http.StatusInternalServerError, err)
		}
	}

	workspaceIDParam := c.Param("id")
	if workspaceIDParam == "" {
		c.Redirect(http.StatusNotFound, "/dashboard/workspaces")
		return
	}

	// Convert the workspace ID from string to int64
	workspaceID, err := ConvertToInt64(workspaceIDParam)
	sendError(err)

	// Fetch the workspace from the service
	workspace, err := w.WorkspaceService.GetWorkspace(workspaceID)
	sendError(err)

	// Render the workspace edit page
	defaultContent := workspaces.Edit(workspace)
	base := layout.Layout(defaultContent, "", nil)

	err = base.Render(c.Request.Context(), c.Writer)
	sendError(err)
}

func (w WorkspaceHandler) Update(c *gin.Context) {
	var body struct {
		Name string `form:"name"`
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

	workspaceIDParam := c.Param("id")
	if workspaceIDParam == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	// Convert the workspace ID from string to int64
	workspaceID, err := ConvertToInt64(workspaceIDParam)
	sendError(err)

	err = w.WorkspaceService.UpdateWorkspace(workspaceID, body.Name)
	sendError(err)

	message := "workspace updated"
	c = ContextWithMessage(c, message)

	w.Edit(c)
}

var redirectURL = "http://%s/dashboard/workspaces/"

func (w WorkspaceHandler) Store(c *gin.Context) {
	var body struct {
		Name string `form:"name"`
	}
	if err := c.Bind(&body); err != nil {
		c = ContextWithError(c, err.Error())
		w.Create(c)
		return
	}
	ownerID := GetUserIDFromContext(c)
	_, err := w.WorkspaceService.CreateWorkspace(body.Name, ownerID)
	if err != nil {
		c = ContextWithError(c, err.Error())
		w.Create(c)
		return
	}

	c.Redirect(http.StatusSeeOther, fmt.Sprintf(redirectURL, c.Request.Host))
}

func (w WorkspaceHandler) SetDefault(c *gin.Context) {
	workspaceIDParam := c.Param("id")
	if workspaceIDParam == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	workspaceID, err := ConvertToInt64(workspaceIDParam)
	ownerID := GetUserIDFromContext(c)
	if err != nil {
		c = ContextWithError(c, err.Error())
		c.Redirect(http.StatusFound, "/dashboard/workspaces")
		return
	}

	err = w.WorkspaceService.SetDefaultWorkspace(workspaceID, ownerID)
	if err != nil {
		c = ContextWithError(c, err.Error())
	}

	c.Redirect(http.StatusFound, "/dashboard/workspaces")
}

func (w WorkspaceHandler) Delete(c *gin.Context) {
	workspaceIDParam := c.Param("id")
	if workspaceIDParam == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	workspaceID, err := ConvertToInt64(workspaceIDParam)
	err = w.WorkspaceService.DeleteWorkspace(workspaceID)
	if err != nil {
		c = ContextWithError(c, err.Error())
	}

	c.Redirect(http.StatusFound, "/dashboard/workspaces")
}

func NewWorkspaceHandler(workspaceService services.WorkspaceService) WorkspaceHandler {
	return WorkspaceHandler{WorkspaceService: workspaceService}
}
