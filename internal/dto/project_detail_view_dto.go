package dto

import (
	"github.com/itsLeonB/catfeinated-time-tracker/internal/model"
)

type ProjectDetailViewDto struct {
	User      model.User
	Project   ProjectResponse
	StartDate string
	EndDate   string
}
