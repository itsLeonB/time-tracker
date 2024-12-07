package service

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/model"
)

type ProjectService interface {
	Create(name string) (*model.Project, error)
	GetAll() ([]*model.Project, error)
	GetByID(id uuid.UUID) (*model.Project, error)
	Update(id uuid.UUID, name string) (*model.Project, error)
	Delete(id uuid.UUID) error
	FirstByQuery(options model.FindProjectOptions) (*model.Project, error)
}

type TaskService interface {
	Create(request *model.NewTaskRequest) (*model.Task, error)
	GetAll() ([]*model.Task, error)
	GetByID(id uuid.UUID) (*model.Task, error)
	GetByNumber(number string) (*model.Task, error)
	Update(id uuid.UUID, name string) (*model.Task, error)
	Delete(id uuid.UUID) error
	Log(id uuid.UUID, action string) (*model.TaskLog, error)
	LogByNumber(number string, action string) (*model.TaskLog, error)
	Find(options *model.QueryOptions) ([]*model.Task, error)
}
