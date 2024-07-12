package inventory

import (
	"sao/battle"
	"sao/types"
	"sao/utils"

	"github.com/google/uuid"
)

type EnduranceSkill struct{}

func (skill EnduranceSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	return nil
}

func (skill EnduranceSkill) GetPath() types.SkillPath {
	return types.PathControl
}

func (skill EnduranceSkill) GetUUID() uuid.UUID {
	return uuid.Nil
}

func (skill EnduranceSkill) IsLevelSkill() bool {
	return true
}

type END_LVL_1 struct {
	EnduranceSkill
	DefaultCost
	DefaultActiveTrigger
	NoEvents
	NoStats
}

func (skill END_LVL_1) GetName() string {
	return "Poziom 1 - wytrzymałość"
}

func (skill END_LVL_1) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	baseIncrease := 10
	baseDuration := 1

	if HasUpgrade(upgrades, 1) {
		baseIncrease = 20
		baseIncrease += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_HP), 3)
		baseIncrease += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AD), 2)
	}

	if HasUpgrade(upgrades, 2) {
		baseDuration++
	}

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Target: target.(battle.Entity).GetUUID(),
		Source: owner.(battle.PlayerEntity).GetUUID(),
		Meta: battle.ActionEffect{
			Effect:   battle.EFFECT_SHIELD,
			Value:    baseIncrease,
			Duration: baseDuration,
			Caster:   owner.(battle.PlayerEntity).GetUUID(),
		},
	})

	return nil
}

func (skill END_LVL_1) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill END_LVL_1) GetCooldown(upgrades int) int {
	baseCD := skill.GetCD()

	if HasUpgrade(upgrades, 1) {
		return baseCD - 1
	}

	return baseCD
}

func (skill END_LVL_1) GetDescription() string {
	return "Daje tarczę o wartości 10 na jedną turę"
}

func (skill END_LVL_1) GetLevel() int {
	return 1
}

func (skill END_LVL_1) GetUpgrades() []PlayerSkillUpgrade {
	return []PlayerSkillUpgrade{
		{
			Name:        "Ulepszenie 1",
			Id:          "Cooldown",
			Events:      nil,
			Description: "Zmniejsza czas odnowienia o 1 turę",
		},
		{
			Name:        "Ulepszenie 2",
			Id:          "Damage",
			Events:      nil,
			Description: "Zwiększa wartość tarczy do 20 + 3%HP + 2%ATK",
		},
		{
			Name:        "Ulepszenie 3",
			Id:          "Duration",
			Events:      nil,
			Description: "Zwiększa czas trwania o 1 turę",
		},
	}
}

type END_LVL_2 struct {
	EnduranceSkill
	NoExecute
	NoStats
	NoTrigger
}

func (skill END_LVL_2) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return map[types.CustomTrigger]func(owner interface{}){
		types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
			owner.(battle.PlayerEntity).SetLevelStat(types.STAT_HP, 15)
		},
	}
}

func (skill END_LVL_2) GetUpgrades() []PlayerSkillUpgrade {
	return []PlayerSkillUpgrade{
		{
			Name: "Ulepszenie 1",
			Id:   "Increase",
			Events: &map[types.CustomTrigger]func(owner battle.PlayerEntity){
				types.CUSTOM_TRIGGER_UNLOCK: func(owner battle.PlayerEntity) {
					owner.SetLevelStat(types.STAT_HP, 20)
				},
			},
			Description: "Zwiększa wartość do 20",
		},
		{
			Name: "Ulepszenie 2",
			Id:   "DEFIncrease",
			Events: &map[types.CustomTrigger]func(owner battle.PlayerEntity){
				types.CUSTOM_TRIGGER_UNLOCK: func(owner battle.PlayerEntity) {
					owner.AppendDerivedStat(types.DerivedStat{
						Base:    types.STAT_HP_PLUS,
						Derived: types.STAT_DEF,
						Percent: 1,
					})
				},
			},
			Description: "Zwiększa wartość pancerza o 1% dodatkowego HP",
		},
		{
			Name: "Ulepszenie 3",
			Id:   "RESIncrease",
			Events: &map[types.CustomTrigger]func(owner battle.PlayerEntity){
				types.CUSTOM_TRIGGER_UNLOCK: func(owner battle.PlayerEntity) {
					owner.AppendDerivedStat(types.DerivedStat{
						Base:    types.STAT_HP_PLUS,
						Derived: types.STAT_MR,
						Percent: 1,
					})
				},
			},
			Description: "Zwiększa wartość odporności na magię o 1% dodatkowego HP",
		},
	}
}

