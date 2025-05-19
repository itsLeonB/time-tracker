package model

import (
	"time"

	"github.com/google/uuid"
)

type UserTaskLog struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserTaskId uuid.UUID
	Action     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (utl *UserTaskLog) TableName() string {
	return "user_task_logs"
}

func (utl *UserTaskLog) IsZero() bool {
	return utl.ID == uuid.Nil || utl.Action == ""
}

type UserTaskLogQueryOptions struct {
	QueryOptions
}
