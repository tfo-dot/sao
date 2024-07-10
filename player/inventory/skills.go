package inventory

import (
	"fmt"
	"sao/battle"
	"sao/types"

	"github.com/google/uuid"
)

type PlayerSkillUpgradable interface {
	types.PlayerSkill

	GetUpgrades() []PlayerSkillUpgrade
	GetCooldown(upgrades int) int
	UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{}
	GetStats(upgrades int) map[types.Stat]int
	GetUpgradableCost(upgrades int) int
}

type PlayerSkillUpgrade struct {
	Name        string
	Description string
	Id          string
	Events      *map[types.CustomTrigger]func(owner battle.PlayerEntity)
}

type PlayerSkillLevel interface {
	PlayerSkillUpgradable

	GetLevel() int
	GetPath() types.SkillPath
}

var AVAILABLE_SKILLS = map[types.SkillPath]map[int][]PlayerSkillLevel{
	types.PathDamage: {
		1: []PlayerSkillLevel{DMG_LVL_1{}},
		2: []PlayerSkillLevel{DMG_LVL_2{}},
		3: []PlayerSkillLevel{DMG_LVL_3{}},
	},
	types.PathEndurance: {
		1: []PlayerSkillLevel{END_LVL_1{}},
		2: []PlayerSkillLevel{END_LVL_2{}},
		3: []PlayerSkillLevel{END_LVL_3{}},
	},
	types.PathControl: {
		1: []PlayerSkillLevel{CON_LVL_1{}},
		2: []PlayerSkillLevel{CON_LVL_2{}},
		3: []PlayerSkillLevel{CON_LVL_3{}},
		4: []PlayerSkillLevel{CON_LVL_4{}},
	},
	types.PathSpecial: {
		1: []PlayerSkillLevel{SPC_LVL_1{}},
		2: []PlayerSkillLevel{SPC_LVL_2{}},
		3: []PlayerSkillLevel{SPC_LVL_3{}},
		4: []PlayerSkillLevel{SPC_LVL_4{}},
	},
}

// LVL to CD
var BaseCooldowns = map[int]int{
	1: 3,
	3: 4,
	4: 4,
}

type BaseLvlSkill struct{}

func (skill BaseLvlSkill) IsLevelSkill() bool {
	return true
}

func (skill BaseLvlSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type:  types.TRIGGER_ACTIVE,
		Event: nil,
	}
}

func (skill BaseLvlSkill) GetUUID() uuid.UUID {
	return uuid.Nil
}

func (skill BaseLvlSkill) GetLevel() int {
	return 1
}

func (skill BaseLvlSkill) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return nil
}

func (skill BaseLvlSkill) GetName() string {
	return fmt.Sprintf("Poziom %v", skill.GetLevel())
}

func (skill BaseLvlSkill) GetCost() int {
	return 1
}

func (skill BaseLvlSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	return nil
}

func (skill BaseLvlSkill) GetStats(upgrades int) map[types.Stat]int {
	return map[types.Stat]int{}
}

func (skill BaseLvlSkill) GetUpgradableCost(upgrades int) int {
	return 1
}

func (skill BaseLvlSkill) GetCooldown(upgrades int) int {
	return 0
}

func (skill BaseLvlSkill) GetCD() int {
	if cd, ok := BaseCooldowns[skill.GetLevel()]; ok {
		return cd
	}

	return 0
}
