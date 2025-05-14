package mapper

import (
	"time"

	"github.com/itsLeonB/catfeinated-time-tracker/internal/constant"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
)

func UserTaskToResponse(userTask model.UserTask) dto.UserTaskResponse {
	logs := make([]dto.UserTaskLogResponse, len(userTask.Logs))
	for i, log := range userTask.Logs {
		logs[i] = UserTaskLogToResponse(log)
	}

	userTaskResponse := dto.UserTaskResponse{
		ID:         userTask.ID,
		UserId:     userTask.UserId,
		TaskId:     userTask.TaskId,
		TaskNumber: userTask.Task.Number,
		TaskName:   userTask.Task.Name,
		ProjectId:  userTask.Task.ProjectID,
		CreatedAt:  userTask.CreatedAt,
		UpdatedAt:  userTask.UpdatedAt,
		TimeSpent:  dto.TimeSpentFromDuration(userTask.GetTotalTime()),
		Logs:       logs,
	}

	if !userTask.Task.IsZero() {
		userTaskResponse.Task = TaskToResponse(userTask.Task)
	}

	return userTaskResponse
}

func PopulateUserTaskWithLogs(userTask model.UserTask, userTaskLogs []dto.UserTaskLogResponse) dto.UserTaskResponse {
	durationSpent := durationSpentFromLogs(userTaskLogs)

	return dto.UserTaskResponse{
		ID:         userTask.ID,
		UserId:     userTask.UserId,
		TaskId:     userTask.TaskId,
		TaskNumber: userTask.Task.Number,
		TaskName:   userTask.Task.Name,
		ProjectId:  userTask.Task.ProjectID,
		CreatedAt:  userTask.CreatedAt,
		UpdatedAt:  userTask.UpdatedAt,
		TimeSpent:  dto.TimeSpentFromDuration(durationSpent),
		Logs:       userTaskLogs,
	}
}

func durationSpentFromLogs(logs []dto.UserTaskLogResponse) time.Duration {
	var totalTime time.Duration

	for i := 0; i < len(logs); i += 2 {
		if i+1 < len(logs) {
			startLog := logs[i+1]
			stopLog := logs[i]

			if startLog.Action == constant.LogAction.Start && stopLog.Action == constant.LogAction.Stop {
				duration := stopLog.CreatedAt.Sub(startLog.CreatedAt)
				totalTime += duration
			}
		}
	}

	return totalTime
}

func QueryParamsToOptions(queryParams dto.UserTaskQueryParams) model.UserTaskQueryOptions {
	return model.UserTaskQueryOptions{
		Filters: map[string]any{
			"user_id": queryParams.UserId,
			"task_id": queryParams.TaskId,
		},
		PreloadRelations: []string{"Task"},
	}
}
