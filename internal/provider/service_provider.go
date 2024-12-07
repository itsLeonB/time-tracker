package provider

import (
	"github.com/itsLeonB/time-tracker/internal/service"
	strategy "github.com/itsLeonB/time-tracker/internal/service/strategy/point"
)

type Services struct {
	Project service.ProjectService
	Task    service.TaskService
}

func ProvideServices(repositories *Repositories) *Services {
	pointStrategy := strategy.NewHourBasedPointStrategy()

	return &Services{
		Project: service.NewProjectService(repositories.Project),
		Task:    service.NewTaskService(repositories.Task, pointStrategy),
	}
}
