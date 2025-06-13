package dto

import (
	"github.com/itsLeonB/time-tracker/internal/model"
)

type HomeViewDto struct {
	User     model.User
	Projects []UserProjectResponse
}
