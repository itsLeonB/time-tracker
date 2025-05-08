package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db}
}

func (ur *userRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := ur.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, eris.Wrap(err, apperror.MsgQueryError)
	}

	return &user, nil
}

func (ur *userRepositoryImpl) Insert(ctx context.Context, user *model.User) error {
	err := ur.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return eris.Wrap(err, apperror.MsgInsertError)
	}

	return nil
}

func (ur *userRepositoryImpl) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	err := ur.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, eris.Wrap(err, apperror.MsgQueryError)
	}

	return &user, nil
}
