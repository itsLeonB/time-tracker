package mapper

import (
	"fmt"

	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/itsLeonB/time-tracker/internal/util"
)

func ProjectToResponse(project model.UserProject) (dto.UserProjectResponse, error) {
	timeSpent := dto.TimeSpent{}
	activeTaskCount := 0
	userTasks := make([]dto.UserTaskResponse, len(project.UserTasks))

	for i, userTask := range project.UserTasks {
		userTaskResponse, err := UserTaskToResponse(userTask)
		if err != nil {
			return dto.UserProjectResponse{}, err
		}

		timeSpent.Add(userTaskResponse.TimeSpent)

		if userTaskResponse.IsActive {
			activeTaskCount++
		}

		userTasks[i] = userTaskResponse
	}

	avgHours := timeSpent.Hours / float64(len(userTasks))

	userProjectResponse := dto.UserProjectResponse{
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

	return userProjectResponse, nil
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
		populatedTask, err := PopulateUserTaskWithLogs(userTask, logsByUserTaskId[userTask.ID])
		if err != nil {
			return dto.UserProjectResponse{}, err
		}

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
