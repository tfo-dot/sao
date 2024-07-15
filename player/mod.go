package player

import (
	"sao/battle"
	"sao/battle/mobs"
	"sao/player/inventory"
	"sao/types"
	"sao/utils"
	"sao/world/fury"
	"strconv"

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
	Fury          *fury.Fury
}

func (pM *PlayerMeta) SerializeFuries() map[string]interface{} {
	if pM.Fury == nil {
		return nil
	}

	return pM.Fury.Serialize()
}

func (pM *PlayerMeta) Serialize() map[string]interface{} {
	party := ""
	if pM.Party != nil {
		party = pM.Party.String()
	}

	return map[string]interface{}{
		"location": []string{pM.Location.FloorName, pM.Location.LocationName},
		"uuid":     pM.OwnUUID.String(),
		"uid":      pM.UserID,
		"fury":     pM.SerializeFuries(),
		"party":    party,
	}
}

type Player struct {
	Name         string
	XP           PlayerXP
	Stats        PlayerStats
	Meta         PlayerMeta
	Inventory    inventory.PlayerInventory
	DynamicStats []types.DerivedStat
	LevelStats   map[types.Stat]int
	DefaultStats map[types.Stat]int
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
		"dynamic_stats": p.DynamicStats,
		"level_stats":   p.LevelStats,
		"default_stats": p.DefaultStats,
		"meta":          p.Meta.Serialize(),
		"inventory":     p.Inventory.Serialize(),
	}
}

func DeserializeDerivedStats(data []interface{}) []types.DerivedStat {
	tempArray := make([]types.DerivedStat, 0)

	for _, stat := range data {
		stat := stat.(map[string]interface{})

		temp := types.DerivedStat{
			Derived: types.Stat(stat["derived"].(float64)),
			Base:    types.Stat(stat["base"].(float64)),
			Percent: int(stat["percent"].(float64)),
			Source:  uuid.MustParse(stat["source"].(string)),
		}

		tempArray = append(tempArray, temp)
	}

	return tempArray
}

func DeserializeLevelStats(data map[string]interface{}) map[types.Stat]int {
	tempMap := make(map[types.Stat]int, 0)

	for key, value := range data {
		val, _ := strconv.Atoi(key)

		tempMap[types.Stat(val)] = int(value.(float64))
	}

	return tempMap
}

func DeserializeDefaultStats(data map[string]interface{}) map[types.Stat]int {
	tempMap := make(map[types.Stat]int, 0)

	for key, value := range data {
		val, _ := strconv.Atoi(key)

		tempMap[types.Stat(val)] = int(value.(float64))
	}

	return tempMap
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
		DeserializeDerivedStats(data["dynamic_stats"].([]interface{})),
		DeserializeLevelStats(data["level_stats"].(map[string]interface{})),
		DeserializeDefaultStats(data["default_stats"].(map[string]interface{})),
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
	var furyData *fury.Fury
	if data["fury"] != nil {
		furyData = fury.Deserialize(data["fury"].(map[string]interface{}))
	}

	var party *uuid.UUID

	if data["party"] != "" {
		temp := uuid.MustParse(data["party"].(string))

		party = &temp
	}

	return &PlayerMeta{
		types.DefaultPlayerLocation(),
		uuid.MustParse(data["uuid"].(string)),
		data["uid"].(string),
		nil,
		party,
		nil,
		furyData,
	}
}

func (p *Player) AppendTempSkill(skill types.WithExpire[types.PlayerSkill]) {
	p.Inventory.AddTempSkill(skill)
}

func (p *Player) GetCurrentMana() int {
	return p.Stats.CurrentMana
}

func (p *Player) GetUUID() uuid.UUID {
	return p.Meta.OwnUUID
}

func (p *Player) Heal(val int) {
	p.Stats.HP += val

	if p.Stats.HP > p.GetStat(types.STAT_HP) {
		p.Stats.HP = p.GetStat(types.STAT_HP)
	}
}

