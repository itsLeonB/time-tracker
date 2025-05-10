package dto

import (
	"time"

	"github.com/google/uuid"
)

type NewTaskRequest struct {
	ProjectID uuid.UUID `json:"projectId" binding:"required"`
	Number    string    `json:"number" binding:"required"`
	Name      string    `json:"name" binding:"required"`
}

type TaskResponse struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"projectId"`
	Number    string    `json:"number"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
