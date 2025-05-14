package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/auth"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/constant"
	"github.com/itsLeonB/catfeinated-time-tracker/templates/auth_pages"
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

		isValid, token := isValidBearerToken(strings.Split(token, " "))
		if !isValid {
			_ = ctx.Error(apperror.UnauthorizedError(eris.Errorf("invalid token")))
			ctx.Abort()
			return
		}

		claims, err := jwt.VerifyToken(token)
		if err != nil {
			_ = ctx.Error(err)
			ctx.Abort()
			return
		}

		ctx.Set(constant.ContextUserID, claims.Data[constant.ContextUserID])
		ctx.Next()
	}
}

func AuthorizeViaCookie(jwt auth.JWT) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookie, err := ctx.Request.Cookie("session_token")
		if err != nil {
			ctx.Abort()
			ctx.HTML(http.StatusForbidden, "", auth_pages.Login("Not logged in", ""))
			return
		}

		claims, err := jwt.VerifyToken(cookie.Value)
		if err != nil {
			ctx.Abort()
			ctx.HTML(http.StatusForbidden, "", auth_pages.Login(err.Error(), ""))
			return
		}

		ctx.Set(constant.ContextUserID, claims.Data[constant.ContextUserID])
		ctx.Next()
	}
}

func isValidBearerToken(splits []string) (bool, string) {
	// First, check the length in a way that takes consistent time
	if len(splits) != 2 {
		return false, ""
	}

	// Use subtle.ConstantTimeCompare for timing-safe comparison
	ok := subtle.ConstantTimeCompare([]byte(splits[0]), []byte("Bearer")) == 1
	if !ok {
		return false, ""
	}

	return true, splits[1]
}
