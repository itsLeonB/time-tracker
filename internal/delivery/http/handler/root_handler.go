package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/time-tracker/internal/dto"
)

type RootHandler struct{}

func (h *RootHandler) Root() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, dto.NewSuccessJSON("time-tracker"))
	}
}

func (h *RootHandler) HealthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, dto.NewSuccessJSON("app is healthy"))
	}
}

func (h *RootHandler) NotFound() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, dto.NewErrorJSON(&dto.ErrorResponse{
			Type:    "RouteNotFoundError",
			Message: "Route is not found",
		}))
	}
}

func (h *RootHandler) TestReadData() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body, err := io.ReadAll(ctx.Request.Body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}

		// Restore the body so it can be read again if needed
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Try to pretty-print the JSON if it's JSON
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, body, "", "  "); err == nil {
			// If it's valid JSON, log the pretty version
			log.Println("Request Body (JSON):\n", prettyJSON.String())
		} else {
			// If not JSON, log the raw body
			log.Println("Request Body (raw):\n", string(body))
		}
	}
}
