package auth

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itsLeonB/catfeinated-time-tracker/internal/dto"
)

type AuthService interface {
	Register(ctx context.Context, request *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Login(ctx context.Context, request *dto.LoginRequest) (*dto.LoginResponse, error)
	CheckToken(ctx context.Context, token string) (bool, error)
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
