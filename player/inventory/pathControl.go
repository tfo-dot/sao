package inventory

import (
	"sao/battle"
	"sao/battle/mobs"
	"sao/types"
	"sao/utils"

	"github.com/google/uuid"
)

type ControlSkill struct{}

func (skill ControlSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	return nil
}

func (skill ControlSkill) GetPath() types.SkillPath {
	return types.PathControl
}

func (skill ControlSkill) GetUUID() uuid.UUID {
	return uuid.Nil
}

func (skill ControlSkill) IsLevelSkill() bool {
	return true
}

func (skill ControlSkill) CanUse(owner interface{}, fightInstance interface{}, upgrades int) bool {
	return true
}

type CON_LVL_1 struct {
	ControlSkill
	DefaultCost
	NoEvents
	NoStats
}

func (skill CON_LVL_1) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE, Target: &types.TargetTrigger{Target: types.TARGET_ENEMY, MaxTargets: 1}}
}

func (skill CON_LVL_1) GetUpgradableTrigger(upgrades int) types.Trigger {
	return skill.GetTrigger()
}

func (skill CON_LVL_1) GetName() string {
	return "Poziom 1 - Kontrola"
}

func (skill CON_LVL_1) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "Cooldown",
			Description: "Zmniejsza czas odnowienia o 1 turę",
		},
		{
			Id:          "Speed",
			Description: "Zwiększa prędkość użytkownika o 10 na 1 turę",
		},
		{
			Id:          "Slow",
			Description: "Spowalnia przeciwnika o 10 po zakończeniu ogłuszenia",
		},
	}
}

func (skill CON_LVL_1) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	if HasUpgrade(upgrades, 2) {
		fightInstance.(*battle.Fight).HandleAction(battle.Action{
			Event:  battle.ACTION_EFFECT,
			Target: target.(battle.Entity).GetUUID(),
			Source: owner.(battle.PlayerEntity).GetUUID(),
			Meta: battle.ActionEffect{
				Effect:   battle.EFFECT_STAT_INC,
				Value:    10,
				Duration: 1,
				Caster:   owner.(battle.PlayerEntity).GetUUID(),
				Meta: battle.ActionEffectStat{
					Stat:      types.STAT_SPD,
					Value:     10,
					IsPercent: false,
				},
			},
		})

	}

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Target: target.(battle.Entity).GetUUID(),
		Source: owner.(battle.PlayerEntity).GetUUID(),
		Meta: battle.ActionEffect{
			Effect:   battle.EFFECT_STUN,
			Value:    1,
			Duration: 1,
			Caster:   owner.(battle.PlayerEntity).GetUUID(),
			OnExpire: func(owner, fightInstance interface{}, meta battle.ActionEffect) {
				if !HasUpgrade(upgrades, 3) {
					return
				}

				fightInstance.(*battle.Fight).HandleAction(battle.Action{
					Event:  battle.ACTION_EFFECT,
					Target: target.(battle.Entity).GetUUID(),
					Source: owner.(battle.PlayerEntity).GetUUID(),
					Meta: battle.ActionEffect{
						Effect:   battle.EFFECT_STAT_DEC,
						Value:    10,
						Duration: 1,
						Caster:   owner.(battle.PlayerEntity).GetUUID(),
						Meta: battle.ActionEffectStat{
							Stat:      types.STAT_SPD,
							Value:     10,
							IsPercent: false,
						},
					},
				})
			},
		},
	})

	return nil
}

func (skill CON_LVL_1) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill CON_LVL_1) GetCooldown(upgrades int) int {
	baseCD := skill.GetCD()

	if HasUpgrade(upgrades, 1) {
		return baseCD - 1
	}

	return baseCD
}

func (skill CON_LVL_1) GetDescription() string {
	return "Ogłusza przeciwnika na jedną turę"
}

func (skill CON_LVL_1) GetUpgradableDescription(upgrades int) string {
	upgradeSegment := []string{"", ""}

	if HasUpgrade(upgrades, 1) {
		upgradeSegment[0] = " Zmniejsza czas odnowienia o 1 turę"
	}

	if HasUpgrade(upgrades, 2) {
		upgradeSegment[1] = "\nPo zakończeniu ogłuszenia spowalnia cel o 10 na turę"
	}

	return "Ogłusza przeciwnika na jedną turę." + upgradeSegment[0] + upgradeSegment[1]
}

func (skill CON_LVL_1) GetLevel() int {
	return 1
}

var CON_LVL_2_UUID = uuid.MustParse("00000000-0001-0000-0000-000000000001")

type CON_LVL_2 struct {
	ControlSkill
	NoEvents
	NoTrigger
	NoCost
	NoCooldown
}

