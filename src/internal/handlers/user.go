package handlers

import (
	"avan/internal/db/sqlc"
	"avan/internal/services"
	"avan/views/layout"
	"avan/views/users"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService services.UserService
}

func (w UserHandler) Index(c *gin.Context) {
	workspaceID := GetWorkspaceIDFromContext(c)
	userList, err := w.UserService.GetUsers(c.Request.Context(), workspaceID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	roleList, err := w.UserService.GetRoles(c.Request.Context(), workspaceID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	workspacesValue, _ := c.Get("workspaces")
	workspaceList, _ := workspacesValue.([]sqlc.Workspace)
	manageable := make([]sqlc.Workspace, 0, len(workspaceList))
	for _, workspace := range workspaceList {
		allowed, err := w.UserService.HasPermission(c.Request.Context(), workspace.ID, GetUserIDFromContext(c), "users.manage")
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if allowed {
			manageable = append(manageable, workspace)
		}
	}
	workspaceList = manageable
	workspaceRoles := make(map[int64][]users.RoleOption, len(workspaceList))
	for _, workspace := range workspaceList {
		roles, err := w.UserService.GetRoles(c.Request.Context(), workspace.ID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		options := make([]users.RoleOption, 0, len(roles))
		for _, role := range roles {
			options = append(options, users.RoleOption{ID: role.ID, Name: role.Name})
		}
		workspaceRoles[workspace.ID] = options
	}
	members := make([]users.Member, 0, len(userList))
	for _, user := range userList {
		member := users.Member{User: user, IsSuper: w.UserService.Repo.Base.Config.AdminEmail != "" && strings.EqualFold(user.Email, w.UserService.Repo.Base.Config.AdminEmail)}
		if !member.IsSuper {
			assigned, err := w.UserService.UserWorkspaceRoles(c.Request.Context(), user.ID)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			assignedRoles := make(map[int64]int64, len(assigned))
			for _, role := range assigned {
				assignedRoles[role.WorkspaceID] = role.RoleID
			}
			for _, workspace := range workspaceList {
				member.Workspaces = append(member.Workspaces, users.WorkspaceAccess{
					Workspace:      workspace,
					Roles:          workspaceRoles[workspace.ID],
					SelectedRoleID: assignedRoles[workspace.ID],
				})
			}
		}
		members = append(members, member)
	}

	page := users.Index(members, roleList)
	base := layout.Layout(page, "", nil)
	if err := base.Render(c.Request.Context(), c.Writer); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

var userPermissions = []string{
	"analysis.read", "api_keys.manage", "developer.read", "messages.read",
	"providers.manage", "settings.manage", "templates.manage", "users.manage", "workspaces.manage",
}

func (w UserHandler) SaveRole(c *gin.Context) {
	workspaceID := GetWorkspaceIDFromContext(c)
	roleID := int64(0)
	if id := c.Param("id"); id != "" {
		var err error
		roleID, err = ConvertToInt64(id)
		if err != nil || roleID <= 0 {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	permissions := c.PostFormArray("permissions")
	allowed := make(map[string]bool, len(userPermissions))
	for _, permission := range userPermissions {
		allowed[permission] = true
	}
	for _, permission := range permissions {
		if !allowed[permission] {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}
	name := strings.TrimSpace(c.PostForm("name"))
	if name == "" || len(name) > 50 {
		c = ContextWithError(c, "Role name is required and must be 50 characters or fewer")
		w.Index(c)
		return
	}
	if err := w.UserService.SaveRole(c.Request.Context(), workspaceID, roleID, name, permissions); err != nil {
		c = ContextWithError(c, err.Error())
		w.Index(c)
		return
	}
	c.Redirect(http.StatusSeeOther, "/dashboard/users")
}

func (w UserHandler) DeleteRole(c *gin.Context) {
	roleID, err := ConvertToInt64(c.Param("id"))
	if err != nil || roleID <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := w.UserService.DeleteRole(c.Request.Context(), GetWorkspaceIDFromContext(c), roleID); err != nil {
		c = ContextWithError(c, err.Error())
		w.Index(c)
		return
	}
	c.Redirect(http.StatusSeeOther, "/dashboard/users")
}

func (w UserHandler) AssignRole(c *gin.Context) {
	userID, err := ConvertToInt64(c.Param("id"))
	roleID, roleErr := ConvertToInt64(c.PostForm("role_id"))
	if err != nil || roleErr != nil || userID <= 0 || roleID < 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if w.userIsSuperAdmin(c, userID) || !w.userInWorkspace(c, userID) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	if err := w.UserService.AssignRole(c.Request.Context(), GetWorkspaceIDFromContext(c), userID, roleID); err != nil {
		c = ContextWithError(c, err.Error())
		w.Index(c)
		return
	}
	c.Redirect(http.StatusSeeOther, "/dashboard/users")
}

func (w UserHandler) userIsSuperAdmin(c *gin.Context, userID int64) bool {
	isAdmin, err := w.UserService.IsSuperAdmin(c.Request.Context(), userID)
	return err != nil || isAdmin
}

func (w UserHandler) userInWorkspace(c *gin.Context, userID int64) bool {
	users, err := w.UserService.GetUsers(c.Request.Context(), GetWorkspaceIDFromContext(c))
	if err != nil {
		return false
	}
	for _, user := range users {
		if user.ID == userID {
			return true
		}
	}
	return false
}

func (w UserHandler) Update(c *gin.Context) {
	userID, err := ConvertToInt64(c.Param("id"))
	if err != nil || userID <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if w.userIsSuperAdmin(c, userID) || !w.userInWorkspace(c, userID) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	name := strings.TrimSpace(c.PostForm("name"))
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")
	parsedEmail, emailErr := mail.ParseAddress(email)
	if name == "" || len(name) > 50 || emailErr != nil || parsedEmail.Address != email || len(email) > 50 || (password != "" && len(password) < 8) {
		c = ContextWithError(c, "Enter a name, valid email, and an optional password of at least 8 characters")
		w.Index(c)
		return
	}
	if err := w.UserService.UpdateUser(c.Request.Context(), userID, name, email, password); err != nil {
		c = ContextWithError(c, err.Error())
		w.Index(c)
		return
	}
	c.Redirect(http.StatusSeeOther, "/dashboard/users")
}

func (w UserHandler) Delete(c *gin.Context) {
	userID, err := ConvertToInt64(c.Param("id"))
	if err != nil || userID <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if w.userIsSuperAdmin(c, userID) || !w.userInWorkspace(c, userID) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	if err := w.UserService.DeleteUser(c.Request.Context(), userID); err != nil {
		c = ContextWithError(c, err.Error())
		w.Index(c)
		return
	}
	c.Redirect(http.StatusSeeOther, "/dashboard/users")
}

func (w UserHandler) AssignWorkspaces(c *gin.Context) {
	userID, err := ConvertToInt64(c.Param("id"))
	if err != nil || userID <= 0 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if w.userIsSuperAdmin(c, userID) || !w.userInWorkspace(c, userID) {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}
	workspacesValue, _ := c.Get("workspaces")
	workspaces, ok := workspacesValue.([]sqlc.Workspace)
	if !ok {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	allowed := make(map[int64]bool, len(workspaces))
	workspaceIDs := make([]int64, 0, len(workspaces))
	for _, workspace := range workspaces {
		canManage, err := w.UserService.HasPermission(c.Request.Context(), workspace.ID, GetUserIDFromContext(c), "users.manage")
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if canManage {
			allowed[workspace.ID] = true
			workspaceIDs = append(workspaceIDs, workspace.ID)
		}
	}
	assignments := make(map[int64]int64)
	for _, rawID := range c.PostFormArray("workspace_ids") {
		workspaceID, err := ConvertToInt64(rawID)
		if err != nil || !allowed[workspaceID] {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		roleID, err := ConvertToInt64(c.PostForm(fmt.Sprintf("role_%d", workspaceID)))
		if err != nil || roleID <= 0 {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		assignments[workspaceID] = roleID
	}
	if err := w.UserService.AssignWorkspaces(c.Request.Context(), userID, workspaceIDs, assignments); err != nil {
		c = ContextWithError(c, err.Error())
		w.Index(c)
		return
	}
	c.Redirect(http.StatusSeeOther, "/dashboard/users")
}

func (w UserHandler) Store(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")
	roleID, err := ConvertToInt64(c.PostForm("role_id"))
	if name == "" || len(name) > 50 || email == "" || len(email) > 50 || len(password) < 8 || roleID <= 0 || err != nil {
		c = ContextWithError(c, "Enter a name, email, password of at least 8 characters, and role")
		w.Index(c)
		return
	}
	err = w.UserService.CreateWorkspaceUser(c.Request.Context(), GetWorkspaceIDFromContext(c), name, email, password, roleID)
	if err != nil {
		c = ContextWithError(c, fmt.Sprint(err))
		w.Index(c)
		return
	}
	c.Redirect(http.StatusSeeOther, "/dashboard/users")
}

func NewUserHandler(userService services.UserService) UserHandler {
	return UserHandler{UserService: userService}
}
