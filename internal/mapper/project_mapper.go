package mapper

import (
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
)

func ProjectToResponse(project model.Project) dto.ProjectResponse {
	timeSpent := dto.TimeSpent{}
	tasks := make([]dto.TaskResponse, len(project.Tasks))

	for i, task := range project.Tasks {
		tasks[i] = TaskToResponse(task)
		timeSpent.Add(tasks[i].TimeSpent)
	}

	return dto.ProjectResponse{
		ID:        project.ID,
		Name:      project.Name,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
		TimeSpent: timeSpent,
		Tasks:     tasks,
	}
}

func PopulateProjectWithLogs(project model.Project, userTaskLogs []dto.UserTaskLogResponse) (dto.ProjectResponse, error) {
	timeSpent := dto.TimeSpent{}
	tasks := make([]dto.TaskResponse, len(project.Tasks))

	for i, task := range project.Tasks {
		populatedTask, err := PopulateTaskWithLogs(task, userTaskLogs)
		if err != nil {
			return dto.ProjectResponse{}, err
		}

		tasks[i] = populatedTask
		timeSpent.Add(tasks[i].TimeSpent)
	}

	projectResponse := dto.ProjectResponse{
		ID:        project.ID,
		Name:      project.Name,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
		TimeSpent: timeSpent,
		Tasks:     tasks,
	}

	return projectResponse, nil
}
