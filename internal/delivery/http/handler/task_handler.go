package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/service"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/util"
	"github.com/rotisserie/eris"
)

type TaskHandler struct {
	taskService         service.TaskService
	externalTaskService service.ExternalTrackerService
}

func NewTaskHandler(
	taskService service.TaskService,
	externalTrackerService service.ExternalTrackerService,
) *TaskHandler {
	return &TaskHandler{
		taskService,
		externalTrackerService,
	}
}

func (th *TaskHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := new(dto.NewTaskRequest)
		err := ctx.ShouldBindJSON(request)
		if err != nil {
			_ = ctx.Error(apperror.BadRequestError(
				eris.Wrap(err, apperror.MsgBindJsonError),
			))
			return
		}

		task, err := th.taskService.Create(ctx, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, dto.NewSuccessJSON(task))
	}
}

func (th *TaskHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tasks, err := th.taskService.GetAll(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, dto.NewSuccessJSON(tasks))
	}
}

func (th *TaskHandler) HandleFind() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		queryParams := constructQueryParams(ctx)

		tasks, err := th.taskService.Find(ctx, queryParams)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, dto.NewSuccessJSON(tasks))
	}
}

func (th *TaskHandler) FindExternal() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		number := ctx.Query("number")

		queryOptions := dto.ExternalQueryOptions{Number: number}

		tasks, err := th.externalTaskService.FindTask(ctx, queryOptions)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, dto.NewSuccessJSON(tasks))
	}
}

func (th *TaskHandler) HandleAddToUserProject() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		projectId, err := util.GetUuidParam(ctx, "id")
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		addToProjectRequest, err := util.BindFormRequest[dto.AddToProjectRequest](ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		addToProjectRequest.ProjectID = projectId

		_, err = th.taskService.AddToUserProject(ctx, addToProjectRequest)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/projects/%s", projectId))
	}
}

func constructQueryParams(ctx *gin.Context) dto.TaskQueryParams {
	var queryParams dto.TaskQueryParams

	number := ctx.Query("number")
	if number != "" {
		queryParams.Number = number
	}

	return queryParams
}
