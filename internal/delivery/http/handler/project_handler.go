package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/constant"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/service"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/util"
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
		request := new(dto.NewProjectRequest)
		err := ctx.ShouldBindJSON(request)
		if err != nil {
			_ = ctx.Error(apperror.BadRequestError(
				eris.Wrap(err, apperror.MsgBindJsonError),
			))
			return
		}

		project, err := ph.projectService.Create(ctx, request.Name)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, dto.NewSuccessJSON(project))
	}
}

func (ph *ProjectHandler) HandleGetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := util.GetUuidFromCtx(ctx, constant.ContextUserID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		projects, err := ph.projectService.FindByUserId(ctx, userId)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, dto.NewSuccessJSON(projects))
	}
}

func (ph *ProjectHandler) HandleGetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := util.GetUuidParam(ctx, "id")
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		userId, err := util.GetUuidFromCtx(ctx, constant.ContextUserID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		params, err := util.GetDatetimeParams(ctx, "start", "end")
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		params.ProjectID = id
		params.UserId = userId

		projectParams := dto.ProjectQueryParams{QueryParams: params}

		project, err := ph.projectService.FirstByQuery(ctx, projectParams)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, dto.NewSuccessJSON(project))
	}
}

func (ph *ProjectHandler) FirstByQuery() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		name := ctx.Query("name")
		options := dto.ProjectQueryParams{Name: name}

		project, err := ph.projectService.FirstByQuery(ctx, options)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, dto.NewSuccessJSON(project))
	}
}
