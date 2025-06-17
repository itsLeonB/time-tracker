package entity

import (
	"time"

	"github.com/google/uuid"
)

type TaskLog struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TaskID    uuid.UUID
	Action    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (tl *TaskLog) TableName() string {
	return "task_logs"
}

func (tl *TaskLog) IsZero() bool {
	return tl.ID == uuid.Nil || tl.Action == ""
}
