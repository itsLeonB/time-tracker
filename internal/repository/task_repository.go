package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type taskRepositoryGorm struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepositoryGorm{db}
}

func (tr *taskRepositoryGorm) Insert(task *model.Task) (*model.Task, error) {
	err := tr.db.Create(task).Error
	if err != nil {
		return nil, eris.Wrap(err, "error inserting")
	}

	return task, nil
}

func (tr *taskRepositoryGorm) GetAll() ([]model.Task, error) {
	var tasks []model.Task

	err := tr.db.Find(&tasks).Error
	if err != nil {
		return nil, eris.Wrap(err, "error querying")
	}

	return tasks, nil
}

func (tr *taskRepositoryGorm) GetByID(id uuid.UUID) (*model.Task, error) {
	var task model.Task

	err := tr.db.First(&task, "id = ?", id).Error
	if err != nil {
		return nil, eris.Wrap(err, "error querying")
	}

	return &task, nil
}

func (tr *taskRepositoryGorm) Update(task *model.Task) (*model.Task, error) {
	err := tr.db.Save(task).Error
	if err != nil {
		return nil, eris.Wrap(err, "error updating")
	}

	return task, nil
}

func (tr *taskRepositoryGorm) Delete(task *model.Task) error {
	err := tr.db.Delete(task).Error
	if err != nil {
		return eris.Wrap(err, "error deleting")
	}

	return nil
}

func (tr *taskRepositoryGorm) Log(task *model.Task, action string) (*model.TaskLog, error) {
	newLog := &model.TaskLog{
		TaskID: task.ID,
		Action: action,
	}

	err := tr.db.Create(newLog).Error
	if err != nil {
		return nil, eris.Wrap(err, "error inserting")
	}

	return newLog, nil
}

func (tr *taskRepositoryGorm) GetLatestLog(task *model.Task) (*model.TaskLog, error) {
	var latestLog model.TaskLog

	err := tr.db.Where("task_id = ?", task.ID).Order("created_at DESC").First(&latestLog).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, eris.Wrap(err, "error querying")
	}

	return &latestLog, nil
}

func (tr *taskRepositoryGorm) GetLogs(task *model.Task) ([]model.TaskLog, error) {
	var logs []model.TaskLog

	err := tr.db.Where("task_id = ?", task.ID).Order("created_at ASC").Find(&logs).Error
	if err != nil {
		return nil, eris.Wrap(err, "error querying")
	}

	return logs, nil
}
