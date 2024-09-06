package inventory

import (
	"fmt"
	"sao/types"
	"sao/utils"

	"github.com/google/uuid"
)

type DamageSkill struct{}

func (skill DamageSkill) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
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

func (skill DamageSkill) CanUse(owner types.PlayerEntity, fightInstance types.FightInstance) bool {
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

func (skill DMG_LVL_1_Effect) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
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

func (skill DMG_LVL_1) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	baseIncrease := 10
	baseDuration := 1

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 2) {
		baseIncrease = 20
		baseIncrease += utils.PercentOf(owner.GetStat(types.STAT_AD), 1)
		baseIncrease += utils.PercentOf(owner.GetStat(types.STAT_AP), 1)
	}

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 3) {
		baseDuration++
	}

	owner.AppendTempSkill(types.WithExpire[types.PlayerSkill]{
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

func (skill DMG_LVL_2) GetEvents() map[types.CustomTrigger]func(owner types.PlayerEntity) {
	return map[types.CustomTrigger]func(owner types.PlayerEntity){
		types.CUSTOM_TRIGGER_UNLOCK: func(owner types.PlayerEntity) {
			owner.SetLevelStat(types.STAT_AD, owner.GetLevelStat(types.STAT_AD)+5)
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
			Events: &map[types.CustomTrigger]func(owner types.PlayerEntity){
				types.CUSTOM_TRIGGER_UNLOCK: func(owner types.PlayerEntity) {
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
			Id: "APStat",
			Events: &map[types.CustomTrigger]func(owner types.PlayerEntity){
				types.CUSTOM_TRIGGER_UNLOCK: func(owner types.PlayerEntity) {
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
			Id: "ADPen",
			Events: &map[types.CustomTrigger]func(owner types.PlayerEntity){
				types.CUSTOM_TRIGGER_UNLOCK: func(owner types.PlayerEntity) {
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

func (skill DMG_LVL_3_Effect_2) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {

	dmgValue := 25 + utils.PercentOf(owner.GetStat(types.STAT_AP), 20)

	for _, entity := range fightInstance.GetAlliesFor(skill.OriginalEntity) {
		fightInstance.HandleAction(types.Action{
			Event:  types.ACTION_DMG,
			Source: owner.GetUUID(),
			Target: entity.GetUUID(),
			Meta:   types.ActionDamage{Damage: []types.Damage{{Value: dmgValue, IsPercent: false, Type: 0}}},
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

func (skill DMG_LVL_3_Effect) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	owner.AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value:      DMG_LVL_3_Effect_2{OriginalEntity: target.GetUUID()},
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

func (skill DMG_LVL_3) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	baseDamage := 25
	ripple := HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 2)

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 1) {
		baseDamage = 30
		baseDamage += utils.PercentOf(owner.GetStat(types.STAT_AD), 2)
		baseDamage += utils.PercentOf(owner.GetStat(types.STAT_AP), 2)
	}

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 3) && len(fightInstance.GetAlliesFor(target.GetUUID())) == 0 {
		baseDamage += utils.PercentOf(baseDamage, 125)
	}

	owner.AppendTempSkill(types.WithExpire[types.PlayerSkill]{
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

func (skill DMG_LVL_4_Effect) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	if target.GetCurrentHP() <= 0 {
		for _, loot := range target.GetLoot() {
			if loot.Type == types.LOOT_ITEM {
				continue
			}

			tempLoot := loot

			tempLoot.Count = utils.PercentOf(loot.Count, skill.IncreaseValue)

			fightInstance.AddAdditionalLoot(tempLoot, owner.GetUUID(), skill.PartyWide)

			if skill.ManaReturn {
				owner.RestoreMana(1)
			}
		}
	}

	return nil
}

func (skill DMG_LVL_4) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	increaseValue := 10

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 1) {
		increaseValue = 20
	}

	owner.AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value: DMG_LVL_4_Effect{IncreaseValue: increaseValue,
			PartyWide:  HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 2),
			ManaReturn: HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 3),
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

func (skill DMG_LVL_5) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	baseDamage := 100
	baseDamage += utils.PercentOf(owner.GetStat(types.STAT_AD), 40)
	baseDamage += utils.PercentOf(owner.GetStat(types.STAT_AP), 30)

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 1) {
		baseDamage += utils.PercentOf(owner.GetStat(types.STAT_AD), 10)
		baseDamage += utils.PercentOf(owner.GetStat(types.STAT_AP), 10)
	}

	validTargets := fightInstance.GetEnemiesFor(owner.GetUUID())
	var targets []types.Entity

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 2) {
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

		if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 3) && entity.GetCurrentHP() > maxHPThreshold {
			baseDamage = utils.PercentOf(baseDamage, 125)
		}

		fightInstance.HandleAction(types.Action{
			Event:  types.ACTION_DMG,
			Source: owner.GetUUID(),
			Target: entity.GetUUID(),
			Meta: types.ActionDamage{
				Damage: []types.Damage{
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

type DMG_LVL_6 struct {
	DamageSkill
	DefaultCost
	NoStats
	NoEvents
}

func (skill DMG_LVL_6) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE, Target: &types.TargetTrigger{Target: types.TARGET_ENEMY, MaxTargets: 1}}
}

func (skill DMG_LVL_6) GetUpgradableTrigger(upgrades int) types.Trigger {
	return skill.GetTrigger()
}

type DMG_LVL_6_EFFECT struct {
	NoCost
	NoCooldown
	NoEvents
	NoLevel
	Owner          uuid.UUID
	DmgFactor      int
	ExecuteUpgrade bool
}

func (effect DMG_LVL_6_EFFECT) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	if target.GetUUID() != effect.Owner {
		return nil
	}

	og := fightInstance.GetEntity(effect.Owner)

	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_DMG,
		Target: effect.Owner,
		Source: effect.Owner,
		Meta: types.ActionDamage{
			Damage: []types.Damage{
				{
					Value:     utils.PercentOf(target.GetCurrentHP(), effect.DmgFactor),
					IsPercent: false,
					Type:      0,
				},
			},
		},
	})

	if effect.ExecuteUpgrade {
		execThreshold := utils.PercentOf(og.GetStat(types.STAT_AD), 10)

		if target.GetCurrentHP() < execThreshold {
			fightInstance.HandleAction(types.Action{
				Event:  types.ACTION_DMG,
				Target: effect.Owner,
				Source: effect.Owner,
				Meta: types.ActionDamage{
					Damage: []types.Damage{
						{
							Value:     target.GetCurrentHP() + 1,
							IsPercent: false,
							Type:      types.DMG_TRUE,
						},
					},
				},
			})
		}
	}

	return nil
}

func (effect DMG_LVL_6_EFFECT) GetDescription() string {
	return "Zadaje obrażenia równe 5% maksymalnego HP celu"
}

func (effect DMG_LVL_6_EFFECT) GetName() string {
	return "Efekt obrażenia - poziom 6"
}

func (effect DMG_LVL_6_EFFECT) GetTrigger() types.Trigger {
	return types.Trigger{
		Type:  types.TRIGGER_PASSIVE,
		Event: types.TRIGGER_ATTACK_GOT_HIT,
	}
}

func (skill DMG_LVL_6) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	dmgFactor := 5

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 1) {
		dmgFactor += 5
	}

	target.AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value: DMG_LVL_6_EFFECT{
			Owner:          owner.GetUUID(),
			DmgFactor:      dmgFactor,
			ExecuteUpgrade: HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 2),
		},
		AfterUsage: false,
		Expire:     1,
	})

	return nil
}

func (skill DMG_LVL_6) GetName() string {
	return "Poziom 6 - obrażenia"
}

func (skill DMG_LVL_6) GetDescription() string {
	return "Oznacza wroga, uderzenie zada mu 5% aktualnego HP"
}

func (skill DMG_LVL_6) GetLevel() int {
	return 6
}

func (skill DMG_LVL_6) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "DmgIncrease",
			Description: "Obrażenia zwiększone o 5%",
		},
		{
			Id:          "Execute",
			Description: "Gdy zdrowie celu po rozbiciu spadnie poniżej 10% ATK użytkownika dobije go.",
		},
	}
}

