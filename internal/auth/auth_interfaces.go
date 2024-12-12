package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itsLeonB/time-tracker/internal/model"
)

type AuthService interface {
	Register(ctx context.Context, request *model.RegisterRequest) (*model.RegisterResponse, error)
	Login(ctx context.Context, request *model.LoginRequest) (*model.LoginResponse, error)
}

type Hasher interface {
	Hash(val string) (string, error)
	CheckHash(hash, val string) (bool, error)
}

type JWT interface {
	CreateToken(data map[string]any) (string, error)
	VerifyToken(token string) (*JWTClaims, error)
}

type JWTClaims struct {
	jwt.RegisteredClaims
	Data map[string]any `json:"data"`
}
