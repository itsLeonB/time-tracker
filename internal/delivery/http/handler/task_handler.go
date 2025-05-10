package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/service"
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

func (th *TaskHandler) Log() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := new(model.NewLogRequest)
		err := ctx.ShouldBindJSON(request)
		if err != nil {
			_ = ctx.Error(apperror.BadRequestError(
				eris.Wrap(err, apperror.MsgBindJsonError),
			))
			return
		}

		parsedId, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			_ = ctx.Error(apperror.BadRequestError(
				eris.Wrap(err, "error parsing UUID"),
			))
			return
		}

		log, err := th.taskService.Log(ctx, parsedId, request.Action)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, dto.NewSuccessJSON(log))
	}
}

func (th *TaskHandler) LogByNumber() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := new(model.LogByNumberRequest)
		err := ctx.ShouldBindJSON(request)
		if err != nil {
			_ = ctx.Error(apperror.BadRequestError(
				eris.Wrap(err, apperror.MsgBindJsonError),
			))
			return
		}

		log, err := th.taskService.LogByNumber(ctx, request.Number, request.Action)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, dto.NewSuccessJSON(log))
	}
}

func (th *TaskHandler) Find() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		queryOptions := constructQueryOptions(ctx)
		tasks, err := th.taskService.Find(ctx, queryOptions)
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

func constructQueryOptions(ctx *gin.Context) *dto.QueryOptions {
	options := new(dto.QueryOptions)

	withLogs := ctx.Query("withLogs")
	if withLogs == "true" {
		options.WithLogs = true
	}

	params := constructQueryParams(ctx)
	if params != nil {
		options.Params = params
	}

	return options
}

func constructQueryParams(ctx *gin.Context) *dto.QueryParams {
	number := ctx.Query("number")
	if number == "" {
		return nil
	}

	return &dto.QueryParams{Number: number}
}
