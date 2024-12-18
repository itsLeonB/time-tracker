package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itsLeonB/time-tracker/internal/apperror"
	"github.com/itsLeonB/time-tracker/internal/config"
	"github.com/rotisserie/eris"
)

type jwtHS256 struct {
	issuer        string
	secretKey     string
	tokenDuration time.Duration
}

func NewJWTHS256(configs *config.Auth) *jwtHS256 {
	return &jwtHS256{
		issuer:        configs.Issuer,
		secretKey:     configs.SecretKey,
		tokenDuration: configs.TokenDuration,
	}
}

func (j *jwtHS256) CreateToken(data map[string]any) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    j.issuer,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			Data: data,
		},
	)

	signed, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", eris.Wrap(err, "error signing token")
	}

	return signed, nil
}

func (j *jwtHS256) VerifyToken(tokenstr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenstr,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.secretKey), nil
		},
		jwt.WithIssuer(j.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, apperror.UnauthorizedError(err)
		}

		return nil, eris.Wrap(err, "error parsing token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, eris.Errorf("error parsing token claims")
	}

	return claims, nil
}
