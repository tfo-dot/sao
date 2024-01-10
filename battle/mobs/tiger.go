package mobs

import (
	"sao/battle"
	"sao/utils"

	"github.com/google/uuid"
)

type Tiger struct {
	HP   int
	UUID uuid.UUID
	Meta MobMeta
}

func NewTiger() *Tiger {
	meta := MobMeta{
		Name:    MOB_TIGER,
		HP:      50,
		SPD:     40,
		ATK:     40,
		Effects: make([]battle.ActionEffect, 0),
	}

	return &Tiger{
		HP:   meta.HP,
		UUID: uuid.New(),
		Meta: meta,
	}
}

func (t *Tiger) GetName() string {
	return string(t.Meta.Name)
}

func (t *Tiger) GetMaxHP() int {
	return t.Meta.HP
}

func (t *Tiger) GetCurrentHP() int {
	return t.HP
}

func (t *Tiger) Heal(val int) {
	t.HP += val

	if t.HP > t.GetMaxHP() {
		t.HP = t.GetMaxHP()
	}
}

func (t *Tiger) GetATK() int {
	return t.Meta.ATK
}

func (t *Tiger) GetSPD() int {
	return t.Meta.SPD
}

func (t *Tiger) IsAuto() bool {
	return true
}

func (t *Tiger) Action(f *battle.Fight) int {
	enemiesList := f.GetEnemiesFor(t.GetUUID())

	if len(enemiesList) == 0 {
		return 0
	}

	enemy := utils.RandomElement[battle.Entity](enemiesList)

	f.ActionChannel <- battle.Action{
		Event:  battle.ACTION_ATTACK,
		Source: t.GetUUID(),
		Target: enemy.GetUUID(),
		Meta: battle.ActionMetaFromList([]battle.Damage{{
			Value: t.GetATK() / 2,
			Type:  battle.DMG_TRUE,
		}, {
			Value: t.GetATK() / 2,
			Type:  battle.DMG_PHYSICAL,
		}}, true),
	}

	return 1
}

func (t *Tiger) TakeDMG(dmg battle.ActionDamage) int {
	currentHP := t.HP

	for _, dmg := range dmg.Damage {
		t.HP -= dmg.Value
	}

	return currentHP - t.HP
}

func (t *Tiger) GetUUID() uuid.UUID {
	return t.UUID
}

func (t *Tiger) GetLoot() []battle.Loot {
	return []battle.Loot{{
		Type: battle.LOOT_EXP, Meta: &map[string]interface{}{"value": 40},
	}}
}

func (t *Tiger) CanDodge() bool {
	return false
}

func (t *Tiger) ApplyEffect(e battle.ActionEffect) {
	t.Meta.Effects = append(t.Meta.Effects, e)
}

func (t *Tiger) HasEffect(e battle.Effect) bool {
	return t.Meta.Effects.HasEffect(e)
}

func (t *Tiger) GetEffect(effect battle.Effect) *battle.ActionEffect {
	return t.Meta.Effects.GetEffect(effect)
}

func (t *Tiger) TriggerAllEffects() {
	t.Meta.Effects = t.Meta.Effects.TriggerAllEffects(t)
}

func (t *Tiger) GetAllEffects() []battle.ActionEffect {
	return t.Meta.Effects
}

func (t *Tiger) GetDEF() int {
	return 0
}

func (t *Tiger) GetMR() int {
	return 0
}

func (t *Tiger) GetDGD() int {
	return 0
}

func (t *Tiger) GetMaxMana() int {
	return 0
}

func (t *Tiger) GetCurrentMana() int {
	return 0
}

func (t *Tiger) GetAP() int {
	return 0
}
