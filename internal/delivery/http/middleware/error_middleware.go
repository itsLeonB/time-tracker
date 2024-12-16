package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/time-tracker/internal/apperror"
	strategy "github.com/itsLeonB/time-tracker/internal/delivery/http/middleware/strategy/error"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/rotisserie/eris"
)

func HandleError(errorStrategyMap *strategy.ErrorStrategyMap) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if err := ctx.Errors.Last(); err != nil {
			var handledError error
			handledError, ok := err.Err.(*apperror.AppError)
			if !ok {
				handledError = eris.Unwrap(err.Err)
			}

			errorStrategy := errorStrategyMap.DetermineStrategy(handledError)
			errorResponse := errorStrategy.HandleError(handledError)
			ctx.AbortWithStatusJSON(errorResponse.Code, model.NewErrorJSON(errorResponse))
			return
		}
	}
}
