package utils

import (
	"math"
	"math/rand"

	"github.com/google/uuid"
)

func RandomElement[v any](slice []v) v {
	return slice[rand.Intn(len(slice))]
}

func RandomNumber(min, max int) int {
	return rand.Intn(max+1-min) + min
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
	return int(math.Round(float64(value) * float64(percent) / 100.0))
}

func Map[v any, r any](slice []v, mapFunc func(v) r) []r {
	result := make([]r, len(slice))

	for i, v := range slice {
		result[i] = mapFunc(v)
	}

	return result
}

func BoolToText(value bool, ifTrue string, ifFalse string) string {
	if value {
		return ifTrue
	} else {
		return ifFalse
	}
}

func SkillUUIDToItemUUID(skillUuid uuid.UUID) uuid.UUID {
	bytes, _ := skillUuid.MarshalBinary()

	copy(bytes[6:8], []byte{0, 0})

	parsedUuid, _ := uuid.FromBytes(bytes)

	return parsedUuid
}
