package model

import (
	"time"

	"github.com/google/uuid"
)

type TaskLog struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"userId"`
	TaskID    uuid.UUID `json:"taskId"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (tl *TaskLog) TableName() string {
	return "task_logs"
}

type NewLogRequest struct {
	Action string `json:"action" binding:"required,oneof=START STOP"`
}

type LogByNumberRequest struct {
	Number string `json:"number" binding:"required"`
	Action string `json:"action" binding:"required,oneof=START STOP"`
}
