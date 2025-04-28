package player

import (
	"errors"
	"fmt"
	"sao/base"
	"sao/data"
	"sao/player/inventory"
	"sao/types"
	"sao/utils"
	"sao/world/party"
	"strconv"

	"github.com/google/uuid"
)

type PlayerStats struct {
	HP          int
	Effects     []types.ActionEffect
	Defending   bool
	CurrentMana int
}

type PlayerXP struct {
	Level int
	Exp   int
}

type PartialParty struct {
	Role         party.PartyRole
	UUID         uuid.UUID
	MembersCount int
}

type PlayerMeta struct {
	OwnUUID       uuid.UUID
	UserID        string
	FightInstance *uuid.UUID
	Party         *PartialParty
	Transaction   *uuid.UUID
	WaitToHeal    bool
}

func (pM *PlayerMeta) Serialize() map[string]any {
	party := ""

	if pM.Party != nil {
		party = pM.Party.UUID.String()
	}

	return map[string]any{
		"uuid":  pM.OwnUUID.String(),
		"uid":   pM.UserID,
		"party": party,
	}
}

type Player struct {
	Name         string
	XP           PlayerXP
	Stats        PlayerStats
	Meta         PlayerMeta
	Inventory    inventory.PlayerInventory
	LevelStats   map[types.Stat]int
	DefaultStats map[types.Stat]int
}

func (p *Player) Serialize() map[string]any {
	return map[string]any{
		"name":          p.Name,
		"xp":            []int{p.XP.Level, p.XP.Exp},
		"stats":         map[string]any{"hp": p.Stats.HP, "current_mana": p.Stats.CurrentMana, "effects": p.Stats.Effects},
		"level_stats":   p.LevelStats,
		"default_stats": p.DefaultStats,
		"meta":          p.Meta.Serialize(),
		"inventory":     p.Inventory.Serialize(),
	}
}

func DeserializeLevelStats(data map[string]any) map[types.Stat]int {
	tempMap := make(map[types.Stat]int, 0)

	for key, value := range data {
		val, _ := strconv.Atoi(key)

		tempMap[types.Stat(val)] = int(value.(float64))
	}

	return tempMap
}

func DeserializeDefaultStats(data map[string]any) map[types.Stat]int {
	tempMap := make(map[types.Stat]int, 0)

	for key, value := range data {
		val, _ := strconv.Atoi(key)

		tempMap[types.Stat(val)] = int(value.(float64))
	}

	return tempMap
}

func Deserialize(data map[string]any) *Player {
	return &Player{
		data["name"].(string),
		PlayerXP{
			Level: int(data["xp"].([]any)[0].(float64)),
			Exp:   int(data["xp"].([]any)[1].(float64)),
		},
		PlayerStats{
			HP:          int(data["stats"].(map[string]any)["hp"].(float64)),
			Effects:     DeserializeEffects(data["stats"].(map[string]any)["effects"].([]any)),
			CurrentMana: int(data["stats"].(map[string]any)["current_mana"].(float64)),
		},
		*DeserializeMeta(data["meta"].(map[string]any)),
		inventory.DeserializeInventory(data["inventory"].(map[string]any)),
		DeserializeLevelStats(data["level_stats"].(map[string]any)),
		DeserializeDefaultStats(data["default_stats"].(map[string]any)),
	}
}

func DeserializeEffects(data []any) []types.ActionEffect {
	temp := make([]types.ActionEffect, len(data))

	if len(data) == 0 {
		return temp
	}

	for idx, effect := range data {
		effect := effect.(map[string]any)

		temp[idx] = types.ActionEffect{
			Effect: types.Effect(effect["effect"].(float64)),
			Value:  int(effect["value"].(float64)),
			Meta:   effect["meta"],
		}
	}

	return temp
}

