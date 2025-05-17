package dto

import (
	"time"

	"github.com/google/uuid"
)

type NewTaskRequest struct {
	Number string `json:"number" binding:"required,min=3"`
	Name   string `json:"name" binding:"required,min=3"`
}

type TaskResponse struct {
	ID        uuid.UUID `json:"id"`
	Number    string    `json:"number"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type TaskQueryParams struct {
	Number string
	Date   time.Time
}
