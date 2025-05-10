package mapper

import (
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
)

func ProjectToResponse(project model.Project) dto.ProjectResponse {
	return dto.ProjectResponse{
		ID:        project.ID,
		Name:      project.Name,
		CreatedAt: project.CreatedAt,
		UpdatedAt: project.UpdatedAt,
		TimeSpent: dto.TimeSpent{},
	}
}
