package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/time-tracker/internal/apperror"
	"github.com/itsLeonB/time-tracker/internal/auth"
	"github.com/itsLeonB/time-tracker/internal/constant"
	"github.com/rotisserie/eris"
)

func Authorize(jwt auth.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			_ = ctx.Error(apperror.UnauthorizedError(eris.Errorf("missing token")))
			ctx.Abort()
			return
		}

		splits := strings.Split(token, " ")
		if len(splits) != 2 || splits[0] != "Bearer" {
			_ = ctx.Error(apperror.UnauthorizedError(eris.Errorf("invalid token")))
			ctx.Abort()
			return
		}

		claims, err := jwt.VerifyToken(splits[1])
		if err != nil {
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Set(constant.ContextUserID, claims.Data[constant.ContextUserID])
		ctx.Next()
	}
}
