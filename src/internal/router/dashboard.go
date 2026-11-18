package router

import (
	"avan/internal/middlewares"
	"avan/views"
	"avan/views/layout"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (r Router) GenerateDashboardRoutes() {
	authMiddleware := middlewares.AuthMiddleware(r.AuthHandler.AuthService, r.WorkspaceHandler.WorkspaceService)
	dashboard := r.Server.Engine.Group("/dashboard", authMiddleware)
	{
		dashboard.GET("/", func(c *gin.Context) {
			// مقدار پیش‌فرض برای محتوای صفحه‌ی داشبورد
			defaultContent := views.Home()
			var items []layout.MenuItem
			base := layout.Layout(defaultContent, "", items)
			err := base.Render(c.Request.Context(), c.Writer)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		})

	}

	r.GenerateAuthRoutes(r.Server.Engine)
	r.GenerateApiRoutes(r.Server.Engine)
	r.GenerateApiKeyRoutes(dashboard)
	r.GenerateUserRoutes(dashboard)
	r.GenerateProviderRoutes(dashboard)
	r.GenerateWorkspaceRoutes(dashboard)
	r.GenerateTemplateRoutes(dashboard)
	r.GenerateMessageRoutes(dashboard)
	r.GenerateAnalysisRoutes(dashboard)
	r.GenerateSettingsRoutes(dashboard)
}
