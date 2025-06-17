package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

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

// generateCSRFToken generates a random CSRF token
func (ah *AuthHandler) generateCSRFToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// validateCSRFToken validates the CSRF token from form and session
func (ah *AuthHandler) validateCSRFToken(ctx *gin.Context) bool {
	formToken := ctx.PostForm("_token")
	sessionToken, err := ctx.Cookie("csrf_token")
	if err != nil || formToken == "" || sessionToken == "" {
		return false
	}
	return formToken == sessionToken
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

		_, err = ah.authService.Login(ctx, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ezutil.NewResponse(appconstant.MsgLoginSuccess),
		)
	}
}

func (ah *AuthHandler) HandleLoginPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if isLoggedIn, err := ah.isLoggedIn(ctx); isLoggedIn && err == nil {
			ctx.Redirect(http.StatusSeeOther, "/")
			return
		}

		// Generate and set CSRF token
		csrfToken := ah.generateCSRFToken()
		ctx.SetCookie("csrf_token", csrfToken, 3600, "/", "", false, true)

		ctx.HTML(http.StatusOK, "", auth_pages.Login("", ""))
	}
}

func (ah *AuthHandler) HandleLoginForm() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Validate CSRF token
		if !ah.validateCSRFToken(ctx) {
			ah.loginError(ctx, "Invalid CSRF token")
			return
		}

		var loginRequest dto.LoginRequest
		if err := ctx.ShouldBind(&loginRequest); err != nil {
			ah.loginError(ctx, err.Error())
			return
		}

		result, err := ah.authService.Login(ctx, loginRequest)
		if err != nil {
			ah.loginError(ctx, err.Error())
			return
		}

		ctx.SetCookie("session_token", result.Token, 24*60*60, "/", "", false, true)
		ctx.Redirect(http.StatusSeeOther, "/")
	}
}

func (ah *AuthHandler) HandleLogout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.SetCookie("session_token", "", -1, "/", "", false, true)
		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func (ah *AuthHandler) HandleRegisterPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if isLoggedIn, err := ah.isLoggedIn(ctx); isLoggedIn && err == nil {
			ctx.Redirect(http.StatusSeeOther, "/")
			return
		}

		// Generate and set CSRF token
		csrfToken := ah.generateCSRFToken()
		ctx.SetCookie("csrf_token", csrfToken, 3600, "/", "", false, true)

		ctx.HTML(http.StatusOK, "", auth_pages.Register(""))
	}
}

func (ah *AuthHandler) HandleRegisterForm() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Validate CSRF token
		if !ah.validateCSRFToken(ctx) {
			ah.registerError(ctx, "Invalid CSRF token")
			return
		}

		var registerRequest dto.RegisterRequest
		if err := ctx.ShouldBind(&registerRequest); err != nil {
			ah.registerError(ctx, err.Error())
			return
		}

		err := ah.authService.Register(ctx, registerRequest)
		if err != nil {
			ah.registerError(ctx, err.Error())
			return
		}

		ctx.Redirect(http.StatusSeeOther, "/login")
	}
}

func (ah *AuthHandler) loginError(ctx *gin.Context, errMsg string) {
	csrfToken := ah.generateCSRFToken()
	ctx.SetCookie("csrf_token", csrfToken, 3600, "/", "", false, true)
	ctx.HTML(http.StatusBadRequest, "", auth_pages.Login(errMsg, ""))
}

func (ah *AuthHandler) registerError(ctx *gin.Context, errMsg string) {
	csrfToken := ah.generateCSRFToken()
	ctx.SetCookie("csrf_token", csrfToken, 3600, "/", "", false, true)
	ctx.HTML(http.StatusBadRequest, "", auth_pages.Register(errMsg))
}

func (ah *AuthHandler) isLoggedIn(ctx *gin.Context) (bool, error) {
	cookie, err := ctx.Request.Cookie("session_token")
	if err != nil {
		return false, err
	}

	isValid, err := ah.authService.CheckToken(ctx, cookie.Value)
	if err != nil {
		return false, err
	}

	return isValid, nil
}
