package inventory

import (
	"sao/types"

	"github.com/google/uuid"
)

var AVAILABLE_SKILLS = map[types.SkillPath]map[int][]types.PlayerSkillUpgradable{
	types.PathDamage: {
		1: []types.PlayerSkillUpgradable{DMG_LVL_1{}},
		2: []types.PlayerSkillUpgradable{DMG_LVL_2{}},
		3: []types.PlayerSkillUpgradable{DMG_LVL_3{}},
		4: []types.PlayerSkillUpgradable{DMG_LVL_4{}},
		5: []types.PlayerSkillUpgradable{DMG_LVL_5{}},
		6: []types.PlayerSkillUpgradable{DMG_LVL_6{}},
		10: []types.PlayerSkillUpgradable{
			DMG_ULT_1{},
			DMG_LVL_2{},
		},
	},
	types.PathEndurance: {
		1: []types.PlayerSkillUpgradable{END_LVL_1{}},
		2: []types.PlayerSkillUpgradable{END_LVL_2{}},
		3: []types.PlayerSkillUpgradable{END_LVL_3{}},
		4: []types.PlayerSkillUpgradable{END_LVL_4{}},
		5: []types.PlayerSkillUpgradable{END_LVL_5{}},
		6: []types.PlayerSkillUpgradable{END_LVL_6{}},
		10: []types.PlayerSkillUpgradable{
			END_ULT_1{},
			END_ULT_2{},
		},
	},
	types.PathControl: {
		1: []types.PlayerSkillUpgradable{CON_LVL_1{}},
		2: []types.PlayerSkillUpgradable{CON_LVL_2{}},
		3: []types.PlayerSkillUpgradable{CON_LVL_3{}},
		4: []types.PlayerSkillUpgradable{CON_LVL_4{}},
		5: []types.PlayerSkillUpgradable{CON_LVL_5{}},
		6: []types.PlayerSkillUpgradable{CON_LVL_6{}},
		10: []types.PlayerSkillUpgradable{
			CON_ULT_1{},
			CON_ULT_2{},
		},
	},
	types.PathSpecial: {
		1: []types.PlayerSkillUpgradable{SPC_LVL_1{}},
		2: []types.PlayerSkillUpgradable{SPC_LVL_2{}},
		3: []types.PlayerSkillUpgradable{SPC_LVL_3{}},
		4: []types.PlayerSkillUpgradable{SPC_LVL_4{}},
		5: []types.PlayerSkillUpgradable{SPC_LVL_5{}},
		6: []types.PlayerSkillUpgradable{SPC_LVL_6{}},
		10: []types.PlayerSkillUpgradable{
			SPC_ULT_1{},
			SPC_ULT_2{},
		},
	},
}

// LVL to CD
var BaseCooldowns = map[int]int{
	1: 3,
	3: 4,
	4: 4,
	5: 4,
	6: 4,
}

func HasUpgrade(upgrades, check int) bool {
	return upgrades&(1<<(check-1)) != 0
}

type NoEvents struct{}

func (n NoEvents) GetEvents() map[types.CustomTrigger]func(owner types.PlayerEntity) {
	return nil
}

type NoStats struct{}

func (n NoStats) GetStats(upgrades int) map[types.Stat]int {
	return map[types.Stat]int{}
}

func (n NoStats) GetDerivedStats(upgrades int) []types.DerivedStat {
	return []types.DerivedStat{}
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

func (d DefaultActiveTrigger) GetUpgradableTrigger(upgrades int) types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE}
}

type NoTrigger struct{}

func (n NoTrigger) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_TYPE_NONE}
}

func (n NoTrigger) GetUpgradableTrigger(upgrades int) types.Trigger {
	return types.Trigger{Type: types.TRIGGER_TYPE_NONE}
}

type NoExecute struct{}

func (n NoExecute) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
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

type Counter struct {
	NoExecute
	NoEvents
	NoStats
	NoTrigger
	NoLevel
}

func (c Counter) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	return nil
}

func (c Counter) GetDescription() string {
	return ""
}

func (c Counter) GetName() string {
	return ""
}

func (c Counter) GetUUID() uuid.UUID {
	return uuid.New()
}