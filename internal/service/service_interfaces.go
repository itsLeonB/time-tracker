package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/itsLeonB/time-tracker/internal/entity"
)

type AuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) error
	Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error)
	CheckToken(ctx context.Context, token string) (bool, error)
}

type UserService interface {
	ValidateUser(ctx context.Context) (entity.User, error)
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)
}

type ProjectService interface {
	Create(ctx context.Context, name string) (dto.ProjectResponse, error)
	FirstByQuery(ctx context.Context, params dto.ProjectQueryParams) (dto.ProjectResponse, error)
	FindByUserId(ctx context.Context, userId uuid.UUID) ([]dto.ProjectResponse, error)
}

type TaskService interface {
	Create(ctx context.Context, request dto.NewTaskRequest) (dto.TaskResponse, error)
	GetAll(ctx context.Context) ([]dto.TaskResponse, error)
	Find(ctx context.Context, params dto.TaskQueryParams) ([]dto.TaskResponse, error)
}

type TaskLogService interface {
	Create(ctx context.Context, request dto.NewTaskLogRequest) (dto.TaskLogResponse, error)
}
