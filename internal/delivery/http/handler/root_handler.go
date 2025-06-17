package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/ezutil"
)

type RootHandler struct{}

func (h *RootHandler) Root() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, ezutil.NewResponse("time-tracker"))
	}
}

func (h *RootHandler) HealthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, ezutil.NewResponse("app is healthy"))
	}
}

func (h *RootHandler) NotFound() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(
			http.StatusNotFound,
			ezutil.NewErrorResponse(ezutil.AppError{
				Type:    "RouteNotFoundError",
				Message: "Route is not found",
			}),
		)
	}
}
