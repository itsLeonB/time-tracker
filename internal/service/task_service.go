package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/constant"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/itsLeonB/time-tracker/internal/repository"
	"github.com/rotisserie/eris"
)

type taskServiceImpl struct {
	taskRepository repository.TaskRepository
}

func NewTaskService(taskRepository repository.TaskRepository) TaskService {
	return &taskServiceImpl{taskRepository}
}

func (ts *taskServiceImpl) Create(projectID uuid.UUID, name string) (*model.Task, error) {
	return ts.taskRepository.Insert(&model.Task{
		ProjectID: projectID,
		Name:      name,
	})
}

func (ts *taskServiceImpl) GetAll() ([]model.Task, error) {
	return ts.taskRepository.GetAll()
}

func (ts *taskServiceImpl) GetByID(id uuid.UUID) (*model.Task, error) {
	return ts.taskRepository.GetByID(id)
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

	latestLog, err := ts.taskRepository.GetLatestLog(task)
	if err != nil {
		return nil, err
	}

	// should always start with START action
	if latestLog == nil && action == constant.LOG_ACTION.STOP {
		return nil, eris.New("task is not started")
	}

	// should alternate actions
	if latestLog != nil && latestLog.Action == action {
		return nil, eris.New("action cannot be the same with last log")
	}

	log, err := ts.taskRepository.Log(task, action)
	if err != nil {
		return nil, err
	}

	return log, nil
}

func (ts *taskServiceImpl) GetTotalTime(id uuid.UUID) (*time.Duration, error) {
	task, err := ts.taskRepository.GetByID(id)
	if err != nil {
		return nil, err
	}

	logs, err := ts.taskRepository.GetLogs(task)
	if err != nil {
		return nil, err
	}

	var totalTime time.Duration

	for i := 0; i < len(logs); i += 2 {
		if i+1 < len(logs) {
			startLog := logs[i]
			stopLog := logs[i+1]

			if startLog.Action == constant.LOG_ACTION.START && stopLog.Action == constant.LOG_ACTION.STOP {
				duration := stopLog.CreatedAt.Sub(startLog.CreatedAt)
				totalTime += duration
			}
		}
	}

	return &totalTime, nil
}
