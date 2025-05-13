package mapper

import (
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
)

func UserTaskToResponse(userTask model.UserTask) dto.UserTaskResponse {
	return dto.UserTaskResponse{
		ID:         userTask.ID,
		UserId:     userTask.UserId,
		TaskId:     userTask.TaskId,
		TaskNumber: userTask.Task.Number,
		TaskName:   userTask.Task.Name,
		ProjectId:  userTask.Task.ProjectID,
		CreatedAt:  userTask.CreatedAt,
		UpdatedAt:  userTask.UpdatedAt,
	}
}

func QueryParamsToOptions(queryParams dto.UserTaskQueryParams) model.UserTaskQueryOptions {
	return model.UserTaskQueryOptions{
		Filters: map[string]any{
			"user_id": queryParams.UserId,
		},
		PreloadRelations: []string{"Task"},
	}
}
