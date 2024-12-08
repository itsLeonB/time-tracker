package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/time-tracker/internal/model"
)

type RootHandler struct{}

func (h *RootHandler) Root() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, model.NewSuccessJSON("time-tracker"))
	}
}

func (h *RootHandler) HealthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, model.NewSuccessJSON("app is healthy"))
	}
}

func (h *RootHandler) NotFound() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, model.NewErrorJSON(&model.ErrorResponse{
			Type:    "RouteNotFoundError",
			Message: "Route is not found",
		}))
	}

}
