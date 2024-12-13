package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/time-tracker/internal/apperror"
	"github.com/itsLeonB/time-tracker/internal/auth"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/rotisserie/eris"
)

type AuthHandler struct {
	authService auth.AuthService
}

func NewAuthHandler(authService auth.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (ah *AuthHandler) HandleRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request model.RegisterRequest
		err := ctx.ShouldBindJSON(&request)
		if err != nil {
			_ = ctx.Error(eris.Wrap(err, apperror.MsgBindJsonError))
			return
		}

		response, err := ah.authService.Register(ctx, &request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, model.NewSuccessJSON(response))
	}
}

func (ah *AuthHandler) HandleLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request model.LoginRequest
		err := ctx.ShouldBindJSON(&request)
		if err != nil {
			_ = ctx.Error(eris.Wrap(err, apperror.MsgBindJsonError))
			return
		}

		response, err := ah.authService.Login(ctx, &request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, model.NewSuccessJSON(response))
	}
}
