package strategy

import (
	"log"
	"net/http"

	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
)

type unhandledErrorStrategy struct{}

func (es *unhandledErrorStrategy) HandleError(err error) *model.ErrorResponse {
	log.Printf("unexpected error type: %T\n", err)

	return &model.ErrorResponse{
		Code:    http.StatusInternalServerError,
		Type:    "InternalServerError",
		Message: "An unexpected error occurred",
	}
}
