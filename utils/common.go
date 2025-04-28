package utils

import (
	"crypto/rand"
	"math"
	"math/big"
)

func RandomElement[v any](slice []v) v {
	number, err := rand.Int(rand.Reader, big.NewInt(int64(len(slice))))

	if err != nil {
		panic(err)
	}

	return slice[int(number.Int64())]
}

func RandomNumber(min, max int) int {
	number, err := rand.Int(rand.Reader, big.NewInt(int64(max+1-min)))

	if err != nil {
		panic(err)
	}

	return int(number.Int64()) + min
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

func PercentOf(value, percent int) int {
	if value == 0 || percent == 0 {
		return 0
	}

	return int(math.Round(float64(value) * float64(percent) / 100.0))
}