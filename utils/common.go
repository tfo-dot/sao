package utils

import (
	"crypto/rand"
	"math"
	"math/big"
	"sao/types"

	"github.com/google/uuid"
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
	return int(math.Round(float64(value) * float64(percent) / 100.0))
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

var StringToStat = map[string]types.Stat{
	"None":       types.STAT_NONE,
	"HP":         types.STAT_HP,
	"SPD":        types.STAT_SPD,
	"AGL":        types.STAT_AGL,
	"ATK":        types.STAT_AD,
	"DEF":        types.STAT_DEF,
	"MR":         types.STAT_MR,
	"MANA":       types.STAT_MANA,
	"AP":         types.STAT_AP,
	"HEAL_SELF":  types.STAT_HEAL_SELF,
	"HEAL_POWER": types.STAT_HEAL_POWER,
	"LETHAL":     types.STAT_LETHAL,
	"LETHAL%":    types.STAT_LETHAL_PERCENT,
	"MAGIC_PEN":  types.STAT_MAGIC_PEN,
	"MAGIC_PEN%": types.STAT_MAGIC_PEN_PERCENT,
	"ADAPTIVE":   types.STAT_ADAPTIVE,
	"ADAPTIVE%":  types.STAT_ADAPTIVE_PERCENT,
	"OMNI_VAMP":  types.STAT_OMNI_VAMP,
	"ATK_VAMP":   types.STAT_ATK_VAMP,
}
