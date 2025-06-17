package repository

import (
	"context"

	"github.com/itsLeonB/time-tracker/internal/entity"
)

type UserRepository interface {
	Insert(ctx context.Context, user entity.User) (entity.User, error)
	FindFirst(ctx context.Context, spec entity.User) (entity.User, error)
}

type ProjectRepository interface {
	Insert(ctx context.Context, project entity.Project) (entity.Project, error)
	FindAll(ctx context.Context, spec entity.ProjectSpecification) ([]entity.Project, error)
	FindFirst(ctx context.Context, spec entity.ProjectSpecification) (entity.Project, error)
	FindFirstPopulated(ctx context.Context, spec entity.ProjectSpecification) (entity.Project, error)
}

type TaskRepository interface {
	Insert(ctx context.Context, task entity.Task) (entity.Task, error)
	FindAll(ctx context.Context, spec entity.Task) ([]entity.Task, error)
	FindFirst(ctx context.Context, spec entity.Task) (entity.Task, error)
	Update(ctx context.Context, task entity.Task) (entity.Task, error)
	Delete(ctx context.Context, task entity.Task) (entity.Task, error)
}

type TaskLogRepository interface {
	Insert(ctx context.Context, taskLog entity.TaskLog) (entity.TaskLog, error)
	FindLatest(ctx context.Context, spec entity.TaskLog) (entity.TaskLog, error)
}
