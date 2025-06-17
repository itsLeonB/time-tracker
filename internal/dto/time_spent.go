package dto

import (
	"math"
	"time"
)

type TimeSpent struct {
	Duration time.Duration `json:"duration"`
	Minutes  float64       `json:"minutes"`
	Hours    float64       `json:"hours"`
	String   string        `json:"string"`
}

func (ts *TimeSpent) Add(timeSpent TimeSpent) {
	ts.Duration += timeSpent.Duration
	ts.Minutes = ts.Duration.Minutes()
	ts.Hours = round(ts.Duration.Hours())
	ts.String = ts.Duration.String()
}

func TimeSpentFromDuration(duration time.Duration) TimeSpent {
	return TimeSpent{
		Duration: duration,
		Minutes:  duration.Minutes(),
		Hours:    round(duration.Hours()),
		String:   duration.String(),
	}
}

func round(val float64) float64 {
	ratio := math.Pow(10, float64(2))

	return math.Ceil(val*ratio) / ratio
}
