package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
			_ = ctx.Error(eris.Wrap(err, "error reading request body"))
			return
		}

		task, err := th.taskService.Create(request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, model.NewSuccessJSON(task))
	}
}

func (th *TaskHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tasks, err := th.taskService.GetAll()
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
			_ = ctx.Error(eris.Wrap(err, "error reading request body"))
			return
		}

		parsedId, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			_ = ctx.Error(eris.Wrap(err, "error parsing UUID"))
			return
		}

		log, err := th.taskService.Log(parsedId, request.Action)
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
			_ = ctx.Error(eris.Wrap(err, "error reading request body"))
			return
		}

		log, err := th.taskService.LogByNumber(request.Number, request.Action)
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
		tasks, err := th.taskService.Find(queryOptions)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, model.NewSuccessJSON(tasks))
	}

}

func constructQueryOptions(ctx *gin.Context) *model.QueryOptions {
	number := ctx.Query("number")
	if number == "" {
		return nil
	}

	return &model.QueryOptions{
		Params: map[string]any{
			"number": number,
		},
	}
}
