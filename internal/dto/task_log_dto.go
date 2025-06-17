package dto

import (
	"time"

	"github.com/google/uuid"
)

type NewTaskLogRequest struct {
	ProjectID uuid.UUID
	TaskID    uuid.UUID `json:"taskId"`
	Action    string    `json:"action" form:"action" binding:"required,oneof=START STOP"`
}

type TaskLogResponse struct {
	Id           uuid.UUID `json:"id"`
	TaskID       uuid.UUID `json:"taskId"`
	Action       string    `json:"action"`
	CreatedAt    time.Time `json:"createdAt"`
	CreatedAtIso string    `json:"createdAtIso"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
