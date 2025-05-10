package model

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Project) TableName() string {
	return "projects"
}
