package mobs

import (
	"sao/battle"
	"sao/utils"
)

type MobType string

const (
	MOB_BEAR   MobType = "Niedźwiedź"
	MOB_TIGER  MobType = "Tygrys"
	MOB_GNAR   MobType = "Gnar"
	MOB_ORC    MobType = "Ork"
	MOB_UNDEAD MobType = "Nieumarły"
)

type MobMeta struct {
	Name    MobType
	HP      int
	SPD     int
	ATK     int
	Effects EffectList
	Props   map[string]interface{}
}

type EffectList []battle.ActionEffect

func (e EffectList) HasEffect(effect battle.Effect) bool {
	for _, eff := range e {
		if eff.Effect == effect {
			return true
		}
	}

	return false
}

func (e EffectList) GetEffect(effect battle.Effect) *battle.ActionEffect {
	for _, eff := range e {
		if eff.Effect == effect {
			return &eff
		}
	}

	return nil
}

func (e EffectList) TriggerAllEffects(en battle.Entity) EffectList {
	effects := make([]battle.ActionEffect, 0)

	for _, effect := range e {
		if effect.Duration > 0 {
			effect.Duration--
		}

		switch effect.Effect {
		case battle.EFFECT_POISON:
			en.TakeDMG(battle.Damage{Value: effect.Value, Type: battle.DMG_TRUE, CanDodge: false}.ToActionMeta())
		case battle.EFFECT_HEAL:
			en.Heal(effect.Value)
		}

		if effect.Duration == 0 {
			continue
		}

		effects = append(effects, effect)
	}

	return effects
}

func MobEncounter(mobType MobType) []battle.Entity {
	switch mobType {
	case MOB_BEAR:
		return []battle.Entity{NewBear()}
	case MOB_GNAR:
		return []battle.Entity{NewGnar()}
	case MOB_ORC:
		num := utils.RandomNumber(0, 4)
		if num < 2 {
			return []battle.Entity{NewOrc()}
		} else {
			group := NewOrcGroup(num)

			entities := make([]battle.Entity, len(group))
			for i, orc := range group {
				entities[i] = &orc
			}

			return entities
		}
	case MOB_TIGER:
		return []battle.Entity{NewTiger()}
	case MOB_UNDEAD:
		return []battle.Entity{NewUndead()}
	}

	return []battle.Entity{}
}
