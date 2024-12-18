package provider

import "github.com/itsLeonB/time-tracker/internal/delivery/http/handler"

type Handlers struct {
	Root    *handler.RootHandler
	Auth    *handler.AuthHandler
	Project *handler.ProjectHandler
	Task    *handler.TaskHandler
}

func ProvideHandlers(services *Services) *Handlers {
	return &Handlers{
		Root:    &handler.RootHandler{},
		Auth:    handler.NewAuthHandler(services.Auth),
		Project: handler.NewProjectHandler(services.Project),
		Task:    handler.NewTaskHandler(services.Task),
	}
}
