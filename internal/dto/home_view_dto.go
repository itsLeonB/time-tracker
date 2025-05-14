package dto

import (
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
)

type HomeViewDto struct {
	User     model.User
	Projects []ProjectResponse
}
