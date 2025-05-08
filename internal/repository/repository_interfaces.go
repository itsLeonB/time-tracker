package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Insert(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
}

type ProjectRepository interface {
	Insert(ctx context.Context, project *model.Project) (*model.Project, error)
	GetAll(ctx context.Context) ([]*model.Project, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Project, error)
	Update(ctx context.Context, project *model.Project) (*model.Project, error)
	Delete(ctx context.Context, project *model.Project) error
	Find(ctx context.Context, options *model.FindProjectOptions) ([]*model.Project, error)
	GetByName(ctx context.Context, name string) (*model.Project, error)
}

type TaskRepository interface {
	Insert(ctx context.Context, task *model.Task) (*model.Task, error)
	GetAll(ctx context.Context) ([]*model.Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error)
	GetByNumber(ctx context.Context, number string) (*model.Task, error)
	Update(ctx context.Context, task *model.Task) (*model.Task, error)
	Delete(ctx context.Context, task *model.Task) error
	Log(ctx context.Context, task *model.Task, action string) (*model.TaskLog, error)
	GetLatestLog(ctx context.Context, task *model.Task) (*model.TaskLog, error)
	GetLogs(ctx context.Context, task *model.Task) ([]*model.TaskLog, error)
	Find(ctx context.Context, options *model.QueryOptions) ([]*model.Task, error)
	GetInProgress(ctx context.Context, projectID uuid.UUID) ([]*model.Task, error)
}
