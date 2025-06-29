package provider

import "github.com/itsLeonB/time-tracker/internal/delivery/http/handler"

type Handlers struct {
	Root    *handler.RootHandler
	Auth    *handler.AuthHandler
	Project *handler.ProjectHandler
	Task    *handler.TaskHandler
	Home    *handler.HomeHandler
	TaskLog *handler.TaskLogHandler
}

func ProvideHandlers(services *Services) *Handlers {
	return &Handlers{
		Root:    &handler.RootHandler{},
		Auth:    handler.NewAuthHandler(services.Auth),
		Project: handler.NewProjectHandler(services.Project, services.User),
		Task:    handler.NewTaskHandler(services.Task),
		Home:    handler.NewHomeHandler(services.User, services.Project),
		TaskLog: handler.NewTaskLogHandler(services.TaskLog),
	}
}
