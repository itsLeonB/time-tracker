package service

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/apperror"
	"github.com/itsLeonB/time-tracker/internal/constant"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/itsLeonB/time-tracker/internal/repository"
	strategy "github.com/itsLeonB/time-tracker/internal/service/strategy/point"
	"github.com/rotisserie/eris"
)

type taskServiceImpl struct {
	taskRepository repository.TaskRepository
	pointStrategy  strategy.PointStrategy
}

func NewTaskService(
	taskRepository repository.TaskRepository,
	pointStrategy strategy.PointStrategy,
) TaskService {
	return &taskServiceImpl{taskRepository, pointStrategy}
}

func (ts *taskServiceImpl) Create(request *model.NewTaskRequest) (*model.Task, error) {
	task, err := ts.taskRepository.GetByNumber(request.Number)
	if err != nil {
		return nil, err
	}
	if task != nil {
		return nil, apperror.ConflictError(
			eris.Errorf("task number: %s already exists", request.Number),
		)
	}

	return ts.taskRepository.Insert(&model.Task{
		ProjectID: request.ProjectID,
		Number:    request.Number,
		Name:      request.Name,
	})
}

func (ts *taskServiceImpl) GetAll() ([]*model.Task, error) {
	return ts.taskRepository.GetAll()
}

func (ts *taskServiceImpl) GetByID(id uuid.UUID) (*model.Task, error) {
	return ts.taskRepository.GetByID(id)
}

func (ts *taskServiceImpl) Find(options *model.QueryOptions) ([]*model.Task, error) {
	tasks, err := ts.taskRepository.Find(options)
	if err != nil {
		return nil, err
	}

	for _, task := range tasks {
		ts.populateAdditionalTaskFields(task)
	}

	return tasks, nil
}

func (ts *taskServiceImpl) GetByNumber(number string) (*model.Task, error) {
	task, err := ts.taskRepository.GetByNumber(number)
	if err != nil {
		return nil, err
	}

	ts.populateAdditionalTaskFields(task)

	return task, nil
}

func (ts *taskServiceImpl) populateAdditionalTaskFields(task *model.Task) {
	task.CalculateTotalTime()
	task.Points = ts.pointStrategy.CalculatePoints(task)
}

func (ts *taskServiceImpl) Update(id uuid.UUID, name string) (*model.Task, error) {
	return ts.taskRepository.Update(&model.Task{
		ID:   id,
		Name: name,
	})
}

func (ts *taskServiceImpl) Delete(id uuid.UUID) error {
	return ts.taskRepository.Delete(&model.Task{ID: id})
}

func (ts *taskServiceImpl) Log(id uuid.UUID, action string) (*model.TaskLog, error) {
	task, err := ts.taskRepository.GetByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, apperror.NotFoundError(eris.Errorf("task ID: %s is not found", id))
	}

	return ts.createLog(task, action)
}

func (ts *taskServiceImpl) LogByNumber(number string, action string) (*model.TaskLog, error) {
	task, err := ts.taskRepository.GetByNumber(number)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, apperror.NotFoundError(eris.Errorf("task number: %s is not found", number))
	}

	return ts.createLog(task, action)
}

func (ts *taskServiceImpl) createLog(task *model.Task, action string) (*model.TaskLog, error) {
	latestLog, err := ts.taskRepository.GetLatestLog(task)
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

	log, err := ts.taskRepository.Log(task, action)
	if err != nil {
		return nil, err
	}

	return log, nil
}
