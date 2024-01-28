package mobs

import (
	"sao/battle"

	"github.com/google/uuid"
)

type MobEntity struct {
	//ID as mob type
	Id      string
	MaxHP   int
	HP      int
	SPD     int
	ATK     int
	Effects EffectList
	UUID    uuid.UUID
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

func (m *MobEntity) GetName() string {
	//get name from db
	return m.Id
}

func (m *MobEntity) GetCurrentHP() int {
	return m.HP
}

func (m *MobEntity) GetMaxHP() int {
	return m.HP
}

func (m *MobEntity) GetATK() int {
	return m.ATK
}

func (m *MobEntity) GetSPD() int {
	return m.SPD
}

func (m *MobEntity) GetDEF() int {
	return 0
}

func (m *MobEntity) GetMR() int {
	return 0
}

func (m *MobEntity) GetAGL() int {
	return 0
}

func (m *MobEntity) GetMaxMana() int {
	return 0
}

func (m *MobEntity) GetCurrentMana() int {
	return 0
}

func (m *MobEntity) GetAP() int {
	return 0
}

func (m *MobEntity) IsAuto() bool {
	return true
}

func (m *MobEntity) Action(f *battle.Fight) int {
	//TODO Generic action implementation

	return 0
}

func (m *MobEntity) TakeDMG(dmg battle.ActionDamage) int {
	currentHP := m.HP

	for _, dmg := range dmg.Damage {
		m.HP -= dmg.Value
	}

	return currentHP - m.HP
}

func (m *MobEntity) GetUUID() uuid.UUID {
	return m.UUID
}

func (m *MobEntity) GetLoot() []battle.Loot {
	return []battle.Loot{{
		Type: battle.LOOT_EXP, Meta: &map[string]interface{}{"value": 55},
	}}
}

func (m *MobEntity) CanDodge() bool {
	return false
}

func (m *MobEntity) ApplyEffect(e battle.ActionEffect) {
	m.Effects = append(m.Effects, e)
}

func (m *MobEntity) HasEffect(e battle.Effect) bool {
	return m.Effects.HasEffect(e)
}

func (m *MobEntity) GetEffect(effect battle.Effect) *battle.ActionEffect {
	return m.Effects.GetEffect(effect)
}

func (m *MobEntity) GetAllEffects() []battle.ActionEffect {
	return m.Effects
}

func (m *MobEntity) Heal(value int) {
	m.HP += value
}

func (m *MobEntity) TriggerAllEffects() {
	m.Effects = m.Effects.TriggerAllEffects(m)
}
