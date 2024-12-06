package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/model"
)

type ProjectService interface {
	Create(name string) (*model.Project, error)
	GetAll() ([]model.Project, error)
	GetByID(id uuid.UUID) (*model.Project, error)
	Update(id uuid.UUID, name string) (*model.Project, error)
	Delete(id uuid.UUID) error
}

type TaskService interface {
	Create(projectID uuid.UUID, name string) (*model.Task, error)
	GetAll() ([]model.Task, error)
	GetByID(id uuid.UUID) (*model.Task, error)
	Update(id uuid.UUID, name string) (*model.Task, error)
	Delete(id uuid.UUID) error
	Log(id uuid.UUID, action string) (*model.TaskLog, error)
	GetTotalTime(id uuid.UUID) (*time.Duration, error)
}
