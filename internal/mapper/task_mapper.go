package mapper

import (
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/util"
)

func TaskToResponse(task model.Task) dto.TaskResponse {
	timeSpent := dto.TimeSpent{}
	userTasks := make([]dto.UserTaskResponse, len(task.UserTasks))

	for i, userTask := range task.UserTasks {
		userTasks[i] = UserTaskToResponse(userTask)
		timeSpent.Add(userTasks[i].TimeSpent)
	}

	return dto.TaskResponse{
		ID:        task.ID,
		ProjectID: task.ProjectID,
		Number:    task.Number,
		Name:      task.Name,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
		TimeSpent: timeSpent,
		UserTasks: userTasks,
	}
}

func PopulateTaskWithLogs(task model.Task, userTaskLogs []dto.UserTaskLogResponse) (dto.TaskResponse, error) {
	timeSpent := dto.TimeSpent{}
	userTasks := make([]dto.UserTaskResponse, len(task.UserTasks))

	logsByUserTaskId, err := util.GroupBy(userTaskLogs, "UserTaskId")
	if err != nil {
		return dto.TaskResponse{}, err
	}

	for i, userTask := range task.UserTasks {
		userTasks[i] = PopulateUserTaskWithLogs(userTask, logsByUserTaskId[userTask.ID])
		timeSpent.Add(userTasks[i].TimeSpent)
	}

	taskResponse := dto.TaskResponse{
		ID:        task.ID,
		ProjectID: task.ProjectID,
		Number:    task.Number,
		Name:      task.Name,
		CreatedAt: task.CreatedAt,
		UpdatedAt: task.UpdatedAt,
		TimeSpent: timeSpent,
		UserTasks: userTasks,
	}

	return taskResponse, nil
}
