package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/model"
)

type UserService interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	ValidateUser(ctx context.Context) (*model.User, error)
}

type ProjectService interface {
	Create(ctx context.Context, name string) (*model.Project, error)
	GetAll(ctx context.Context) ([]*model.Project, error)
	GetByID(ctx context.Context, options *model.QueryOptions) (*model.Project, error)
	Update(ctx context.Context, id uuid.UUID, name string) (*model.Project, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FirstByQuery(ctx context.Context, options *model.FindProjectOptions) (*model.Project, error)
	GetInProgressTasks(ctx context.Context, id uuid.UUID) ([]*model.Task, error)
}

type TaskService interface {
	Create(ctx context.Context, request *model.NewTaskRequest) (*model.Task, error)
	GetAll(ctx context.Context) ([]*model.Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error)
	GetByNumber(ctx context.Context, number string) (*model.Task, error)
	Update(ctx context.Context, id uuid.UUID, name string) (*model.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Log(ctx context.Context, id uuid.UUID, action string) (*model.TaskLog, error)
	LogByNumber(ctx context.Context, number string, action string) (*model.TaskLog, error)
	Find(ctx context.Context, options *model.QueryOptions) ([]*model.Task, error)
}
