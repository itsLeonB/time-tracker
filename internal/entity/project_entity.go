package entity

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserId    uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Tasks     []Task `gorm:"foreignKey:ProjectID"`
}

func (p *Project) TableName() string {
	return "projects"
}

func (p *Project) IsZero() bool {
	return p.ID == uuid.Nil ||
		p.UserId == uuid.Nil ||
		p.Name == ""
}

type ProjectSpecification struct {
	QuerySpecification
	Project
	Start time.Time
	End   time.Time
}
