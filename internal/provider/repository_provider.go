package provider

import (
	"github.com/itsLeonB/time-tracker/internal/repository"
	"gorm.io/gorm"
)

type Repositories struct {
	Project repository.ProjectRepository
}

func ProvideRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Project: repository.NewProjectRepository(db),
	}
}