func (skill CON_LVL_2) GetName() string {
	return "Poziom 2 - Kontrola"
}

func (skill CON_LVL_2) GetUpgradableDescription(upgrades int) string {
	upgradeSegment := []string{"5 SPD i AGL.", "", "", ""}

	if HasUpgrade(upgrades, 2) {
		upgradeSegment[0] = "10 SPD i AGL."
	}

	if HasUpgrade(upgrades, 1) {
		upgradeSegment[1] = "\nPo trafieniu przeciwnika spowolnia go o 10."
	}

	if HasUpgrade(upgrades, 3) {
		upgradeSegment[2] = "\nPo trafieniu przeciwnika zmniejsza jego AGL o 10."
	}

	if HasUpgrade(upgrades, 1) && HasUpgrade(upgrades, 3) {
		upgradeSegment[1] = ""
		upgradeSegment[2] = ""
		upgradeSegment[3] = "\nPo trafieniu przeciwnika zmniejsza jego SPD i AGL o 10."
	}

	return "Dostajesz " + upgradeSegment[0] + upgradeSegment[1] + upgradeSegment[2] + upgradeSegment[3]
}

func (skill CON_LVL_2) GetDescription() string {
	return "Dostajesz 5 SPD i AGL"
}

func (skill CON_LVL_2) GetLevel() int {
	return 2
}

func (skill CON_LVL_2) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "OnHit",
			Description: "Po trafieniu spowalnia przeciwnika o 10 SPD",
		},
		{
			Id:          "Increase",
			Description: "Zwiększa wartości dwukrotnie",
		},
		{
			Id:          "OnHit",
			Description: "Po trafieniu zmniejsza statystyki przeciwnika o 10 DGD",
		},
	}
}

func (skill CON_LVL_2) GetStats(upgrades int) map[types.Stat]int {
	stats := map[types.Stat]int{
		types.STAT_SPD: 5,
		types.STAT_AGL: 5,
	}

	if HasUpgrade(upgrades, 2) {
		stats[types.STAT_SPD] = 10
		stats[types.STAT_AGL] = 10
	}

	return stats
}

func (skill CON_LVL_2) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	if HasUpgrade(upgrades, 1) {
		fightInstance.(*battle.Fight).HandleAction(battle.Action{
			Event:  battle.ACTION_EFFECT,
			Target: target.(battle.Entity).GetUUID(),
			Source: owner.(battle.PlayerEntity).GetUUID(),
			Meta: battle.ActionEffect{
				Effect:   battle.EFFECT_STAT_DEC,
				Value:    10,
				Duration: 1,
				Caster:   owner.(battle.PlayerEntity).GetUUID(),
				Meta: battle.ActionEffectStat{
					Stat:      types.STAT_SPD,
					Value:     10,
					IsPercent: false,
				},
			},
		})
	}

	if HasUpgrade(upgrades, 3) {
		fightInstance.(*battle.Fight).HandleAction(battle.Action{
			Event:  battle.ACTION_EFFECT,
			Target: target.(battle.Entity).GetUUID(),
			Source: owner.(battle.PlayerEntity).GetUUID(),
			Meta: battle.ActionEffect{
				Effect:   battle.EFFECT_STAT_DEC,
				Value:    10,
				Duration: 1,
				Caster:   owner.(battle.PlayerEntity).GetUUID(),
				Meta: battle.ActionEffectStat{
					Stat:      types.STAT_AGL,
					Value:     10,
					IsPercent: false,
				},
			},
		})
	}

	return nil

}

type CON_LVL_3 struct {
	ControlSkill
	NoEvents
	NoStats
}

func (skill CON_LVL_3) GetName() string {
	return "Poziom 3 - Kontrola"
}

func (skill CON_LVL_3) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	target.(battle.PlayerEntity).Cleanse()

	return nil
}

func (skill CON_LVL_3) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill CON_LVL_3) GetCooldown(upgrades int) int {
	baseCD := skill.GetCD()

	if HasUpgrade(upgrades, 2) {
		return baseCD - 1
	}

	return baseCD
}

func (skill CON_LVL_3) GetDescription() string {
	return "Usuwa wszystkie negatywne efekty"
}

func (skill CON_LVL_3) GetLevel() int {
	return 3
}

func (skill CON_LVL_3) GetCost() int {
	return 1
}

func (skill CON_LVL_3) GetUpgradableCost(upgrades int) int {
	base := skill.GetCost()

	if HasUpgrade(upgrades, 3) {
		return 0
	}

	return base
}

func (skill CON_LVL_3) GetUpgradableDescription(upgrades int) string {
	upgradeSegment := ""

	if HasUpgrade(upgrades, 1) {
		upgradeSegment = " Można użyć na sojuszniku"
	}

	return "Usuwa wszystkie negatywne efekty." + upgradeSegment
}

