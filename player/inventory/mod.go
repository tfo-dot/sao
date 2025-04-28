package inventory

import (
	"errors"
	"sao/data"
	"sao/types"
	"strconv"

	"github.com/google/uuid"
	"slices"
)

type PlayerInventory struct {
	Gold        int
	TempSkills  []*types.WithExpire[types.PlayerSkill]
	Items       []*types.PlayerItem
	ItemCD      map[uuid.UUID]int
	LevelSkills map[int]*LevelSkillInfo
}

type LevelSkillInfo struct {
	CD       int
	Choice   int
	Upgrades int
	Skill    types.PlayerSkillUpgradable
	Meta     any
}

func (inv *PlayerInventory) AddTempSkill(skill types.WithExpire[types.PlayerSkill]) {
	inv.TempSkills = append(inv.TempSkills, &skill)
}

func (inv *PlayerInventory) Serialize() map[string]any {
	items := make([]map[string]any, 0)

	for _, item := range inv.Items {
		items = append(items, map[string]any{"uuid": item.UUID.String(), "count": item.Count})
	}

	lvlSkills := make(map[int]any)

	for key, info := range inv.LevelSkills {
		lvlSkills[key] = map[string]any{
			"path":     info.Skill.GetPath(),
			"choice":   info.Choice,
			"upgrades": info.Upgrades,
		}
	}

	return map[string]any{
		"gold":        inv.Gold,
		"items":       items,
		"levelSkills": lvlSkills,
	}
}

func DeserializeInventory(rawData map[string]any) PlayerInventory {
	inv := GetDefaultInventory()

	inv.Gold = int(rawData["gold"].(float64))

	if rawItemData, okay := rawData["items"].([]map[string]any); okay {
		for _, item := range rawItemData {
			uuid, _ := uuid.Parse(item["uuid"].(string))
			copy := data.Items[uuid]

			copy.Count = item["count"].(int)

			inv.Items = append(inv.Items, &copy)
		}
	}

	if _, okay := rawData["levelSkills"].(map[string]any); okay {
		for key, lvlData := range rawData["levelSkills"].(map[string]any) {
			parsed, _ := strconv.Atoi(key)
			data := lvlData.(map[string]any)

			inv.LevelSkills[parsed] = &LevelSkillInfo{
				CD:       0,
				Choice:   int(data["choice"].(float64)),
				Upgrades: int(data["upgrades"].(float64)),
				Skill:    AVAILABLE_SKILLS[types.SkillPath(data["path"].(float64))][parsed][int(data["choice"].(float64))],
			}
		}
	}

	return inv
}

func (inv *PlayerInventory) AddItem(item *types.PlayerItem) {
	for _, invItem := range inv.Items {
		if invItem.UUID == item.UUID && invItem.Stacks && invItem.Count < invItem.MaxCount {
			if invItem.Count+item.Count > invItem.MaxCount {
				item.Count -= invItem.MaxCount - invItem.Count
				invItem.Count = invItem.MaxCount
			} else {
				invItem.Count += item.Count
				return
			}
		}
	}

	//This will bite me in the ass later
	inv.Items = append(inv.Items, item)
}

func (inv PlayerInventory) GetStat(stat types.Stat) int {
	value := 0

	for _, item := range inv.Items {
		val, exists := item.Stats[stat]

		if exists {
			value += val
		}
	}

	for key, skillInfo := range inv.LevelSkills {
		skillStats := skillInfo.Skill.GetStats(inv.LevelSkills[key].Upgrades)

		val, exists := skillStats[stat]

		if exists {
			value += val
		}
	}

	return value
}

func (inv PlayerInventory) GetDerivedStats() []types.DerivedStat {
	statList := make([]types.DerivedStat, 0)

	for _, item := range inv.Items {
		statList = append(statList, item.DerivedStats...)
	}

	for key, skillInfo := range inv.LevelSkills {
		statList = append(statList, skillInfo.Skill.GetDerivedStats(inv.LevelSkills[key].Upgrades)...)
	}

	return statList
}

func (inv *PlayerInventory) UseItem(itemUuid uuid.UUID, owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance) {
	for i, item := range inv.Items {
		if item.UUID == itemUuid {
			if item.Consume && item.Count <= 0 {
				return
			}

			if item.Consume {
				item.Count--
			}

			inv.Items[i].UseItem(owner, target, fightInstance)

			if item.Count == 0 && item.Consume {
				inv.Items = slices.Delete(inv.Items, i, i+1)
			}
		}
	}
}

func (inv *PlayerInventory) UpgradeSkill(lvl int, upgradeIdx int) error {
	skillInfo, exists := inv.LevelSkills[lvl]

	if !exists {
		return errors.New("SKILL_NOT_UNLOCKED")
	}

	if upgradeIdx == -1 {
		return errors.New("UPGRADE_NOT_FOUND")
	}

	if HasUpgrade(skillInfo.Upgrades, upgradeIdx) {
		return errors.New("UPGRADE_ALREADY_UNLOCKED")
	}

	inv.LevelSkills[lvl].Upgrades = inv.LevelSkills[lvl].Upgrades & (1 << upgradeIdx)

	return nil
}

func (inv *PlayerInventory) UnlockSkill(lvl int, path types.SkillPath, choice int) error {
	if _, exists := inv.LevelSkills[lvl]; exists {
		return errors.New("SKILL_ALREADY_UNLOCKED")
	}

	skillChoices := AVAILABLE_SKILLS[path][lvl]

	if choice > len(skillChoices) {
		return errors.New("INVALID_CHOICE")
	}

	inv.LevelSkills[lvl] = &LevelSkillInfo{
		CD:       0,
		Choice:   choice,
		Upgrades: 0,
		Skill:    skillChoices[choice],
		Meta:     nil,
	}

	return nil
}

func GetDefaultInventory() PlayerInventory {
	return PlayerInventory{
		TempSkills:  make([]*types.WithExpire[types.PlayerSkill], 0),
		Items:       make([]*types.PlayerItem, 0),
		LevelSkills: make(map[int]*LevelSkillInfo),
	}
}
