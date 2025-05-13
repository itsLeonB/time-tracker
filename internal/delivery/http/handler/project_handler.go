package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (ph *ProjectHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		parsedId, err := uuid.Parse(id)
		if err != nil {
			_ = ctx.Error(apperror.BadRequestError(
				eris.Wrap(err, "error parsing UUID"),
			))
			return
		}

		options := dto.QueryOptions{
			Params: &dto.QueryParams{
				ProjectID: parsedId,
			},
		}

		date := ctx.Query("date")
		if date == "today" {
			options.Params.Date = time.Now()
		} else if date != "" {
			timestamp, err := time.Parse(time.DateOnly, date)
			if err != nil {
				_ = ctx.Error(apperror.BadRequestError(
					eris.Wrap(err, "error parsing date"),
				))
				return
			}

			options.Params.Date = timestamp
		}

		project, err := ph.projectService.GetByID(ctx, &options)
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
		options := &dto.FindProjectOptions{Name: name}

		project, err := ph.projectService.FirstByQuery(ctx, options)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, dto.NewSuccessJSON(project))
	}
}
