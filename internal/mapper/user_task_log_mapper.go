package mapper

import (
	"time"

	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
)

func UserTaskLogToResponse(userTaskLog model.UserTaskLog) dto.UserTaskLogResponse {
	return dto.UserTaskLogResponse{
		Id:           userTaskLog.ID,
		UserTaskId:   userTaskLog.UserTaskId,
		Action:       userTaskLog.Action,
		CreatedAt:    userTaskLog.CreatedAt,
		CreatedAtIso: userTaskLog.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    userTaskLog.UpdatedAt,
	}
}
