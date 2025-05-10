package strategy

import (
	"log"
	"net/http"

	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
)

type unhandledErrorStrategy struct{}

func (es *unhandledErrorStrategy) HandleError(err error) *dto.ErrorResponse {
	log.Printf("unexpected error type: %T\n", err)

	return &dto.ErrorResponse{
		Code:    http.StatusInternalServerError,
		Type:    "InternalServerError",
		Message: "An unexpected error occurred",
	}
}
