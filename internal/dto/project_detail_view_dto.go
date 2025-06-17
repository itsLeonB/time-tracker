package dto

import "github.com/itsLeonB/time-tracker/internal/entity"

type ProjectDetailViewDto struct {
	User      entity.User
	Project   ProjectResponse
	StartDate string
	EndDate   string
}
