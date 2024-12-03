package provider

import "github.com/itsLeonB/time-tracker/internal/delivery/http/handler"

type Handlers struct {
	Project *handler.ProjectHandler
}

func ProvideHandlers(services *Services) *Handlers {
	return &Handlers{
		Project: handler.NewProjectHandler(services.Project),
	}
}
