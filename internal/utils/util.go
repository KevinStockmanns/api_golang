package utils

import (
	"math"
	"regexp"
)

func IsInteger(num string) bool {
	matched, _ := regexp.MatchString("^[0-9]+$", num)
	return matched
}

func RoundDecimal(decimal float64) float64 {
	return (math.Ceil(decimal*10) / 10)
}
