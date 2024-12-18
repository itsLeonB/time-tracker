package model

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID  `json:"userId"`
	Name        string     `json:"name"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	TotalPoints float64    `json:"totalPoints" gorm:"-"`
	TimeSpent   *TimeSpent `json:"timeSpent,omitempty" gorm:"-"`
}

func (p *Project) TableName() string {
	return "projects"
}

type NewProjectRequest struct {
	Name string `json:"name" binding:"required"`
}

type FindProjectOptions struct {
	Name string
}
