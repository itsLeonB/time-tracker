package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Insert(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type ProjectRepository interface {
	Insert(ctx context.Context, project *model.UserProject) (*model.UserProject, error)
	GetAll(ctx context.Context) ([]model.UserProject, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.UserProject, error)
	Update(ctx context.Context, project *model.UserProject) (*model.UserProject, error)
	Delete(ctx context.Context, project *model.UserProject) error
	Find(ctx context.Context, options dto.UserProjectQueryParams) ([]model.UserProject, error)
	GetByName(ctx context.Context, name string) (*model.UserProject, error)
	First(ctx context.Context, options model.ProjectQueryOptions) (model.UserProject, error)
}

type TaskRepository interface {
	Insert(ctx context.Context, task *model.Task) (*model.Task, error)
	GetAll(ctx context.Context) ([]model.Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error)
	GetByNumber(ctx context.Context, number string) (*model.Task, error)
	Update(ctx context.Context, task *model.Task) (*model.Task, error)
	Delete(ctx context.Context, task *model.Task) error
	Find(ctx context.Context, options model.TaskQueryOptions) ([]model.Task, error)
}

type UserTaskRepository interface {
	Insert(ctx context.Context, userTask model.UserTask) (model.UserTask, error)
	FindAll(ctx context.Context, options model.UserTaskQueryOptions) ([]model.UserTask, error)
	FindById(ctx context.Context, id uuid.UUID) (model.UserTask, error)
}

type UserTaskLogRepository interface {
	Insert(ctx context.Context, userTaskLog model.UserTaskLog) (model.UserTaskLog, error)
	FindLatest(ctx context.Context, options model.UserTaskLogQueryOptions) (model.UserTaskLog, error)
	FindAll(ctx context.Context, options model.UserTaskLogQueryOptions) ([]model.UserTaskLog, error)
}
