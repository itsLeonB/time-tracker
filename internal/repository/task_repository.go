package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/apperror"
	"github.com/itsLeonB/time-tracker/internal/constant"
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

func (tr *taskRepositoryGorm) Insert(ctx context.Context, task *model.Task) (*model.Task, error) {
	err := tr.db.WithContext(ctx).Create(task).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgInsertError))
	}

	return task, nil
}

func (tr *taskRepositoryGorm) GetAll(ctx context.Context) ([]*model.Task, error) {
	var tasks []*model.Task

	err := tr.db.WithContext(ctx).Find(&tasks, "user_id = ?", ctx.Value(constant.ContextUserID)).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return tasks, nil
}

func (tr *taskRepositoryGorm) GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	var task model.Task

	err := tr.db.WithContext(ctx).First(&task, "id = ? AND user_id = ?", id, ctx.Value(constant.ContextUserID)).Error
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

	err := tr.db.WithContext(ctx).Preload("Logs").First(&task,
		"number = ? AND user_id = ?",
		number, ctx.Value(constant.ContextUserID),
	).Error
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
	newLog := &model.TaskLog{
		TaskID: task.ID,
		Action: action,
	}

	err := tr.db.WithContext(ctx).Create(newLog).Error
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
	}).Where("user_id = ?", ctx.Value(constant.ContextUserID))

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
