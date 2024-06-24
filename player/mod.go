package player

import (
	"sao/battle"
	"sao/battle/mobs"
	"sao/player/inventory"
	"sao/types"
	"sao/utils"
	"sao/world/fury"

	"github.com/google/uuid"
)

type PlayerStats struct {
	HP          int
	Effects     mobs.EffectList
	Defending   bool
	CurrentMana int
}

type PlayerXP struct {
	Level int
	Exp   int
}

type PlayerMeta struct {
	Location      types.PlayerLocation
	OwnUUID       uuid.UUID
	UserID        string
	FightInstance *uuid.UUID
	Party         *uuid.UUID
	Transaction   *uuid.UUID
	//Array just in case someone will have multiple furies
	Fury []*fury.Fury
}

func (pM *PlayerMeta) SerializeFuries() []map[string]interface{} {
	furies := make([]map[string]interface{}, 0)

	for _, fury := range pM.Fury {
		furies = append(furies, fury.Serialize())
	}

	return furies
}

func (pM *PlayerMeta) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"location": []string{pM.Location.FloorName, pM.Location.LocationName},
		"uuid":     pM.OwnUUID.String(),
		"uid":      pM.UserID,
		"fury":     pM.SerializeFuries(),
	}
}

type Player struct {
	Name         string
	XP           PlayerXP
	Stats        PlayerStats
	Meta         PlayerMeta
	Inventory    inventory.PlayerInventory
	DynamicStats []types.DerivedStat
}

func (p *Player) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"name": p.Name,
		"xp":   []int{p.XP.Level, p.XP.Exp},
		"stats": map[string]interface{}{
			"hp":           p.Stats.HP,
			"current_mana": p.Stats.CurrentMana,
			"effects":      p.Stats.Effects,
		},
		"meta":      p.Meta.Serialize(),
		"inventory": p.Inventory.Serialize(),
	}
}

func Deserialize(data map[string]interface{}) *Player {
	return &Player{
		data["name"].(string),
		PlayerXP{
			Level: int(data["xp"].([]interface{})[0].(float64)),
			Exp:   int(data["xp"].([]interface{})[1].(float64)),
		},
		PlayerStats{
			int(data["stats"].(map[string]interface{})["hp"].(float64)),
			DeserializeEffects(data["stats"].(map[string]interface{})["effects"].([]interface{})),
			false,
			int(data["stats"].(map[string]interface{})["current_mana"].(float64)),
		},
		*DeserializeMeta(data["meta"].(map[string]interface{})),
		inventory.DeserializeInventory(data["inventory"].(map[string]interface{})),
		make([]types.DerivedStat, 0),
	}
}

func DeserializeEffects(data []interface{}) mobs.EffectList {
	temp := make(mobs.EffectList, 0)

	if len(data) == 0 {
		return temp
	}

	for _, effect := range data {
		effect := effect.(map[string]interface{})

		temp = append(temp, battle.ActionEffect{
			Effect: battle.Effect(effect["effect"].(float64)),
			Value:  int(effect["value"].(float64)),
			Meta:   effect["meta"],
		})
	}

	return temp
}

func DeserializeMeta(data map[string]interface{}) *PlayerMeta {
	deserializedFuries := make([]*fury.Fury, 0)

	for _, furyData := range data["fury"].([]interface{}) {
		deserializedFuries = append(deserializedFuries, fury.Deserialize(furyData.(map[string]interface{})))
	}

	return &PlayerMeta{
		types.DefaultPlayerLocation(),
		uuid.MustParse(data["uuid"].(string)),
		data["uid"].(string),
		nil,
		nil,
		nil,
		deserializedFuries,
	}
}

func (p *Player) GetCurrentMana() int {
	return p.Stats.CurrentMana
}

func (p *Player) GetFuryStat(stat types.Stat) int {
	value := 0

	for _, fury := range p.Meta.Fury {
		statValue, ok := fury.GetStats()[stat]
		if ok {
			value += statValue
		}
	}

	return value
}

func (p *Player) GetMaxMana() int {
	return 10 + p.Inventory.GetStat(types.STAT_MANA) + p.GetFuryStat(types.STAT_MANA)
}

func (p *Player) GetUUID() uuid.UUID {
	return p.Meta.OwnUUID
}

func (p *Player) GetCurrentHP() int {
	return p.Stats.HP
}

func (p *Player) GetMaxHP() int {
	return 100 + ((p.XP.Level - 1) * 10) + p.Inventory.GetStat(types.STAT_HP) + p.GetFuryStat(types.STAT_HP)
}

