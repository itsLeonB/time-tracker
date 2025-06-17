package route

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil"
	"github.com/itsLeonB/time-tracker/internal/appconstant"
	"github.com/itsLeonB/time-tracker/internal/provider"
	"github.com/itsLeonB/time-tracker/internal/service"
	"github.com/rotisserie/eris"
)

func setupApiRoutes(router *gin.Engine, handlers *provider.Handlers, services *provider.Services) {
	// Middlewares
	tokenCheckFunc := newTokenCheckFunc(services.JWT, services.User)
	authMiddleware := ezutil.NewAuthMiddleware("Bearer", tokenCheckFunc)
	errorMiddleware := ezutil.NewErrorMiddleware()

	routeConfigs := []ezutil.RouteConfig{
		{
			Group: "/api",
			Handlers: []gin.HandlerFunc{
				errorMiddleware,
			},
			Versions: []ezutil.RouteVersionConfig{
				{
					Version:  1,
					Handlers: []gin.HandlerFunc{},
					Groups: []ezutil.RouteGroupConfig{
						{
							Group:    "/",
							Handlers: []gin.HandlerFunc{},
							Endpoints: []ezutil.EndpointConfig{
								{
									Method:   "GET",
									Endpoint: "/",
									Handlers: []gin.HandlerFunc{handlers.Root.Root()},
								},
								{
									Method:   "GET",
									Endpoint: "/health",
									Handlers: []gin.HandlerFunc{handlers.Root.HealthCheck()},
								},
							},
						},
						{
							Group:    "/auth",
							Handlers: []gin.HandlerFunc{},
							Endpoints: []ezutil.EndpointConfig{
								{
									Method:   "POST",
									Endpoint: "/register",
									Handlers: []gin.HandlerFunc{
										handlers.Auth.HandleRegister(),
									},
								},
								{
									Method:   "POST",
									Endpoint: "/login",
									Handlers: []gin.HandlerFunc{
										handlers.Auth.HandleLogin(),
									},
								},
							},
						},
						{
							Group:    "/",
							Handlers: []gin.HandlerFunc{authMiddleware},
							Endpoints: []ezutil.EndpointConfig{
								{
									Method:   "POST",
									Endpoint: "/projects",
									Handlers: []gin.HandlerFunc{handlers.Project.Create()},
								},
								{
									Method:   "GET",
									Endpoint: "/projects",
									Handlers: []gin.HandlerFunc{handlers.Project.HandleGetAll()},
								},
								{
									Method:   "GET",
									Endpoint: "/projects/:id",
									Handlers: []gin.HandlerFunc{handlers.Project.HandleGetById()},
								},
								{
									Method:   "GET",
									Endpoint: "/projects/first",
									Handlers: []gin.HandlerFunc{handlers.Project.FirstByQuery()},
								},
								{
									Method:   "POST",
									Endpoint: "/tasks",
									Handlers: []gin.HandlerFunc{handlers.Task.Create()},
								},
								{
									Method:   "GET",
									Endpoint: "/tasks",
									Handlers: []gin.HandlerFunc{handlers.Task.HandleFind()},
								},
							},
						},
					},
				},
			},
		},
	}

	ezutil.SetupRoutes(router, routeConfigs)
}

func newTokenCheckFunc(jwtService ezutil.JWTService, userService service.UserService) func(ctx *gin.Context, token string) (bool, map[string]any, error) {
	return func(ctx *gin.Context, token string) (bool, map[string]any, error) {
		claims, err := jwtService.VerifyToken(token)
		if err != nil {
			return false, nil, err
		}

		tokenUserId, exists := claims.Data[appconstant.ContextUserID]
		if !exists {
			return false, nil, eris.New("missing user ID from token")
		}
		stringUserID, ok := tokenUserId.(string)
		if !ok {
			return false, nil, eris.New("error asserting userID, is not a string")
		}
		userID, err := ezutil.Parse[uuid.UUID](stringUserID)
		if err != nil {
			return false, nil, err
		}

		exists, err = userService.ExistsByID(ctx, userID)
		if err != nil {
			return false, nil, err
		}
		if !exists {
			return false, nil, ezutil.UnauthorizedError(appconstant.MsgAuthUserNotFound)
		}

		authData := map[string]any{
			appconstant.ContextUserID: userID,
		}

		return true, authData, nil
	}
}
