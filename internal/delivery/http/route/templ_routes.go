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

	router.Static("/static", "./static")

	router.GET("/login", handlers.Auth.HandleLoginPage())
	router.GET("/register", handlers.Auth.HandleRegisterPage())

	router.POST("/auth/login", handlers.Auth.HandleLoginForm())
	router.POST("/auth/register", handlers.Auth.HandleRegisterForm())

	protectedRoutes := router.Group("/", authMiddleware)

	protectedRoutes.GET("/", handlers.Home.HandleHomePage())
	protectedRoutes.GET("/logout", handlers.Auth.HandleLogout())

	protectedRoutes.POST("/user-projects", handlers.Project.HandleNewProjectForm())
	protectedRoutes.GET("/user-projects/:id", handlers.Project.HandleProjectDetailPage())
	protectedRoutes.POST("/user-projects/:id/user-tasks", handlers.UserTask.HandleCreateForm())
	protectedRoutes.POST("/user-projects/:id/user-tasks/:userTaskId/:action", handlers.UserTask.HandleLogPost())
}
