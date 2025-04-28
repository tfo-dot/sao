package base

import (
	"sao/types"
	"sao/utils"

	"github.com/google/uuid"
)

func TakeDMGOrDodge[T types.Entity](dmg types.ActionDamage, entity T) (map[types.DamageType]int, bool) {
	if utils.RandomNumber(0, 100) <= entity.GetStat(types.STAT_AGL) && dmg.CanDodge {
		return make(map[types.DamageType]int), true
	}

	return entity.TakeDMG(dmg), false
}

func TakeDMG[T types.Entity](dmg types.ActionDamage, entity T) map[types.DamageType]int {
	resistMapRaw := make(map[types.DamageType]int)
	resistMapPercent := make(map[types.DamageType]int)

	dmgStats := make(map[types.DamageType]int)
	resistFull := make(map[types.DamageType]bool)

	for _, effect := range entity.GetAllEffects() {
		if effect.Effect == types.EFFECT_RESIST {
			meta := effect.Meta.(types.ActionEffectResist)

			if meta.All {
				if meta.DmgType == 4 {
					return make(map[types.DamageType]int)
				} else {
					resistFull[types.DamageType(meta.DmgType)] = true
				}

				continue
			}

			if meta.IsPercent {
				if meta.DmgType == 4 {
					resistMapPercent[types.DMG_PHYSICAL] += effect.Value
					resistMapPercent[types.DMG_MAGICAL] += effect.Value
					resistMapPercent[types.DMG_TRUE] += effect.Value
				} else {
					resistMapPercent[types.DamageType(meta.DmgType)] += effect.Value
				}
			} else {
				if meta.DmgType == 4 {
					resistMapRaw[types.DMG_PHYSICAL] += effect.Value
					resistMapRaw[types.DMG_MAGICAL] += effect.Value
					resistMapRaw[types.DMG_TRUE] += effect.Value
				} else {
					resistMapRaw[types.DamageType(meta.DmgType)] += effect.Value
				}
			}
		}
	}

	for _, dmg := range dmg.Damage {
		if resistFull[dmg.Type] || resistMapPercent[dmg.Type] >= 100 {
			continue
		}

		if resistMapRaw[dmg.Type] > 0 {
			dmg.Value -= resistMapRaw[dmg.Type]
		}

		if dmg.Value <= 0 {
			continue
		}

		if resistMapPercent[dmg.Type] > 0 {
			dmg.Value = utils.PercentOf(dmg.Value, 100-resistMapPercent[dmg.Type])
		}

		if dmg.Value <= 0 {
			continue
		}

		defValue := 0

		if dmg.Type == types.DMG_PHYSICAL {
			defValue = entity.GetStat(types.STAT_DEF)
		} else if dmg.Type == types.DMG_MAGICAL {
			defValue = entity.GetStat(types.STAT_MR)
		}

		if defValue != 0 {
			dmg.Value = utils.CalcReducedDamage(dmg.Value, defValue)
		}

		if dmg.Value <= 0 {
			continue
		}

		if dmg.Type != types.DMG_TRUE {
			dmg.Value = entity.DamageShields(dmg.Value)
		}

		dmgStats[dmg.Type] += dmg.Value

		entity.ChangeHP(-dmg.Value)
	}

	return dmgStats
}

func DamageShields[T types.Entity](dmg int, entity T) ([]types.ActionEffect, int) {
	leftOverDmg := dmg

	validEffects := make([]types.ActionEffect, 0)

	for _, effect := range entity.GetAllEffects() {
		if effect.Effect != types.EFFECT_SHIELD {
			validEffects = append(validEffects, effect)
			continue
		}

		newShieldValue := effect.Value - leftOverDmg

		if newShieldValue <= 0 {
			leftOverDmg = newShieldValue * -1
		} else {
			effect.Value = newShieldValue
			leftOverDmg = 0
		}
	}

	return validEffects, leftOverDmg
}

func Cleanse[T types.Entity](entity T) []types.ActionEffect {
	keepList := make([]types.ActionEffect, 0)

	for _, effect := range entity.GetAllEffects() {
		switch effect.Effect {
		case types.EFFECT_DOT:
			continue
		case types.EFFECT_STUN:
			continue
		case types.EFFECT_TAUNTED:
			continue
		case types.EFFECT_STAT_DEC:
			continue
		}

		keepList = append(keepList, effect)
	}

	return keepList
}

func TriggerAllEffects[T types.Entity](entity T) ([]types.ActionEffect) {
	effects := make([]types.ActionEffect, 0)

	for _, effect := range entity.GetAllEffects() {
		if effect.Duration > 0 {
			effect.Duration--
		}

		switch effect.Effect {
		case types.EFFECT_DOT:
			entity.TakeDMG(types.ActionDamage{Damage: []types.Damage{{Value: effect.Value, Type: types.DMG_TRUE}}})
		case types.EFFECT_HEAL:
			entity.Heal(effect.Value)
		case types.EFFECT_MANA_RESTORE:
			entity.RestoreMana(effect.Value)
		}

		effects = append(effects, effect)
	}

	return effects
}

func GetEffectByType[T types.Entity](effect types.Effect, entity T) *types.ActionEffect {
	for _, eff := range entity.GetAllEffects() {
		if eff.Effect == effect {
			return &eff
		}
	}

	return nil
}

func GetEffectByUUID[T types.Entity](eUuid uuid.UUID, entity T) *types.ActionEffect {
	for _, eff := range entity.GetAllEffects() {
		if eff.Uuid == eUuid {
			return &eff
		}
	}

	return nil
}

func RemoveEffect[T types.Entity](eUuid uuid.UUID, entity T) []types.ActionEffect {
	effects := make([]types.ActionEffect, 0)

	for _, eff := range entity.GetAllEffects() {
		if eff.Uuid == eUuid {
			continue
		}

		effects = append(effects, eff)
	}

	return effects
}

func Heal[T types.Entity](entity T, value int) {
	if entity.GetStat(types.STAT_HEAL_POWER) != 0 {
		value = utils.PercentOf(value, 100+entity.GetStat(types.STAT_HEAL_POWER))
	}

	entity.ChangeHP(value)
}

func DefaultAction[T types.Entity](f types.FightInstance, entity T) []types.Action {
	enemies := f.GetEnemiesFor(entity.GetUUID())

	if len(enemies) == 0 {
		return []types.Action{}
	}

	tauntEffect := entity.GetEffectByType(types.EFFECT_TAUNTED)

	if tauntEffect != nil {
		return []types.Action{{
			Event: types.ACTION_ATTACK, Source: entity.GetUUID(), Target: tauntEffect.Meta.(uuid.UUID),
		}}
	}

	return []types.Action{{
		Event: types.ACTION_ATTACK, Source: entity.GetUUID(), Target: utils.RandomElement(enemies).GetUUID(),
	}}
}
