package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/util"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type taskRepositoryGorm struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepositoryGorm{db}
}

func (tr *taskRepositoryGorm) Insert(ctx context.Context, task *model.Task) (*model.Task, error) {
	err := tr.db.WithContext(ctx).Create(task).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgInsertError))
	}

	return task, nil
}

func (tr *taskRepositoryGorm) GetAll(ctx context.Context) ([]*model.Task, error) {
	var tasks []*model.Task

	err := tr.db.WithContext(ctx).Find(&tasks).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return tasks, nil
}

func (tr *taskRepositoryGorm) GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	var task model.Task

	err := tr.db.WithContext(ctx).First(&task, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return &task, nil
}

func (tr *taskRepositoryGorm) GetByNumber(ctx context.Context, number string) (*model.Task, error) {
	var task model.Task

	err := tr.db.WithContext(ctx).Preload("Logs").First(&task, "number = ?", number).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return &task, nil
}

func (tr *taskRepositoryGorm) Update(ctx context.Context, task *model.Task) (*model.Task, error) {
	err := tr.db.WithContext(ctx).Save(task).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgUpdateError))
	}

	return task, nil
}

func (tr *taskRepositoryGorm) Delete(ctx context.Context, task *model.Task) error {
	err := tr.db.WithContext(ctx).Delete(task).Error
	if err != nil {
		return apperror.InternalServerError(eris.Wrap(err, apperror.MsgDeleteError))
	}

	return nil
}

func (tr *taskRepositoryGorm) Log(ctx context.Context, task *model.Task, action string) (*model.TaskLog, error) {
	userId, err := util.GetUserID(ctx)
	if err != nil {
		return nil, err
	}

	newLog := &model.TaskLog{
		UserID: userId,
		TaskID: task.ID,
		Action: action,
	}

	err = tr.db.WithContext(ctx).Create(newLog).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgInsertError))
	}

	return newLog, nil
}

func (tr *taskRepositoryGorm) GetLatestLog(ctx context.Context, task *model.Task) (*model.TaskLog, error) {
	var latestLog model.TaskLog

	err := tr.db.WithContext(ctx).Where("task_id = ?", task.ID).Order("created_at DESC").First(&latestLog).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return &latestLog, nil
}

func (tr *taskRepositoryGorm) GetLogs(ctx context.Context, task *model.Task) ([]*model.TaskLog, error) {
	var logs []*model.TaskLog

	err := tr.db.WithContext(ctx).Where("task_id = ?", task.ID).Order("created_at ASC").Find(&logs).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return logs, nil
}

func (tr *taskRepositoryGorm) Find(ctx context.Context, options *model.QueryOptions) ([]*model.Task, error) {
	var tasks []*model.Task

	query := tr.db.WithContext(ctx).Preload("Logs", func(db *gorm.DB) *gorm.DB {
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

func (tr *taskRepositoryGorm) GetInProgress(ctx context.Context, projectID uuid.UUID) ([]*model.Task, error) {
	var tasks []*model.Task
	subQuery := tr.db.WithContext(ctx).
		Model(&model.TaskLog{}).
		Select("task_id, MAX(created_at) as last_log").
		Group("task_id")

	err := tr.db.WithContext(ctx).
		Model(&model.Task{}).
		Joins("JOIN (?) as last_logs ON tasks.id = last_logs.task_id", subQuery).
		Joins("JOIN task_logs ON task_logs.task_id = tasks.id AND task_logs.created_at = last_logs.last_log").
		Where("tasks.project_id = ? AND task_logs.action = ?", projectID, "START").
		Find(&tasks).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return tasks, nil
}
