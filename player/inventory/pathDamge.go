package inventory

import (
	"fmt"
	"sao/battle"
	"sao/types"
	"sao/utils"

	"github.com/google/uuid"
)

type DamageSkill struct{}

func (skill DamageSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	return nil
}

func (skill DamageSkill) GetPath() types.SkillPath {
	return types.PathDamage
}

func (skill DamageSkill) GetUUID() uuid.UUID {
	return uuid.Nil
}

func (skill DamageSkill) IsLevelSkill() bool {
	return true
}

func (skill DamageSkill) CanUse(owner interface{}, fightInstance interface{}, upgrades int) bool {
	return true
}

type DMG_LVL_1 struct {
	DamageSkill
	DefaultCost
	NoEvents
	NoStats
}

func (skill DMG_LVL_1) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE, Flags: types.FLAG_INSTANT_SKILL}
}

func (skill DMG_LVL_1) GetUpgradableTrigger(upgrades int) types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE, Flags: types.FLAG_INSTANT_SKILL}
}

func (skill DMG_LVL_1) GetName() string {
	return "Poziom 1 - obrażenia"
}

type DMG_LVL_1_Effect struct {
	NoEvents
	Damage int
}

func (skill DMG_LVL_1_Effect) GetName() string {
	return "Obrażenia 1 - Efekt"
}

func (skill DMG_LVL_1_Effect) GetDescription() string {
	return "Zwiększa kolejnego ataku o 10 na 1 turę"
}

func (skill DMG_LVL_1) GetUpgradableDescription(upgrades int) string {
	upgradeSegment := []string{"10", "1 turę"}

	if HasUpgrade(upgrades, 2) {
		upgradeSegment[0] = "20 + 1%ATK + 1%AP"
	}

	if HasUpgrade(upgrades, 3) {
		upgradeSegment[1] = "2 tury"
	}

	return fmt.Sprintf("Zwiększa kolejnego ataku o %s na %s.", upgradeSegment[0], upgradeSegment[1])
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
		Type:  types.TRIGGER_PASSIVE,
		Event: types.TRIGGER_ATTACK_BEFORE,
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

func (skill DMG_LVL_1) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	baseIncrease := 10
	baseDuration := 1

	if HasUpgrade(upgrades, 2) {
		baseIncrease = 20
		baseIncrease += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AD), 1)
		baseIncrease += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AP), 1)
	}

	if HasUpgrade(upgrades, 3) {
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

	if HasUpgrade(upgrades, 1) {
		return baseCD - 1
	}

	return baseCD
}

func (skill DMG_LVL_1) GetDescription() string {
	return "Zwiększa obrażenia kolejnego ataku o 10 na jedną turę"
}

func (skill DMG_LVL_1) GetLevel() int {
	return 1
}

func (skill DMG_LVL_1) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "Cooldown",
			Description: "Zmniejsza czas odnowienia o 1 turę",
		},
		{
			Id:          "Damage",
			Description: "Wartość zwiększenia wynosi 20 + 1%ATK + 1%AP",
		},
		{
			Id:          "Duration",
			Description: "Zwiększa czas trwania o 1 turę",
		},
	}
}

type DMG_LVL_2 struct {
	DamageSkill
	NoExecute
	NoStats
	NoEvents
	NoTrigger
}

func (skill DMG_LVL_2) GetName() string {
	return "Poziom 2 - obrażenia"
}

func (skill DMG_LVL_2) GetLevel() int {
	return 2
}

func (skill DMG_LVL_2) GetDescription() string {
	return "Zwiększa otrzymywany atak co poziom o 5"
}

func (skill DMG_LVL_2) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return map[types.CustomTrigger]func(owner interface{}){
		types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
			owner.(battle.PlayerEntity).SetLevelStat(types.STAT_AD, owner.(battle.PlayerEntity).GetLevelStat(types.STAT_AD)+5)
		},
	}
}

func (skill DMG_LVL_2) GetUpgradableDescription(upgrades int) string {

	upgradeSegments := []string{"", "", ""}

	if HasUpgrade(upgrades, 1) {
		upgradeSegments[0] = "\nOtrzymujesz 1% AP jako przebicie magiczne."
	}

	if HasUpgrade(upgrades, 2) {
		upgradeSegments[1] = "\nOtrzymujesz 10% ATK jako AP"
	}

	if HasUpgrade(upgrades, 3) {
		upgradeSegments[2] = "\nOtrzymujesz 1% ATK jako przebicie pancerza"
	}

	return "Zwiększa otrzymywany atak co poziom o 5." + upgradeSegments[0] + upgradeSegments[1] + upgradeSegments[2]
}

