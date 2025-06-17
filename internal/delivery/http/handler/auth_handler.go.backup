package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/itsLeonB/ezutil"
	"github.com/itsLeonB/time-tracker/internal/appconstant"
	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/itsLeonB/time-tracker/internal/service"
	"github.com/itsLeonB/time-tracker/templates/auth_pages"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (ah *AuthHandler) HandleRegister() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request, err := ezutil.BindRequest[dto.RegisterRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		err = ah.authService.Register(ctx, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusCreated,
			ezutil.NewResponse(appconstant.MsgRegisterSuccess),
		)
	}
}

func (ah *AuthHandler) HandleLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request, err := ezutil.BindRequest[dto.LoginRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		response, err := ah.authService.Login(ctx, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ezutil.NewResponse(appconstant.MsgLoginSuccess).WithData(response),
		)
	}
}

func (ah *AuthHandler) HandleLoginPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if isLoggedIn, err := ah.isLoggedIn(ctx); isLoggedIn && err == nil {
			ctx.Redirect(http.StatusSeeOther, "/")
			return
		}

		ctx.HTML(http.StatusOK, "", auth_pages.Login("", ""))
	}
}

func (ah *AuthHandler) HandleLoginForm() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request, err := ezutil.BindRequest[dto.LoginRequest](ctx, binding.Form)
		if err != nil {
			loginError(ctx, err)
			return
		}

		response, err := ah.authService.Login(ctx, request)
		if err != nil {
			loginError(ctx, err)
			return
		}

		ctx.SetCookie("session_token", response.Token, int((24 * time.Hour).Nanoseconds()), "/", "", true, true)
		ctx.Redirect(http.StatusSeeOther, "/")
	}
}

func (ah *AuthHandler) HandleLogout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.SetCookie("session_token", "", -1, "/", "", true, true)
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func (ah *AuthHandler) HandleRegisterPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if isLoggedIn, err := ah.isLoggedIn(ctx); isLoggedIn && err == nil {
			ctx.Redirect(http.StatusSeeOther, "/")
			return
		}

		ctx.HTML(http.StatusOK, "", auth_pages.Register(""))
	}
}

func (ah *AuthHandler) HandleRegisterForm() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request, err := ezutil.BindRequest[dto.RegisterRequest](ctx, binding.Form)
		if err != nil {
			registerError(ctx, err)
			return
		}

		err = ah.authService.Register(ctx, request)
		if err != nil {
			registerError(ctx, err)
			return
		}

		ctx.HTML(http.StatusOK, "", auth_pages.Login("", appconstant.MsgRegisterSuccess))
	}
}

func loginError(ctx *gin.Context, err error) {
	ctx.HTML(http.StatusForbidden, "", auth_pages.Login(err.Error(), ""))
}

func registerError(ctx *gin.Context, err error) {
	ctx.HTML(http.StatusInternalServerError, "", auth_pages.Register(err.Error()))
}

func (ah *AuthHandler) isLoggedIn(ctx *gin.Context) (bool, error) {
	cookie, err := ctx.Request.Cookie("session_token")
	if err == nil {
		return ah.authService.CheckToken(ctx, cookie.Value)
	}

	return false, nil
}
