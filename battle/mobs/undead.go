package mobs

import (
	"sao/battle"
	"sao/utils"

	"github.com/google/uuid"
)

type Undead struct {
	HP   int
	UUID uuid.UUID
	Meta MobMeta
	Turn int
}

func NewUndead() *Undead {
	meta := MobMeta{
		Name:    MOB_UNDEAD,
		HP:      70,
		SPD:     40,
		ATK:     50,
		Effects: make([]battle.ActionEffect, 0),
	}

	return &Undead{
		HP:   meta.HP,
		UUID: uuid.New(),
		Meta: meta,
	}
}

func (u *Undead) GetName() string {
	return string(u.Meta.Name)
}

func (u *Undead) GetCurrentHP() int {
	return u.HP
}

func (u *Undead) GetMaxHP() int {
	return u.Meta.HP
}

func (u *Undead) GetATK() int {
	return u.Meta.ATK
}

func (u *Undead) GetSPD() int {
	return u.Meta.SPD
}

func (u *Undead) IsAuto() bool {
	return true
}

func (u *Undead) Action(f *battle.Fight) int {
	enemiesList := f.GetEnemiesFor(u.GetUUID())

	if len(enemiesList) == 0 {
		return 0
	}

	u.Turn++

	enemy := utils.RandomElement[battle.Entity](enemiesList)

	f.ActionChannel <- battle.Action{
		Event:  battle.ACTION_ATTACK,
		Source: u.GetUUID(),
		Target: enemy.GetUUID(),
		Meta:   battle.Damage{Value: u.GetATK(), Type: battle.DMG_PHYSICAL, CanDodge: true}.ToActionMeta(),
	}

	if u.Turn == 4 {
		f.ActionChannel <- battle.Action{
			Event:  battle.ACTION_EFFECT,
			Source: u.GetUUID(),
			Target: enemy.GetUUID(),
			Meta: battle.ActionEffect{
				Effect:   battle.EFFECT_POISON,
				Value:    10,
				Duration: 3,
			},
		}
		u.Turn = 0

		return 2
	}

	return 1
}

func (u *Undead) TakeDMG(dmg battle.ActionDamage) int {
	currentHP := u.HP

	for _, dmg := range dmg.Damage {
		u.HP -= dmg.Value
	}

	return currentHP - u.HP
}

func (u *Undead) GetUUID() uuid.UUID {
	return u.UUID
}

func (u *Undead) GetLoot() []battle.Loot {
	return []battle.Loot{{
		Type: battle.LOOT_EXP, Meta: &map[string]interface{}{"value": 55},
	}}
}

func (u *Undead) CanDodge() bool {
	return false
}

func (u *Undead) ApplyEffect(e battle.ActionEffect) {
	u.Meta.Effects = append(u.Meta.Effects, e)
}

func (u *Undead) HasEffect(e battle.Effect) bool {
	return u.Meta.Effects.HasEffect(e)
}

func (u *Undead) GetEffect(effect battle.Effect) *battle.ActionEffect {
	return u.Meta.Effects.GetEffect(effect)
}

func (u *Undead) TriggerAllEffects() {
	u.Meta.Effects = u.Meta.Effects.TriggerAllEffects(u)
}

func (u *Undead) Heal(value int) {
	u.HP += value
}

func (u *Undead) GetAllEffects() []battle.ActionEffect {
	return u.Meta.Effects
}

func (u *Undead) GetDEF() int {
	return 0
}

func (u *Undead) GetMR() int {
	return 0
}

func (u *Undead) GetDGD() int {
	return 0
}

func (u *Undead) GetMaxMana() int {
	return 0
}

func (u *Undead) GetCurrentMana() int {
	return 0
}

func (u *Undead) GetAP() int {
	return 0
}