func (skill DMG_LVL_2) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id: "APPen",
			Events: &map[types.CustomTrigger]func(owner interface{}){
				types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
					owner.(battle.PlayerEntity).AppendDerivedStat(types.DerivedStat{
						Base:    types.STAT_AP,
						Derived: types.STAT_AP,
						Percent: 1,
					})
				},
			},
			Description: "Otrzymujesz 1% AP jako przebicie magiczne",
		},
		{
			Id: "APStat",
			Events: &map[types.CustomTrigger]func(owner interface{}){
				types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
					owner.(battle.PlayerEntity).AppendDerivedStat(types.DerivedStat{
						Base:    types.STAT_AD,
						Derived: types.STAT_AP,
						Percent: 10,
					})
				},
			},
			Description: "Otrzymujesz 10% ATK jako AP",
		},
		{
			Id: "ADPen",
			Events: &map[types.CustomTrigger]func(owner interface{}){
				types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
					owner.(battle.PlayerEntity).AppendDerivedStat(types.DerivedStat{
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

type DMG_LVL_3 struct {
	DamageSkill
	DefaultCost
	NoStats
	NoEvents
}

func (skill DMG_LVL_3) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE, Flags: types.FLAG_INSTANT_SKILL}
}

func (skill DMG_LVL_3) GetUpgradableTrigger(upgrades int) types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE, Flags: types.FLAG_INSTANT_SKILL}
}

func (skill DMG_LVL_3) GetName() string {
	return "Poziom 3 - obrażenia"
}

func (skill DMG_LVL_3) GetDescription() string {
	return "Kolejny atak zada dodatkowe 25 obrażeń"
}

func (skill DMG_LVL_3) GetLevel() int {
	return 3
}

func (skill DMG_LVL_3) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill DMG_LVL_3) GetCooldown(upgrades int) int {
	return skill.GetCD()
}

type DMG_LVL_3_Effect struct {
	NoEvents
	NoCost
	NoLevel
	NoCooldown
	Damage int
	Ripple bool
}

type DMG_LVL_3_Effect_2 struct {
	NoCost
	NoCooldown
	NoEvents
	NoLevel
	OriginalEntity uuid.UUID
}

func (skill DMG_LVL_3_Effect_2) GetName() string {
	return "Obrażenia 3 - Efekt 2"
}

func (skill DMG_LVL_3_Effect_2) GetDescription() string {
	return "Zadaje 25 obrażeń wszystkim wrogom (poza celem ataku)"
}

func (skill DMG_LVL_3_Effect_2) Execute(owner, target, fightInstance, meta interface{}) interface{} {

	dmgValue := 25 + utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AP), 20)

	for _, entity := range fightInstance.(*battle.Fight).GetAlliesFor(skill.OriginalEntity) {
		fightInstance.(*battle.Fight).HandleAction(battle.Action{
			Event:  battle.ACTION_DMG,
			Source: owner.(battle.Entity).GetUUID(),
			Target: entity.GetUUID(),
			Meta:   battle.ActionDamage{Damage: []battle.Damage{{Value: dmgValue, IsPercent: false, Type: 0}}},
		})
	}

	return nil
}

func (skill DMG_LVL_3_Effect_2) GetTrigger() types.Trigger {
	return types.Trigger{
		Type:  types.TRIGGER_PASSIVE,
		Event: types.TRIGGER_ATTACK_HIT,
	}
}

func (skill DMG_LVL_3_Effect) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	owner.(battle.PlayerEntity).AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value:      DMG_LVL_3_Effect_2{OriginalEntity: target.(battle.Entity).GetUUID()},
		AfterUsage: false,
		Expire:     1,
	})

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

func (skill DMG_LVL_3_Effect) GetDescription() string {
	return "Kolejny atak zada dodatkowe 25 obrażeń"
}

func (skill DMG_LVL_3_Effect) GetName() string {
	return "Obrażenia 3 - Efekt"
}

func (skill DMG_LVL_3_Effect) GetTrigger() types.Trigger {
	return types.Trigger{
		Type:  types.TRIGGER_PASSIVE,
		Event: types.TRIGGER_ATTACK_BEFORE,
	}
}

func (skill DMG_LVL_3) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	baseDamage := 25
	ripple := HasUpgrade(upgrades, 2)

	if HasUpgrade(upgrades, 1) {
		baseDamage = 30
		baseDamage += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AD), 2)
		baseDamage += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AP), 2)
	}

	if HasUpgrade(upgrades, 3) && len(fightInstance.(*battle.Fight).GetAlliesFor(target.(battle.Entity).GetUUID())) == 0 {
		baseDamage += utils.PercentOf(baseDamage, 125)
	}

	owner.(battle.PlayerEntity).AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value:      DMG_LVL_3_Effect{Damage: baseDamage, Ripple: ripple},
		AfterUsage: true,
		Expire:     1,
	})

	return nil
}

