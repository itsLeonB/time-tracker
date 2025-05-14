package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/constant"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/service"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/util"
)

type UserTaskHandler struct {
	userTaskService    service.UserTaskService
	userTaskLogService service.UserTaskLogService
}

func NewUserTaskHandler(
	userTaskService service.UserTaskService,
	userTaskLogService service.UserTaskLogService,
) *UserTaskHandler {
	return &UserTaskHandler{
		userTaskService,
		userTaskLogService,
	}
}

func (uth *UserTaskHandler) HandleCreate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request, err := util.BindJsonRequest[dto.NewUserTaskRequest](ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		userId, err := util.GetUuidFromCtx(ctx, constant.ContextUserID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.UserId = userId

		response, err := uth.userTaskService.Create(ctx, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, dto.NewSuccessJSON(response))
	}
}

func (uth *UserTaskHandler) HandleFindAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId, err := util.GetUuidFromCtx(ctx, constant.ContextUserID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		queryParams := dto.UserTaskQueryParams{
			UserId: userId,
		}

		response, err := uth.userTaskService.FindAll(ctx, queryParams)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, dto.NewSuccessJSON(response))
	}
}

func (uth *UserTaskHandler) HandleGetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// userId, err := util.GetUuidFromCtx(ctx, constant.ContextUserID)
		// if err != nil {
		// 	_ = ctx.Error(err)
		// 	return
		// }

		id, err := util.GetUuidParam(ctx, "id")
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		// queryParams := dto.UserTaskQueryParams{
		// 	UserId: userId,
		// }

		response, err := uth.userTaskService.GetById(ctx, id)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, dto.NewSuccessJSON(response))
	}
}

func (uth *UserTaskHandler) HandleLog() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request, err := util.BindJsonRequest[dto.NewUserTaskLogRequest](ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		userTaskId, err := util.GetUuidParam(ctx, "id")
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.UserTaskId = userTaskId

		response, err := uth.userTaskLogService.Create(ctx, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, dto.NewSuccessJSON(response))
	}
}
