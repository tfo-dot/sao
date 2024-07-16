package inventory

import (
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
		4: []PlayerSkillLevel{DMG_LVL_4{}},
		5: []PlayerSkillLevel{DMG_LVL_5{}},
	},
	types.PathEndurance: {
		1: []PlayerSkillLevel{END_LVL_1{}},
		2: []PlayerSkillLevel{END_LVL_2{}},
		3: []PlayerSkillLevel{END_LVL_3{}},
		4: []PlayerSkillLevel{END_LVL_4{}},
		5: []PlayerSkillLevel{END_LVL_5{}},
	},
	types.PathControl: {
		1: []PlayerSkillLevel{CON_LVL_1{}},
		2: []PlayerSkillLevel{CON_LVL_2{}},
		3: []PlayerSkillLevel{CON_LVL_3{}},
		4: []PlayerSkillLevel{CON_LVL_4{}},
		5: []PlayerSkillLevel{CON_LVL_5{}},
	},
	types.PathSpecial: {
		1: []PlayerSkillLevel{SPC_LVL_1{}},
		2: []PlayerSkillLevel{SPC_LVL_2{}},
		3: []PlayerSkillLevel{SPC_LVL_3{}},
		4: []PlayerSkillLevel{SPC_LVL_4{}},
		5: []PlayerSkillLevel{SPC_LVL_5{}},
	},
}

// LVL to CD
var BaseCooldowns = map[int]int{
	1: 3,
	3: 4,
	4: 4,
	5: 4,
}

func HasUpgrade(upgrades, check int) bool {
	return upgrades&(1<<(check-1)) != 0
}

type NoEvents struct{}

func (n NoEvents) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return nil
}

type NoStats struct{}

func (n NoStats) GetStats(upgrades int) map[types.Stat]int {
	return map[types.Stat]int{}
}

type DefaultCost struct{}

func (d DefaultCost) GetCost() int {
	return 1
}

func (d DefaultCost) GetUpgradableCost(upgrades int) int {
	return 1
}

type DefaultActiveTrigger struct{}

func (d DefaultActiveTrigger) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE}
}

type NoTrigger struct{}

func (n NoTrigger) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_TYPE_NONE}
}

type NoExecute struct{}

func (n NoExecute) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	return nil
}

func (n NoExecute) GetCD() int {
	return 0
}

func (n NoExecute) GetCooldown(upgrades int) int {
	return 0
}

func (n NoExecute) GetCost() int {
	return 0
}

func (n NoExecute) GetUpgradableCost(upgrades int) int {
	return 0
}

type NoCost struct{}

func (n NoCost) GetCost() int {
	return 0
}

func (n NoCost) GetUpgradableCost(upgrades int) int {
	return 0
}

type NoLevel struct{}

func (n NoLevel) IsLevelSkill() bool {
	return false
}

func (n NoLevel) GetUUID() uuid.UUID {
	return uuid.New()
}

type NoCooldown struct{}

func (n NoCooldown) GetCD() int {
	return 0
}

func (n NoCooldown) GetCooldown(upgrades int) int {
	return 0
}
