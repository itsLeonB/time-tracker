package route

import (
	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/time-tracker/internal/provider"
)

func SetupRoutes(router *gin.Engine, handlers *provider.Handlers) *gin.Engine {
	router.HandleMethodNotAllowed = true
	router.ContextWithFallback = true

	router.NoRoute(handlers.Root.NotFound())
	router.GET("", handlers.Root.Root())
	router.GET("/health", handlers.Root.HealthCheck())

	projectRoutes := router.Group("/projects")
	projectRoutes.POST("", handlers.Project.Create())
	projectRoutes.GET("", handlers.Project.GetAll())
	projectRoutes.GET("/:id", handlers.Project.GetByID())
	projectRoutes.GET("/first", handlers.Project.FirstByQuery())

	taskRoutes := router.Group("/tasks")
	taskRoutes.POST("", handlers.Task.Create())
	taskRoutes.GET("", handlers.Task.Find())
	taskRoutes.POST("/:id/logs", handlers.Task.Log())
	taskRoutes.POST("/log-by-number", handlers.Task.LogByNumber())

	return router
}
