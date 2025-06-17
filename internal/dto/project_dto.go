package dto

import (
	"time"

	"github.com/google/uuid"
)

type NewProjectRequest struct {
	Name string `json:"name" form:"name" binding:"required"`
}

type ProjectQueryParams struct {
	QueryParams
	ID   uuid.UUID
	Name string
}

type ProjectResponse struct {
	ID              uuid.UUID      `json:"id"`
	UserId          uuid.UUID      `json:"userId"`
	Name            string         `json:"name"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	TimeSpent       TimeSpent      `json:"timeSpent"`
	Tasks           []TaskResponse `json:"userTasks"`
	ActiveTaskCount int            `json:"activeTaskCount"`
	AverageTaskTime string         `json:"averageTaskTime"`
}
