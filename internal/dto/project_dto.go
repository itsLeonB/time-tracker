package dto

import (
	"time"

	"github.com/google/uuid"
)

type NewProjectRequest struct {
	Name string `json:"name" binding:"required"`
}

type FindProjectOptions struct {
	Name string
	Ids  []uuid.UUID
}

type ProjectResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	TimeSpent TimeSpent `json:"timeSpent"`
}
