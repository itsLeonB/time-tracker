package dto

import (
	"github.com/itsLeonB/time-tracker/internal/model"
)

type ProjectDetailViewDto struct {
	User      model.User
	Project   UserProjectResponse
	StartDate string
	EndDate   string
}
