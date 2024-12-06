package provider

import "github.com/itsLeonB/time-tracker/internal/service"

type Services struct {
	Project service.ProjectService
	Task    service.TaskService
}

func ProvideServices(repositories *Repositories) *Services {
	return &Services{
		Project: service.NewProjectService(repositories.Project),
		Task:    service.NewTaskService(repositories.Task),
	}
}