func (skill END_LVL_2) GetDescription() string {
	return "Zwiększa zdrowie zyskiwane co poziom do 15"
}

func (skill END_LVL_2) GetLevel() int {
	return 2
}

func (skill END_LVL_2) GetName() string {
	return "Poziom 2 - wytrzymałość"
}

type END_LVL_3 struct {
	EnduranceSkill
	DefaultCost
	NoEvents
	NoStats
	DefaultActiveTrigger
}

type END_LVL_3_EFFECT struct {
	NoCost
	NoCooldown
	NoEvents
	NoLevel
	Owner      uuid.UUID
	HealFactor int
	OnlyOwner  bool
}

func (effect END_LVL_3_EFFECT) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	if target.(battle.Entity).GetUUID() != effect.Owner && !effect.OnlyOwner {
		return nil
	}

	og := fightInstance.(*battle.Fight).Entities[effect.Owner]

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Target: effect.Owner,
		Source: effect.Owner,
		Meta: battle.ActionEffect{
			Effect: battle.EFFECT_HEAL_SELF,
			Value:  utils.PercentOf(og.Entity.(battle.PlayerEntity).GetStat(types.STAT_HP), effect.HealFactor),
		},
	})

	return nil
}

func (effect END_LVL_3_EFFECT) GetDescription() string {
	return "Leczy o x% HP"
}

func (effect END_LVL_3_EFFECT) GetName() string {
	return "Efekt wytrzymałości - poziom 3"
}

func (effect END_LVL_3_EFFECT) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType: types.TRIGGER_ATTACK_GOT_HIT,
		},
	}
}

func (skill END_LVL_3) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	healFactor := 5

	if HasUpgrade(upgrades, 3) {
		healFactor += 5
	}

	target.(battle.Entity).AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value: END_LVL_3_EFFECT{
			Owner:      owner.(battle.PlayerEntity).GetUUID(),
			HealFactor: healFactor,
			OnlyOwner:  HasUpgrade(upgrades, 1),
		},
		AfterUsage: false,
		Expire:     1,
	})

	return nil
}

func (skill END_LVL_3) GetName() string {
	return "Poziom 3 - wytrzymałość"
}

func (skill END_LVL_3) GetLevel() int {
	return 3
}

func (skill END_LVL_3) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill END_LVL_3) GetCooldown(upgrades int) int {
	baseCD := skill.GetCD()

	if HasUpgrade(upgrades, 2) {
		return baseCD - 1
	}

	return baseCD
}

func (skill END_LVL_3) GetDescription() string {
	return "Oznacza cel na turę, atak rozbije oznaczenie i wyleczy"
}

func (skill END_LVL_3) GetUpgrades() []PlayerSkillUpgrade {
	return []PlayerSkillUpgrade{
		{
			Name:        "Ulepszenie 1",
			Id:          "Ally",
			Events:      nil,
			Description: "Sojusznik może rozbić",
		},
		{
			Name:        "Ulepszenie 2",
			Id:          "Cooldown",
			Events:      nil,
			Description: "Zmniejsza czas odnowienia o 1 turę",
		},
		{
			Name:        "Ulepszenie 3",
			Id:          "Increase",
			Events:      nil,
			Description: "Zwiększa leczenie o X",
		},
	}
}
