package route

import (
	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/delivery/http/middleware"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/delivery/http/renderer"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/provider"
)

func setupTemplRoutes(router *gin.Engine, handlers *provider.Handlers, services *provider.Services) {
	htmlRenderer := router.HTMLRender
	router.HTMLRender = &renderer.HTMLTemplRenderer{FallbackHtmlRenderer: htmlRenderer}

	authMiddleware := middleware.AuthorizeViaCookie(services.JWT)

	router.SetTrustedProxies(nil)

	router.GET("/register", handlers.Auth.HandleRegisterPage())
	router.POST("/auth/register", handlers.Auth.HandleRegisterForm())

	router.GET("/login", handlers.Auth.HandleLoginPage())
	router.POST("/auth/login", handlers.Auth.HandleLoginForm())

	protectedRoutes := router.Group("/", authMiddleware)

	protectedRoutes.GET("/", handlers.Home.HandleHomePage())
	protectedRoutes.GET("/logout", handlers.Auth.HandleLogout())

	protectedRoutes.GET("/projects/:id", handlers.Project.HandleProjectDetailPage())
	protectedRoutes.POST("/projects/:id/tasks", handlers.Task.HandleAddToUserProject())
}
