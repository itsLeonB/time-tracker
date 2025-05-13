package mapper

import (
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
)

func UserTaskLogToResponse(userTaskLog model.UserTaskLog) dto.UserTaskLogResponse {
	return dto.UserTaskLogResponse{
		Id:         userTaskLog.ID,
		UserTaskId: userTaskLog.UserTaskId,
		Action:     userTaskLog.Action,
		CreatedAt:  userTaskLog.CreatedAt,
		UpdatedAt:  userTaskLog.UpdatedAt,
	}
}
