package provider

import (
	"github.com/itsLeonB/ezutil"
	"github.com/itsLeonB/time-tracker/internal/service"
)

type Services struct {
	JWT     ezutil.JWTService
	User    service.UserService
	Auth    service.AuthService
	Project service.ProjectService
	Task    service.TaskService
	TaskLog service.TaskLogService
}

func ProvideServices(configs *ezutil.Config, repositories *Repositories) *Services {
	userService := service.NewUserService(repositories.User)
	hashService := ezutil.NewHashService(10)
	jwtService := ezutil.NewJwtService(configs.Auth)
	authService := service.NewAuthService(hashService, jwtService, repositories.User, repositories.Transactor)
	userTaskLogService := service.NewTaskLogService(repositories.UserTaskLog)
	taskService := service.NewTaskService(repositories.Task, userService)
	projectService := service.NewProjectService(repositories.Project, userService)

	return &Services{
		JWT:     jwtService,
		User:    userService,
		Auth:    authService,
		Project: projectService,
		Task:    taskService,
		TaskLog: userTaskLogService,
	}
}
