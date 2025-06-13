package mapper

import (
	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/itsLeonB/time-tracker/internal/model"
)

func TaskToResponse(task model.Task) dto.TaskResponse {
	return dto.TaskResponse{
		ID:        task.ID,
		Number:    task.Number,
		Name:      task.Name,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}
