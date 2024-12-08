package strategy

import (
	"github.com/itsLeonB/time-tracker/internal/model"
	"github.com/itsLeonB/time-tracker/internal/util"
)

type hourBasedPointStrategy struct{}

func (ps *hourBasedPointStrategy) CalculatePoints(task *model.Task) float64 {
	if task.TimeSpent == nil {
		return 0
	}

	return util.Round(task.TimeSpent.Hours, 1)
}
