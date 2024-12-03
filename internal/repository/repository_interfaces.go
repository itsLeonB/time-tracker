package repository

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/model"
)

type ProjectRepository interface {
	Insert(project *model.Project) (*model.Project, error)
	GetAll() ([]model.Project, error)
	GetByID(id uuid.UUID) (*model.Project, error)
	Update(project *model.Project) (*model.Project, error)
	Delete(project *model.Project) error
}
