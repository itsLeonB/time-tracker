package repository

import (
	"context"

	"github.com/itsLeonB/ezutil"
	"github.com/itsLeonB/time-tracker/internal/appconstant"
	"github.com/itsLeonB/time-tracker/internal/entity"
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

func (tr *taskRepositoryGorm) Insert(ctx context.Context, task entity.Task) (entity.Task, error) {
	db, err := tr.getGormInstance(ctx)
	if err != nil {
		return entity.Task{}, err
	}

	err = db.Create(&task).Error
	if err != nil {
		return entity.Task{}, eris.Wrap(err, appconstant.MsgInsertError)
	}

	return task, nil
}

func (tr *taskRepositoryGorm) FindAll(ctx context.Context, spec entity.Task) ([]entity.Task, error) {
	var tasks []entity.Task

	db, err := tr.getGormInstance(ctx)
	if err != nil {
		return nil, err
	}

	err = db.
		Scopes(util.DefaultOrdering(), ezutil.WhereBySpec(spec)).
		Find(&tasks).
		Error

	if err != nil {
		return nil, eris.Wrap(err, appconstant.MsgQueryError)
	}

	return tasks, nil
}

func (tr *taskRepositoryGorm) FindFirst(ctx context.Context, spec entity.Task) (entity.Task, error) {
	var task entity.Task

	db, err := tr.getGormInstance(ctx)
	if err != nil {
		return entity.Task{}, err
	}

	err = db.
		Scopes(util.DefaultOrdering(), ezutil.WhereBySpec(spec)).
		First(&task).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Task{}, nil
		}

		return entity.Task{}, eris.Wrap(err, appconstant.MsgQueryError)
	}

	return task, nil
}

func (tr *taskRepositoryGorm) Update(ctx context.Context, task entity.Task) (entity.Task, error) {
	db, err := tr.getGormInstance(ctx)
	if err != nil {
		return entity.Task{}, err
	}

	err = db.Save(&task).Error
	if err != nil {
		return entity.Task{}, eris.Wrap(err, appconstant.MsgUpdateError)
	}

	return task, nil
}

func (tr *taskRepositoryGorm) Delete(ctx context.Context, task entity.Task) (entity.Task, error) {
	db, err := tr.getGormInstance(ctx)
	if err != nil {
		return entity.Task{}, err
	}

	err = db.Delete(&task).Error
	if err != nil {
		return entity.Task{}, eris.Wrap(err, appconstant.MsgDeleteError)
	}

	return task, nil
}

func (tr *taskRepositoryGorm) getGormInstance(ctx context.Context) (*gorm.DB, error) {
	tx, err := ezutil.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx, nil
	}

	return tr.db.WithContext(ctx), nil
}
