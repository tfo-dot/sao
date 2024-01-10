package mobs

import (
	"sao/battle"
	"sao/utils"

	"github.com/google/uuid"
)

type Bear struct {
	HP   int
	UUID uuid.UUID
	Meta MobMeta
	Turn int
}

func NewBear() *Bear {
	meta := MobMeta{
		Name: MOB_BEAR,
		HP:   100,
		SPD:  40,
		ATK:  65,
		Effects: append(make([]battle.ActionEffect, 0), battle.ActionEffect{
			Effect:   battle.EFFECT_VAMP,
			Value:    5,
			Duration: -1,
		}),
	}

	return &Bear{
		HP:   meta.HP,
		UUID: uuid.New(),
		Meta: meta,
	}
}

func (b *Bear) GetName() string {
	return string(b.Meta.Name)
}

func (b *Bear) GetCurrentHP() int {
	return b.HP
}

func (b *Bear) GetMaxHP() int {
	return b.Meta.HP
}

func (b *Bear) GetATK() int {
	return b.Meta.ATK
}

func (b *Bear) GetSPD() int {
	return b.Meta.SPD
}

func (b *Bear) GetDEF() int {
	return 0
}

func (b *Bear) GetMR() int {
	return 0
}

func (b *Bear) GetDGD() int {
	return 0
}

func (b *Bear) GetMaxMana() int {
	return 0
}

func (b *Bear) GetCurrentMana() int {
	return 0
}

func (b *Bear) GetAP() int {
	return 0
}

func (b *Bear) IsAuto() bool {
	return true
}

func (b *Bear) Action(f *battle.Fight) int {
	enemiesList := f.GetEnemiesFor(b.GetUUID())

	if len(enemiesList) == 0 {
		return 0
	}

	b.Turn++

	enemy := utils.RandomElement[battle.Entity](enemiesList)

	dmg := b.GetATK()

	if enemy.HasEffect(battle.EFFECT_FEAR) {
		dmg = int(float64(dmg) * 1.1)
	}

	f.ActionChannel <- battle.Action{
		Event:  battle.ACTION_ATTACK,
		Source: b.GetUUID(),
		Target: enemy.GetUUID(),
		Meta:   battle.Damage{Value: dmg, Type: battle.DMG_PHYSICAL, CanDodge: true}.ToActionMeta(),
	}

	if b.Turn == 3 {
		f.ActionChannel <- battle.Action{
			Event:  battle.ACTION_EFFECT,
			Source: b.GetUUID(),
			Target: enemy.GetUUID(),
			Meta: battle.ActionEffect{
				Effect:   battle.EFFECT_FEAR,
				Value:    -1,
				Duration: 1,
			},
		}
		b.Turn = 0

		return 2
	}

	return 1
}

func (b *Bear) TakeDMG(dmg battle.ActionDamage) int {
	currentHP := b.HP

	for _, dmg := range dmg.Damage {
		b.HP -= dmg.Value
	}

	return currentHP - b.HP
}

func (b *Bear) GetUUID() uuid.UUID {
	return b.UUID
}

func (b *Bear) GetLoot() []battle.Loot {
	return []battle.Loot{{
		Type: battle.LOOT_EXP, Meta: &map[string]interface{}{"value": 55},
	}}
}

func (b *Bear) CanDodge() bool {
	return false
}

func (b *Bear) ApplyEffect(e battle.ActionEffect) {
	b.Meta.Effects = append(b.Meta.Effects, e)
}

func (b *Bear) HasEffect(e battle.Effect) bool {
	return b.Meta.Effects.HasEffect(e)
}

func (b *Bear) GetEffect(effect battle.Effect) *battle.ActionEffect {
	return b.Meta.Effects.GetEffect(effect)
}

func (b *Bear) GetAllEffects() []battle.ActionEffect {
	return b.Meta.Effects
}

func (b *Bear) Heal(value int) {
	b.HP += value
}

func (b *Bear) TriggerAllEffects() {
	b.Meta.Effects = b.Meta.Effects.TriggerAllEffects(b)
}
