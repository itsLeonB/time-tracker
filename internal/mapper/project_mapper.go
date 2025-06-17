package mapper

import (
	"fmt"

	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/itsLeonB/time-tracker/internal/entity"
)

func ProjectToResponse(project entity.Project) (dto.ProjectResponse, error) {
	timeSpent := dto.TimeSpent{}
	activeTaskCount := 0
	var avgHours float64 = 0
	userTasks := make([]dto.TaskResponse, len(project.Tasks))

	if len(project.Tasks) > 0 {
		for i, task := range project.Tasks {
			taskResponse, err := TaskToResponseWithError(task)
			if err != nil {
				return dto.ProjectResponse{}, err
			}

			timeSpent.Add(taskResponse.TimeSpent)

			if taskResponse.IsActive {
				activeTaskCount++
			}

			userTasks[i] = taskResponse
		}

		avgHours = timeSpent.Hours / float64(len(userTasks))
	}

	userProjectResponse := dto.ProjectResponse{
		ID:              project.ID,
		UserId:          project.UserID,
		Name:            project.Name,
		CreatedAt:       project.CreatedAt,
		UpdatedAt:       project.UpdatedAt,
		TimeSpent:       timeSpent,
		Tasks:           userTasks,
		ActiveTaskCount: activeTaskCount,
		AverageTaskTime: fmt.Sprintf("%.1f", avgHours),
	}

	return userProjectResponse, nil
}
