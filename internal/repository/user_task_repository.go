package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/util"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type userTaskRepositoryImpl struct {
	db *gorm.DB
}

func NewUserTaskRepository(db *gorm.DB) UserTaskRepository {
	return &userTaskRepositoryImpl{db}
}

func (utr *userTaskRepositoryImpl) Insert(ctx context.Context, userTask model.UserTask) (model.UserTask, error) {
	err := utr.db.WithContext(ctx).Create(&userTask).Error
	if err != nil {
		return userTask, eris.Wrap(err, apperror.MsgInsertError)
	}

	return userTask, nil
}

func (utr *userTaskRepositoryImpl) FindAll(ctx context.Context, options model.UserTaskQueryOptions) ([]model.UserTask, error) {
	var userTasks []model.UserTask

	query := utr.db.WithContext(ctx).Scopes(
		util.PreloadRelations(options.PreloadRelations),
		util.FilterByColumns(options.Filters),
		util.DefaultOrdering(),
	)

	err := query.Find(&userTasks).Error
	if err != nil {
		return nil, eris.Wrap(err, apperror.MsgQueryError)
	}

	return userTasks, nil
}

func (utr *userTaskRepositoryImpl) FindById(ctx context.Context, id uuid.UUID) (model.UserTask, error) {
	var userTask model.UserTask

	err := utr.db.WithContext(ctx).First(&userTask, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return userTask, nil
		}

		return userTask, eris.Wrap(err, apperror.MsgQueryError)
	}

	return userTask, nil
}
