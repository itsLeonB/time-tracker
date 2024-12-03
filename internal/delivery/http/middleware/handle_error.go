package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/rotisserie/eris"
)

func HandleError() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
		if err := ctx.Errors.Last(); err != nil {
			log.Println(eris.ToString(eris.Unwrap(err), true))
			abortWithError(ctx, http.StatusInternalServerError, err)
			return
		}
	}
}

func abortWithError(ctx *gin.Context, statusCode int, err error) {
	ctx.AbortWithStatusJSON(statusCode, model.NewErrorJSON(err))
}