func (p *Player) Heal(val int) {
	p.Stats.HP += val

	if p.Stats.HP > p.GetMaxHP() {
		p.Stats.HP = p.GetMaxHP()
	}
}

func (p *Player) GetAdaptiveType() types.AdaptiveType {
	atk := p.GetATKWithoutAdaptive()
	ap := p.GetAPWithoutAdaptive()

	if atk > ap || atk == ap {
		return types.ADAPTIVE_ATK
	} else {
		return types.ADAPTIVE_AP
	}
}

func (p *Player) GetATK() int {
	adaptiveAtk := p.Inventory.GetStat(types.STAT_ADAPTIVE)
	if p.GetAdaptiveType() != types.ADAPTIVE_ATK {
		adaptiveAtk = 0
	}

	return 40 + ((p.XP.Level - 1) * 15) + p.GetATKWithoutAdaptive() + adaptiveAtk
}

func (p *Player) GetATKWithoutAdaptive() int {
	return p.Inventory.GetStat(types.STAT_AD) + p.GetFuryStat(types.STAT_AD)
}

func (p *Player) GetAPWithoutAdaptive() int {
	return p.Inventory.GetStat(types.STAT_AP) + p.GetFuryStat(types.STAT_AP)
}

func (p *Player) GetAP() int {
	adaptiveAp := p.Inventory.GetStat(types.STAT_ADAPTIVE)
	if p.GetAdaptiveType() != types.ADAPTIVE_AP {
		adaptiveAp = 0
	}

	return p.Inventory.GetStat(types.STAT_AP) + p.GetFuryStat(types.STAT_AP) + adaptiveAp
}

func (p *Player) GetSPD() int {
	return 40 + p.Inventory.GetStat(types.STAT_SPD) + p.GetFuryStat(types.STAT_SPD)
}

func (p *Player) IsAuto() bool {
	return false
}

func (p *Player) GetDEF() int {
	return p.Inventory.GetStat(types.STAT_DEF) + p.GetFuryStat(types.STAT_DEF)
}

func (p *Player) GetMR() int {
	return p.Inventory.GetStat(types.STAT_MR) + p.GetFuryStat(types.STAT_MR)
}

func (p *Player) GetAGL() int {
	return 50 + p.Inventory.GetStat(types.STAT_AGL) + p.GetFuryStat(types.STAT_AGL)
}

func (p *Player) SetDefendingState(state bool) {
	p.Stats.Defending = state
}

func (p *Player) GetDefendingState() bool {
	return p.Stats.Defending
}

func (p *Player) Action(f *battle.Fight) []battle.Action { return []battle.Action{} }

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

func (p *Player) AddEXP(maxFloor, value int) {
	p.XP.Exp += value

	if p.XP.Level >= maxFloor*5 {

		println("Max level reached 1")

		p.XP.Level = maxFloor * 5
		p.XP.Exp = 0
		return
	}

	for _, fury := range p.Meta.Fury {
		fury.AddXP(utils.PercentOf(value, 20))
	}

	for p.XP.Exp >= ((p.XP.Level * 100) + 100) {
		if p.XP.Level >= maxFloor*5 {

			println("Max level reached 3")

			p.XP.Level = maxFloor * 5
			p.XP.Exp = 0
			return
		}

		p.XP.Exp -= ((p.XP.Level * 100) + 100)
		p.LevelUP()
	}
}

func (p *Player) LevelUP() {
	p.XP.Level++

	if p.Stats.HP < p.GetMaxHP() {
		p.Stats.HP = utils.PercentOf(p.GetMaxHP(), 20) + p.GetCurrentHP()
	}

	if p.Stats.HP > p.GetMaxHP() {
		p.Stats.HP = p.GetMaxHP()
	}
}

func (p *Player) AddGold(value int) {
	p.Inventory.Gold += value
}

func (p *Player) GetLoot() []battle.Loot {
	return nil
}

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) CanDodge() bool {
	return !p.Stats.Defending
}

func (p *Player) ApplyEffect(e battle.ActionEffect) {
	p.Stats.Effects = append(p.Stats.Effects, e)
}

func (p *Player) GetEffectByType(effect battle.Effect) *battle.ActionEffect {
	return p.Stats.Effects.GetEffectByType(effect)
}

func (p *Player) GetEffectByUUID(uuid uuid.UUID) *battle.ActionEffect {
	return p.Stats.Effects.GetEffectByUUID(uuid)
}

func (p *Player) TriggerAllEffects() []battle.ActionEffect {
	effects, expiredEffects := p.Stats.Effects.TriggerAllEffects(p)

	p.Stats.Effects = effects

	return expiredEffects
}

