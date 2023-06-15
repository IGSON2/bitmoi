package common

import "math"

const (
	Decimal = 10000
)

// TODO :
func FloorDecimal(f float64) float64 {
	return math.Floor(f*Decimal) / Decimal
}
