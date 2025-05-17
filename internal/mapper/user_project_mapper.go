package mapper

import (
	"fmt"

	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/util"
)

func ProjectToResponse(project model.UserProject) dto.UserProjectResponse {
	timeSpent := dto.TimeSpent{}
	activeTaskCount := 0
	userTasks := make([]dto.UserTaskResponse, len(project.UserTasks))

	for i, userTask := range project.UserTasks {
		userTaskResponse := UserTaskToResponse(userTask)
		timeSpent.Add(userTaskResponse.TimeSpent)

		if userTaskResponse.IsActive {
			activeTaskCount++
		}

		userTasks[i] = userTaskResponse
	}

	avgHours := timeSpent.Hours / float64(len(userTasks))

	return dto.UserProjectResponse{
		ID:              project.ID,
		UserId:          project.UserId,
		Name:            project.Name,
		CreatedAt:       project.CreatedAt,
		UpdatedAt:       project.UpdatedAt,
		TimeSpent:       timeSpent,
		UserTasks:       userTasks,
		ActiveTaskCount: activeTaskCount,
		AverageTaskTime: fmt.Sprintf("%.1f", avgHours),
	}
}

func PopulateProjectWithLogs(project model.UserProject, userTaskLogs []dto.UserTaskLogResponse) (dto.UserProjectResponse, error) {
	timeSpent := dto.TimeSpent{}
	activeTaskCount := 0
	userTasks := make([]dto.UserTaskResponse, len(project.UserTasks))
	logsByUserTaskId, err := util.GroupBy(userTaskLogs, "UserTaskId")
	if err != nil {
		return dto.UserProjectResponse{}, err
	}

	for i, userTask := range project.UserTasks {
		populatedTask := PopulateUserTaskWithLogs(userTask, logsByUserTaskId[userTask.ID])
		timeSpent.Add(populatedTask.TimeSpent)

		if populatedTask.IsActive {
			activeTaskCount++
		}

		userTasks[i] = populatedTask
	}

	avgHours := timeSpent.Hours / float64(len(userTasks))

	projectResponse := dto.UserProjectResponse{
		ID:              project.ID,
		UserId:          project.UserId,
		Name:            project.Name,
		CreatedAt:       project.CreatedAt,
		UpdatedAt:       project.UpdatedAt,
		TimeSpent:       timeSpent,
		UserTasks:       userTasks,
		ActiveTaskCount: activeTaskCount,
		AverageTaskTime: fmt.Sprintf("%.1f", avgHours),
	}

	return projectResponse, nil
}
