package model

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (p *Project) TableName() string {
	return "projects"
}

type NewProjectRequest struct {
	Name string `json:"name" binding:"required"`
}
