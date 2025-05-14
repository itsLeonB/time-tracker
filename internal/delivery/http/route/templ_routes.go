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

	router.GET("/", authMiddleware, handlers.Home.HandleHomePage())
	router.GET("/logout", authMiddleware, handlers.Auth.HandleLogout())
}
