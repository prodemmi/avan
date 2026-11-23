package handlers

import (
	"avan/internal/services"
	"avan/views/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService services.AuthService
}

func (w AuthHandler) LoginPage(c *gin.Context) {
	base := auth.Login("/auth/login")
	err := base.Render(c.Request.Context(), c.Writer)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
}

func (w AuthHandler) Login(c *gin.Context) {
	var body struct {
		Email    string `form:"email"`
		Password string `form:"password"`
	}
	ctx := c.Request.Context()
	if err := c.Bind(&body); err != nil {
		c = ContextWithError(c, err.Error())
		w.LoginPage(c)
		return
	}

	user, err := w.AuthService.LoginUser(ctx, body.Email, body.Password)
	if err != nil {
		c = ContextWithError(c, err.Error())
		w.LoginPage(c)
		return
	}

	token, err := w.AuthService.GenerateJWT(user.ID)
	if err != nil {
		c = ContextWithError(c, err.Error())
		w.LoginPage(c)
		return
	}
	secure := RequestIsHTTPS(c)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth_token", token, 3600, "/", "", secure, true)

	c.Redirect(http.StatusFound, "/dashboard")
}

func (w AuthHandler) Logout(c *gin.Context) {
	secure := RequestIsHTTPS(c)
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("auth_token", "", 3600, "/", "", secure, true)

	c.Redirect(http.StatusFound, "/auth/login")
}

func NewAuthHandler(authService services.AuthService) AuthHandler {
	return AuthHandler{AuthService: authService}
}
