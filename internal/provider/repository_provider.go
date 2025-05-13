package provider

import (
	"github.com/itsLeonB/catfeinated-time-tracker/internal/config"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/repository"
	"gorm.io/gorm"
)

type Repositories struct {
	User        repository.UserRepository
	Project     repository.ProjectRepository
	Task        repository.TaskRepository
	Youtrack    *repository.YoutrackRepository
	UserTask    repository.UserTaskRepository
	UserTaskLog repository.UserTaskLogRepository
}

func ProvideRepositories(
	db *gorm.DB,
	configs *config.Config,
) *Repositories {
	return &Repositories{
		User:        repository.NewUserRepository(db),
		Project:     repository.NewProjectRepository(db),
		Task:        repository.NewTaskRepository(db),
		Youtrack:    repository.NewYoutrackRepository(configs.Youtrack),
		UserTask:    repository.NewUserTaskRepository(db),
		UserTaskLog: repository.NewUserTaskLogRepository(db),
	}
}
