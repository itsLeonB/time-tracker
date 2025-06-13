package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/itsLeonB/time-tracker/internal/service"
	"github.com/itsLeonB/time-tracker/templates/pages"
)

type HomeHandler struct {
	userService    service.UserService
	projectService service.ProjectService
}

func NewHomeHandler(
	userService service.UserService,
	projectService service.ProjectService,
) *HomeHandler {
	return &HomeHandler{
		userService,
		projectService,
	}
}

func (hh *HomeHandler) HandleHomePage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := hh.userService.ValidateUser(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		projects, err := hh.projectService.FindByUserId(ctx, user.ID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		homeViewDto := dto.HomeViewDto{
			User:     *user,
			Projects: projects,
		}

		ctx.HTML(http.StatusOK, "", pages.Home(homeViewDto))
	}
}