func (skill DMG_LVL_6) GetUpgradableDescription(upgrades int) string {
	dmgFactor := 5

	if HasUpgrade(upgrades, 1) {
		dmgFactor = 10
	}

	executeUpgrade := ""

	if HasUpgrade(upgrades, 2) {
		executeUpgrade = " Gdy zdrowie celu po rozbiciu spadnie poniżej 10% ATK użytkownika dobije go."
	}

	return fmt.Sprintf("Oznacza wroga, uderzenie zada mu %d%%.%s", dmgFactor, executeUpgrade)
}

func (skill DMG_LVL_6) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill DMG_LVL_6) GetCooldown(upgrades int) int {
	return skill.GetCD()
}

type DMG_ULT_1 struct {
	DamageSkill
	DefaultCost
	NoStats
	NoEvents
	DefaultActiveTrigger
}

func (skill DMG_ULT_1) GetName() string {
	return "Poziom 10 - obrażenia"
}

func (skill DMG_ULT_1) GetDescription() string {
	return "Na 2 tury prowokuje wszystkich wrogów niwelując wszystkie otrzymane obrażenia.\nPo zakończeniu prowokacji leczy się całkowicie, zyskując jedną turę odporności na efekty kontroli tłumu, a wszyscy wrogowie zostają ogłuszeni.\nNastępnie przez 8 tur zyskuje 10% kradzieży życia oraz tarczę równą 50% przyjętych obrażeń. Jego AD i AP jest zwiększone o 120% i dodatkowo o 1% przyjętych obrażeń"
}

