package dto

import (
	"time"

	"github.com/google/uuid"
)

type NewUserTaskRequest struct {
	UserProjectId uuid.UUID `json:"userProjectId"`
	TaskId        uuid.UUID `json:"taskId" binding:"required"`
}

type UserTaskResponse struct {
	ID            uuid.UUID             `json:"id"`
	UserProjectId uuid.UUID             `json:"userProjectId"`
	TaskId        uuid.UUID             `json:"taskId"`
	CreatedAt     time.Time             `json:"createdAt"`
	UpdatedAt     time.Time             `json:"updatedAt"`
	TaskNumber    string                `json:"taskNumber"`
	TaskName      string                `json:"taskName"`
	TimeSpent     TimeSpent             `json:"timeSpent"`
	Logs          []UserTaskLogResponse `json:"logs"`
	IsActive      bool                  `json:"isActive"`
	StartTime     string                `json:"startTime"`
}

type UserTaskQueryParams struct {
	UserProjectId uuid.UUID
	TaskId        uuid.UUID
}
