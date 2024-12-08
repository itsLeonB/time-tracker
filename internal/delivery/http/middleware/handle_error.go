package middleware

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/time-tracker/internal/apperror"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/rotisserie/eris"
)

func HandleError() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if err := ctx.Errors.Last(); err != nil {
			appError := err.Err.(*apperror.AppError)
			log.Println(eris.ToString(appError.Err, true))

			errorResponse := model.ErrorResponse{
				Code:    appError.HttpStatusCode,
				Type:    appError.Type,
				Message: appError.Message,
				Details: appError.Err.Error(),
			}

			abortWithError(ctx, &errorResponse)
			return
		}
	}
}

func abortWithError(ctx *gin.Context, err *model.ErrorResponse) {
	ctx.AbortWithStatusJSON(err.Code, model.NewErrorJSON(err))
}