func DeserializeMeta(data map[string]any) *PlayerMeta {
	var partyTemp *PartialParty = nil

	if data["party"] != "" {
		//TODO fetch party? XDDD
		partyTemp = &PartialParty{
			Role: party.None, UUID: uuid.MustParse(data["party"].(string)), MembersCount: 0,
		}
	}

	return &PlayerMeta{
		OwnUUID: uuid.MustParse(data["uuid"].(string)),
		UserID:  data["uid"].(string),
		Party:   partyTemp,
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

func (p *Player) Action(f types.FightInstance) []types.Action { return []types.Action{} }

func (p *Player) TakeDMG(dmgList types.ActionDamage) map[types.DamageType]int {
	return base.TakeDMG(dmgList, p)
}

func (p *Player) TakeDMGOrDodge(dmg types.ActionDamage) (map[types.DamageType]int, bool) {
	return base.TakeDMGOrDodge(dmg, p)
}

func (p *Player) DamageShields(dmg int) int {
	keepEffects, leftover := base.DamageShields(dmg, p)

	p.Stats.Effects = keepEffects

	return leftover
}

func (p *Player) AddEXP(value int) {
	p.XP.Exp += value

	maxLevel := (data.FloorMap.GetUnlockedFloorCount() * 5) - 1

	if p.XP.Level >= maxLevel {
		p.XP.Level = maxLevel
		p.XP.Exp = 0
		return
	}

	for p.XP.Exp >= ((p.XP.Level * 100) + 100) {
		if p.XP.Level >= maxLevel {
			p.XP.Level = maxLevel
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

func (p *Player) GetLoot() []types.Loot {
	return nil
}

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) CanDodge() bool {
	return !p.Stats.Defending
}

func (p *Player) ApplyEffect(e types.ActionEffect) {
	p.Stats.Effects = append(p.Stats.Effects, e)
}

func (p *Player) GetEffectByType(effect types.Effect) *types.ActionEffect {
	return base.GetEffectByType(effect, p)
}

func (p *Player) ChangeHP(value int) {
	p.Stats.HP += value
}

func (p *Player) GetEffectByUUID(uuid uuid.UUID) *types.ActionEffect {
	return base.GetEffectByUUID(uuid, p)
}

func (p *Player) TriggerAllEffects() {
	p.Stats.Effects = base.TriggerAllEffects(p)
}

func (p *Player) RemoveEffect(uuid uuid.UUID) {
	p.Stats.Effects = base.RemoveEffect(uuid, p)
}

func (p *Player) GetAllEffects() []types.ActionEffect {
	temporaryEffects := make([]types.ActionEffect, 0)

	if p.GetDefendingState() {
		temporaryEffects = append(temporaryEffects, types.ActionEffect{
			Effect:   types.EFFECT_STAT_INC,
			Value:    20,
			Duration: -1,
			Meta:     types.ActionEffectStat{Stat: types.STAT_DEF, IsPercent: true},
		}, types.ActionEffect{
			Effect:   types.EFFECT_STAT_INC,
			Duration: -1,
			Value:    20,
			Meta:     types.ActionEffectStat{Stat: types.STAT_MR, IsPercent: true},
		})
	}

	if p.Meta.Party != nil {
		switch p.Meta.Party.Role {
		case party.DPS:
			increase := 10 + (p.Meta.Party.MembersCount-1)*5

			temporaryEffects = append(temporaryEffects, types.ActionEffect{
				Effect:   types.EFFECT_STAT_INC,
				Duration: -1,
				Value:    increase,
				Meta:     types.ActionEffectStat{Stat: types.STAT_AD, IsPercent: true},
			}, types.ActionEffect{
				Effect:   types.EFFECT_STAT_INC,
				Duration: -1,
				Value:    increase,
				Meta:     types.ActionEffectStat{Stat: types.STAT_AP, IsPercent: true},
			})
		case party.Tank:
			temporaryEffects = append(temporaryEffects, types.ActionEffect{
				Effect:   types.EFFECT_STAT_INC,
				Duration: -1,
				Value:    25,
				Meta:     types.ActionEffectStat{Stat: types.STAT_DEF},
			}, types.ActionEffect{
				Effect:   types.EFFECT_STAT_INC,
				Duration: -1,
				Value:    25,
				Meta:     types.ActionEffectStat{Stat: types.STAT_MR},
			}, types.ActionEffect{
				Effect:   types.EFFECT_STAT_INC,
				Duration: -1,
				Value:    (p.Meta.Party.MembersCount - 1) * 5,
				Meta:     types.ActionEffectStat{Stat: types.STAT_HP, IsPercent: true},
			}, types.ActionEffect{
				Effect:   types.EFFECT_STAT_INC,
				Duration: -1,
				Value:    (p.Meta.Party.MembersCount - 1) * 5,
				Meta:     types.ActionEffectStat{Stat: types.STAT_DEF, IsPercent: true},
			}, types.ActionEffect{
				Effect:   types.EFFECT_STAT_INC,
				Duration: -1,
				Value:    (p.Meta.Party.MembersCount - 1) * 5,
				Meta:     types.ActionEffectStat{Stat: types.STAT_MR, IsPercent: true},
			})

			if (p.Meta.Party.MembersCount - 1) > 2 {
				temporaryEffects = append(temporaryEffects, types.ActionEffect{Effect: types.EFFECT_TAUNT, Duration: -1})
			}
		case party.Support:
			temporaryEffects = append(temporaryEffects, types.ActionEffect{
				Effect:   types.EFFECT_STAT_INC,
				Duration: -1,
				Value:    15 + (p.Meta.Party.MembersCount-1)*5,
				Meta:     types.ActionEffectStat{Stat: types.STAT_HEAL_POWER, IsPercent: true},
			})
		}
	}

	return append(temporaryEffects, p.Stats.Effects...)
}

func (p *Player) CanAttack() bool {
	return p.GetEffectByType(types.EFFECT_STUN) == nil
}

func (p *Player) CanDefend() bool {
	return p.GetEffectByType(types.EFFECT_STUN) == nil
}

func (p *Player) CanUseSkill(skill types.PlayerSkill) bool {
	if skill.IsLevelSkill() {
		lvl := skill.(types.PlayerSkillUpgradable).GetLevel()

		skillTrigger := skill.(types.PlayerSkillUpgradable).GetUpgradableTrigger(p.Inventory.LevelSkills[lvl].Upgrades)

		if skillTrigger.Type == types.TRIGGER_PASSIVE || skillTrigger.Type == types.TRIGGER_TYPE_NONE {
			return false
		}

		if p.GetEffectByType(types.EFFECT_STUN) != nil {
			return skillTrigger.Flags&types.FLAG_IGNORE_CC != 0
		}

		if info := p.Inventory.LevelSkills[lvl]; info.CD != 0 {
			return false
		}
	} else {
		skillTrigger := skill.GetTrigger()

		if skillTrigger.Type == types.TRIGGER_PASSIVE || skillTrigger.Type == types.TRIGGER_TYPE_NONE {
			return false
		}

		if p.GetEffectByType(types.EFFECT_STUN) != nil {
			return skillTrigger.Flags&types.FLAG_IGNORE_CC != 0
		}

		if currentCD, onCooldown := p.Inventory.ItemCD[skill.GetUUID()]; onCooldown && currentCD != 0 {
			return false
		}
	}

	if p.GetCurrentMana() < skill.GetCost() {
		return false
	}

	return true
}

func (p *Player) GetAllItems() []*types.PlayerItem {
	return p.Inventory.Items
}

func (p *Player) GetSkills() []types.PlayerSkill {
	arr := make([]types.PlayerSkill, 0)

	for _, skill := range p.Inventory.LevelSkills {
		arr = append(arr, skill.Skill)
	}

	return arr
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

func (p *Player) GetDerivedStats() []types.DerivedStat {
	return p.Inventory.GetDerivedStats()
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
		if effect.Effect == types.EFFECT_STAT_INC {

			if value, ok := effect.Meta.(types.ActionEffectStat); ok {
				if value.Stat != stat {
					continue
				}

				if value.IsPercent {
					percentValue += effect.Value
				} else {
					statValue += effect.Value
				}
			}
		}

		if effect.Effect == types.EFFECT_STAT_DEC {

			if value, ok := effect.Meta.(types.ActionEffectStat); ok {
				if value.Stat != stat {
					continue
				}

				if value.IsPercent {
					percentValue -= effect.Value
				} else {
					statValue -= effect.Value
				}
			}
		}
	}

	if value, ok := p.LevelStats[stat]; ok {
		statValue += ((p.XP.Level - 1) * value)
	}

	statValue += p.Inventory.GetStat(stat)

	for _, effect := range p.GetDerivedStats() {
		if effect.Derived == stat {
			if effect.Base == effect.Derived {
				statValue += utils.PercentOf(statValue, effect.Percent)
			} else {
				statValue += utils.PercentOf(p.GetStat(effect.Base), effect.Percent)
			}
		}
	}

	return statValue + (statValue * percentValue / 100)
}

func (p *Player) Cleanse() {
	p.Stats.Effects = base.Cleanse(p)
}

func (p *Player) GetUpgrades(lvl int) int {
	return p.Inventory.LevelSkills[lvl].Upgrades
}

func (p *Player) GetLvlSkill(lvl int) types.PlayerSkill {
	skill, skillExists := p.Inventory.LevelSkills[lvl]

	if !skillExists {
		return nil
	}

	return skill.Skill
}

func (p *Player) GetUID() string {
	return p.Meta.UserID
}

func (p *Player) GetLvlCD(lvl int) int {
	return p.Inventory.LevelSkills[lvl].CD
}

func (p *Player) SetLvlCD(lvl int, value int) {
	p.Inventory.LevelSkills[lvl].CD = value
}

func (p *Player) GetLevelSkillsCD() map[int]int {
	mapTemp := make(map[int]int, 0)

	for skill, cd := range p.Inventory.LevelSkills {
		if cd.CD == 0 {
			continue
		}

		mapTemp[skill] = cd.CD
	}

	return mapTemp
}

func (p *Player) SetLevelStat(stat types.Stat, value int) {
	p.LevelStats[stat] = value
}

func (p *Player) GetLevelStat(stat types.Stat) int {
	if value, ok := p.LevelStats[stat]; ok {
		return value
	}

	return 0
}

func (p *Player) ReduceCooldowns(event types.SkillTrigger) {
	for skill := range p.Inventory.ItemCD {
		rawSUID, _ := skill.MarshalBinary()

		copy(rawSUID[6:8], []byte{0, 0})

		itemUuid, _ := uuid.FromBytes(rawSUID)

		for _, effect := range data.Items[itemUuid].Effects {
			cdMeta := effect.GetTrigger().Cooldown

			if cdMeta == nil && event == types.TRIGGER_TURN {
				p.Inventory.ItemCD[skill]--
			}

			if cdMeta != nil && event != cdMeta.PassEvent {
				continue
			}
		}
	}

	for skillLevel := range p.Inventory.LevelSkills {
		skillData := p.Inventory.LevelSkills[skillLevel]
		cdMeta := skillData.Skill.GetUpgradableTrigger(skillData.Upgrades).Cooldown

		if cdMeta == nil && event == types.TRIGGER_TURN {
			p.Inventory.LevelSkills[skillLevel].CD--
		}

		if cdMeta != nil && event != cdMeta.PassEvent {
			continue
		}
	}
}

func (p *Player) GetLvl() int {
	return p.XP.Level
}

func (p *Player) TriggerEvent(event types.SkillTrigger, data types.EventData, meta any) []any {
	returnMeta := make([]any, 0)

	for _, item := range p.Inventory.Items {
		for _, effect := range item.Effects {
			if cd := effect.GetCD(); cd != 0 {
				currentCD, onCooldown := p.Inventory.ItemCD[effect.GetUUID()]

				if onCooldown {
					if currentCD == 0 {
						p.Inventory.ItemCD[effect.GetUUID()] = cd
					} else {
						continue
					}
				}
			}

			trigger := effect.GetTrigger()

			if trigger.Type == types.TRIGGER_PASSIVE && trigger.Event == event {
				if cost := effect.GetCost(); cost != 0 {
					if cost > p.GetCurrentMana() {
						continue
					} else {
						p.Stats.CurrentMana -= cost
					}
				}

				temp := effect.Execute(p, data.Target, data.Fight, meta)

				if temp != nil {
					returnMeta = append(returnMeta, temp)
				}
			}
		}
	}

	for skillLevel, skillStruct := range p.Inventory.LevelSkills {
		trigger := skillStruct.Skill.GetUpgradableTrigger(skillStruct.Upgrades)

		if cd := skillStruct.Skill.GetCooldown(skillStruct.Upgrades); cd != 0 {
			if currentCD := p.Inventory.LevelSkills[skillLevel].CD; currentCD == 0 {
				p.Inventory.LevelSkills[skillLevel].CD = cd
			} else {
				continue
			}
		}

		if trigger.Type == types.TRIGGER_PASSIVE && trigger.Event == event {
			if cost := skillStruct.Skill.GetCost(); cost != 0 {
				if cost > p.GetCurrentMana() {
					continue
				} else {
					p.Stats.CurrentMana -= cost
				}
			}

			if temp := skillStruct.Skill.Execute(p, data.Target, data.Fight, meta); temp != nil {
				returnMeta = append(returnMeta, temp)
			}
		}
	}

	for _, effect := range p.Inventory.TempSkills {
		trigger := effect.Value.GetTrigger()

		if trigger.Type == types.TRIGGER_PASSIVE && trigger.Event == event {
			if cost := effect.Value.GetCost(); cost != 0 {
				if cost > p.GetCurrentMana() {
					continue
				} else {
					p.Stats.CurrentMana -= cost
				}
			}

			temp := effect.Value.Execute(p, data.Target, data.Fight, meta)

			if temp != nil {
				returnMeta = append(returnMeta, temp)
			}
		}
	}

	return returnMeta
}

func (p *Player) UnlockSkill(path types.SkillPath, lvl, choice int) error {
	if lvl > p.GetLvl() {
		return errors.New("PLAYER_LVL_TOO_LOW")
	}

	if p.GetAvailableSkillActions() < 1 {
		return errors.New("NO_ACTIONS_AVAILABLE")
	}

	if err := p.Inventory.UnlockSkill(lvl, path, choice); err != nil {
		return err
	}

	skillEvents := p.Inventory.LevelSkills[lvl].Skill.GetEvents()

	if effect, effectExists := skillEvents[types.CUSTOM_TRIGGER_UNLOCK]; effectExists {
		effect(p)
	}

	return nil
}

func (p *Player) UpgradeSkill(lvl int, upgradeIdx int) error {
	if p.GetAvailableSkillActions() <= 0 {
		return errors.New("NO_ACTIONS_AVAILABLE")
	}

	skillUpgrade := p.Inventory.LevelSkills[lvl].Skill.GetUpgrades()[upgradeIdx]

	if err := p.Inventory.UpgradeSkill(lvl, upgradeIdx); err != nil {
		return err
	}

	if effect, effectExists := skillUpgrade.Events[types.CUSTOM_TRIGGER_UNLOCK]; effectExists {
		effect(p)
	}

	return nil
}

func (p *Player) TriggerTempSkills() {
	list := make([]*types.WithExpire[types.PlayerSkill], 0)

	for _, skill := range p.Inventory.TempSkills {
		if !skill.AfterUsage || skill.Either {
			skill.Expire--

			if skill.Expire > 0 {
				list = append(list, skill)
			} else {
				continue
			}
		} else {
			list = append(list, skill)
		}
	}

	p.Inventory.TempSkills = list
}

func (p *Player) ClearFight() {
	p.Meta.FightInstance = nil
}

func (p *Player) GetAvailableSkillActions() int {
	overall := 0

	if p.XP.Level <= 6 {
		overall = p.XP.Level
	}

	{
		temp := p.XP.Level - 6

		for temp >= 2 {
			temp -= 2
			overall++
		}
	}

	used := len(p.Inventory.LevelSkills)

	for _, skill := range p.Inventory.LevelSkills {
		for i := range skill.Skill.GetUpgrades() {
			if inventory.HasUpgrade(skill.Upgrades, i) {
				used++
			}
		}
	}

	return overall - used
}

func (p *Player) SetLevelSkillMeta(lvl int, meta any) {
	p.Inventory.LevelSkills[lvl].Meta = meta
}

func (p *Player) GetLevelSkillMeta(lvl int) any {
	return p.Inventory.LevelSkills[lvl].Meta
}

func (p *Player) UseItem(item uuid.UUID, target types.Entity, fight types.FightInstance) {
	p.Inventory.UseItem(item, p, target, fight)
}

func (p *Player) GetSkillPath() types.SkillPath {
	counts := make(map[types.SkillPath]int)

	for lvl, skill := range p.Inventory.LevelSkills {
		if lvl%10 == 0 {
			continue
		}

		counts[skill.Skill.GetPath()]++

		if skill.Upgrades != 0 {
			for i := range skill.Skill.GetUpgrades() {
				if inventory.HasUpgrade(skill.Upgrades, i) {
					counts[skill.Skill.GetPath()]++
				}
			}
		}
	}

	max := types.PathControl
	maxCount := 0

	for path, count := range counts {
		fmt.Printf("%d %d", path, count)
		if count > maxCount {
			max = path
			maxCount = count
		}
	}

	return max
}

func NewPlayer(name string, uid string) Player {
	return Player{
		name,
		PlayerXP{Level: 1},
		PlayerStats{
			HP:          data.PlayerDefaults.Stats[types.STAT_HP],
			Effects:     make([]types.ActionEffect, 0),
			CurrentMana: data.PlayerDefaults.Stats[types.STAT_MANA],
		},
		PlayerMeta{
			OwnUUID: uuid.New(),
			UserID:  uid,
		},
		inventory.GetDefaultInventory(),
		data.PlayerDefaults.Level,
		data.PlayerDefaults.Stats,
	}
}
