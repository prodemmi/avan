package handlers

import (
	"avan/internal/services"
	"avan/views/layout"
	"avan/views/users"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService services.UserService
}

func (w UserHandler) Index(c *gin.Context) {
	userList, err := w.UserService.GetUsers(c.Request.Context(), GetWorkspaceIDFromContext(c))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	page := users.Index(userList)
	base := layout.Layout(page, "", nil)
	if err := base.Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func NewUserHandler(userService services.UserService) UserHandler {
	return UserHandler{UserService: userService}
}
