package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProjectID uuid.UUID
	Number    string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserTasks []UserTask `gorm:"foreignKey:TaskId"`
}

func (t *Task) TableName() string {
	return "tasks"
}

func (t *Task) IsZero() bool {
	return t.ID == uuid.Nil
}

type TaskQueryOptions struct {
	Number    string
	ProjectID uuid.UUID
	Date      time.Time
}
