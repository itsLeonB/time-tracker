package route

import (
	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/time-tracker/internal/delivery/http/middleware"
	strategy "github.com/itsLeonB/time-tracker/internal/delivery/http/middleware/strategy/error"
	"github.com/itsLeonB/time-tracker/internal/provider"
)

func setupApiRoutes(router *gin.Engine, handlers *provider.Handlers, services *provider.Services, path string) {
	errorStrategyMap := strategy.NewErrorStrategyMap()
	apiRoutes := router.Group(path, middleware.HandleError(errorStrategyMap))

	apiRoutes.GET("", handlers.Root.Root())
	apiRoutes.GET("/health", handlers.Root.HealthCheck())

	apiRoutes.POST("/test", handlers.Root.TestReadData())

	authRoutes := apiRoutes.Group("/auth")
	authRoutes.POST("/register", handlers.Auth.HandleRegister())
	authRoutes.POST("/login", handlers.Auth.HandleLogin())

	authenticatedRoutes := apiRoutes.Group("", middleware.Authorize(services.JWT))

	projectRoutes := authenticatedRoutes.Group("/projects")
	projectRoutes.POST("", handlers.Project.Create())
	projectRoutes.GET("", handlers.Project.HandleGetAll())
	projectRoutes.GET("/:id", handlers.Project.HandleGetById())
	projectRoutes.GET("/first", handlers.Project.FirstByQuery())

	taskRoutes := authenticatedRoutes.Group("/tasks")
	taskRoutes.POST("", handlers.Task.Create())
	taskRoutes.GET("", handlers.Task.HandleFind())

	externalTrackerRoutes := authenticatedRoutes.Group("/external")
	externalTaskRoutes := externalTrackerRoutes.Group("/tasks")
	externalTaskRoutes.GET("", handlers.Task.FindExternal())

	userTaskRoutes := authenticatedRoutes.Group("/user-tasks")
	userTaskRoutes.POST("", handlers.UserTask.HandleCreate())
	userTaskRoutes.GET("", handlers.UserTask.HandleFindAll())
	userTaskRoutes.GET("/:id", handlers.UserTask.HandleGetById())
	userTaskRoutes.POST("/:id/logs", handlers.UserTask.HandleLog())
}
