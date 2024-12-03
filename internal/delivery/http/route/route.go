package route

import (
	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/time-tracker/internal/provider"
)

func SetupRoutes(router *gin.Engine, handlers *provider.Handlers) *gin.Engine {
	router.HandleMethodNotAllowed = true
	router.ContextWithFallback = true

	projectRoutes := router.Group("/projects")
	projectRoutes.POST("", handlers.Project.Create())

	return router
}
