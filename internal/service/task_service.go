package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/mapper"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/repository"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/util"
	"github.com/rotisserie/eris"
)

type taskServiceImpl struct {
	taskRepository         repository.TaskRepository
	userService            UserService
	externalTrackerService ExternalTrackerService
}

func NewTaskService(
	taskRepository repository.TaskRepository,
	userService UserService,
	externalTrackerService ExternalTrackerService,
) TaskService {
	return &taskServiceImpl{
		taskRepository,
		userService,
		externalTrackerService,
	}
}

func (ts *taskServiceImpl) Create(ctx context.Context, request *dto.NewTaskRequest) (dto.TaskResponse, error) {
	var taskResponse dto.TaskResponse

	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return taskResponse, err
	}

	task, err := ts.taskRepository.GetByNumber(ctx, request.Number)
	if err != nil {
		return taskResponse, err
	}
	if task != nil {
		return taskResponse, apperror.ConflictError(eris.Errorf("task number: %s already exists", request.Number))
	}

	newTask := &model.Task{
		Number: request.Number,
		Name:   request.Name,
	}

	insertedTask, err := ts.taskRepository.Insert(ctx, newTask)
	if err != nil {
		return taskResponse, err
	}

	return mapper.TaskToResponse(*insertedTask), nil
}

func (ts *taskServiceImpl) GetAll(ctx context.Context) ([]dto.TaskResponse, error) {
	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	tasks, err := ts.taskRepository.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return util.MapSlice(tasks, mapper.TaskToResponse), nil
}

func (ts *taskServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (dto.TaskResponse, error) {
	var taskResponse dto.TaskResponse

	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return taskResponse, err
	}

	task, err := ts.taskRepository.GetByID(ctx, id)
	if err != nil {
		return taskResponse, err
	}

	return mapper.TaskToResponse(*task), nil
}

func (ts *taskServiceImpl) Find(ctx context.Context, queryParams dto.TaskQueryParams) ([]dto.TaskResponse, error) {
	tasks, _, err := ts.findTask(ctx, queryParams)
	if err != nil {
		return nil, err
	}

	if len(tasks) > 0 {
		return util.MapSlice(tasks, mapper.TaskToResponse), nil
	}

	// If no tasks are found, check the external tracker
	return ts.findAndInsertFromExternal(ctx, queryParams.Number)
}

func (ts *taskServiceImpl) findAndInsertFromExternal(ctx context.Context, number string) ([]dto.TaskResponse, error) {
	externalQueryOptions := dto.ExternalQueryOptions{
		Number: number,
	}

	externalTasks, err := ts.externalTrackerService.FindTask(ctx, externalQueryOptions)
	if err != nil {
		return nil, err
	}

	insertedTasks := make([]dto.TaskResponse, 0, len(externalTasks))

	for _, externalTask := range externalTasks {
		insertedTask, err := ts.insertFromExternal(ctx, externalTask)
		if err != nil {
			return nil, err
		}

		insertedTasks = append(insertedTasks, mapper.TaskToResponse(insertedTask))
	}

	return insertedTasks, nil
}

func (ts *taskServiceImpl) insertFromExternal(ctx context.Context, externalTask model.ExternalTask) (model.Task, error) {
	newTask := model.Task{
		Number: externalTask.Number,
		Name:   externalTask.Name,
	}

	insertedTask, err := ts.taskRepository.Insert(ctx, &newTask)
	if err != nil {
		return newTask, err
	}

	return *insertedTask, nil
}

func (ts *taskServiceImpl) GetByNumber(ctx context.Context, number string) (dto.TaskResponse, error) {
	var taskResponse dto.TaskResponse

	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return taskResponse, err
	}

	task, err := ts.getTaskByNumber(ctx, number)
	if err != nil {
		return taskResponse, err
	}

	return mapper.TaskToResponse(*task), nil
}

func (ts *taskServiceImpl) Update(ctx context.Context, id uuid.UUID, name string) (dto.TaskResponse, error) {
	var taskResponse dto.TaskResponse

	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return taskResponse, err
	}

	task, err := ts.getTask(ctx, id)
	if err != nil {
		return taskResponse, err
	}

	task.Name = name

	updatedTask, err := ts.taskRepository.Update(ctx, task)
	if err != nil {
		return taskResponse, err
	}

	return mapper.TaskToResponse(*updatedTask), nil
}

func (ts *taskServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return err
	}

	task, err := ts.getTask(ctx, id)
	if err != nil {
		return err
	}

	return ts.taskRepository.Delete(ctx, task)
}

func (ts *taskServiceImpl) getTask(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	task, err := ts.taskRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, apperror.NotFoundError(eris.Errorf("task ID: %s is not found", id))
	}

	return task, nil
}

func (ts *taskServiceImpl) getTaskByNumber(ctx context.Context, number string) (*model.Task, error) {
	task, err := ts.taskRepository.GetByNumber(ctx, number)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, apperror.NotFoundError(eris.Errorf("task number: %s is not found", number))
	}

	return task, nil
}

func (ts *taskServiceImpl) findTask(ctx context.Context, params dto.TaskQueryParams) ([]model.Task, model.User, error) {
	user, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return nil, model.User{}, err
	}

	queryOptions := model.TaskQueryOptions{
		Number: params.Number,
		Date:   params.Date,
	}

	tasks, err := ts.taskRepository.Find(ctx, queryOptions)
	if err != nil {
		return nil, model.User{}, err
	}

	return tasks, *user, nil
}
