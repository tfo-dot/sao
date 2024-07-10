package inventory

import (
	"sao/battle"
	"sao/types"
	"sao/utils"

	"github.com/google/uuid"
)

type DMG_LVL_1 struct{ BaseLvlSkill }

type DMG_LVL_1_Effect struct {
	Damage int
}

func (skill DMG_LVL_1_Effect) GetName() string {
	return "Obrażenia 1 - Efekt"
}

func (skill DMG_LVL_1_Effect) GetDescription() string {
	return "Zwiększa obrażenia o 10 na turę"
}

func (skill DMG_LVL_1_Effect) GetUUID() uuid.UUID {
	return uuid.New()
}

func (skill DMG_LVL_1_Effect) GetCD() int {
	return 0
}

func (skill DMG_LVL_1_Effect) GetCost() int {
	return 0
}

func (skill DMG_LVL_1_Effect) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType: types.TRIGGER_ATTACK_BEFORE,
		},
	}
}

func (skill DMG_LVL_1_Effect) IsLevelSkill() bool {
	return false
}

func (skill DMG_LVL_1_Effect) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	return types.AttackTriggerMeta{
		Effects: []types.DamagePartial{
			{
				Value:   skill.Damage,
				Percent: false,
				Type:    0,
			},
		},
	}
}

func (skill DMG_LVL_1_Effect) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return nil
}

func (skill DMG_LVL_1) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	baseIncrease := 10
	baseDuration := 1

	if upgrades&(1<<1) == 1 {
		baseIncrease = 20
		baseIncrease += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AD), 1)
		baseIncrease += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AP), 1)
	}

	if upgrades&(1<<2) == 1 {
		baseDuration++
	}

	owner.(battle.PlayerEntity).AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value:      DMG_LVL_1_Effect{Damage: baseIncrease},
		AfterUsage: false,
		Expire:     baseDuration,
	})

	return nil
}

func (skill DMG_LVL_1) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill DMG_LVL_1) GetCooldown(upgrades int) int {
	baseCD := skill.GetCD()

	if upgrades&(1<<0) == 1 {
		return baseCD - 1
	}

	return baseCD
}

func (skill DMG_LVL_1) GetDescription() string {
	return "Zwiększa obrażenia o 10 na jedną turę"
}

func (skill DMG_LVL_1) GetPath() types.SkillPath {
	return types.PathDamage
}

func (skill DMG_LVL_1) GetLevel() int {
	return 1
}

func (skill DMG_LVL_1) GetUpgrades() []PlayerSkillUpgrade {
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
			Description: "Zwiększa obrażenia do 20 + 1%ATK + 1%AP",
		},
		{
			Name:        "Ulepszenie 3",
			Id:          "Duration",
			Events:      nil,
			Description: "Zwiększa czas trwania o 1 turę",
		},
	}
}

type DMG_LVL_2 struct{ BaseLvlSkill }

func (skill DMG_LVL_2) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	return nil
}

func (skill DMG_LVL_2) GetDescription() string {
	return "Zwiększa otrzymywany atak co poziom o 20"
}

func (skill DMG_LVL_2) GetPath() types.SkillPath {
	return types.PathDamage
}

func (skill DMG_LVL_2) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return map[types.CustomTrigger]func(owner interface{}){
		types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
			owner.(battle.PlayerEntity).SetLevelStat(types.STAT_AD, 20)
		},
	}
}

