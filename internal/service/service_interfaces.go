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
	FirstByQuery(ctx context.Context, options dto.ProjectQueryParams) (dto.ProjectResponse, error)
	GetOrCreate(ctx context.Context, name string) (dto.ProjectResponse, error)
	FindByUserId(ctx context.Context, userId uuid.UUID) ([]dto.ProjectResponse, error)
}

type TaskService interface {
	Create(ctx context.Context, request *dto.NewTaskRequest) (dto.TaskResponse, error)
	GetAll(ctx context.Context) ([]dto.TaskResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (dto.TaskResponse, error)
	GetByNumber(ctx context.Context, number string) (dto.TaskResponse, error)
	Update(ctx context.Context, id uuid.UUID, name string) (dto.TaskResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Find(ctx context.Context, queryParams dto.TaskQueryParams) ([]dto.TaskResponse, error)
}

type ExternalTrackerService interface {
	FindTask(ctx context.Context, queryOptions dto.ExternalQueryOptions) ([]model.ExternalTask, error)
}

type UserTaskService interface {
	Create(ctx context.Context, request dto.NewUserTaskRequest) (dto.UserTaskResponse, error)
	FindAll(ctx context.Context, options dto.UserTaskQueryParams) ([]dto.UserTaskResponse, error)
	GetById(ctx context.Context, id uuid.UUID) (dto.UserTaskResponse, error)
}

type UserTaskLogService interface {
	Create(ctx context.Context, request dto.NewUserTaskLogRequest) (dto.UserTaskLogResponse, error)
	FindAll(ctx context.Context, params dto.UserTaskLogParams) ([]dto.UserTaskLogResponse, error)
}
