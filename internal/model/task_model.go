package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID uuid.UUID `json:"project_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TaskLogs  []TaskLog `json:"task_logs,omitempty"`
}

func (t *Task) TableName() string {
	return "tasks"
}

type NewTaskRequest struct {
	ProjectID uuid.UUID `json:"project_id" binding:"required"`
	Name      string    `json:"name" binding:"required"`
}

type TaskLog struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TaskID    uuid.UUID `json:"task_id"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (tl *TaskLog) TableName() string {
	return "task_logs"
}

type NewLogRequest struct {
	Action string `json:"action" binding:"required,oneof=START STOP"`
}

type TotalTimeResponse struct {
	Duration time.Duration `json:"duration"`
	Minutes  float64       `json:"minutes"`
	Hours    float64       `json:"hours"`
	String   string        `json:"string"`
}
