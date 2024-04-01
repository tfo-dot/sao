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

type PlayerStats struct {
	HP          int
	SPD         int
	AGL         int
	Effects     mobs.EffectList
	Defending   bool
	CurrentMana int
}

type PlayerMeta struct {
	Location      location.PlayerLocation
	OwnUUID       uuid.UUID
	UserID        string
	Navel         *navel.Navel
	FightInstance *uuid.UUID
	Party         *uuid.UUID
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

func (p *Player) GetAGL() int {
	return p.Stats.AGL + p.Inventory.GetStat(battle.STAT_AGL)
}

func (p *Player) GetAP() int {
	return p.Inventory.GetStat(battle.STAT_AP)
}

func (p *Player) SetDefendingState(state bool) {
	p.Stats.Defending = state
}

func (p *Player) GetDefendingState() bool {
	return p.Stats.Defending
}

func (p *Player) Action(f *battle.Fight) int {
	return 0
}

func (p *Player) TakeDMG(dmgList battle.ActionDamage) int {
	startingHP := p.Stats.HP

	for _, dmg := range dmgList.Damage {
		//Skip shield and such
		if dmg.Type == battle.DMG_TRUE {
			p.Stats.HP -= dmg.Value
			continue
		}

		rawDmg := dmg.Value

		switch dmg.Type {
		case battle.DMG_PHYSICAL:
			rawDmg = utils.CalcReducedDamage(dmg.Value, p.GetDEF())
		case battle.DMG_MAGICAL:
			rawDmg = utils.CalcReducedDamage(dmg.Value, p.GetMR())
		}

		p.Stats.HP -= p.DamageShields(rawDmg)
	}

	return startingHP - p.Stats.HP
}

func (p *Player) TakeDMGOrDodge(dmg battle.ActionDamage) (int, bool) {
	if utils.RandomNumber(0, 100) <= p.GetAGL() && dmg.CanDodge {
		return 0, true
	}

	return p.TakeDMG(dmg), false
}

func (p *Player) DamageShields(dmg int) int {
	leftOverDmg := dmg
	idxToRemove := make([]int, 0)

	for idx, effect := range p.Stats.Effects {
		if effect.Effect == battle.EFFECT_SHIELD {
			newShieldValue := effect.Value - leftOverDmg

			if newShieldValue <= 0 {
				leftOverDmg = newShieldValue * -1

				idxToRemove = append(idxToRemove, idx)
			} else {
				effect.Value = newShieldValue
				leftOverDmg = 0
			}
		}
	}

	for _, idx := range idxToRemove {
		p.Stats.Effects = append(p.Stats.Effects[:idx], p.Stats.Effects[idx+1:]...)
	}

	return leftOverDmg
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
	// If player is defending, he can't dodge (counterattack is possible)
	return !p.Stats.Defending
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

func (p *Player) TriggerAllEffects() []battle.ActionEffect {
	effects, expiredEffects := p.Stats.Effects.TriggerAllEffects(p)

	p.Stats.Effects = effects

	return expiredEffects
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

	if p.GetCurrentMana() < skill.Cost {
		return false
	}

	return true
}

func (p *Player) CanUseLvlSkill(skill inventory.PlayerSkill) bool {
	if skill.Trigger.Type == types.TRIGGER_PASSIVE {
		return false
	}

	if p.HasEffect(battle.EFFECT_SILENCE) {
		return false
	}

	if p.Inventory.LevelSkillsCDS[skill.ForLevel] > 0 {
		return false
	}

	if p.GetCurrentMana() < skill.Cost {
		return false
	}

	return true
}

func (p *Player) GetAllSkills() []types.PlayerSkill {
	tempArr := p.Inventory.Skills

	for _, item := range p.Inventory.Items {
		for _, effect := range item.Effects {

			tempArr = append(tempArr, &types.PlayerSkill{
				Name:    effect.Name,
				Trigger: effect.Trigger,
				Cost:    0,
				UUID:    item.UUID,
				Action: func(source interface{}, target interface{}, fight interface{}) {
					effect.Execute(source, target, fight.(*interface{}))
				},
			})
		}
	}

	for _, skill := range p.Inventory.LevelSkills {
		tempArr = append(tempArr, &types.PlayerSkill{
			Name:        skill.Name,
			Description: skill.Description,
			Trigger:     *skill.Trigger,
			Cost:        skill.Cost,
			UUID:        uuid.Nil,
			Action: func(source interface{}, target interface{}, fight interface{}) {
				skill.Execute(source.(battle.PlayerEntity), target.(battle.Entity), fight.(*battle.Fight))
			},
			CD: skill.CD.Calc(skill, p.GetUpgrades(skill.ForLevel)),
		})
	}

	return []types.PlayerSkill{}
}

func (p *Player) AddItem(item *types.PlayerItem) {
	p.Inventory.Items = append(p.Inventory.Items, item)
}

func (p *Player) GetAllItems() []*types.PlayerItem {
	return p.Inventory.Items
}

func (p *Player) RemoveItem(item int) {
	p.Inventory.Items = append(p.Inventory.Items[:item], p.Inventory.Items[item+1:]...)
}

func (p *Player) RestoreMana(value int) {
	p.Stats.CurrentMana += value

	if p.Stats.CurrentMana > p.GetMaxMana() {
		p.Stats.CurrentMana = p.GetMaxMana()
	}
}

func (p *Player) GetStat(stat battle.Stat) int {
	statValue := 0
	percentValue := 0

	for _, effect := range p.GetAllEffects() {
		if effect.Effect == battle.EFFECT_STAT_INC {

			if value, ok := effect.Meta.(battle.ActionEffectStat); ok {
				if value.Stat != stat {
					continue
				}

				if value.IsPercent {
					percentValue += value.Value
				} else {
					statValue += value.Value
				}
			}
		}

		if effect.Effect == battle.EFFECT_STAT_DEC {

			if value, ok := effect.Meta.(battle.ActionEffectStat); ok {
				if value.Stat != stat {
					continue
				}

				if value.IsPercent {
					percentValue -= value.Value
				} else {
					statValue -= value.Value
				}
			}
		}
	}

	tempValue := statValue

	switch stat {
	case battle.STAT_HP:
		tempValue += p.GetMaxHP()
	case battle.STAT_AD:
		tempValue += p.GetATK()
	case battle.STAT_SPD:
		tempValue += p.GetSPD()
	case battle.STAT_AGL:
		tempValue += p.GetAGL()
	case battle.STAT_AP:
		tempValue += p.GetAP()
	case battle.STAT_DEF:
		tempValue += p.GetDEF()
	case battle.STAT_MR:
		tempValue += p.GetMR()
	case battle.STAT_MANA:
		tempValue += p.GetCurrentMana()
	}

	return tempValue + (tempValue * percentValue / 100)
}

func (p *Player) Cleanse() {
	p.Stats.Effects.Cleanse()
}

func (p *Player) GetUpgrades(lvl int) []string {
	return p.Inventory.LevelSkillsUpgrades[lvl]
}

func (p *Player) GetLvlSkill(lvl int) *types.PlayerSkill {
	for _, skill := range p.Inventory.LevelSkills {
		if skill.ForLevel == lvl {

			return &types.PlayerSkill{
				Name:        skill.Name,
				Trigger:     *skill.Trigger,
				Description: skill.Description,
				Cost:        skill.Cost,
				UUID:        uuid.Nil,
				Action: func(source interface{}, target interface{}, fight interface{}) {
					skill.Execute(source.(battle.PlayerEntity), target.(battle.Entity), fight.(*battle.Fight))
				},
				CD: skill.CD.Calc(skill, p.GetUpgrades(lvl)),
			}
		}
	}

	return nil
}

func (p *Player) GetSkill(skillUUID uuid.UUID) *types.PlayerSkill {
	for _, skill := range p.Inventory.Skills {
		if skill.UUID == skillUUID {
			return skill
		}
	}

	return nil
}

func (p *Player) GetUID() string {
	return p.Meta.UserID
}

func (p *Player) SetCD(skillUUID uuid.UUID, value int) {
	p.Inventory.CDS[skillUUID] = value
}

func (p *Player) GetCD(skillUUID uuid.UUID) int {
	return p.Inventory.CDS[skillUUID]
}

func (p *Player) GetParty() *uuid.UUID {
	return p.Meta.Party
}

func NewPlayer(name string, uid string) Player {
	return Player{
		name,
		xp.PlayerXP{Level: 1, Exp: 0},
		PlayerStats{100, 40, 50, make(mobs.EffectList, 0), false, 10},
		PlayerMeta{location.DefaultLocation(), uuid.New(), uid, nil, nil, nil},
		inventory.GetDefaultInventory(),
	}
}
