package provider

import (
	"github.com/itsLeonB/catfeinated-time-tracker/internal/repository"
	"gorm.io/gorm"
)

type Repositories struct {
	User    repository.UserRepository
	Project repository.ProjectRepository
	Task    repository.TaskRepository
}

func ProvideRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:    repository.NewUserRepository(db),
		Project: repository.NewProjectRepository(db),
		Task:    repository.NewTaskRepository(db),
	}
}
