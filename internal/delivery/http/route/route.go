package route

import (
	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/delivery/http/middleware"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/provider"
)

func SetupRoutes(router *gin.Engine, handlers *provider.Handlers, services *provider.Services) {
	router.HandleMethodNotAllowed = true
	router.ContextWithFallback = true

	router.NoRoute(handlers.Root.NotFound())
	router.GET("", handlers.Root.Root())
	router.GET("/health", handlers.Root.HealthCheck())

	router.POST("/test", handlers.Root.TestReadData())

	authRoutes := router.Group("/auth")
	authRoutes.POST("/register", handlers.Auth.HandleRegister())
	authRoutes.POST("/login", handlers.Auth.HandleLogin())

	authenticatedRoutes := router.Group("", middleware.Authorize(services.JWT))

	projectRoutes := authenticatedRoutes.Group("/projects")
	projectRoutes.POST("", handlers.Project.Create())
	projectRoutes.GET("", handlers.Project.GetAll())
	projectRoutes.GET("/:id", handlers.Project.GetByID())
	projectRoutes.GET("/first", handlers.Project.FirstByQuery())

	taskRoutes := authenticatedRoutes.Group("/tasks")
	taskRoutes.POST("", handlers.Task.Create())
	taskRoutes.GET("", handlers.Task.Find())
	taskRoutes.POST("/:id/logs", handlers.Task.Log())
	taskRoutes.POST("/log-by-number", handlers.Task.LogByNumber())

	externalTrackerRoutes := authenticatedRoutes.Group("/external")
	externalTaskRoutes := externalTrackerRoutes.Group("/tasks")
	externalTaskRoutes.GET("", handlers.Task.FindExternal())
}
