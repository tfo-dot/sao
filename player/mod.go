package player

import (
	"sao/battle"
	"sao/battle/mobs"
	"sao/navel"
	"sao/player/inventory"
	"sao/player/location"
	"sao/player/xp"
	"sao/types"
	"sao/utils"

	"github.com/google/uuid"
)

type PlayerGender byte

const (
	MALE PlayerGender = iota
	FEMALE
)

type PlayerStats struct {
	HP          int
	SPD         int
	DGD         int
	Effects     mobs.EffectList
	CurrentMana int
}

type PlayerMeta struct {
	Gender        PlayerGender
	Location      location.PlayerLocation
	OwnUUID       uuid.UUID
	UserID        string
	Navel         *navel.Navel
	FightInstance *uuid.UUID
}

type Player struct {
	Name      string
	XP        xp.PlayerXP
	Stats     PlayerStats
	Meta      PlayerMeta
	Inventory inventory.PlayerInventory
}

func (p *Player) GetCurrentMana() int {
	return p.Stats.CurrentMana
}

func (p *Player) GetMaxMana() int {
	return 10 + p.Inventory.GetStat(battle.STAT_MANA)
}

func (p *Player) GetUUID() uuid.UUID {
	return p.Meta.OwnUUID
}

func (p *Player) GetCurrentHP() int {
	return p.Stats.HP
}

func (p *Player) GetMaxHP() int {
	return 100 + ((p.XP.Level - 1) * 10) + p.Inventory.GetStat(battle.STAT_HP)
}

func (p *Player) Heal(val int) {
	p.Stats.HP += val

	if p.Stats.HP > p.GetMaxHP() {
		p.Stats.HP = p.GetMaxHP()
	}
}

func (p *Player) GetATK() int {
	return 40 + ((p.XP.Level - 1) * 15) + p.Inventory.GetStat(battle.STAT_AD)
}

func (p *Player) GetSPD() int {
	return p.Stats.SPD + p.Inventory.GetStat(battle.STAT_SPD)
}

func (p *Player) IsAuto() bool {
	return false
}

func (p *Player) GetDEF() int {
	return p.Inventory.GetStat(battle.STAT_DEF)
}

func (p *Player) GetMR() int {
	return p.Inventory.GetStat(battle.STAT_MR)
}

func (p *Player) GetDGD() int {
	return p.Stats.DGD + p.Inventory.GetStat(battle.STAT_DGD)
}

func (p *Player) GetAP() int {
	return p.Inventory.GetStat(battle.STAT_AP)
}

func (p *Player) getDGD() int {
	return p.Inventory.GetStat(battle.STAT_DGD) + p.Stats.DGD
}

func (p *Player) Action(f *battle.Fight) int {
	return 0
}

func (p *Player) TakeDMG(dmgList battle.ActionDamage) int {
	currentHP := p.Stats.HP

	for _, dmg := range dmgList.Damage {
		switch dmg.Type {
		case battle.DMG_PHYSICAL:
			dmg.Value -= utils.CalcReducedDamage(dmg.Value, p.GetDEF())
		case battle.DMG_MAGICAL:
			dmg.Value -= utils.CalcReducedDamage(dmg.Value, p.GetMR())
		}

		p.Stats.HP -= dmg.Value
	}

	//DMG TAKEN NOT THE SAME AS DMG DEALT
	return currentHP - p.Stats.HP
}

func (p *Player) TakeDMGOrDodge(dmg battle.ActionDamage) (int, bool) {
	if utils.RandomNumber(0, 100) <= p.getDGD() && dmg.CanDodge {
		return 0, true
	}

	return p.TakeDMG(dmg), false
}

func (p *Player) AddEXP(value int) {
	p.XP.Exp += value

	for p.XP.Exp >= ((p.XP.Level * 100) + 100) {
		p.XP.Exp -= ((p.XP.Level * 100) + 100)
		p.LevelUP()
	}
}

func (p *Player) LevelUP() {
	p.XP.Level++

	if p.Stats.HP < p.GetMaxHP() {
		p.Stats.HP = p.GetMaxHP() - (p.Stats.HP / 5)
	}
}

func (p *Player) AddGold(value int) {
	p.Inventory.Gold += value
}

