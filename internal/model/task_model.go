package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid"`
	ProjectID uuid.UUID  `json:"project_id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	TaskLogs  []*TaskLog `json:"task_logs"`
}

func (t *Task) TableName() string {
	return "tasks"
}

type TaskLog struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key"`
	TaskID    uuid.UUID `json:"task_id"`
	Action    Action    `json:"action"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (tl *TaskLog) TableName() string {
	return "task_logs"
}

type Action string

const (
	Start Action = "START"
	Stop  Action = "STOP"
)
