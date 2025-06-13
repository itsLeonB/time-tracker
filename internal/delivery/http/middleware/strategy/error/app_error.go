package strategy

import (
	"log"

	"github.com/itsLeonB/time-tracker/internal/apperror"
	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/rotisserie/eris"
)

type appErrorStrategy struct{}

func (es *appErrorStrategy) HandleError(err error) *dto.ErrorResponse {
	appError := err.(*apperror.AppError)

	if appError.HttpStatusCode == 500 {
		log.Println(eris.ToString(appError.Err, true))
	}

	return &dto.ErrorResponse{
		Code:    appError.HttpStatusCode,
		Type:    appError.Type,
		Message: appError.Message,
		Details: appError.Err.Error(),
	}
}
