package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/apperror"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/itsLeonB/time-tracker/internal/service"
	"github.com/rotisserie/eris"
)

type TaskHandler struct {
	taskService service.TaskService
}

func NewTaskHandler(taskService service.TaskService) *TaskHandler {
	return &TaskHandler{taskService}
}

func (th *TaskHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request := new(model.NewTaskRequest)
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

		ctx.JSON(http.StatusCreated, model.NewSuccessJSON(task))
	}
}

func (th *TaskHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tasks, err := th.taskService.GetAll(ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, model.NewSuccessJSON(tasks))
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

		ctx.JSON(http.StatusCreated, model.NewSuccessJSON(log))
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

		ctx.JSON(http.StatusCreated, model.NewSuccessJSON(log))
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

		ctx.JSON(http.StatusOK, model.NewSuccessJSON(tasks))
	}

}

func constructQueryOptions(ctx *gin.Context) *model.QueryOptions {
	options := new(model.QueryOptions)

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

func constructQueryParams(ctx *gin.Context) *model.QueryParams {
	number := ctx.Query("number")
	if number == "" {
		return nil
	}

	return &model.QueryParams{Number: number}
}
