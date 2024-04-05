package inventory

import (
	"sao/battle"
	"sao/types"
)

var AVAILABLE_SKILLS = map[types.SkillPath]map[int]PlayerSkillLevel{}

type PlayerSkillUpgradable interface {
	types.PlayerSkill

	GetUpgrades() []PlayerSkillUpgrade
	HasUpgrade(upgradeId string) bool

	UnlockUpgrade(upgradeId string)
}

type PlayerSkillUpgrade interface {
	GetName() string
	GetDescription() string
	GetId() string

	Execute(owner battle.PlayerEntity, target battle.Entity, fightInstance *battle.Fight)
	GetEvents() *map[types.CustomTrigger]func(owner *battle.PlayerEntity)
}

// Lvl
type PlayerSkillLevel interface {
	PlayerSkillUpgradable

	GetLevel() int
	GetPath() types.SkillPath
}

type PSkillUpgrade struct {
	Name        string
	Description string
	Id          string
	Execute     func(owner battle.PlayerEntity, target battle.Entity, fightInstance *battle.Fight)
	OnEvent     *map[types.CustomTrigger]func(owner *battle.PlayerEntity)
}
