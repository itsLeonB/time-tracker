package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/model"
)

type UserService interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
}

type ProjectService interface {
	Create(name string) (*model.Project, error)
	GetAll() ([]*model.Project, error)
	GetByID(options *model.QueryOptions) (*model.Project, error)
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
