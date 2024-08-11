package inventory

import (
	"fmt"
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
	return types.PathEndurance
}

func (skill EnduranceSkill) GetUUID() uuid.UUID {
	return uuid.Nil
}

func (skill EnduranceSkill) IsLevelSkill() bool {
	return true
}

func (skill EnduranceSkill) CanUse(owner interface{}, fightInstance interface{}, upgrades int) bool {
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

func (skill END_LVL_1) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "Cooldown",
			Description: "Zmniejsza czas odnowienia o 1 turę",
		},
		{
			Id:          "Damage",
			Description: "Zwiększa wartość tarczy do 20 + 3%HP + 2%ATK",
		},
		{
			Id:          "Duration",
			Description: "Zwiększa czas trwania o 1 turę",
		},
	}
}

func (skill END_LVL_1) GetUpgradableDescription(upgrades int) string {
	baseIncrease := "10"
	baseDuration := "1"

	if HasUpgrade(upgrades, 1) {
		baseIncrease = "20 + 3%HP + 2%ATK"
	}

	if HasUpgrade(upgrades, 2) {
		baseDuration = "2"
	}

	return fmt.Sprintf("Daje tarczę o wartości %s na %s tur", baseIncrease, baseDuration)
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

func (skill END_LVL_2) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id: "Increase",
			Events: &map[types.CustomTrigger]func(owner interface{}){
				types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
					owner.(battle.PlayerEntity).SetLevelStat(types.STAT_HP, 20)
				},
			},
			Description: "Zwiększa wartość do 20",
		},
		{
			Id: "DEFIncrease",
			Events: &map[types.CustomTrigger]func(owner interface{}){
				types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
					owner.(battle.PlayerEntity).AppendDerivedStat(types.DerivedStat{
						Base:    types.STAT_HP_PLUS,
						Derived: types.STAT_DEF,
						Percent: 1,
					})
				},
			},
			Description: "Zwiększa wartość pancerza o 1% dodatkowego HP",
		},
		{
			Id: "RESIncrease",
			Events: &map[types.CustomTrigger]func(owner interface{}){
				types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
					owner.(battle.PlayerEntity).AppendDerivedStat(types.DerivedStat{
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

func (skill END_LVL_2) GetUpgradableDescription(upgrades int) string {
	baseIncrease := "15"

	if HasUpgrade(upgrades, 1) {
		baseIncrease = "20"
	}

	defUpgrade := ""

	if HasUpgrade(upgrades, 2) {
		defUpgrade = "\nOtrzymujesz pancerz równy 1% dodatkowego HP."
	}

	mrUpgrade := ""

	if HasUpgrade(upgrades, 3) {
		mrUpgrade = "\nOtrzymujesz odporność na magię równą 1% dodatkowego HP."
	}

	return fmt.Sprintf("Zwiększa zdrowie zyskiwane co poziom do %s.%s%s", baseIncrease, defUpgrade, mrUpgrade)
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
		Type:  types.TRIGGER_PASSIVE,
		Event: types.TRIGGER_ATTACK_GOT_HIT,
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

func (skill END_LVL_3) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "Ally",
			Description: "Sojusznik może rozbić",
		},
		{
			Id:          "Cooldown",
			Description: "Zmniejsza czas odnowienia o 1 turę",
		},
		{
			Id:          "Increase",
			Description: "Zwiększa leczenie o X",
		},
	}
}

func (skill END_LVL_3) GetUpgradableDescription(upgrades int) string {
	baseHeal := "5"

	if HasUpgrade(upgrades, 3) {
		baseHeal = "10"
	}

	allyUpgrade := ""

	if HasUpgrade(upgrades, 1) {
		allyUpgrade = " (także sojusznika)"
	}

	return fmt.Sprintf("Oznacza cel na turę, kolejny atak%s rozbije oznaczenie i wyleczy o %s%%HP", allyUpgrade, baseHeal)
}

type END_LVL_4 struct {
	EnduranceSkill
	DefaultCost
	DefaultActiveTrigger
	NoEvents
	NoStats
}

type END_LVL_4_EFFECT struct {
	NoCost
	NoCooldown
	NoEvents
	NoLevel
	Reduce int
	Event  types.SkillTrigger
}

func (skill END_LVL_4_EFFECT) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	if skill.Event == types.TRIGGER_ATTACK_GOT_HIT {
		tempMeta := meta.(types.AttackTriggerMeta)

		for idx, dmg := range tempMeta.Effects {
			tempMeta.Effects[idx].Value = -utils.PercentOf(dmg.Value, skill.Reduce)
		}

		return tempMeta
	}

	tempMeta := meta.(types.DamageTriggerMeta)

	for idx, dmg := range tempMeta.Effects {
		tempMeta.Effects[idx].Value = -utils.PercentOf(dmg.Value, skill.Reduce)
	}

	return tempMeta
}

func (skill END_LVL_4_EFFECT) GetDescription() string {
	return "Zmniejsza obrażenia o x%"
}

func (skill END_LVL_4_EFFECT) GetName() string {
	return "Efekt wytrzymałości - poziom 4"
}

func (skill END_LVL_4_EFFECT) GetTrigger() types.Trigger {
	return types.Trigger{
		Type:  types.TRIGGER_PASSIVE,
		Event: skill.Event,
	}
}

func (skill END_LVL_4) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	baseReduce := 20

	if HasUpgrade(upgrades, 1) {
		baseReduce += 5
	}

	baseEvent := types.TRIGGER_ATTACK_GOT_HIT

	if HasUpgrade(upgrades, 3) {
		baseEvent = types.TRIGGER_DAMAGE_BEFORE
	}

	owner.(battle.PlayerEntity).AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value: END_LVL_4_EFFECT{
			Reduce: baseReduce,
			Event:  baseEvent,
		},
		AfterUsage: true,
		Expire:     1,
		Either:     HasUpgrade(upgrades, 3),
	})

	return nil
}

