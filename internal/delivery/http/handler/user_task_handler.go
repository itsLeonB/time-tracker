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
	userTaskService service.UserTaskService
}

func NewUserTaskHandler(userTaskService service.UserTaskService) *UserTaskHandler {
	return &UserTaskHandler{
		userTaskService: userTaskService,
	}
}

func (uth *UserTaskHandler) HandleCreate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		request, err := util.BindRequest[dto.NewUserTaskRequest](ctx)
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

		ctx.JSON(http.StatusOK, dto.NewSuccessJSON(response))
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
