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

	taskService := service.NewTaskService(repositories.Task, pointStrategy)
	projectService := service.NewProjectService(repositories.Project, taskService)

	return &Services{
		Project: projectService,
		Task:    taskService,
	}
}
