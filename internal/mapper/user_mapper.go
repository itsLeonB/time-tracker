package mapper

import (
	"github.com/itsLeonB/time-tracker/internal/appconstant"
	"github.com/itsLeonB/time-tracker/internal/dto"
	"github.com/itsLeonB/time-tracker/internal/entity"
)

func UserToAuthData(user entity.User) map[string]any {
	return map[string]any{
		appconstant.ContextUserID: user.ID,
	}
}

func UserToResponse(user entity.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
