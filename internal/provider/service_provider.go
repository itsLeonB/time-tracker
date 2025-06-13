package provider

import (
	"github.com/itsLeonB/time-tracker/internal/auth"
	"github.com/itsLeonB/time-tracker/internal/config"
	"github.com/itsLeonB/time-tracker/internal/service"
)

type Services struct {
	User            service.UserService
	Hasher          auth.Hasher
	JWT             auth.JWT
	Auth            auth.AuthService
	Project         service.ProjectService
	Task            service.TaskService
	ExternalTracker service.ExternalTrackerService
	UserTask        service.UserTaskService
	UserTaskLog     service.UserTaskLogService
}

func ProvideServices(configs *config.Config, repositories *Repositories) *Services {
	userService := service.NewUserService(repositories.User)
	hasher := auth.NewHasherBcrypt(10)
	jwt := auth.NewJWTHS256(configs.Auth)
	authService := auth.NewAuthService(hasher, jwt, userService)
	externalTrackerService := service.NewYoutrackService(repositories.Youtrack)
	userTaskLogService := service.NewUserTaskLogService(repositories.UserTaskLog)
	taskService := service.NewTaskService(repositories.Task, userService, externalTrackerService)
	userTaskService := service.NewUserTaskService(repositories.UserTask, taskService)
	projectService := service.NewProjectService(repositories.Project, userService, userTaskService, userTaskLogService)

	return &Services{
		User:            userService,
		Hasher:          hasher,
		JWT:             jwt,
		Auth:            authService,
		Project:         projectService,
		Task:            taskService,
		ExternalTracker: externalTrackerService,
		UserTask:        userTaskService,
		UserTaskLog:     userTaskLogService,
	}
}
