package repository

import (
	"context"
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

func (tr *taskRepositoryGorm) GetAll(ctx context.Context) ([]model.Task, error) {
	var tasks []model.Task

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

func (tr *taskRepositoryGorm) Find(ctx context.Context, options model.TaskQueryOptions) ([]model.Task, error) {
	var tasks []model.Task

	query := tr.db.WithContext(ctx)

	if options.Number != "" {
		query = query.Where("number = ?", options.Number)
	}

	if options.ProjectID != uuid.Nil {
		query = query.Where("project_id = ?", options.ProjectID)
	}

	if options.Date != (time.Time{}) {
		query = query.Where("created_at >= ? AND created_at <= ?",
			util.StartOfDay(options.Date),
			util.EndOfDay(options.Date),
		)
	}

	err := query.Find(&tasks).Error
	if err != nil {
		return nil, apperror.InternalServerError(eris.Wrap(err, apperror.MsgQueryError))
	}

	return tasks, nil
}
