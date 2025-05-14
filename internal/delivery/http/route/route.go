package route

import (
	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/provider"
)

func SetupRoutes(router *gin.Engine, handlers *provider.Handlers, services *provider.Services) {
	router.HandleMethodNotAllowed = true
	router.ContextWithFallback = true

	router.NoRoute(handlers.Root.NotFound())

	setupApiRoutes(router, handlers, services, "/api")
	setupTemplRoutes(router, handlers, services)
}
