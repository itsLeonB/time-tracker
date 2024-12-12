package provider

import (
	"github.com/itsLeonB/time-tracker/internal/auth"
	"github.com/itsLeonB/time-tracker/internal/config"
	"github.com/itsLeonB/time-tracker/internal/service"
	strategy "github.com/itsLeonB/time-tracker/internal/service/strategy/point"
)

type Services struct {
	User    service.UserService
	Hasher  auth.Hasher
	JWT     auth.JWT
	Auth    auth.AuthService
	Project service.ProjectService
	Task    service.TaskService
}

func ProvideServices(configs *config.Config, repositories *Repositories) *Services {
	pointStrategy := strategy.NewHourBasedPointStrategy()

	userService := service.NewUserService(repositories.User)
	hasher := auth.NewHasherBcrypt(10)
	jwt := auth.NewJWTHS256(configs.Auth)
	authService := auth.NewAuthService(hasher, jwt, userService)
	taskService := service.NewTaskService(repositories.Task, pointStrategy)
	projectService := service.NewProjectService(repositories.Project, taskService)

	return &Services{
		User:    userService,
		Hasher:  hasher,
		JWT:     jwt,
		Auth:    authService,
		Project: projectService,
		Task:    taskService,
	}
}