func (p *Player) GetLoot() []battle.Loot {
	return nil
}

func (p *Player) ReceiveLoot(loot battle.Loot) {
	switch loot.Type {
	case battle.LOOT_GOLD:
		p.AddGold((*loot.Meta)["value"].(int))
	case battle.LOOT_EXP:
		p.AddEXP((*loot.Meta)["value"].(int))
	}
}

func (p *Player) ReceiveMultipleLoot(loot []battle.Loot) {
	for _, l := range loot {
		p.ReceiveLoot(l)
	}
}

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) CanDodge() bool {
	return true
}

func (p *Player) ApplyEffect(e battle.ActionEffect) {
	p.Stats.Effects = append(p.Stats.Effects, e)
}

func (p *Player) HasEffect(e battle.Effect) bool {
	return p.Stats.Effects.HasEffect(e)
}

func (p *Player) GetEffect(effect battle.Effect) *battle.ActionEffect {
	return p.Stats.Effects.GetEffect(effect)
}

func (p *Player) TriggerAllEffects() {
	p.Stats.Effects = p.Stats.Effects.TriggerAllEffects(p)
}

func (p *Player) GetAllEffects() []battle.ActionEffect {
	return p.Stats.Effects
}

func (p *Player) CanAttack() bool {
	return !(p.HasEffect(battle.EFFECT_DISARM) || p.HasEffect(battle.EFFECT_STUN))
}

func (p *Player) CanDodgeNow() bool {
	return !(p.HasEffect(battle.EFFECT_STUN) || p.HasEffect(battle.EFFECT_ROOT) || p.HasEffect(battle.EFFECT_GROUND) || p.HasEffect(battle.EFFECT_BLIND))
}

func (p *Player) CanDefend() bool {
	return !(p.HasEffect(battle.EFFECT_STUN) || p.HasEffect(battle.EFFECT_ROOT))
}

func (p *Player) CanUseSkill(skill types.PlayerSkill) bool {
	if skill.Trigger.Type == types.TRIGGER_PASSIVE {
		return false
	}

	if p.HasEffect(battle.EFFECT_SILENCE) {
		return false
	}

	if p.Inventory.CDS[skill.UUID] > 0 {
		return false
	}

	switch skill.Cost.Resource {
	case types.ManaResource:
		if p.GetCurrentMana() < skill.Cost.Cost {
			return false
		}
	}

	return true
}

func (p *Player) GetAvailableActions() []battle.ActionPartial {
	actions := make([]battle.ActionPartial, 0)

	if p.CanAttack() {
		actions = append(actions, battle.ActionPartial{Event: battle.ACTION_ATTACK, Meta: nil})
	}

	if p.CanDodgeNow() {
		actions = append(actions, battle.ActionPartial{Event: battle.ACTION_DODGE, Meta: nil})
	}

	if p.CanDefend() {
		actions = append(actions, battle.ActionPartial{Event: battle.ACTION_DEFEND, Meta: nil})
	}

	for _, skill := range p.Inventory.Skills {
		if skill.Trigger.Type == types.TRIGGER_ACTIVE && p.CanUseSkill(skill) {
			actions = append(actions, battle.ActionPartial{Event: battle.ACTION_SKILL, Meta: &skill.UUID})
		}
	}

	return actions
}

func (p *Player) GetAllSkills() []types.PlayerSkill {

	tempArr := p.Inventory.Skills

	for _, item := range p.Inventory.Items {
		for _, effect := range item.Effects {

			tempArr = append(tempArr, types.PlayerSkill{
				Name:    effect.Name,
				Trigger: effect.Trigger,
				Cost:    types.SkillCost{Cost: 0, Resource: types.ManaResource},
				UUID:    item.UUID,
				Grade:   types.GradeCommon,
				Action:  effect.Execute,
			})
		}
	}

	return tempArr
}

func NewPlayer(gender PlayerGender, name string, uid string) Player {
	return Player{
		name,
		xp.PlayerXP{Level: 1, Exp: 0},
		PlayerStats{100, 40, 50, make(mobs.EffectList, 0), 0},
		PlayerMeta{gender, location.DefaultLocation(), uuid.New(), uid, nil, nil},
		inventory.GetDefaultInventory(),
	}
}
