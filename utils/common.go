package utils

import (
	"math"
	"math/rand"
)

// Random element from slice
func RandomElement[v any](slice []v) v {
	return slice[rand.Intn(len(slice))]
}

func CalcReducedDamage(atk, def int) int {
	if def < 0 {
		return int(float32(atk) * float32(2.0-float32(100/(100-def))))
	} else {
		return int(float32(atk) * float32(100/(100+def)))
	}
}

// Random number
func RandomNumber(min, max int) int {
	return rand.Intn(max-min) + min
}

func PercentOf(value, percent int) int {
	return int(math.Round(float64(value) * float64(percent) / 100.0))
}

func ReadStringWithOffset(offset int, buf []byte) (int, string) {
	strLen := int(buf[offset])

	return offset + 1 + strLen, string(buf[offset+1 : offset+1+strLen])
}

func WriteStringWithOffset(buf []byte, offset int, str string) int {
	buf[offset] = byte(len(str))
	offset++

	copy(buf[offset:offset+len(str)], str)

	return offset + len(str)
}
