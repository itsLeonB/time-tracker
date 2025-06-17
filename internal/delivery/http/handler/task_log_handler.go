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

type TaskLogHandler struct {
	taskLogService service.TaskLogService
}

func NewTaskLogHandler(
	taskLogService service.TaskLogService,
) *TaskLogHandler {
	return &TaskLogHandler{
		taskLogService,
	}
}

func (tlh *TaskLogHandler) HandleCreate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		projectID, err := ezutil.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextProjectID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		taskID, err := ezutil.GetRequiredPathParam[uuid.UUID](ctx, appconstant.ContextTaskID)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request, err := ezutil.BindRequest[dto.NewTaskLogRequest](ctx, binding.Form)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		request.ProjectID = projectID
		request.TaskID = taskID

		_, err = tlh.taskLogService.Create(ctx, request)
		if err != nil {
			_ = ctx.Error(err)
			return
		}

		ctx.Redirect(http.StatusSeeOther, fmt.Sprintf("/projects/%s", projectID))
	}
}
