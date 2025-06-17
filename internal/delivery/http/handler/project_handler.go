package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil"
	"github.com/itsLeonB/time-tracker/internal/appconstant"
	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/itsLeonB/time-tracker/internal/service"
	"github.com/itsLeonB/time-tracker/templates/pages"
)

type ProjectHandler struct {
	projectService service.ProjectService
	userService    service.UserService
}

func NewProjectHandler(
	projectService service.ProjectService,
	userService service.UserService,
) *ProjectHandler {
	return &ProjectHandler{
		projectService,
		userService,
	}
}

func (ph *ProjectHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request, err := ezutil.BindRequest[dto.NewProjectRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		project, err := ph.projectService.Create(ctx, request.Name)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusCreated,
			ezutil.NewResponse(appconstant.MsgInsertData).WithData(project),
		)
	}
}

func (ph *ProjectHandler) HandleGetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := ezutil.GetAndParseFromContext[uuid.UUID](ctx, appconstant.ContextUserID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		projects, err := ph.projectService.FindByUserId(ctx, userId)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ezutil.NewResponse(appconstant.MsgGetData).WithData(projects),
		)
	}
}

func (ph *ProjectHandler) HandleGetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := ezutil.GetAndParseFromContext[uuid.UUID](ctx, appconstant.ContextUserID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		project, _, err := ph.getProjectById(ctx, userId)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ezutil.NewResponse(appconstant.MsgGetData).WithData(project),
		)
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

		ctx.JSON(
			http.StatusOK,
			ezutil.NewResponse(appconstant.MsgInsertData).WithData(project),
		)
	}
}

func (ph *ProjectHandler) HandleProjectDetailPage() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := ph.userService.ValidateUser(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		project, params, err := ph.getProjectById(ctx, user.ID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		viewDto := dto.ProjectDetailViewDto{
			User:      user,
			Project:   project,
			StartDate: ezutil.FormatTimeNullable(params.StartDatetime, time.DateOnly),
			EndDate:   ezutil.FormatTimeNullable(params.EndDatetime, time.DateOnly),
		}

		ctx.HTML(http.StatusOK, "", pages.ProjectDetail(viewDto))
	}
}

func (ph *ProjectHandler) HandleNewProjectForm() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := ph.userService.ValidateUser(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request, err := ezutil.BindRequest[dto.NewProjectRequest](ctx, binding.Form)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		_, err = ph.projectService.Create(ctx, request.Name)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.Redirect(http.StatusSeeOther, "/")
	}
}

func (ph *ProjectHandler) getProjectById(ctx *gin.Context, userId uuid.UUID) (dto.ProjectResponse, dto.QueryParams, error) {
	id, err := ezutil.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextProjectID)
	if err != nil {
		return dto.ProjectResponse{}, dto.QueryParams{}, err
	}

	params, err := ezutil.BindRequest[dto.QueryParams](ctx, binding.Query)
	if err != nil {
		return dto.ProjectResponse{}, dto.QueryParams{}, err
	}

	params.UserId = userId
	projectParams := dto.ProjectQueryParams{QueryParams: params}
	projectParams.ID = id

	project, err := ph.projectService.FirstByQuery(ctx, projectParams)
	if err != nil {
		return dto.ProjectResponse{}, dto.QueryParams{}, err
	}

	return project, params, nil
}
