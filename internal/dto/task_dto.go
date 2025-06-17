package dto

import (
	"time"

	"github.com/google/uuid"
)

type NewTaskRequest struct {
	ProjectID uuid.UUID `json:"projectId"`
	Number    string    `json:"number" form:"number" binding:"required,min=3"`
	Name      string    `json:"name" form:"name" binding:"required,min=3"`
}

type TaskResponse struct {
	ID        uuid.UUID         `json:"id"`
	ProjectID uuid.UUID         `json:"projectId"`
	Number    string            `json:"number"`
	Name      string            `json:"name"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
	TimeSpent TimeSpent         `json:"timeSpent"`
	Logs      []TaskLogResponse `json:"logs"`
	IsActive  bool              `json:"isActive"`
	StartTime string            `json:"startTime"`
}

type TaskQueryParams struct {
	QueryParams
	Number string
}
