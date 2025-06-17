package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/itsLeonB/ezutil"
	"github.com/itsLeonB/time-tracker/internal/appconstant"
	"github.com/itsLeonB/time-tracker/internal/entity"
	"github.com/itsLeonB/time-tracker/internal/repository"
	"github.com/rotisserie/eris"
)

type userServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userServiceImpl{userRepository}
}

func (us *userServiceImpl) ValidateUser(ctx context.Context) (entity.User, error) {
	userID, err := ezutil.GetAndParseFromContext[uuid.UUID](ctx.(*gin.Context), appconstant.ContextUserID)
	if err != nil {
		return entity.User{}, err
	}

	spec := entity.User{ID: userID}

	user, err := us.userRepository.FindFirst(ctx, spec)
	if err != nil {
		return entity.User{}, err
	}
	if user.IsZero() {
		return entity.User{}, ezutil.NotFoundError(eris.Errorf("user not found"))
	}

	return user, nil
}

func (us *userServiceImpl) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	spec := entity.User{ID: id}

	user, err := us.userRepository.FindFirst(ctx, spec)
	if err != nil {
		return false, err
	}

	return !user.IsZero(), nil
}
