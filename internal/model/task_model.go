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
}

func (t *Task) TableName() string {
	return "tasks"
}
