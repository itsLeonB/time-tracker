package mapper

import (
	"time"

	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/itsLeonB/time-tracker/internal/entity"
)

func TaskLogToResponse(taskLog entity.TaskLog) dto.TaskLogResponse {
	return dto.TaskLogResponse{
		Id:           taskLog.ID,
		TaskID:       taskLog.TaskID,
		Action:       taskLog.Action,
		CreatedAt:    taskLog.CreatedAt,
		CreatedAtIso: taskLog.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    taskLog.UpdatedAt,
	}
}
