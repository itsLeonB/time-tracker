package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil"
	"github.com/itsLeonB/time-tracker/internal/appconstant"
	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/itsLeonB/time-tracker/internal/service"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(
	taskService service.TaskService,
) *TaskHandler {
	return &TaskHandler{
		taskService,
	}
}

func (th *TaskHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request, err := ezutil.BindRequest[dto.NewTaskRequest](ctx, binding.JSON)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		task, err := th.taskService.Create(ctx, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusCreated,
			ezutil.NewResponse(appconstant.MsgInsertData).WithData(task),
		)
	}
}

func (th *TaskHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tasks, err := th.taskService.GetAll(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(
			http.StatusOK,
			ezutil.NewResponse(appconstant.MsgGetData).WithData(tasks),
		)
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

		ctx.JSON(
			http.StatusOK,
			ezutil.NewResponse(appconstant.MsgGetData).WithData(tasks),
		)
	}
}

func (th *TaskHandler) HandleCreate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		projectID, exists, err := ezutil.GetPathParam[uuid.UUID](ctx, appconstant.ContextProjectID)
		if err != nil || !exists {
			_ = ctx.Error(err)
			return
		}

		request, err := ezutil.BindRequest[dto.NewTaskRequest](ctx, binding.Form)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.ProjectID = projectID

		_, err = th.taskService.Create(ctx, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/projects/%s", projectID))
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
