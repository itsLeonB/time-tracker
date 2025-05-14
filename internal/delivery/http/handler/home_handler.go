package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/service"
	"github.com/itsLeonB/catfeinated-time-tracker/templates/pages"
)

type HomeHandler struct {
	userService service.UserService
}

func NewHomeHandler(userService service.UserService) *HomeHandler {
	return &HomeHandler{userService}
}

func (hh *HomeHandler) HandleHomePage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := hh.userService.ValidateUser(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.HTML(http.StatusOK, "", pages.Home(user))
	}
}
