package model

import (
	"time"

	"github.com/google/uuid"
)

type UserProject struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserId    uuid.UUID
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserTasks []UserTask `gorm:"foreignKey:UserProjectId"`
}

func (p *UserProject) TableName() string {
	return "user_projects"
}

func (p *UserProject) IsZero() bool {
	return p.ID == uuid.Nil ||
		p.UserId == uuid.Nil ||
		p.Name == ""
}

type ProjectQueryOptions struct {
	QueryOptions
}
