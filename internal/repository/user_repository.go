package repository

import (
	"context"

	"github.com/itsLeonB/ezutil"
	"github.com/itsLeonB/time-tracker/internal/appconstant"
	"github.com/itsLeonB/time-tracker/internal/entity"
	"github.com/rotisserie/eris"
	"gorm.io/gorm"
)

type userRepositoryGorm struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryGorm{db}
}

func (ur *userRepositoryGorm) Insert(ctx context.Context, user entity.User) (entity.User, error) {
	db, err := ur.getGormInstance(ctx)
	if err != nil {
		return entity.User{}, err
	}

	err = db.Create(&user).Error
	if err != nil {
		return entity.User{}, eris.Wrap(err, appconstant.MsgInsertError)
	}

	return user, nil
}

func (ur *userRepositoryGorm) FindFirst(ctx context.Context, spec entity.User) (entity.User, error) {
	var user entity.User

	db, err := ur.getGormInstance(ctx)
	if err != nil {
		return entity.User{}, err
	}

	err = db.Scopes(ezutil.WhereBySpec(spec)).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.User{}, nil
		}

		return entity.User{}, eris.Wrap(err, appconstant.MsgQueryError)
	}

	return user, nil
}

func (us *userRepositoryGorm) getGormInstance(ctx context.Context) (*gorm.DB, error) {
	tx, err := ezutil.GetTxFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if tx != nil {
		return tx, nil
	}

	return us.db.WithContext(ctx), nil
}
