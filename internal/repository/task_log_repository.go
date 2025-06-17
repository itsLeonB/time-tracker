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

type taskLogRepositoryGorm struct {
	db *gorm.DB
}

func NewTaskLogRepository(db *gorm.DB) TaskLogRepository {
	return &taskLogRepositoryGorm{db}
}

func (tlr *taskLogRepositoryGorm) Insert(ctx context.Context, taskLog entity.TaskLog) (entity.TaskLog, error) {
	db, err := tlr.getGormInstance(ctx)
	if err != nil {
		return entity.TaskLog{}, err
	}

	err = db.Create(&taskLog).Error
	if err != nil {
		return entity.TaskLog{}, eris.Wrap(err, appconstant.MsgInsertError)
	}

	return taskLog, nil
}

func (tlr *taskLogRepositoryGorm) FindLatest(ctx context.Context, spec entity.TaskLog) (entity.TaskLog, error) {
	var taskLog entity.TaskLog

	db, err := tlr.getGormInstance(ctx)
	if err != nil {
		return entity.TaskLog{}, err
	}

	err = db.
		Scopes(ezutil.WhereBySpec(spec), util.DefaultOrdering()).
		Take(&taskLog).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.TaskLog{}, nil
		}

		return entity.TaskLog{}, eris.Wrap(err, appconstant.MsgQueryError)
	}

	return taskLog, nil
}

func (tlr *taskLogRepositoryGorm) getGormInstance(ctx context.Context) (*gorm.DB, error) {
	tx, err := ezutil.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx, nil
	}

	return tlr.db.WithContext(ctx), nil
}
