package route

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/ezutil"
	"github.com/itsLeonB/time-tracker/internal/appconstant"
	"github.com/itsLeonB/time-tracker/internal/delivery/http/renderer"
	"github.com/itsLeonB/time-tracker/internal/provider"
	"github.com/itsLeonB/time-tracker/templates/auth_pages"
)

func setupTemplRoutes(router *gin.Engine, handlers *provider.Handlers, services *provider.Services) {
	htmlRenderer := router.HTMLRender
	router.HTMLRender = &renderer.HTMLTemplRenderer{FallbackHtmlRenderer: htmlRenderer}

	authMiddleware := authorizeViaCookie(services.JWT)

	_ = router.SetTrustedProxies(nil)

	router.Static("/static", "./static")

	router.GET("/login", handlers.Auth.HandleLoginPage())
	router.GET("/register", handlers.Auth.HandleRegisterPage())

	router.POST("/auth/login", handlers.Auth.HandleLoginForm())
	router.POST("/auth/register", handlers.Auth.HandleRegisterForm())

	protectedRoutes := router.Group("/", authMiddleware)

	protectedRoutes.GET("/", handlers.Home.HandleHomePage())
	protectedRoutes.GET("/logout", handlers.Auth.HandleLogout())

	protectedRoutes.POST("/projects", handlers.Project.HandleNewProjectForm())

	protectedRoutes.GET(
		fmt.Sprintf("/projects/:%s", appconstant.ContextProjectID),
		handlers.Project.HandleProjectDetailPage(),
	)

	protectedRoutes.POST(
		fmt.Sprintf("/projects/:%s/tasks", appconstant.ContextProjectID),
		handlers.Task.HandleCreate(),
	)

	protectedRoutes.POST(
		fmt.Sprintf(
			"/projects/:%s/tasks/:%s/logs",
			appconstant.ContextProjectID,
			appconstant.ContextTaskID,
		),
		handlers.TaskLog.HandleCreate(),
	)
}

func authorizeViaCookie(jwt ezutil.JWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Request.Cookie("session_token")
		if err != nil {
			ctx.Abort()
			ctx.HTML(http.StatusForbidden, "", auth_pages.Login("Not logged in", ""))
			return
		}

		claims, err := jwt.VerifyToken(cookie.Value)
		if err != nil {
			ctx.Abort()
			ctx.HTML(http.StatusForbidden, "", auth_pages.Login(err.Error(), ""))
			return
		}

		for key, val := range claims.Data {
			ctx.Set(key, val)
		}

		ctx.Next()
	}
}
