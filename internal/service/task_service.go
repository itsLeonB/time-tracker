package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/constant"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/repository"
	strategy "github.com/itsLeonB/catfeinated-time-tracker/internal/service/strategy/point"
	"github.com/rotisserie/eris"
)

type taskServiceImpl struct {
	taskRepository repository.TaskRepository
	pointStrategy  strategy.PointStrategy
	userService    UserService
}

func NewTaskService(
	taskRepository repository.TaskRepository,
	pointStrategy strategy.PointStrategy,
	userService UserService,
) TaskService {
	return &taskServiceImpl{taskRepository, pointStrategy, userService}
}

func (ts *taskServiceImpl) Create(ctx context.Context, request *model.NewTaskRequest) (*model.Task, error) {
	user, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	task, err := ts.taskRepository.GetByNumber(ctx, request.Number)
	if err != nil {
		return nil, err
	}
	if task != nil {
		return nil, apperror.ConflictError(
			eris.Errorf("task number: %s already exists", request.Number),
		)
	}

	newTask := &model.Task{
		UserID:    user.ID,
		ProjectID: request.ProjectID,
		Number:    request.Number,
		Name:      request.Name,
	}

	return ts.taskRepository.Insert(ctx, newTask)
}

func (ts *taskServiceImpl) GetAll(ctx context.Context) ([]*model.Task, error) {
	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	return ts.taskRepository.GetAll(ctx)
}

func (ts *taskServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	return ts.taskRepository.GetByID(ctx, id)
}

func (ts *taskServiceImpl) Find(ctx context.Context, options *model.QueryOptions) ([]*model.Task, error) {
	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	tasks, err := ts.taskRepository.Find(ctx, options)
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		ts.populateAdditionalTaskFields(task, options)
	}

	return tasks, nil
}

func (ts *taskServiceImpl) GetByNumber(ctx context.Context, number string) (*model.Task, error) {
	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	task, err := ts.getTaskByNumber(ctx, number)
	if err != nil {
		return nil, err
	}

	ts.populateAdditionalTaskFields(task, nil)

	return task, nil
}

func (ts *taskServiceImpl) populateAdditionalTaskFields(
	task *model.Task,
	options *model.QueryOptions,
) {
	task.DetermineProgress()
	task.CalculateTotalTime()
	task.Points = ts.pointStrategy.CalculatePoints(task)

	if options != nil && !options.WithLogs {
		task.Logs = nil
	}
}

func (ts *taskServiceImpl) Update(ctx context.Context, id uuid.UUID, name string) (*model.Task, error) {
	_, err := ts.userService.ValidateUser(ctx)
	if err != nil {
		return nil, err
	}

	task, err := ts.getTask(ctx, id)
	if err != nil {
		return nil, err
	}

	task.Name = name

	return ts.taskRepository.Update(ctx, task)
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