func (p *Player) GetFlags() types.EntityFlag {
	return 0
}

func (p *Player) SetDefendingState(state bool) {
	p.Stats.Defending = state
}

func (p *Player) GetDefendingState() bool {
	return p.Stats.Defending
}

func (p *Player) Action(f *battle.Fight) []battle.Action { return []battle.Action{} }

func (p *Player) TakeDMG(dmgList battle.ActionDamage) []battle.Damage {
	dmgStats := []battle.Damage{
		{Value: 0, Type: battle.DMG_PHYSICAL},
		{Value: 0, Type: battle.DMG_MAGICAL},
		{Value: 0, Type: battle.DMG_TRUE},
	}

	for _, dmg := range dmgList.Damage {
		//Skip shield and such
		if dmg.Type == battle.DMG_TRUE {
			p.Stats.HP -= dmg.Value
			dmgStats[2].Value += dmg.Value
			continue
		}

		rawDmg := dmg.Value

		switch dmg.Type {
		case battle.DMG_PHYSICAL:
			rawDmg = utils.CalcReducedDamage(dmg.Value, p.GetStat(types.STAT_DEF))
		case battle.DMG_MAGICAL:
			rawDmg = utils.CalcReducedDamage(dmg.Value, p.GetStat(types.STAT_MR))
		}

		actualDmg := p.DamageShields(rawDmg)

		dmgStats[dmg.Type].Value += actualDmg

		p.Stats.HP -= actualDmg
	}

	return dmgStats
}

