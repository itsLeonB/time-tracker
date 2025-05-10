package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
)

type UserService interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	ValidateUser(ctx context.Context) (*model.User, error)
}

type ProjectService interface {
	Create(ctx context.Context, name string) (dto.ProjectResponse, error)
	GetAll(ctx context.Context) ([]dto.ProjectResponse, error)
	GetByID(ctx context.Context, options *dto.QueryOptions) (dto.ProjectResponse, error)
	Update(ctx context.Context, id uuid.UUID, name string) (dto.ProjectResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FirstByQuery(ctx context.Context, options *dto.FindProjectOptions) (dto.ProjectResponse, error)
}

type TaskService interface {
	Create(ctx context.Context, request *dto.NewTaskRequest) (*model.Task, error)
	GetAll(ctx context.Context) ([]*model.Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error)
	GetByNumber(ctx context.Context, number string) (*model.Task, error)
	Update(ctx context.Context, id uuid.UUID, name string) (*model.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Log(ctx context.Context, id uuid.UUID, action string) (*model.TaskLog, error)
	LogByNumber(ctx context.Context, number string, action string) (*model.TaskLog, error)
	Find(ctx context.Context, options *dto.QueryOptions) ([]*model.Task, error)
}

type ExternalTrackerService interface {
	FindTask(ctx context.Context, queryOptions dto.ExternalQueryOptions) ([]model.ExternalTask, error)
}
