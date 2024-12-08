package util

import "math"

func Round(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))

	return math.Ceil(val*ratio) / ratio
}
