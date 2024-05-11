package utils

import (
	"math"
)

func RoundToTwoDigits(f float64) float64 {
	return math.Round(f*100) / 100
}
