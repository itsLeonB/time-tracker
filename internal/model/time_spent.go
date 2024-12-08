package model

import "time"

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
