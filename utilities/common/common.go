package common

import "math"

const (
	Decimal       = 100000
	EnvProduction = "production"
	EnvDevelop    = "develop"
)

func FloorDecimal(f float64) float64 {
	return math.Floor(f*Decimal) / Decimal
}

func RoundDecimal(f float64) float64 {
	return math.Round(f*Decimal) / Decimal
}

func CeilDecimal(f float64) float64 {
	return math.Ceil(f*Decimal) / Decimal
}
