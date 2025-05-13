package repository

import (
	"context"

	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/util"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type gormUserTaskLogRepository struct {
	db *gorm.DB
}

func NewUserTaskLogRepository(db *gorm.DB) UserTaskLogRepository {
	return &gormUserTaskLogRepository{db}
}

func (r *gormUserTaskLogRepository) Insert(ctx context.Context, userTaskLog model.UserTaskLog) (model.UserTaskLog, error) {
	err := r.db.WithContext(ctx).Create(&userTaskLog).Error
	if err != nil {
		return userTaskLog, eris.Wrap(err, "failed to insert user task log")
	}

	return userTaskLog, nil
}

func (utlr *gormUserTaskLogRepository) FindLatest(ctx context.Context, options model.UserTaskLogQueryOptions) (model.UserTaskLog, error) {
	var userTaskLog model.UserTaskLog

	query := utlr.db.WithContext(ctx).Scopes(
		util.PreloadRelations(options.PreloadRelations),
		util.FilterByColumns(options.Filters),
	)

	err := query.Order("created_at DESC").First(&userTaskLog).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return userTaskLog, nil
		}

		return userTaskLog, eris.Wrap(err, "failed to find latest user task log")
	}

	return userTaskLog, nil
}
