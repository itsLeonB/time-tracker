package route

import (
	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/time-tracker/internal/provider"
)

func SetupRoutes(router *gin.Engine, handlers *provider.Handlers, services *provider.Services) {
	setupApiRoutes(router, handlers, services)
	setupTemplRoutes(router, handlers, services)
}
