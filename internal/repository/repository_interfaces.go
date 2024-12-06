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

type TaskRepository interface {
	Insert(task *model.Task) (*model.Task, error)
	GetAll() ([]model.Task, error)
	GetByID(id uuid.UUID) (*model.Task, error)
	Update(task *model.Task) (*model.Task, error)
	Delete(task *model.Task) error
	Log(task *model.Task, action string) (*model.TaskLog, error)
	GetLatestLog(task *model.Task) (*model.TaskLog, error)
	GetLogs(task *model.Task) ([]model.TaskLog, error)
}