func (skill DMG_LVL_3) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "Damage",
			Description: "Wartość zwiększenia wynosi 30 + 2%ATK + 2%AP",
		},
		{
			Id:          "Ripple",
			Description: "Zadaje dodatkowe 25+20% AP obrażeń wszystkim wrogom (poza celem ataku)",
		},
		{
			Id:          "Isolate",
			Description: "Zadaje dodatkowe 25% obrażeń gdy jesteś w walce z jednym (żywym) przeciwnikiem",
		},
	}
}

func (skill DMG_LVL_3) GetUpgradableDescription(upgrades int) string {
	upgradeSegments := []string{"25", "", ""}

	if HasUpgrade(upgrades, 1) {
		upgradeSegments[0] = "30 + 2%ATK + 2%AP"
	}

	if HasUpgrade(upgrades, 2) {
		upgradeSegments[1] = " Zadaje wszystkim wrogom (poza celem ataku) 25+20% AP obrażeń."
	}

	if HasUpgrade(upgrades, 3) {
		upgradeSegments[2] = " Zadaje dodatkowe 25% obrażeń gdy jesteś w walce z jednym (żywym) przeciwnikiem"
	}

	return fmt.Sprintf("Kolejny atak zada dodatkowe %s obrażeń.%s %s", upgradeSegments[0], upgradeSegments[2], upgradeSegments[1])
}

type DMG_LVL_4 struct {
	DamageSkill
	DefaultCost
	NoEvents
	NoStats
}

func (skill DMG_LVL_4) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE, Flags: types.FLAG_INSTANT_SKILL}
}

func (skill DMG_LVL_4) GetUpgradableTrigger(upgrades int) types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE, Flags: types.FLAG_INSTANT_SKILL}
}

func (skill DMG_LVL_4) GetLevel() int {
	return 4
}

func (skill DMG_LVL_4) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "Increase",
			Description: "Wartość zwiększenia wynosi 20.",
		},
		{
			Id:          "PartyWide",
			Description: "Efekt zwiększenia zasobów działa na całą drużynę",
		},
		{
			Id:          "ManaReturn",
			Description: "Zabicie przywróci punkt many",
		},
	}
}

func (skill DMG_LVL_4) GetName() string {
	return "Poziom 4 - obrażenia"
}

func (skill DMG_LVL_4) GetDescription() string {
	return "Zwiększa zasoby (doświadczenie i złoto) o 10% po zabiciu przeciwnika"
}

func (skill DMG_LVL_4) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill DMG_LVL_4) GetCooldown(upgrades int) int {
	return skill.GetCD()
}

type DMG_LVL_4_Effect struct {
	NoCost
	NoCooldown
	NoEvents
	NoLevel
	IncreaseValue int
	PartyWide     bool
	ManaReturn    bool
}

func (skill DMG_LVL_4_Effect) GetName() string {
	return "Obrażenia 4 - Efekt"
}

func (skill DMG_LVL_4_Effect) GetDescription() string {
	return "Kolejny zabity wróg (w ciągu tury) da 10% więcej zasobów"
}

func (skill DMG_LVL_4_Effect) GetTrigger() types.Trigger {
	return types.Trigger{
		Type:  types.TRIGGER_PASSIVE,
		Event: types.TRIGGER_ATTACK_HIT,
	}
}

func (skill DMG_LVL_4_Effect) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	if target.(battle.Entity).GetCurrentHP() <= 0 {
		for _, loot := range target.(battle.Entity).GetLoot() {
			if loot.Type == battle.LOOT_ITEM {
				continue
			}

			tempLoot := loot

			tempLoot.Count = utils.PercentOf(loot.Count, skill.IncreaseValue)

			fightInstance.(*battle.Fight).AddAdditionalLoot(tempLoot, owner.(battle.Entity).GetUUID(), skill.PartyWide)

			if skill.ManaReturn {
				owner.(battle.PlayerEntity).RestoreMana(1)
			}
		}
	}

	return nil
}

func (skill DMG_LVL_4) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	increaseValue := 10

	if HasUpgrade(upgrades, 1) {
		increaseValue = 20
	}

	owner.(battle.PlayerEntity).AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value: DMG_LVL_4_Effect{IncreaseValue: increaseValue,
			PartyWide:  HasUpgrade(upgrades, 2),
			ManaReturn: HasUpgrade(upgrades, 3),
		},
		AfterUsage: true,
		Expire:     1,
	})

	return nil
}

