package utils

import (
	"math"
	"math/rand"
)

func RandomElement[v any](slice []v) v {
	return slice[rand.Intn(len(slice))]
}

func CalcReducedDamage(atk, reductionValue int) int {
	if reductionValue == 0 {
		return atk
	}

	if reductionValue < 0 {
		return int(float32(atk) * float32(2.0-float32(100/(100-reductionValue))))
	} else {
		return int(float32(atk) * float32(100/(100+reductionValue)))
	}
}

// Random number
func RandomNumber(min, max int) int {
	return rand.Intn(max+1-min) + min
}

func PercentOf(value, percent int) int {
	return int(math.Round(float64(value) * float64(percent) / 100.0))
}

func Map[v any, r any](slice []v, mapFunc func(v) r) []r {
	result := make([]r, len(slice))

	for i, v := range slice {
		result[i] = mapFunc(v)
	}

	return result
}
