package dto

import (
	"github.com/itsLeonB/time-tracker/internal/entity"
)

type HomeViewDto struct {
	User     entity.User
	Projects []ProjectResponse
}
