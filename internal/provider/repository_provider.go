package provider

import (
	"github.com/itsLeonB/time-tracker/internal/repository"
	"gorm.io/gorm"
)

type Repositories struct {
	Project repository.ProjectRepository
	Task    repository.TaskRepository
}

func ProvideRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Project: repository.NewProjectRepository(db),
		Task:    repository.NewTaskRepository(db),
	}
}