func (skill CON_LVL_3) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "Ally",
			Description: "Może być użyte na sojuszniku",
		},
		{
			Id:          "Cooldown",
			Description: "Zmniejsza czas odnowienia o 1 turę",
		},
		{
			Id:          "Cost",
			Description: "Nie kosztuje many",
		},
	}
}

func (skill CON_LVL_3) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_TYPE_NONE}
}

func (skill CON_LVL_3) GetUpgradableTrigger(upgrades int) types.Trigger {
	baseTarget := types.TARGET_SELF

	if HasUpgrade(upgrades, 1) {
		baseTarget |= types.TARGET_ALLY
	}

	return types.Trigger{
		Type:  types.TRIGGER_ACTIVE,
		Flags: types.FLAG_IGNORE_CC,
		Target: &types.TargetTrigger{
			Target:     baseTarget,
			MaxTargets: 1,
		},
	}
}

type CON_LVL_4 struct {
	ControlSkill
	DefaultActiveTrigger
	DefaultCost
	NoEvents
	NoStats
}

func (skill CON_LVL_4) GetName() string {
	return "Poziom 4 - Kontrola"
}

type CON_LVL_4_EFFECT struct {
	NoCooldown
	NoCost
	NoLevel
	NoEvents
	Ripple  bool
	CanMiss bool
}

func (skill CON_LVL_4_EFFECT) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	owner.(battle.PlayerEntity).ReduceCooldowns(types.TRIGGER_TURN)

	if skill.Ripple {
		dmgValue := 0

		for _, damage := range meta.(types.AttackTriggerMeta).Effects {
			dmgValue += damage.Value
		}

		randomTarget := utils.RandomElement(fightInstance.(*battle.Fight).GetEnemiesFor(owner.(battle.PlayerEntity).GetUUID()))

		fightInstance.(*battle.Fight).HandleAction(battle.Action{
			Event:  battle.ACTION_DMG,
			Target: randomTarget.GetUUID(),
			Source: owner.(battle.PlayerEntity).GetUUID(),
			Meta: battle.ActionDamage{
				Damage: []battle.Damage{
					{
						Type:  types.DMG_PHYSICAL,
						Value: utils.PercentOf(dmgValue, 20),
					},
				},
			},
		})

	}

	return nil
}

func (skill CON_LVL_4_EFFECT) GetTrigger() types.Trigger {
	return types.Trigger{
		Type:  types.TRIGGER_PASSIVE,
		Event: types.TRIGGER_ATTACK_HIT,
	}
}

func (skill CON_LVL_4_EFFECT) GetName() string {
	return "Kontrola 4 - Efekt"
}

func (skill CON_LVL_4_EFFECT) GetDescription() string {
	return "Po trafieniu ataku zmniejsza CD wszystkich umiejętności o 1"
}

func (skill CON_LVL_4) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	owner.(battle.PlayerEntity).AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value: CON_LVL_4_EFFECT{
			Ripple:  HasUpgrade(upgrades, 1),
			CanMiss: HasUpgrade(upgrades, 2),
		},
		Expire:     1,
		AfterUsage: true,
	})

	return nil
}

func (skill CON_LVL_4) GetUpgradableDescription(upgrades int) string {
	upgradeSegment := []string{" trafionym", ""}

	if HasUpgrade(upgrades, 1) {
		upgradeSegment[1] = "Dodatkowo losowy przeciwnik dostanie obrażenia w wysokości 20% ataku"
	}

	if !HasUpgrade(upgrades, 2) {
		upgradeSegment[0] = ""
	}

	return "Po " + upgradeSegment[0] + "ataku zmniejsza CD wszystkich umiejętności o 1." + upgradeSegment[1]
}

func (skill CON_LVL_4) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill CON_LVL_4) GetCooldown(upgrades int) int {
	return skill.GetCD()
}

func (skill CON_LVL_4) GetDescription() string {
	return "Po trafieniu ataku zmniejsza CD wszystkich umiejętności o 1"
}

func (skill CON_LVL_4) GetLevel() int {
	return 4
}

func (skill CON_LVL_4) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "Ripple",
			Description: "Zadaje 20% obrażeń losowemu przeciwnikowi",
		},
		{
			Id:          "CanMiss",
			Description: "Może nie trafić",
		},
	}
}

type CON_LVL_5 struct {
	ControlSkill
	DefaultActiveTrigger
	DefaultCost
	NoEvents
	NoStats
}

var CON_LVL_5_ENTITY_TYPE = uuid.MustParse("00000000-ffff-ffff-0000-000000000002")

func (skill CON_LVL_5) GetName() string {
	return "Poziom 5 - Kontrola"
}

