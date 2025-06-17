package mapper

import (
	"time"

	"github.com/itsLeonB/time-tracker/internal/appconstant"
	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/itsLeonB/time-tracker/internal/entity"
	"github.com/rotisserie/eris"
)

func TaskToResponse(task entity.Task) dto.TaskResponse {
	return dto.TaskResponse{
		ID:        task.ID,
		Number:    task.Number,
		Name:      task.Name,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
	}
}

func TaskToResponseWithError(task entity.Task) (dto.TaskResponse, error) {
	logs := make([]dto.TaskLogResponse, len(task.Logs))
	for i, log := range task.Logs {
		logs[i] = TaskLogToResponse(log)
	}

	totalUserTaskTime, err := durationSpentFromLogs(logs)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	userTaskResponse := dto.TaskResponse{
		ID:        task.ID,
		ProjectID: task.ProjectID,
		Number:    task.Number,
		Name:      task.Name,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
		TimeSpent: dto.TimeSpentFromDuration(totalUserTaskTime),
		Logs:      logs,
		IsActive:  isActive(logs),
		StartTime: getStartTime(logs),
	}

	return userTaskResponse, nil
}

func isActive(logs []dto.TaskLogResponse) bool {
	return len(logs) > 0 && logs[0].Action == appconstant.LogAction.Start
}

func durationSpentFromLogs(logs []dto.TaskLogResponse) (time.Duration, error) {
	// Handle empty logs
	if len(logs) == 0 {
		return 0, nil
	}

	// Handle single log cases
	if len(logs) == 1 {
		switch logs[0].Action {
		case appconstant.LogAction.Stop:
			return 0, eris.New("Invalid log: single STOP entry without matching START")
		case appconstant.LogAction.Start:
			return 0, nil
		}
	}

	var totalTime time.Duration
	var pendingStart *time.Time

	// Process from last to first (oldest to newest) since input is descending
	for i := len(logs) - 1; i >= 0; i-- {
		log := logs[i]

		switch log.Action {
		case appconstant.LogAction.Start:
			if pendingStart != nil {
				return 0, eris.New("Invalid log sequence: consecutive START entries")
			}
			pendingStart = &log.CreatedAt

		case appconstant.LogAction.Stop:
			if pendingStart == nil {
				return 0, eris.New("Invalid log sequence: STOP without matching START")
			}
			// Calculate duration from START to STOP
			totalTime += log.CreatedAt.Sub(*pendingStart)
			pendingStart = nil
		}
	}

	// If there's an unpaired START at the end, count time until now
	if pendingStart != nil {
		totalTime += time.Since(*pendingStart)
	}

	return totalTime, nil
}

func getStartTime(logs []dto.TaskLogResponse) string {
	if isActive(logs) {
		return logs[0].CreatedAt.Format(time.RFC3339)
	}

	return ""
}
