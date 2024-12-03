package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/itsLeonB/time-tracker/internal/service"
	"github.com/rotisserie/eris"
)

type ProjectHandler struct {
	projectService service.ProjectService
}

func NewProjectHandler(projectService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService}
}

func (ph *ProjectHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := new(model.NewProjectRequest)
		err := ctx.ShouldBindJSON(request)
		if err != nil {
			_ = ctx.Error(eris.Wrap(err, "error reading request body"))
			return
		}

		project, err := ph.projectService.Create(request.Name)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, model.NewSuccessJSON(project))
	}
}