func (skill CON_LVL_5) GetDescription() string {
	//DEF, MR, SPD, AD, AP, HP
	return "Tworzy klona który ma 25% statystyk użytkownika, może tylko atakować i bronić"
}

func (skill CON_LVL_5) GetUpgradableDescription(upgrades int) string {
	upgradeSegment := []string{" 25% ", "", ""}

	if HasUpgrade(upgrades, 1) {
		upgradeSegment[1] = " Jest 20% szans że użyje umiejętności zamiast ataku."
	}

	if HasUpgrade(upgrades, 2) {
		upgradeSegment[0] = " 50% "
	}

	if HasUpgrade(upgrades, 3) {
		upgradeSegment[2] = "Możesz mieć 2 klony, a po przyzwaniu klon prowokuje wszystkich przeciwników"
	} else {
		upgradeSegment[2] = "Możesz mieć 1 klona."
	}

	return "Tworzy klona który ma" + upgradeSegment[0] + "statystyk użytkownika." + upgradeSegment[0] + upgradeSegment[1] + upgradeSegment[2]
}

func (skill CON_LVL_5) GetLevel() int {
	return 5
}

func (skill CON_LVL_5) CanUseSkill(owner interface{}, fightInstance interface{}, upgrades int) bool {
	maxNumber := 1

	if HasUpgrade(upgrades, 3) {
		maxNumber = 2
	}

	return fightInstance.(*battle.Fight).CanSummon(CON_LVL_5_ENTITY_TYPE, maxNumber)
}

func (skill CON_LVL_5) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	percentValue := 25

	if HasUpgrade(upgrades, 2) {
		percentValue = 50
	}

	var customAction func(self *mobs.SummonEntity, f *battle.Fight) []battle.Action = nil

	if HasUpgrade(upgrades, 1) {
		customAction = func(self *mobs.SummonEntity, f *battle.Fight) []battle.Action {
			actions := make([]battle.Action, 0)

			if utils.RandomNumber(1, 100) <= 20 {
				self.AppendTempSkill(types.WithExpire[types.PlayerSkill]{
					Value:      utils.RandomElement(owner.(battle.PlayerEntity).GetSkills()),
					Expire:     1,
					AfterUsage: true,
					Either:     true,
				})
			}

			return actions
		}
	}

	var onSummon func(f *battle.Fight, summonEntity *mobs.SummonEntity) = nil

	if HasUpgrade(upgrades, 3) {
		onSummon = func(f *battle.Fight, summonEntity *mobs.SummonEntity) {
			f.HandleAction(battle.Action{
				Event:  battle.ACTION_EFFECT,
				Source: summonEntity.GetUUID(),
				Target: summonEntity.GetUUID(),
				Meta: battle.ActionEffect{
					Effect:   battle.EFFECT_TAUNT,
					Value:    1,
					Duration: 1,
					Caster:   summonEntity.GetUUID(),
				},
			})
		}
	}

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_SUMMON,
		Target: owner.(battle.Entity).GetUUID(),
		Source: owner.(battle.PlayerEntity).GetUUID(),
		Meta: battle.ActionSummon{
			Flags:       battle.SUMMON_FLAG_EXPIRE,
			ExpireTimer: 5,
			EntityType:  CON_LVL_5_ENTITY_TYPE,
			Entity: &mobs.SummonEntity{
				Owner:     owner.(battle.PlayerEntity).GetUUID(),
				UUID:      uuid.New(),
				Name:      "Klon " + owner.(battle.PlayerEntity).GetName(),
				CurrentHP: utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_HP), percentValue),
				Stats: map[types.Stat]int{
					types.STAT_DEF: utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_DEF), percentValue),
					types.STAT_MR:  utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_MR), percentValue),
					types.STAT_SPD: utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_SPD), percentValue) + 40,
					types.STAT_AD:  utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AD), percentValue),
					types.STAT_AP:  utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AP), percentValue),
				},
				TempSkill:    make([]*types.WithExpire[types.PlayerSkill], 0),
				Effects:      make(mobs.EffectList, 0),
				CustomAction: customAction,
				OnSummon:     onSummon,
			},
		},
	})

	return nil
}

func (skill CON_LVL_5) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill CON_LVL_5) GetCooldown(upgrades int) int {
	// return skill.GetCD()
	return 0
}

func (skill CON_LVL_5) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "Actions",
			Description: "20% szans że użyje umiejętności",
		},
		{
			Id:          "Stats",
			Description: "Klon ma 50% statystyk",
		},
		{
			Id:          "MaxCount",
			Description: "Możesz mieć 2 klony. Po przyzwaniu klon prowokuje wszystkich przeciwników",
		},
	}
}
