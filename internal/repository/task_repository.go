package repository

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/apperror"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/itsLeonB/time-tracker/internal/util"
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
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgInsertError))
	}

	return task, nil
}

func (tr *taskRepositoryGorm) GetAll() ([]*model.Task, error) {
	var tasks []*model.Task

	err := tr.db.Find(&tasks).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return tasks, nil
}

func (tr *taskRepositoryGorm) GetByID(id uuid.UUID) (*model.Task, error) {
	var task model.Task

	err := tr.db.First(&task, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return &task, nil
}

func (tr *taskRepositoryGorm) GetByNumber(number string) (*model.Task, error) {
	var task model.Task

	err := tr.db.Preload("Logs").First(&task, "number = ?", number).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return &task, nil
}

func (tr *taskRepositoryGorm) Update(task *model.Task) (*model.Task, error) {
	err := tr.db.Save(task).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgUpdateError))
	}

	return task, nil
}

func (tr *taskRepositoryGorm) Delete(task *model.Task) error {
	err := tr.db.Delete(task).Error
	if err != nil {
		return apperror.InternalServerError(eris.Wrap(err, apperror.MsgDeleteError))
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
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgInsertError))
	}

	return newLog, nil
}

func (tr *taskRepositoryGorm) GetLatestLog(task *model.Task) (*model.TaskLog, error) {
	var latestLog model.TaskLog

	err := tr.db.Where("task_id = ?", task.ID).Order("created_at DESC").First(&latestLog).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return &latestLog, nil
}

func (tr *taskRepositoryGorm) GetLogs(task *model.Task) ([]*model.TaskLog, error) {
	var logs []*model.TaskLog

	err := tr.db.Where("task_id = ?", task.ID).Order("created_at ASC").Find(&logs).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return logs, nil
}

func (tr *taskRepositoryGorm) Find(options *model.QueryOptions) ([]*model.Task, error) {
	var tasks []*model.Task

	query := tr.db.Preload("Logs", func(db *gorm.DB) *gorm.DB {
		return db.Order("task_logs.created_at ASC")
	})

	if options != nil {
		if options.Params != nil {
			if options.Params.Number != "" {
				query = query.Where("number ILIKE ?", fmt.Sprintf("%%%s%%", options.Params.Number))
			}

			if options.Params.ProjectID != uuid.Nil {
				query = query.Where("project_id = ?", options.Params.ProjectID)
			}

			if options.Params.Date != (time.Time{}) {
				query = query.Where("created_at >= ? AND created_at <= ?",
					util.StartOfDay(options.Params.Date),
					util.EndOfDay(options.Params.Date),
				)
			}
		}
	}

	err := query.Find(&tasks).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return tasks, nil
}
