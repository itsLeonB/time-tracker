package strategy

import (
	"math"

	"github.com/itsLeonB/time-tracker/internal/model"
)

type hourBasedPointStrategy struct{}

func NewHourBasedPointStrategy() PointStrategy {
	return &hourBasedPointStrategy{}
}

func (ps *hourBasedPointStrategy) CalculatePoints(task *model.Task) int {
	if task.TimeSpent == nil {
		return 0
	}

	return int(math.Ceil(task.TimeSpent.Hours))
}
