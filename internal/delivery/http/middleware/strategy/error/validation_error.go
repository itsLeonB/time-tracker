package strategy

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/itsLeonB/time-tracker/internal/model"
)

type validationErrorStrategy struct{}

func (es *validationErrorStrategy) HandleError(err error) *model.ErrorResponse {
	validationErrors := err.(validator.ValidationErrors)
	var errMap = make(map[string]string, len(validationErrors))

	for _, validationError := range validationErrors {
		errMap[validationError.Field()] = validationError.Error()
	}

	return &model.ErrorResponse{
		Code:    http.StatusBadRequest,
		Type:    "BadRequestError",
		Message: "Error validating request body",
		Details: errMap,
	}
}
