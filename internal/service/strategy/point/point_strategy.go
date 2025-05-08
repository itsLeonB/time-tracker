package strategy

import "github.com/itsLeonB/catfeinated-time-tracker/internal/model"

type PointStrategy interface {
	CalculatePoints(task *model.Task) float64
}

// converts hours spent to points,
// rounding up with 1 significant digit
func NewHourBasedPointStrategy() PointStrategy {
	return &hourBasedPointStrategy{}
}
