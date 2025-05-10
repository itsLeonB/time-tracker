package dto

import (
	"time"

	"github.com/google/uuid"
)

type NewProjectRequest struct {
	Name string `json:"name" binding:"required"`
}

type FindProjectOptions struct {
	Name string
}

type ProjectResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	TimeSpent TimeSpent `json:"timeSpent"`
}

type TimeSpent struct {
	Duration time.Duration `json:"duration"`
	Minutes  float64       `json:"minutes"`
	Hours    float64       `json:"hours"`
	String   string        `json:"string"`
}

func (ts *TimeSpent) Add(timeSpent *TimeSpent) {
	ts.Duration += timeSpent.Duration
	ts.Minutes = ts.Duration.Minutes()
	ts.Hours = ts.Duration.Hours()
	ts.String = ts.Duration.String()
}
