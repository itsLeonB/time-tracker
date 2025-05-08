package strategy

import (
	"log"

	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/rotisserie/eris"
)

type appErrorStrategy struct{}

func (es *appErrorStrategy) HandleError(err error) *model.ErrorResponse {
	appError := err.(*apperror.AppError)

	if appError.HttpStatusCode == 500 {
		log.Println(eris.ToString(appError.Err, true))
	}

	return &model.ErrorResponse{
		Code:    appError.HttpStatusCode,
		Type:    appError.Type,
		Message: appError.Message,
		Details: appError.Err.Error(),
	}
}
