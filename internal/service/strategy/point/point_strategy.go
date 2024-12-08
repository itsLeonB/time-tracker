package strategy

import "github.com/itsLeonB/time-tracker/internal/model"

type PointStrategy interface {
	CalculatePoints(task *model.Task) int
}
