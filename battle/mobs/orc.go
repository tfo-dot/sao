package mobs

import (
	"sao/battle"
	"sao/utils"

	"github.com/google/uuid"
)

type Orc struct {
	HP      int
	UUID    uuid.UUID
	Meta    MobMeta
	isGroup bool
}

func NewOrc() *Orc {
	meta := MobMeta{
		Name:    MOB_ORC,
		HP:      150,
		SPD:     40,
		ATK:     60,
		Effects: make([]battle.ActionEffect, 0),
	}

	return &Orc{
		HP:      meta.HP,
		UUID:    uuid.New(),
		Meta:    meta,
		isGroup: false,
	}
}

func NewOrcGroup(num int) []Orc {
	if num < 2 {
		num = 2
	}

	group := make([]Orc, num)

	for i := 0; i < num; i++ {
		meta := MobMeta{
			Name:    MOB_ORC,
			HP:      150 + ((num - 2) * 10),
			SPD:     40,
			ATK:     60 + ((num - 2) * 10),
			Effects: make([]battle.ActionEffect, 0),
		}

		group[i] = Orc{
			HP:      meta.HP,
			UUID:    uuid.New(),
			Meta:    meta,
			isGroup: true,
		}
	}

	return group
}

func (o *Orc) GetName() string {
	return string(o.Meta.Name)
}

func (o *Orc) GetMaxHP() int {
	return o.Meta.HP
}

func (o *Orc) GetCurrentHP() int {
	return o.HP
}

func (o *Orc) Heal(val int) {
	o.HP += val

	if o.HP > o.GetMaxHP() {
		o.HP = o.GetMaxHP()
	}
}

func (o *Orc) GetATK() int {
	return o.Meta.ATK
}

func (o *Orc) GetSPD() int {
	return o.Meta.SPD
}

func (o *Orc) IsAuto() bool {
	return true
}

func (o *Orc) Action(f *battle.Fight) int {
	enemiesList := f.GetEnemiesFor(o.GetUUID())

	if len(enemiesList) == 0 {
		return 0
	}

	enemy := utils.RandomElement[battle.Entity](enemiesList)

	f.ActionChannel <- battle.Action{
		Event:  battle.ACTION_ATTACK,
		Source: o.GetUUID(),
		Target: enemy.GetUUID(),
		Meta: battle.ActionMetaFromList([]battle.Damage{{
			Value: o.GetATK() / 2,
			Type:  battle.DMG_TRUE,
		}, {
			Value: o.GetATK() / 2,
			Type:  battle.DMG_PHYSICAL,
		}}, true),
	}

	return 1
}

func (o *Orc) TakeDMG(dmg battle.ActionDamage) int {
	currentHP := o.HP

	for _, dmg := range dmg.Damage {
		o.HP -= dmg.Value
	}

	return currentHP - o.HP
}

func (o *Orc) GetUUID() uuid.UUID {
	return o.UUID
}

func (o *Orc) GetLoot() []battle.Loot {
	base := 50

	if o.isGroup {
		return []battle.Loot{{
			Type: battle.LOOT_EXP, Meta: &map[string]interface{}{"value": base * 2},
		}}
	}

	return []battle.Loot{{
		Type: battle.LOOT_EXP, Meta: &map[string]interface{}{"value": base},
	}}
}

func (o *Orc) CanDodge() bool {
	return false
}

func (o *Orc) ApplyEffect(e battle.ActionEffect) {
	o.Meta.Effects = append(o.Meta.Effects, e)
}

func (o *Orc) HasEffect(e battle.Effect) bool {
	return o.Meta.Effects.HasEffect(e)
}

func (o *Orc) GetEffect(effect battle.Effect) *battle.ActionEffect {
	return o.Meta.Effects.GetEffect(effect)
}

func (o *Orc) TriggerAllEffects() {
	o.Meta.Effects = o.Meta.Effects.TriggerAllEffects(o)
}

func (o *Orc) GetAllEffects() []battle.ActionEffect {
	return o.Meta.Effects
}

func (o *Orc) GetDEF() int {
	return 0
}

func (o *Orc) GetMR() int {
	return 0
}

func (o *Orc) GetDGD() int {
	return 0
}

func (o *Orc) GetMaxMana() int {
	return 0
}

func (o *Orc) GetCurrentMana() int {
	return 0
}

func (o *Orc) GetAP() int {
	return 0
}
