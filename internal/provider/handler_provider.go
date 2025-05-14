package provider

import "github.com/itsLeonB/catfeinated-time-tracker/internal/delivery/http/handler"

type Handlers struct {
	Root     *handler.RootHandler
	Auth     *handler.AuthHandler
	Project  *handler.ProjectHandler
	Task     *handler.TaskHandler
	UserTask *handler.UserTaskHandler
	Home     *handler.HomeHandler
}

func ProvideHandlers(services *Services) *Handlers {
	return &Handlers{
		Root:     &handler.RootHandler{},
		Auth:     handler.NewAuthHandler(services.Auth),
		Project:  handler.NewProjectHandler(services.Project, services.User),
		Task:     handler.NewTaskHandler(services.Task, services.ExternalTracker),
		UserTask: handler.NewUserTaskHandler(services.UserTask, services.UserTaskLog),
		Home:     handler.NewHomeHandler(services.User, services.Project),
	}
}
