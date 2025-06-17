package provider

import (
	"github.com/itsLeonB/ezutil"
	"github.com/itsLeonB/time-tracker/internal/repository"
)

type Repositories struct {
	Transactor  ezutil.Transactor
	User        repository.UserRepository
	Project     repository.ProjectRepository
	Task        repository.TaskRepository
	UserTaskLog repository.TaskLogRepository
}

func ProvideRepositories(configs *ezutil.Config) *Repositories {
	return &Repositories{
		Transactor:  ezutil.NewTransactor(configs.GORM),
		User:        repository.NewUserRepository(configs.GORM),
		Project:     repository.NewProjectRepository(configs.GORM),
		Task:        repository.NewTaskRepository(configs.GORM),
		UserTaskLog: repository.NewTaskLogRepository(configs.GORM),
	}
}
