package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/constant"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/service"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/util"
	"github.com/rotisserie/eris"
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
		response, err := uth.userTaskService.FindAll(ctx, dto.UserTaskQueryParams{})
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

func (uth *UserTaskHandler) HandleCreateForm() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userProjectId, err := util.GetUuidParam(ctx, "id")
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request, err := util.BindFormRequest[dto.NewUserTaskRequest](ctx)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.UserProjectId = userProjectId

		_, err = uth.userTaskService.Create(ctx, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/user-projects/%s", userProjectId))
	}
}

func (uth *UserTaskHandler) HandleLogPost() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		projectId, err := util.GetUuidParam(ctx, "id")
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		userTaskId, err := util.GetUuidParam(ctx, "userTaskId")
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		action := ctx.Param("action")
		if action != constant.LogAction.Start && action != constant.LogAction.Stop {
			_ = ctx.Error(eris.Errorf("%s is not a valid action", action))
			return
		}

		newLogRequest := dto.NewUserTaskLogRequest{
			UserTaskId: userTaskId,
			Action:     action,
		}

		_, err = uth.userTaskLogService.Create(ctx, newLogRequest)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/user-projects/%s", projectId))
	}
}
