package auth

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/apperror"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/constant"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/service"
	"github.com/rotisserie/eris"
)

type authServiceImpl struct {
	hasher      Hasher
	jwt         JWT
	userService service.UserService
}

func NewAuthService(hasher Hasher, jwt JWT, userService service.UserService) AuthService {
	return &authServiceImpl{hasher, jwt, userService}
}

func (as *authServiceImpl) Register(ctx context.Context, request *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	existingUser, err := as.userService.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, apperror.ConflictError(eris.Errorf("email: %s already registered", request.Email))
	}

	hash, err := as.hasher.Hash(request.Password)
	if err != nil {
		return nil, err
	}

	newUser := &model.User{
		Email:    request.Email,
		Password: hash,
	}

	err = as.userService.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return &dto.RegisterResponse{
		Message: "register success, please login",
	}, nil
}

func (as *authServiceImpl) Login(ctx context.Context, request *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := as.userService.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperror.BadRequestError(eris.Errorf("wrong combination of email/password"))
	}

	ok, err := as.hasher.CheckHash(user.Password, request.Password)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, apperror.BadRequestError(eris.Errorf("wrong combination of email/password"))
	}

	token, err := as.jwt.CreateToken(gin.H{constant.ContextUserID: user.ID})
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Type:  "Bearer",
		Token: token,
	}, nil
}
