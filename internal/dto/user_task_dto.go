package dto

import (
	"time"

	"github.com/google/uuid"
)

type NewUserTaskRequest struct {
	TaskId uuid.UUID `json:"taskId" binding:"required"`
	UserId uuid.UUID `json:"userId"`
}

type UserTaskResponse struct {
	ID         uuid.UUID             `json:"id"`
	UserId     uuid.UUID             `json:"userId"`
	TaskId     uuid.UUID             `json:"taskId"`
	TaskNumber string                `json:"taskNumber"`
	TaskName   string                `json:"taskName"`
	ProjectId  uuid.UUID             `json:"projectId"`
	CreatedAt  time.Time             `json:"createdAt"`
	UpdatedAt  time.Time             `json:"updatedAt"`
	TimeSpent  TimeSpent             `json:"timeSpent"`
	Logs       []UserTaskLogResponse `json:"logs"`
}

type UserTaskQueryParams struct {
	UserId uuid.UUID
}
