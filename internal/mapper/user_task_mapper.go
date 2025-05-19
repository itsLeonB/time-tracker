package mapper

import (
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/constant"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/rotisserie/eris"
)

func UserTaskToResponse(userTask model.UserTask) (dto.UserTaskResponse, error) {
	logs := make([]dto.UserTaskLogResponse, len(userTask.Logs))
	for i, log := range userTask.Logs {
		logs[i] = UserTaskLogToResponse(log)
	}

	totalUserTaskTime, err := durationSpentFromLogs(logs, userTask.ID)
	if err != nil {
		return dto.UserTaskResponse{}, err
	}

	userTaskResponse := dto.UserTaskResponse{
		ID:            userTask.ID,
		UserProjectId: userTask.UserProjectId,
		TaskId:        userTask.TaskId,
		CreatedAt:     userTask.CreatedAt,
		UpdatedAt:     userTask.UpdatedAt,
		TaskNumber:    userTask.Task.Number,
		TaskName:      userTask.Task.Name,
		TimeSpent:     dto.TimeSpentFromDuration(totalUserTaskTime),
		Logs:          logs,
		IsActive:      isActive(logs),
		StartTime:     getStartTime(logs),
	}

	return userTaskResponse, nil
}

func PopulateUserTaskWithLogs(userTask model.UserTask, userTaskLogs []dto.UserTaskLogResponse) (dto.UserTaskResponse, error) {
	durationSpent, err := durationSpentFromLogs(userTaskLogs, userTask.ID)
	if err != nil {
		return dto.UserTaskResponse{}, err
	}

	userTaskResponse := dto.UserTaskResponse{
		ID:            userTask.ID,
		UserProjectId: userTask.UserProjectId,
		TaskId:        userTask.TaskId,
		CreatedAt:     userTask.CreatedAt,
		UpdatedAt:     userTask.UpdatedAt,
		TaskNumber:    userTask.Task.Number,
		TaskName:      userTask.Task.Name,
		TimeSpent:     dto.TimeSpentFromDuration(durationSpent),
		Logs:          userTaskLogs,
		IsActive:      isActive(userTaskLogs),
		StartTime:     getStartTime(userTaskLogs),
	}

	return userTaskResponse, nil
}

func isActive(logs []dto.UserTaskLogResponse) bool {
	return len(logs) > 0 && logs[len(logs)-1].Action == constant.LogAction.Start
}
func durationSpentFromLogs(logs []dto.UserTaskLogResponse, userTaskId uuid.UUID) (time.Duration, error) {
	var totalTime time.Duration
	var pendingStart *time.Time

	for _, log := range logs {
		switch log.Action {
		case constant.LogAction.Start:
			if pendingStart != nil {
				// Unpaired START: Count time from previous START until now
				totalTime += time.Since(*pendingStart)
			}
			pendingStart = &log.CreatedAt

		case constant.LogAction.Stop:
			if pendingStart == nil {
				return totalTime, eris.Errorf("Orphaned STOP log (no matching START) for task ID: %s", userTaskId)
			}
			totalTime += log.CreatedAt.Sub(*pendingStart)
			pendingStart = nil
		}
	}

	// If a START is left unpaired, count time until now
	if pendingStart != nil {
		totalTime += time.Since(*pendingStart)
	}

	return totalTime, nil
}

func QueryParamsToOptions(queryParams dto.UserTaskQueryParams) model.UserTaskQueryOptions {
	return model.UserTaskQueryOptions{
		Filters: map[string]any{
			"user_project_id": queryParams.UserProjectId,
			"task_id":         queryParams.TaskId,
		},
		PreloadRelations: []string{"Task"},
	}
}

func getStartTime(logs []dto.UserTaskLogResponse) string {
	if isActive(logs) {
		return logs[len(logs)-1].CreatedAt.Format(time.RFC3339)
	}

	return ""
}
