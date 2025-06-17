package service

import (
	"context"

	"github.com/itsLeonB/ezutil"
	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/itsLeonB/time-tracker/internal/entity"
	"github.com/itsLeonB/time-tracker/internal/mapper"
	"github.com/itsLeonB/time-tracker/internal/repository"
	"github.com/rotisserie/eris"
)

type taskServiceImpl struct {
	taskRepository repository.TaskRepository
	userService    UserService
}

func NewTaskService(
	taskRepository repository.TaskRepository,
	userService UserService,
) TaskService {
	return &taskServiceImpl{
		taskRepository,
		userService,
	}
}

func (ts *taskServiceImpl) Create(ctx context.Context, request dto.NewTaskRequest) (dto.TaskResponse, error) {
	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	spec := entity.Task{
		ProjectID: request.ProjectID,
		Number:    request.Number,
	}

	task, err := ts.taskRepository.FindFirst(ctx, spec)
	if err != nil {
		return dto.TaskResponse{}, err
	}
	if !task.IsZero() {
		return dto.TaskResponse{}, ezutil.ConflictError(eris.Errorf("task number: %s already exists", request.Number))
	}

	newTask := entity.Task{
		ProjectID: request.ProjectID,
		Number:    request.Number,
		Name:      request.Name,
	}

	insertedTask, err := ts.taskRepository.Insert(ctx, newTask)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	return mapper.TaskToResponse(insertedTask), nil
}

func (ts *taskServiceImpl) GetAll(ctx context.Context) ([]dto.TaskResponse, error) {
	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	tasks, err := ts.taskRepository.FindAll(ctx, entity.Task{})
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(tasks, mapper.TaskToResponse), nil
}

func (ts *taskServiceImpl) Find(ctx context.Context, params dto.TaskQueryParams) ([]dto.TaskResponse, error) {
	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	spec := entity.Task{Number: params.Number}

	tasks, err := ts.taskRepository.FindAll(ctx, spec)
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(tasks, mapper.TaskToResponse), nil
}
