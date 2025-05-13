package dto

import (
	"time"

	"github.com/google/uuid"
)

type NewUserTaskLogRequest struct {
	UserTaskId uuid.UUID `json:"userTaskId"`
	Action     string    `json:"action" binding:"required,oneof=START STOP"`
}

type UserTaskLogResponse struct {
	Id         uuid.UUID `json:"id"`
	UserTaskId uuid.UUID `json:"userTaskId"`
	Action     string    `json:"action"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