func (skill END_LVL_4) GetName() string {
	return "Poziom 4 - wytrzymałość"
}

func (skill END_LVL_4) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill END_LVL_4) GetCooldown(upgrades int) int {
	return skill.GetCD()
}

func (skill END_LVL_4) GetDescription() string {
	return "Zmniejsza kolejny otrzymany atak o 20%"
}

func (skill END_LVL_4) GetLevel() int {
	return 4
}

func (skill END_LVL_4) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "Reduce",
			Description: "Zwiększa redukcję do 25%",
		},
		{
			Id:          "Duration",
			Description: "Działa przez całą ture",
		},
		{
			Id:          "Type",
			Description: "Działa na obrażenia",
		},
	}
}

func (skill END_LVL_4) GetUpgradableDescription(upgrades int) string {
	baseReduce := "20%"

	if HasUpgrade(upgrades, 1) {
		baseReduce = "25%"
	}

	baseDuration := "kolejny"

	if HasUpgrade(upgrades, 2) {
		baseDuration = "wszystkie"
	}

	baseEvent := "otrzymany atak"

	if HasUpgrade(upgrades, 3) {
		if HasUpgrade(upgrades, 2) {
			baseDuration = "kolejne"
		}
		baseEvent = "otrzymane obrażenia"
	}

	return fmt.Sprintf("Przez jedną turę zmniejsza %s %s o %s", baseDuration, baseEvent, baseReduce)
}

type END_LVL_5 struct {
	EnduranceSkill
	DefaultCost
	NoEvents
	NoStats
	DefaultActiveTrigger
}

func (skill END_LVL_5) GetName() string {
	return "Poziom 5 - wytrzymałość"
}

func (skill END_LVL_5) GetDescription() string {
	return "Dostaje 25 tarczy za każdego przeciwnika prowokując ich wszystkich"
}

func (skill END_LVL_5) GetLevel() int {
	return 5
}

func (skill END_LVL_5) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill END_LVL_5) GetCooldown(upgrades int) int {
	return skill.GetCD()
}

func (skill END_LVL_5) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "Duration",
			Description: "Działa turę dłużej",
		},
		{
			Id:          "Increase",
			Description: "Zwiększa wartość tarczy o 10% AP",
		},
		{
			Id:          "Heal",
			Description: "Po zakończeniu leczy o 50% pozostałej tarczy",
		},
	}
}

func (skill END_LVL_5) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	enemiesCount := len(fightInstance.(*battle.Fight).GetEnemiesFor(owner.(battle.PlayerEntity).GetUUID()))

	baseShield := 25

	if HasUpgrade(upgrades, 2) {
		baseShield += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AP), 10)
	}

	duration := 1

	if HasUpgrade(upgrades, 1) {
		duration++
	}

	shieldUuid := uuid.New()

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Target: owner.(battle.Entity).GetUUID(),
		Source: owner.(battle.PlayerEntity).GetUUID(),
		Meta: battle.ActionEffect{
			Effect:   battle.EFFECT_SHIELD,
			Value:    enemiesCount * baseShield,
			Duration: duration,
			Caster:   owner.(battle.PlayerEntity).GetUUID(),
			Uuid:     shieldUuid,
			OnExpire: func(owner, fight interface{}, meta battle.ActionEffect) {
				if HasUpgrade(upgrades, 3) {
					if meta.Value > 0 {
						healValue := utils.PercentOf(meta.Value, 50)

						fight.(*battle.Fight).HandleAction(battle.Action{
							Event:  battle.ACTION_EFFECT,
							Target: owner.(battle.PlayerEntity).GetUUID(),
							Source: owner.(battle.PlayerEntity).GetUUID(),
							Meta: battle.ActionEffect{
								Effect:   battle.EFFECT_HEAL_SELF,
								Value:    healValue,
								Duration: 0,
							},
						})
					}
				}
			},
		},
	})

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Target: owner.(battle.Entity).GetUUID(),
		Source: owner.(battle.PlayerEntity).GetUUID(),
		Meta: battle.ActionEffect{
			Effect:   battle.EFFECT_TAUNT,
			Value:    0,
			Duration: duration,
			Caster:   owner.(battle.PlayerEntity).GetUUID(),
		},
	})

	return nil
}

func (skill END_LVL_5) GetUpgradableDescription(upgrades int) string {
	baseShield := "25"

	if HasUpgrade(upgrades, 2) {
		baseShield = "25 + 10% AP"
	}

	duration := "1"

	if HasUpgrade(upgrades, 1) {
		duration = "2"
	}

	heal := ""

	if HasUpgrade(upgrades, 3) {
		heal = "\nPo zakończeniu leczy o 50% pozostałej tarczy"
	}

	return fmt.Sprintf("Dostaje %s tarczy za każdego przeciwnika prowokując ich wszystkich na %s tur.%s", baseShield, duration, heal)
}
