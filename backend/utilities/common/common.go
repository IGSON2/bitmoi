package common

import "math"

const (
	Decimal = 10000
)

func FloorDecimal(f float64) float64 {
	return math.Floor(f*Decimal) / Decimal
}

func RoundDecimal(f float64) float64 {
	return math.Round(f*Decimal) / Decimal
}