func (p *Player) TakeDMGOrDodge(dmg battle.ActionDamage) ([]battle.Damage, bool) {
	if utils.RandomNumber(0, 100) <= p.GetStat(types.STAT_AGL) && dmg.CanDodge {
		return []battle.Damage{
			{Value: 0, Type: battle.DMG_PHYSICAL},
			{Value: 0, Type: battle.DMG_MAGICAL},
			{Value: 0, Type: battle.DMG_TRUE},
		}, true
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
		p.XP.Level = maxFloor * 5
		p.XP.Exp = 0
		return
	}

	p.Meta.Fury.AddXP(utils.PercentOf(value, 20))

	for p.XP.Exp >= ((p.XP.Level * 100) + 100) {
		if p.XP.Level >= maxFloor*5 {
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

	if p.Stats.HP < p.GetStat(types.STAT_HP) {
		p.Stats.HP = utils.PercentOf(p.GetStat(types.STAT_HP), 20) + p.GetCurrentHP()
	}

	if p.Stats.HP > p.GetStat(types.STAT_HP) {
		p.Stats.HP = p.GetStat(types.STAT_HP)
	}
}

func (p *Player) GetCurrentHP() int {
	return p.Stats.HP
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

// TODO rewrite it so it knows what type of skill what used (uuid moment)
func (p *Player) GetAllSkills() []types.PlayerSkill {
	tempArr := make([]types.PlayerSkill, 0)

	for _, item := range p.Inventory.Items {
		tempArr = append(tempArr, item.Effects...)
	}

	for _, skill := range p.Inventory.LevelSkills {
		tempArr = append(tempArr, skill)
	}

	if p.Meta.Fury != nil {
		tempArr = append(tempArr, p.Meta.Fury.GetSkills()...)
	}

	return tempArr
}

func (p *Player) AddItem(item *types.PlayerItem) {

	for _, effect := range item.Effects {
		effectEvents := effect.GetEvents()

		effectEvents[types.CUSTOM_TRIGGER_UNLOCK](p)
	}

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

	if p.Stats.CurrentMana > p.GetStat(types.STAT_MANA) {
		p.Stats.CurrentMana = p.GetStat(types.STAT_MANA)
	}
}

func (p *Player) GetDefaultStat(stat types.Stat) int {
	if value, ok := p.DefaultStats[stat]; ok {
		return value
	}
	return 0
}

func (p *Player) GetStat(stat types.Stat) int {
	switch stat {
	case types.STAT_MANA_PLUS:
		return p.GetDefaultStat(types.STAT_MANA) - p.GetStat(types.STAT_MANA)
	case types.STAT_HP_PLUS:
		return p.GetDefaultStat(types.STAT_HP) + ((p.XP.Level - 1) * 10) - p.GetStat(types.STAT_HP)
	}

	statValue := p.GetDefaultStat(stat)
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

	if value, ok := p.LevelStats[stat]; ok {
		statValue += ((p.XP.Level - 1) * value)
	}

	statValue += p.Inventory.GetStat(stat)

	if p.Meta.Fury != nil {
		statValue += p.Meta.Fury.GetStat(stat)
	}

	for _, effect := range p.DynamicStats {
		if effect.Derived == stat {
			statValue += utils.PercentOf(p.GetStat(effect.Base), effect.Percent)
		}
	}

	if stat == types.STAT_AD || stat == types.STAT_AP {
		adaptive := p.GetStat(types.STAT_ADAPTIVE)

		if adaptive > 0 {
			adaptiveType := p.GetAdaptiveAttackType()

			if adaptiveType == types.ADAPTIVE_ATK && stat == types.STAT_AD {
				statValue += adaptive
			}

			if adaptiveType == types.ADAPTIVE_AP && stat == types.STAT_AP {
				statValue += adaptive
			}
		}
	}

	baseStat := statValue + (statValue * percentValue / 100)

	if stat == types.STAT_AD || stat == types.STAT_AP {
		adaptive := p.GetStat(types.STAT_ADAPTIVE_PERCENT)

		if adaptive > 0 {
			adaptiveType := p.GetAdaptiveAttackType()

			if adaptiveType == types.ADAPTIVE_ATK && stat == types.STAT_AD {
				baseStat += utils.PercentOf(baseStat, adaptive)
			}

			if adaptiveType == types.ADAPTIVE_AP && stat == types.STAT_AP {
				baseStat += utils.PercentOf(baseStat, adaptive)
			}
		}
	}
	return baseStat
}

func (p *Player) GetRawStat(stat types.Stat) int {
	switch stat {
	case types.STAT_MANA_PLUS:
		return p.GetDefaultStat(types.STAT_MANA) - p.GetStat(types.STAT_MANA)
	case types.STAT_HP_PLUS:
		return p.GetDefaultStat(types.STAT_HP) + ((p.XP.Level - 1) * 10) - p.GetStat(types.STAT_HP)
	}

	statValue := p.GetDefaultStat(stat)
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

	if value, ok := p.LevelStats[stat]; ok {
		statValue += ((p.XP.Level - 1) * value)
	}

	statValue += p.Inventory.GetStat(stat)

	if p.Meta.Fury != nil {
		statValue += p.Meta.Fury.GetStat(stat)
	}

	for _, effect := range p.DynamicStats {
		if effect.Derived == stat {
			statValue += utils.PercentOf(p.GetStat(effect.Base), effect.Percent)
		}
	}

	return statValue + (statValue * percentValue / 100)
}

func (p *Player) GetAdaptiveAttackType() types.AdaptiveAttackType {
	adStat := p.GetRawStat(types.STAT_AD)
	apStat := p.GetRawStat(types.STAT_AP)

	if apStat > adStat {
		return types.ADAPTIVE_AP
	}

	return types.ADAPTIVE_ATK
}

func (p *Player) Cleanse() {
	p.Stats.Effects.Cleanse()
}

func (p *Player) GetUpgrades(lvl int) int {
	return p.Inventory.LevelSkillsUpgrades[lvl]
}

func (p *Player) GetLvlSkill(lvl int) types.PlayerSkill {
	skill, skillExists := p.Inventory.LevelSkills[lvl]

	if !skillExists {
		return nil
	}

	return skill
}

func (p *Player) GetSkill(uuid uuid.UUID) types.PlayerSkill {
	for _, skill := range p.GetAllSkills() {
		if skill.GetUUID() == uuid {
			return skill
		}
	}

	return nil
}

func (p *Player) GetUID() string {
	return p.Meta.UserID
}

func (p *Player) GetCD(skill uuid.UUID) int {
	return p.Inventory.Cooldowns[skill]
}

func (p *Player) SetCD(skill uuid.UUID, value int) {
	if value == 0 {
		delete(p.Inventory.Cooldowns, skill)
	} else {
		p.Inventory.Cooldowns[skill] = value
	}
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

func (p *Player) GetLevelSkillsCD() map[int]int {
	mapTemp := make(map[int]int, 0)

	for skill, cd := range p.Inventory.LevelSkillsCDS {
		mapTemp[skill] = cd
	}

	return mapTemp
}

func (p *Player) GetSkillsCD() map[uuid.UUID]int {
	mapTemp := make(map[uuid.UUID]int, 0)

	for skill, cd := range p.Inventory.Cooldowns {
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

func (p *Player) SetLevelStat(stat types.Stat, value int) {
	p.LevelStats[stat] = value
}

func (p *Player) ReduceCooldowns(event types.SkillTrigger) {
	for skill, cd := range p.Inventory.Cooldowns {
		//TODO check if skill has different cd pass event trigger
		//TODO check if this is even used
		if cd > 0 {
			p.Inventory.Cooldowns[skill]--
		}
	}

	for skill, cd := range p.Inventory.LevelSkillsCDS {
		skillData := p.Inventory.LevelSkills[skill]
		cdMeta := skillData.GetTrigger().Cooldown

		if cdMeta == nil || cdMeta.PassEvent != event {
			continue
		}

		if cd > 0 {
			p.Inventory.LevelSkillsCDS[skill]--
		}
	}
}

// TODO Supply all data to the events
// TODO add cooldowns
// TODO returns meta effect (used for example in TRIGGER_ATTACK_ATTEMPT)
// TODO skill costs
func (p *Player) TriggerEvent(event types.SkillTrigger, meta interface{}) {
	for _, item := range p.Inventory.Items {
		for _, effect := range item.Effects {
			trigger := effect.GetTrigger()

			if trigger.Type == types.TRIGGER_PASSIVE && trigger.Event.TriggerType == event {
				effect.Execute(p, nil, nil, meta)
			}
		}
	}

	for _, skill := range p.Inventory.LevelSkills {
		trigger := skill.GetTrigger()

		if trigger.Type == types.TRIGGER_PASSIVE && trigger.Event.TriggerType == event {
			skill.Execute(p, nil, nil, meta)
		}
	}

	if p.Meta.Fury != nil {
		for _, skill := range p.Meta.Fury.GetSkills() {
			trigger := skill.GetTrigger()

			if trigger.Type == types.TRIGGER_PASSIVE && trigger.Event.TriggerType == event {
				skill.Execute(p, nil, nil, meta)
			}
		}
	}
}

func NewPlayer(name string, uid string) Player {
	return Player{
		name,
		PlayerXP{Level: 1, Exp: 0},
		PlayerStats{100, make(mobs.EffectList, 0), false, 10},
		PlayerMeta{types.DefaultPlayerLocation(), uuid.New(), uid, nil, nil, nil, nil},
		inventory.GetDefaultInventory(),
		make([]types.DerivedStat, 0),
		map[types.Stat]int{
			types.STAT_HP: 15,
			types.STAT_AD: 15,
		},
		map[types.Stat]int{
			types.STAT_HP:   100,
			types.STAT_AD:   40,
			types.STAT_SPD:  40,
			types.STAT_AGL:  50,
			types.STAT_MANA: 10,
		},
	}
}
