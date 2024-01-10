package mobs

import (
	"sao/battle"
	"sao/utils"

	"github.com/google/uuid"
)

type Gnar struct {
	HP   int
	UUID uuid.UUID
	Meta MobMeta
	Turn int
}

func NewGnar() *Gnar {
	meta := MobMeta{
		Name:    MOB_GNAR,
		HP:      150,
		SPD:     40,
		ATK:     60,
		Effects: make([]battle.ActionEffect, 0),
	}

	return &Gnar{
		HP:   meta.HP,
		UUID: uuid.New(),
		Meta: meta,
	}
}

func (g *Gnar) GetName() string {
	return string(g.Meta.Name)
}

func (g *Gnar) GetMaxHP() int {
	return g.Meta.HP
}

func (g *Gnar) GetCurrentHP() int {
	return g.HP
}

func (g *Gnar) Heal(val int) {
	g.HP += val

	if g.HP > g.GetMaxHP() {
		g.HP = g.GetMaxHP()
	}
}

func (g *Gnar) GetATK() int {
	return g.Meta.ATK
}

func (g *Gnar) GetSPD() int {
	return g.Meta.SPD
}

func (g *Gnar) IsAuto() bool {
	return true
}

func (g *Gnar) Action(f *battle.Fight) int {
	enemiesList := f.GetEnemiesFor(g.GetUUID())

	if len(enemiesList) == 0 {
		return 0
	}

	enemy := utils.RandomElement[battle.Entity](enemiesList)

	g.Turn++

	dmg := g.GetATK()

	if g.Turn%3 == 0 {
		dmg += 15 + int(float64(enemy.GetMaxHP())*0.05)
	}

	f.ActionChannel <- battle.Action{
		Event:  battle.ACTION_ATTACK,
		Source: g.GetUUID(),
		Target: enemy.GetUUID(),
		Meta: battle.Damage{
			Value: dmg,
			Type:  battle.DMG_PHYSICAL,
		}.ToActionMeta(),
	}

	return 1
}

func (g *Gnar) TakeDMG(dmg battle.ActionDamage) int {
	currentHP := g.HP

	for _, dmg := range dmg.Damage {
		actualDmg := dmg.Value

		//Passive skill (not trigger on poison or sth else)
		if dmg.Value > 75 && dmg.Type == battle.DMG_PHYSICAL {
			actualDmg -= actualDmg / 2
		}

		g.HP -= actualDmg
	}

	return currentHP - g.HP
}

func (g *Gnar) GetUUID() uuid.UUID {
	return g.UUID
}

func (g *Gnar) GetLoot() []battle.Loot {
	return []battle.Loot{{
		Type: battle.LOOT_EXP, Meta: &map[string]interface{}{"value": 90},
	}}
}

func (g *Gnar) CanDodge() bool {
	return false
}

func (g *Gnar) ApplyEffect(e battle.ActionEffect) {
	g.Meta.Effects = append(g.Meta.Effects, e)
}

func (g *Gnar) HasEffect(e battle.Effect) bool {
	return g.Meta.Effects.HasEffect(e)
}

func (g *Gnar) GetEffect(effect battle.Effect) *battle.ActionEffect {
	return g.Meta.Effects.GetEffect(effect)
}

func (g *Gnar) TriggerAllEffects() {
	g.Meta.Effects = g.Meta.Effects.TriggerAllEffects(g)
}

func (g *Gnar) GetAllEffects() []battle.ActionEffect {
	return g.Meta.Effects
}

func (g *Gnar) GetDEF() int {
	return 0
}

func (g *Gnar) GetMR() int {
	return 0
}

func (g *Gnar) GetDGD() int {
	return 0
}

func (g *Gnar) GetMaxMana() int {
	return 0
}

func (g *Gnar) GetCurrentMana() int {
	return 0
}

func (g *Gnar) GetAP() int {
	return 0
}
