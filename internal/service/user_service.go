package service

import (
	"context"

	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/itsLeonB/time-tracker/internal/repository"
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
