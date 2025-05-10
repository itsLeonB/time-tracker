package dto

import "github.com/google/uuid"

type NewTaskRequest struct {
	ProjectID uuid.UUID `json:"projectId" binding:"required"`
	Number    string    `json:"number" binding:"required"`
	Name      string    `json:"name" binding:"required"`
}
