package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/constant"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/mapper"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/repository"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/util"
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
	return &taskServiceImpl{taskRepository, userService}
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
		ProjectID: request.ProjectID,
		Number:    request.Number,
		Name:      request.Name,
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

func (ts *taskServiceImpl) Find(ctx context.Context, options *dto.QueryOptions) ([]dto.TaskResponse, error) {
	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	tasks, err := ts.taskRepository.Find(ctx, options)
	if err != nil {
		return nil, err
	}

	return util.MapSlice(tasks, mapper.TaskToResponse), nil
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

func (ts *taskServiceImpl) Log(ctx context.Context, id uuid.UUID, action string) (*model.TaskLog, error) {
	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	task, err := ts.getTask(ctx, id)
	if err != nil {
		return nil, err
	}

	return ts.createLog(ctx, task, action)
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

func (ts *taskServiceImpl) LogByNumber(ctx context.Context, number string, action string) (*model.TaskLog, error) {
	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	task, err := ts.getTaskByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	return ts.createLog(ctx, task, action)
}

func (ts *taskServiceImpl) createLog(ctx context.Context, task *model.Task, action string) (*model.TaskLog, error) {
	latestLog, err := ts.taskRepository.GetLatestLog(ctx, task)
	if err != nil {
		return nil, err
	}

	// should always start with START action
	if latestLog == nil && action == constant.LogAction.Stop {
		return nil, apperror.BadRequestError(
			eris.Errorf("task ID: %s is not yet started", task.ID),
		)
	}

	// should alternate actions
	if latestLog != nil && latestLog.Action == action {
		return nil, apperror.BadRequestError(
			eris.Errorf("task ID: %s is already %s", task.ID, action),
		)
	}

	log, err := ts.taskRepository.Log(ctx, task, action)
	if err != nil {
		return nil, err
	}

	return log, nil
}