func (skill DMG_LVL_2) GetUpgrades() []PlayerSkillUpgrade {
	return []PlayerSkillUpgrade{
		{
			Name: "Ulepszenie 1",
			Id:   "APPen",
			Events: &map[types.CustomTrigger]func(owner battle.PlayerEntity){
				types.CUSTOM_TRIGGER_UNLOCK: func(owner battle.PlayerEntity) {
					owner.AppendDerivedStat(types.DerivedStat{
						Base:    types.STAT_AP,
						Derived: types.STAT_AP,
						Percent: 1,
					})
				},
			},
			Description: "Otrzymujesz 1% AP jako przebicie magiczne",
		},
		{
			Name: "Ulepszenie 2",
			Id:   "APStat",
			Events: &map[types.CustomTrigger]func(owner battle.PlayerEntity){
				types.CUSTOM_TRIGGER_UNLOCK: func(owner battle.PlayerEntity) {
					owner.AppendDerivedStat(types.DerivedStat{
						Base:    types.STAT_AD,
						Derived: types.STAT_AP,
						Percent: 10,
					})
				},
			},
			Description: "Otrzymujesz 10% ATK jako AP",
		},
		{
			Name: "Ulepszenie 3",
			Id:   "ADPen",
			Events: &map[types.CustomTrigger]func(owner battle.PlayerEntity){
				types.CUSTOM_TRIGGER_UNLOCK: func(owner battle.PlayerEntity) {
					owner.AppendDerivedStat(types.DerivedStat{
						Base:    types.STAT_AD,
						Derived: types.STAT_LETHAL,
						Percent: 1,
					})
				},
			},
			Description: "Otrzymujesz 1% ATK jako przebicie pancerza",
		},
	}
}

// TODO ripple effect
type DMG_LVL_3 struct{ BaseLvlSkill }

func (skill DMG_LVL_3) GetDescription() string {
	return "Zadaje dodatkowe 25 obrażeń"
}

func (skill DMG_LVL_3) GetPath() types.SkillPath {
	return types.PathDamage
}

type DMG_LVL_3_Effect struct {
	Damage int
	Ripple bool
}

func (skill DMG_LVL_3_Effect) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	return types.AttackTriggerMeta{
		Effects: []types.DamagePartial{
			{
				Value:   skill.Damage,
				Percent: false,
				Type:    0,
			},
		},
	}
}

func (skill DMG_LVL_3_Effect) GetCD() int {
	return 0
}

func (skill DMG_LVL_3_Effect) GetCost() int {
	return 0
}

func (skill DMG_LVL_3_Effect) GetDescription() string {
	return "Zadaje dodatkowe 25 obrażeń"
}

func (skill DMG_LVL_3_Effect) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return nil
}

func (skill DMG_LVL_3_Effect) GetName() string {
	return "Obrażenia 3 - Efekt"
}

func (skill DMG_LVL_3_Effect) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_ATTACK_BEFORE,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
			TargetCount:   -1,
			Meta:          nil,
		},
	}
}

func (skill DMG_LVL_3_Effect) IsLevelSkill() bool {
	return false
}

func (skill DMG_LVL_3_Effect) GetUUID() uuid.UUID {
	return uuid.New()
}

func (skill DMG_LVL_3) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	baseDamage := 25
	ripple := upgrades&(1<<1) == 1

	if upgrades&(1<<0) == 1 {
		baseDamage = 30
		baseDamage += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AD), 2)
		baseDamage += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AP), 2)
	}

	if upgrades&(1<<2) == 1 && len(fightInstance.(*battle.Fight).GetAlliesFor(target.(battle.Entity).GetUUID())) == 0 {
		baseDamage += utils.PercentOf(baseDamage, 125)
	}

	owner.(battle.PlayerEntity).AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value:      DMG_LVL_3_Effect{Damage: baseDamage, Ripple: ripple},
		AfterUsage: true,
		Expire:     1,
	})

	return nil
}

func (skill DMG_LVL_3) GetUpgrades() []PlayerSkillUpgrade {
	return []PlayerSkillUpgrade{
		{
			Name:        "Ulepszenie 1",
			Id:          "Damage",
			Events:      nil,
			Description: "Zwiększa obrażenia do 30 + 2%ATK + 2%AP",
		},
		{
			Name:        "Ulepszenie 2",
			Id:          "Ripple",
			Events:      nil,
			Description: "Zadaje dodatkowe 25 obrażeń sąsiadom",
		},
		{
			Name:        "Ulepszenie 3",
			Id:          "Isolate",
			Events:      nil,
			Description: "Zadaje dodatkowe 125% obrażeń wyizolowanym celom",
		},
	}
}

// TODO rest of it lmao
type DMG_LVL_4 struct{ BaseLvlSkill }

type DMG_LVL_4_Effect struct {
	IncreaseValue int
	PartyWide     bool
}

func (skill DMG_LVL_4) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	return nil
}
