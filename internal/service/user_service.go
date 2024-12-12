package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/time-tracker/internal/apperror"
	"github.com/itsLeonB/time-tracker/internal/constant"
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/itsLeonB/time-tracker/internal/repository"
	"github.com/rotisserie/eris"
)

type userServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userServiceImpl{userRepository}
}

func (us *userServiceImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	return us.userRepository.FindByEmail(ctx, email)
}

func (us *userServiceImpl) Create(ctx context.Context, user *model.User) error {
	return us.userRepository.Insert(ctx, user)
}

func (us *userServiceImpl) ValidateUser(ctx context.Context) (*model.User, error) {
	id := ctx.Value(constant.ContextUserID)
	userID, err := uuid.Parse(id.(string))
	if err != nil {
		return nil, eris.Wrap(err, "error parsing UUID")
	}

	user, err := us.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.NotFoundError(eris.Wrap(err, "user not found"))
	}

	return user, nil
}