func (skill DMG_LVL_4) GetUpgradableDescription(upgrades int) string {
	upgradeSegments := []string{"10%", "", ""}

	if HasUpgrade(upgrades, 1) {
		upgradeSegments[0] = "20%"
	}

	if HasUpgrade(upgrades, 2) {
		upgradeSegments[1] = " Efekt zwiększenia zasobów działa na całą drużynę,"
	}

	if HasUpgrade(upgrades, 3) {
		upgradeSegments[2] = " Dodatkowo przy aktywacji przywróci tobie punkt many,"
	}

	return fmt.Sprintf("Kolejny zabity przeciwnik (w ciągu tury) da zwiększone zasoby (doświadczenie i złoto) o %s.%s%s", upgradeSegments[0], upgradeSegments[2], upgradeSegments[1])
}

type DMG_LVL_5 struct {
	DamageSkill
	DefaultCost
	NoEvents
	NoStats
}

func (skill DMG_LVL_5) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE, Target: &types.TargetTrigger{Target: types.TARGET_ENEMY, MaxTargets: 3}}
}

func (skill DMG_LVL_5) GetUpgradableTrigger(upgrades int) types.Trigger {

	baseTrigger := skill.GetTrigger()

	if HasUpgrade(upgrades, 2) {
		baseTrigger.Target.MaxTargets = 6
	}

	return baseTrigger
}

func (skill DMG_LVL_5) GetName() string {
	return "Poziom 5 - obrażenia"
}

func (skill DMG_LVL_5) GetDescription() string {
	return "Zadaje 3 wrogom 100 + 40% ATK + 30% AP"
}

func (skill DMG_LVL_5) GetLevel() int {
	return 5
}

func (skill DMG_LVL_5) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill DMG_LVL_5) GetCooldown(upgrades int) int {
	return skill.GetCD()
}

func (skill DMG_LVL_5) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "DamageRat",
			Description: "Obrażenia wynoszą 100 + 50% ATK + 40% AP",
		},
		{
			Id:          "MaxCount",
			Description: "Maksymalna ilość celów wynosi 6",
		},
		{
			Id:          "DamageHP",
			Description: "Obrażenia umiejętności są zwiększone o 25% jeśli cel ma więcej niż 70% HP",
		},
	}
}

func (skill DMG_LVL_5) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	baseDamage := 100
	baseDamage += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AD), 40)
	baseDamage += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AP), 30)

	if HasUpgrade(upgrades, 1) {
		baseDamage += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AD), 10)
		baseDamage += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AP), 10)
	}

	validTargets := fightInstance.(*battle.Fight).GetEnemiesFor(owner.(battle.Entity).GetUUID())
	var targets []battle.Entity

	if HasUpgrade(upgrades, 2) {
		if len(validTargets) > 6 {
			targets = validTargets[:6]
		} else {
			targets = validTargets
		}
	} else {
		if len(validTargets) > 3 {
			targets = validTargets[:3]
		} else {
			targets = validTargets
		}
	}

	for _, entity := range targets {

		maxHPThreshold := utils.PercentOf(entity.GetStat(types.STAT_HP), 70)

		if HasUpgrade(upgrades, 3) && entity.GetCurrentHP() > maxHPThreshold {
			baseDamage = utils.PercentOf(baseDamage, 125)
		}

		fightInstance.(*battle.Fight).HandleAction(battle.Action{
			Event:  battle.ACTION_DMG,
			Source: owner.(battle.Entity).GetUUID(),
			Target: entity.GetUUID(),
			Meta: battle.ActionDamage{
				Damage: []battle.Damage{
					{
						Value:     baseDamage,
						IsPercent: false,
						Type:      0,
					},
				},
			},
		})
	}

	return nil
}

func (skill DMG_LVL_5) GetUpgradableDescription(upgrades int) string {
	upgradeSegments := []string{"40% ATK + 30% AP", ""}

	if HasUpgrade(upgrades, 1) {
		upgradeSegments[0] = "50% ATK + 40% AP"
	}

	targetCount := 3

	if HasUpgrade(upgrades, 2) {
		targetCount = 6
	}

	if HasUpgrade(upgrades, 3) {
		upgradeSegments[1] = " Obrażenia (celu) są zwiększone o 25% jeśli cel ma więcej niż 70% HP."
	}

	return fmt.Sprintf("Zadaje %d wrogom 100 + %s.%s", targetCount, upgradeSegments[0], upgradeSegments[1])
}