func (p *Player) RemoveEffect(uuid uuid.UUID) {
	p.Stats.Effects = p.Stats.Effects.RemoveEffect(uuid)
}

func (p *Player) GetAllEffects() []battle.ActionEffect {
	return p.Stats.Effects
}

func (p *Player) CanAttack() bool {
	return !(p.GetEffectByType(battle.EFFECT_DISARM) != nil || p.GetEffectByType(battle.EFFECT_STUN) != nil)
}

func (p *Player) CanDodgeNow() bool {
	return !(p.GetEffectByType(battle.EFFECT_STUN) != nil || p.GetEffectByType(battle.EFFECT_ROOT) != nil || p.GetEffectByType(battle.EFFECT_GROUND) != nil || p.GetEffectByType(battle.EFFECT_BLIND) != nil)
}

func (p *Player) CanDefend() bool {
	return !(p.GetEffectByType(battle.EFFECT_STUN) != nil || p.GetEffectByType(battle.EFFECT_ROOT) != nil)
}

func (p *Player) CanUseSkill(skill types.PlayerSkill) bool {
	if skill.GetTrigger().Type == types.TRIGGER_PASSIVE {
		return false
	}

	if p.GetEffectByType(battle.EFFECT_SILENCE) != nil {
		return false
	}

	if skill.IsLevelSkill() {
		if p.GetLvlCD(skill.(types.PlayerSkillLevel).GetLevel()) > 0 {
			return false
		}
	}

	if p.GetCurrentMana() < skill.GetCost() {
		return false
	}

	return true
}

func (p *Player) CanUseLvlSkill(skill inventory.PlayerSkillLevel) bool {
	if skill.GetTrigger().Type == types.TRIGGER_PASSIVE {
		return false
	}

	if p.GetEffectByType(battle.EFFECT_SILENCE) != nil {
		return false
	}

	if p.Inventory.LevelSkillsCDS[skill.GetLevel()] > 0 {
		return false
	}

	if p.GetCurrentMana() < skill.GetCost() {
		return false
	}

	return true
}

func (p *Player) GetAllSkills() []types.PlayerSkill {
	tempArr := make([]types.PlayerSkill, 0)

	for _, item := range p.Inventory.Items {
		tempArr = append(tempArr, item.Effects...)
	}

	for _, skill := range p.Inventory.LevelSkills {
		tempArr = append(tempArr, skill)
	}

	for _, skill := range p.Meta.Fury {
		tempArr = append(tempArr, skill.GetSkills()...)
	}

	return tempArr
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

func (p *Player) GetStat(stat types.Stat) int {
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
	case types.STAT_HP:
		tempValue += p.GetMaxHP()
	case types.STAT_AD:
		tempValue += p.GetATK()
	case types.STAT_SPD:
		tempValue += p.GetSPD()
	case types.STAT_AGL:
		tempValue += p.GetAGL()
	case types.STAT_AP:
		tempValue += p.GetAP()
	case types.STAT_DEF:
		tempValue += p.GetDEF()
	case types.STAT_MR:
		tempValue += p.GetMR()
	case types.STAT_MANA:
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

func (p *Player) GetLvlSkill(lvl int) types.PlayerSkill {

	skill, skillExists := p.Inventory.LevelSkills[lvl]

	if !skillExists {
		return nil
	}

	return skill
}

func (p *Player) GetUID() string {
	return p.Meta.UserID
}

func (p *Player) GetLvlCD(lvl int) int {
	return p.Inventory.LevelSkillsCDS[lvl]
}

func (p *Player) SetLvlCD(lvl int, value int) {
	if value == 0 {
		delete(p.Inventory.LevelSkillsCDS, lvl)
	} else {
		p.Inventory.LevelSkillsCDS[lvl] = value

	}
}

func (p *Player) GetSkillsCD() map[any]int {
	mapTemp := make(map[any]int, 0)

	for skill, cd := range p.Inventory.LevelSkillsCDS {
		mapTemp[skill] = cd
	}

	return mapTemp
}

func (p *Player) GetParty() *uuid.UUID {
	return p.Meta.Party
}

func (p *Player) AppendDerivedStat(stat types.DerivedStat) {
	p.DynamicStats = append(p.DynamicStats, stat)
}

func NewPlayer(name string, uid string) Player {
	return Player{
		name,
		PlayerXP{Level: 1, Exp: 0},
		PlayerStats{100, make(mobs.EffectList, 0), false, 10},
		PlayerMeta{types.DefaultPlayerLocation(), uuid.New(), uid, nil, nil, nil, nil},
		inventory.GetDefaultInventory(),
		make([]types.DerivedStat, 0),
	}
}