func (skill DMG_ULT_1) GetLevel() int {
	return 10
}

func (skill DMG_ULT_1) GetCD() int {
	return 10
}

func (skill DMG_ULT_1) GetCooldown(upgrades int) int {
	return skill.GetCD()
}

func (skill DMG_ULT_1) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{}
}

func (skill DMG_ULT_1) GetUpgradableDescription(upgrades int) string {
	return skill.GetDescription()
}

func (skill DMG_ULT_1) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	return skill.Execute(owner, target, fightInstance, meta)
}

type DMG_ULT_1_Effect struct {
	NoCost
	NoCooldown
	NoEvents
	NoLevel
	NoStats
}

func (skill DMG_ULT_1_Effect) GetName() string {
	return "Obrażenia 10 - Efekt"
}

func (skill DMG_ULT_1_Effect) GetDescription() string {
	return ""
}

func (skill DMG_ULT_1_Effect) GetTrigger() types.Trigger {
	return types.Trigger{
		Type:  types.TRIGGER_PASSIVE,
		Event: types.TRIGGER_DAMAGE_GOT_HIT,
	}
}

func (skill DMG_ULT_1_Effect) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	return nil
}

func (skill DMG_ULT_1) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	handlerUuid := fightInstance.AppendEventHandler(
		owner.GetUUID(),
		types.TRIGGER_DAMAGE_GOT_HIT,
		func(owner, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
			for _, event := range meta.(types.DamageTriggerMeta).Effects {
				if event.Percent {
					continue
				}

				dmgTaken := target.(types.PlayerEntity).GetLevelSkillMeta(skill.GetLevel())
				target.(types.PlayerEntity).SetLevelSkillMeta(skill.GetLevel(), dmgTaken.(int)+event.Value)
			}

			return nil
		},
	)

	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_EFFECT,
		Target: owner.GetUUID(),
		Source: owner.GetUUID(),
		Meta: types.ActionEffect{
			Effect:   types.EFFECT_RESIST,
			Value:    0,
			Duration: 2,
			Uuid:     uuid.New(),
			Meta: types.ActionEffectResist{
				All:       true,
				Value:     0,
				IsPercent: false,
				DmgType:   4,
			},
			Caster: owner.GetUUID(),
			Target: owner.GetUUID(),
			Source: types.SOURCE_ND,
		},
	})

	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_EFFECT,
		Target: owner.GetUUID(),
		Source: owner.GetUUID(),
		Meta: types.ActionEffect{
			Effect:   types.EFFECT_TAUNT,
			Value:    0,
			Duration: 2,
			Uuid:     uuid.New(),
			Meta:     nil,
			Caster:   owner.GetUUID(),
			Target:   owner.GetUUID(),
			Source:   types.SOURCE_ND,
			OnExpire: func(owner types.Entity, fight types.FightInstance, meta types.ActionEffect) {
				owner.Heal(owner.GetStat(types.STAT_HP))
				fight.RemoveEventHandler(handlerUuid)

				dmgTaken := owner.(types.PlayerEntity).GetLevelSkillMeta(skill.GetLevel()).(int)

				for _, entity := range fight.GetEnemiesFor(owner.GetUUID()) {
					fight.HandleAction(types.Action{
						Event:  types.ACTION_EFFECT,
						Target: entity.GetUUID(),
						Source: owner.GetUUID(),
						Meta: types.ActionEffect{
							Effect:   types.EFFECT_STUN,
							Value:    0,
							Duration: 1,
							Uuid:     uuid.New(),
							Meta:     nil,
							Caster:   owner.GetUUID(),
							Target:   entity.GetUUID(),
							Source:   types.SOURCE_ND,
						},
					})
				}

				fight.HandleAction(types.Action{
					Event:  types.ACTION_EFFECT,
					Target: owner.GetUUID(),
					Source: owner.GetUUID(),
					Meta: types.ActionEffect{
						Effect:   types.EFFECT_STAT_INC,
						Value:    0,
						Duration: 8,
						Uuid:     uuid.New(),
						Meta: types.ActionEffectStat{
							Stat:      types.STAT_ATK_VAMP,
							Value:     10,
							IsPercent: false,
						},
						Caster: owner.GetUUID(),
						Target: owner.GetUUID(),
						Source: types.SOURCE_ND,
					},
				})

				fight.HandleAction(types.Action{
					Event:  types.ACTION_EFFECT,
					Target: owner.GetUUID(),
					Source: owner.GetUUID(),
					Meta: types.ActionEffect{
						Effect:   types.EFFECT_SHIELD,
						Value:    utils.PercentOf(dmgTaken, 50),
						Duration: 8,
						Uuid:     uuid.New(),
						Meta:     nil,
						Caster:   owner.GetUUID(),
						Target:   owner.GetUUID(),
						Source:   types.SOURCE_ND,
					},
				})

				fight.HandleAction(types.Action{
					Event:  types.ACTION_EFFECT,
					Target: owner.GetUUID(),
					Source: owner.GetUUID(),
					Meta: types.ActionEffect{
						Effect:   types.EFFECT_STAT_INC,
						Value:    0,
						Duration: 8,
						Uuid:     uuid.New(),
						Meta: types.ActionEffectStat{
							Stat:      types.STAT_AD,
							Value:     utils.PercentOf(dmgTaken, 1),
							IsPercent: false,
						},
						Caster: owner.GetUUID(),
						Target: owner.GetUUID(),
						Source: types.SOURCE_ND,
					},
				})

				fight.HandleAction(types.Action{
					Event:  types.ACTION_EFFECT,
					Target: owner.GetUUID(),
					Source: owner.GetUUID(),
					Meta: types.ActionEffect{
						Effect:   types.EFFECT_STAT_INC,
						Value:    0,
						Duration: 8,
						Uuid:     uuid.New(),
						Meta: types.ActionEffectStat{
							Stat:      types.STAT_AP,
							Value:     utils.PercentOf(dmgTaken, 1),
							IsPercent: false,
						},
						Caster: owner.GetUUID(),
						Target: owner.GetUUID(),
						Source: types.SOURCE_ND,
					},
				})

				fight.HandleAction(types.Action{
					Event:  types.ACTION_EFFECT,
					Target: owner.GetUUID(),
					Source: owner.GetUUID(),
					Meta: types.ActionEffect{
						Effect:   types.EFFECT_STAT_INC,
						Value:    0,
						Duration: 8,
						Uuid:     uuid.New(),
						Meta: types.ActionEffectStat{
							Stat:      types.STAT_AD,
							Value:     20,
							IsPercent: true,
						},
						Caster: owner.GetUUID(),
						Target: owner.GetUUID(),
						Source: types.SOURCE_ND,
					},
				})

				fight.HandleAction(types.Action{
					Event:  types.ACTION_EFFECT,
					Target: owner.GetUUID(),
					Source: owner.GetUUID(),
					Meta: types.ActionEffect{
						Effect:   types.EFFECT_STAT_INC,
						Value:    0,
						Duration: 8,
						Uuid:     uuid.New(),
						Meta: types.ActionEffectStat{
							Stat:      types.STAT_AP,
							Value:     20,
							IsPercent: true,
						},
						Caster: owner.GetUUID(),
						Target: owner.GetUUID(),
						Source: types.SOURCE_ND,
					},
				})
			},
		},
	})

	return nil
}
