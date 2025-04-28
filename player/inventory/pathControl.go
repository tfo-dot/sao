package inventory

import (
	"fmt"
	"sao/battle/mobs"
	"sao/types"
	"sao/utils"
	"slices"

	"github.com/google/uuid"
)

type ControlSkill struct{}

func (skill ControlSkill) Execute(_ types.PlayerEntity, _ types.Entity, _ types.FightInstance, _ any) any {
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

func (skill ControlSkill) CanUse(_ types.PlayerEntity, _ types.FightInstance) bool {
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
		{Id: "Cooldown", Description: "Zmniejsza czas odnowienia o 1 turę"},
		{Id: "Speed", Description: "Zwiększa SPD użytkownika o 10 na 1 turę"},
		{Id: "Slow", Description: "Po zakończeniu ogłuszenia zmniejsza SPD przeciwnika o 10 na 1 turę."},
	}
}

func (skill CON_LVL_1) GetDescription() string {
	return "Ogłusza przeciwnika na jedną turę"
}

func (skill CON_LVL_1) GetUpgradableDescription(upgrades int) string {
	upgradeSegment := []string{"", ""}

	if HasUpgrade(upgrades, 1) {
		upgradeSegment[0] = "\nPo użyciu daje 10 SPD na 1 turę."
	}

	if HasUpgrade(upgrades, 2) {
		upgradeSegment[1] = "\nPo zakończeniu ogłuszenia zmniejsza SPD celu o 10 przez 1 turę"
	}

	return "Ogłusza przeciwnika na jedną turę." + upgradeSegment[0] + upgradeSegment[1]
}

func (skill CON_LVL_1) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 2) {
		fightInstance.HandleAction(types.Action{
			Event:  types.ACTION_EFFECT,
			Target: target.GetUUID(),
			Source: owner.GetUUID(),
			Meta: types.ActionEffect{
				Effect:   types.EFFECT_STAT_INC,
				Value:    10,
				Duration: 1,
				Meta:     types.ActionEffectStat{Stat: types.STAT_SPD},
			},
		})
	}

	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_EFFECT,
		Target: target.GetUUID(),
		Source: owner.GetUUID(),
		Meta: types.ActionEffect{
			Effect:   types.EFFECT_STUN,
			Duration: 1,
			OnExpire: func(owner types.Entity, fightInstance types.FightInstance, _ types.ActionEffect) {
				if !HasUpgrade(owner.(types.PlayerEntity).GetUpgrades(skill.GetLevel()), 3) {
					return
				}

				fightInstance.HandleAction(types.Action{
					Event:  types.ACTION_EFFECT,
					Target: target.GetUUID(),
					Source: owner.GetUUID(),
					Meta: types.ActionEffect{
						Effect:   types.EFFECT_STAT_DEC,
						Value:    10,
						Duration: 1,
						Meta:     types.ActionEffectStat{Stat: types.STAT_SPD},
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
	upgradeSegment := []string{"", "", ""}

	statValue := 5

	if HasUpgrade(upgrades, 2) {
		statValue = 10
	}

	if HasUpgrade(upgrades, 1) {
		upgradeSegment[0] = "\nPo trafieniu przeciwnika zmniejsza jego SPD o 10 na 1 turę."
	}

	if HasUpgrade(upgrades, 3) {
		upgradeSegment[1] = "\nPo trafieniu przeciwnika zmniejsza jego AGL o 10 na 1 turę."
	}

	if HasUpgrade(upgrades, 1) && HasUpgrade(upgrades, 3) {
		upgradeSegment[0] = ""
		upgradeSegment[1] = ""
		upgradeSegment[2] = "\nPo trafieniu przeciwnika zmniejsza jego SPD i AGL o 10 na 1 turę."
	}

	return fmt.Sprintf("Otrzymujesz %d SPD i AGL. %s%s%s", statValue, upgradeSegment[0], upgradeSegment[1], upgradeSegment[2])
}

func (skill CON_LVL_2) GetDescription() string {
	return "Otrzymujesz 5 SPD i AGL"
}

func (skill CON_LVL_2) GetLevel() int {
	return 2
}

func (skill CON_LVL_2) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{Id: "OnHit", Description: "Po trafieniu ataku zmniejsza SPD przeciwnika o 10 na 1 turę"},
		{Id: "Increase", Description: "Zwiększa wartości pasywne do 10 SPD i AGL"},
		{Id: "OnHit", Description: "Po trafieniu ataku zmniejsza AGL przeciwnika o 10 na 1 turę"},
	}
}

func (skill CON_LVL_2) GetStats(upgrades int) map[types.Stat]int {
	if HasUpgrade(upgrades, 2) {
		return map[types.Stat]int{types.STAT_SPD: 10, types.STAT_AGL: 10}
	}

	return map[types.Stat]int{types.STAT_SPD: 5, types.STAT_AGL: 5}
}

func (skill CON_LVL_2) GetDerivedStats(upgrades int) []types.DerivedStat {
	return []types.DerivedStat{}
}

func (skill CON_LVL_2) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 1) {
		fightInstance.HandleAction(types.Action{
			Event:  types.ACTION_EFFECT,
			Target: target.GetUUID(),
			Source: owner.GetUUID(),
			Meta: types.ActionEffect{
				Effect:   types.EFFECT_STAT_DEC,
				Value:    10,
				Duration: 1,
				Meta:     types.ActionEffectStat{Stat: types.STAT_SPD, IsPercent: false},
			},
		})
	}

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 3) {
		fightInstance.HandleAction(types.Action{
			Event:  types.ACTION_EFFECT,
			Target: target.GetUUID(),
			Source: owner.GetUUID(),
			Meta: types.ActionEffect{
				Effect:   types.EFFECT_STAT_DEC,
				Value:    10,
				Duration: 1,
				Meta:     types.ActionEffectStat{Stat: types.STAT_AGL, IsPercent: false},
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

func (skill CON_LVL_3) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	target.Cleanse()

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
		{Id: "Ally", Description: "Celem umiejętności może być sojusznik"},
		{Id: "Cooldown", Description: "Zmniejsza czas odnowienia o 1 turę"},
		{Id: "Cost", Description: "Umiejętność nie kosztuje many"},
	}
}

func (skill CON_LVL_3) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE}
}

func (skill CON_LVL_3) GetUpgradableTrigger(upgrades int) types.Trigger {
	baseTarget := types.TARGET_SELF

	if HasUpgrade(upgrades, 1) {
		baseTarget |= types.TARGET_ALLY
	}

	return types.Trigger{
		Type:   types.TRIGGER_ACTIVE,
		Flags:  types.FLAG_IGNORE_CC,
		Target: &types.TargetTrigger{Target: baseTarget, MaxTargets: 1},
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

func (skill CON_LVL_4_EFFECT) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta any) any {
	owner.ReduceCooldowns(types.TRIGGER_TURN)

	if skill.Ripple {
		dmgValue := 0

		for _, damage := range meta.(types.AttackTriggerMeta).Effects {
			dmgValue += damage.Value
		}

		fightInstance.HandleAction(types.Action{
			Event:  types.ACTION_DMG,
			Target: utils.RandomElement(fightInstance.GetEnemiesFor(owner.GetUUID())).GetUUID(),
			Source: owner.GetUUID(),
			Meta: types.ActionDamage{
				Damage: []types.Damage{{
					Type:  types.DMG_PHYSICAL,
					Value: utils.PercentOf(dmgValue, 20),
				}},
			},
		})
	}

	return nil
}

func (skill CON_LVL_4_EFFECT) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_PASSIVE, Event: types.TRIGGER_ATTACK_HIT}
}

func (skill CON_LVL_4_EFFECT) GetName() string {
	return "Kontrola 4 - Efekt"
}

func (skill CON_LVL_4_EFFECT) GetDescription() string {
	return "Kolejny trafiony atak zmniejszy CD wszystkich umiejętności o 1"
}

func (skill CON_LVL_4) GetUpgradableDescription(upgrades int) string {
	upgradeSegment := []string{" trafiony", ""}

	if HasUpgrade(upgrades, 1) {
		upgradeSegment[1] = "Dodatkowo losowy przeciwnik dostanie obrażenia w wysokości 20% ataku"
	}

	if !HasUpgrade(upgrades, 2) {
		upgradeSegment[0] = ""
	}

	return "Kolejny" + upgradeSegment[0] + " atak zmniejszy CD wszystkich umiejętności o 1." + upgradeSegment[1]
}

func (skill CON_LVL_4) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	owner.AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value: CON_LVL_4_EFFECT{
			Ripple:  HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 1),
			CanMiss: HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 2),
		},
		Expire:     1,
		AfterUsage: true,
	})

	return nil
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
		{Id: "Ripple", Description: "Przy aktywacji efektu dodatkowo zada 20% obrażeń ataku losowemu przeciwnikowi"},
		{Id: "CanMiss", Description: "Zmienia warunek aktywacji na 'przy kolejnym ataku' (wcześniej 'przy kolejnym trafionym ataku')"},
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
	return "Tworzy klona który ma 25% statystyk (DEF, RES, SPD, ATK, AP, HP) użytkownika i utrzymuje się 5 tur."
}

func (skill CON_LVL_5) GetUpgradableDescription(upgrades int) string {
	ratio := 25

	if HasUpgrade(upgrades, 2) {
		ratio = 50
	}

	cloneCount := 1
	summonEffect := ""

	if HasUpgrade(upgrades, 3) {
		cloneCount = 2
		summonEffect = "Po przyzwaniu klon prowokuje wszystkich przeciwników"
	}

	skillUsage := ""

	if HasUpgrade(upgrades, 1) {
		skillUsage = " Klon ma 20% szans na użycie umiejętności."
	}

	return fmt.Sprintf("Tworzy klona który ma %d%% statystyk (DEF, RES, SPD, ATK, AP, HP) użytkownika i utrzymuje się 5 tur.%s Max ilość klonów %d.  %s", ratio, skillUsage, cloneCount, summonEffect)
}

func (skill CON_LVL_5) GetLevel() int {
	return 5
}

func (skill CON_LVL_5) CanUseSkill(owner types.PlayerEntity, fightInstance types.FightInstance) bool {
	maxNumber := 1

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 3) {
		maxNumber = 2
	}

	return fightInstance.CanSummon(CON_LVL_5_ENTITY_TYPE, maxNumber)
}

func (skill CON_LVL_5) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	percentValue := 25

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 2) {
		percentValue = 50
	}

	var customAction func(self *mobs.SummonEntity, f types.FightInstance) []types.Action = nil

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 1) {
		customAction = func(self *mobs.SummonEntity, f types.FightInstance) []types.Action {
			actions := make([]types.Action, 0)

			if utils.RandomNumber(1, 100) <= 20 {
				self.AppendTempSkill(types.WithExpire[types.PlayerSkill]{
					Value:      utils.RandomElement(owner.GetSkills()),
					Expire:     1,
					AfterUsage: true,
					Either:     true,
				})
			}

			return actions
		}
	}

	var onSummon func(f types.FightInstance, summonEntity *mobs.SummonEntity) = nil

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 3) {
		onSummon = func(f types.FightInstance, summonEntity *mobs.SummonEntity) {
			f.HandleAction(types.Action{
				Event:  types.ACTION_EFFECT,
				Source: summonEntity.GetUUID(),
				Target: summonEntity.GetUUID(),
				Meta:   types.ActionEffect{Effect: types.EFFECT_TAUNT, Value: 1, Duration: 1},
			})
		}
	}

	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_SUMMON,
		Source: owner.GetUUID(),
		Meta: types.ActionSummon{
			Flags:       types.SUMMON_FLAG_EXPIRE,
			ExpireTimer: 5,
			EntityType:  CON_LVL_5_ENTITY_TYPE,
			Entity: &mobs.SummonEntity{
				Owner:     owner.GetUUID(),
				UUID:      uuid.New(),
				Name:      "Klon " + owner.GetName(),
				CurrentHP: utils.PercentOf(owner.GetStat(types.STAT_HP), percentValue),
				Stats: map[types.Stat]int{
					types.STAT_DEF: utils.PercentOf(owner.GetStat(types.STAT_DEF), percentValue),
					types.STAT_MR:  utils.PercentOf(owner.GetStat(types.STAT_MR), percentValue),
					types.STAT_SPD: utils.PercentOf(owner.GetStat(types.STAT_SPD), percentValue) + 40,
					types.STAT_AD:  utils.PercentOf(owner.GetStat(types.STAT_AD), percentValue),
					types.STAT_AP:  utils.PercentOf(owner.GetStat(types.STAT_AP), percentValue),
				},
				TempSkill:    make([]*types.WithExpire[types.PlayerSkill], 0),
				Effects:      make([]types.ActionEffect, 0),
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
	return skill.GetCD()
}

func (skill CON_LVL_5) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{Id: "Actions", Description: "Klon ma 20% szans na użycie umiejętności"},
		{Id: "Stats", Description: "Klon ma 50% statystyk"},
		{Id: "MaxCount", Description: "Limit klonów zwiększa się do 2. Po przyzwaniu klon prowokuje wszystkich przeciwników"},
	}
}

type CON_LVL_6 struct {
	ControlSkill
	DefaultCost
	NoEvents
	NoStats
}

func (skill CON_LVL_6) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE, Target: &types.TargetTrigger{Target: types.TARGET_ENEMY, MaxTargets: 1}}
}

func (skill CON_LVL_6) GetUpgradableTrigger(upgrades int) types.Trigger {
	return skill.GetTrigger()
}

func (skill CON_LVL_6) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill CON_LVL_6) GetCooldown(upgrades int) int {
	return skill.GetCD()
}

func (skill CON_LVL_6) GetName() string {
	return "Poziom 6 - Kontrola"
}

func (skill CON_LVL_6) GetDescription() string {
	return "Ty i wybrany wróg zostajecie ogłuszeni na jedną turę."
}

func (skill CON_LVL_6) GetUpgradableDescription(upgrades int) string {
	damage := ""

	if HasUpgrade(upgrades, 1) {
		damage = ", dodatkowo przeciwnik otrzymuje obrażenia magiczne w wysokości 100% AP."
	}

	shield := ""

	if HasUpgrade(upgrades, 2) {
		shield = "Ty oraz sojusznik o najmniejszym zdrowiu otrzymuje tarczę o wartości 50% AP."
	}

	return fmt.Sprintf("Ty i wybrany wróg zostajecie ogłuszeni na jedną turę%s. %s", damage, shield)
}

func (skill CON_LVL_6) GetLevel() int {
	return 6
}

func (skill CON_LVL_6) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{Id: "Damage", Description: "Dodatkowo zadaje 100% AP obrażeń przeciwnikowi."},
		{Id: "Shield", Description: "Ty oraz sojusznik o najmniejszym zdrowiu otrzymuje tarczę o wartości 50% AP"},
	}
}

func (skill CON_LVL_6) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_EFFECT,
		Target: target.GetUUID(),
		Source: owner.GetUUID(),
		Meta:   types.ActionEffect{Effect: types.EFFECT_STUN, Duration: 1},
	})

	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_EFFECT,
		Source: owner.GetUUID(),
		Meta:   types.ActionEffect{Effect: types.EFFECT_STUN, Duration: 1},
	})

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 1) {
		fightInstance.HandleAction(types.Action{
			Event:  types.ACTION_DMG,
			Target: target.GetUUID(),
			Source: owner.GetUUID(),
			Meta: types.ActionDamage{
				Damage:   []types.Damage{{Type: types.DMG_MAGICAL, Value: owner.GetStat(types.STAT_AP)}},
				CanDodge: false,
			},
		})
	}

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 2) {
		allies := fightInstance.GetAlliesFor(owner.GetUUID())

		if len(allies) > 0 {
			slices.SortFunc(allies, func(entPrev, entNext types.Entity) int {
				return entPrev.GetStat(types.STAT_HP) - entNext.GetStat(types.STAT_HP)
			})

			fightInstance.HandleAction(types.Action{
				Event:  types.ACTION_EFFECT,
				Target: allies[0].GetUUID(),
				Source: owner.GetUUID(),
				Meta: types.ActionEffect{
					Effect: types.EFFECT_SHIELD, Value: utils.PercentOf(owner.GetStat(types.STAT_AP), 50), Duration: 1,
				},
			})
		}

		fightInstance.HandleAction(types.Action{
			Event:  types.ACTION_EFFECT,
			Target: owner.GetUUID(),
			Source: owner.GetUUID(),
			Meta: types.ActionEffect{
				Effect: types.EFFECT_SHIELD, Value: utils.PercentOf(owner.GetStat(types.STAT_AP), 50), Duration: 1,
			},
		})
	}

	return nil
}

type CON_ULT_1 struct {
	ControlSkill
	DefaultCost
	DefaultActiveTrigger
	NoEvents
	NoStats
}

func (skill CON_ULT_1) GetName() string {
	return "Poziom 10 - kontrola"
}

type CON_ULT_1_EFFECT_1 struct {
	NoCooldown
	NoCost
	NoEvents
	NoStats
	NoLevel
}

func (skill CON_ULT_1_EFFECT_1) GetName() string {
	return "Kontrola 10 - Efekt"
}

func (skill CON_ULT_1_EFFECT_1) GetDescription() string {
	return "Co ture zadaje wszystkim wrogom 1% maks HP."
}

func (skill CON_ULT_1_EFFECT_1) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_PASSIVE, Event: types.TRIGGER_TURN}
}

func (skill CON_ULT_1_EFFECT_1) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	for _, enemy := range fightInstance.GetEnemiesFor(owner.GetUUID()) {
		fightInstance.HandleAction(types.Action{
			Event:  types.ACTION_DMG,
			Target: enemy.GetUUID(),
			Source: owner.GetUUID(),
			Meta: types.ActionDamage{
				Damage: []types.Damage{{
					Type:  types.DMG_TRUE,
					Value: utils.PercentOf(enemy.GetStat(types.STAT_HP), 1),
				}},
				CanDodge: false,
			},
		})
	}

	return nil
}

type CON_ULT_1_EFFECT_2 struct {
	NoCooldown
	NoCost
	NoEvents
	NoStats
	NoLevel
}

func (skill CON_ULT_1_EFFECT_2) GetName() string {
	return "Kontrola 10 - Efekt"
}

func (skill CON_ULT_1_EFFECT_2) GetDescription() string {
	return "Zaatakowanie wroga zmniejsza jego SPD o 10"
}

func (skill CON_ULT_1_EFFECT_2) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_PASSIVE, Event: types.TRIGGER_ATTACK_HIT}
}

func (skill CON_ULT_1_EFFECT_2) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_EFFECT,
		Target: target.GetUUID(),
		Source: owner.GetUUID(),
		Meta: types.ActionEffect{
			Effect:   types.EFFECT_STAT_DEC,
			Value:    10,
			Duration: 1,
			Meta:     types.ActionEffectStat{Stat: types.STAT_SPD, IsPercent: false},
		},
	})

	return nil
}

type CON_ULT_1_EFFECT_3 struct {
	NoCooldown
	NoCost
	NoEvents
	NoLevel
}

func (skill CON_ULT_1_EFFECT_3) GetStats() map[types.Stat]int {
	return map[types.Stat]int{types.STAT_SPD: 10}
}

func (skill CON_ULT_1_EFFECT_3) GetName() string {
	return "Kontrola 10 - Efekt"
}

func (skill CON_ULT_1_EFFECT_3) GetDescription() string {
	return "Daje 10 SPD i leczy o 1% maksymalnego HP na turę"
}

func (skill CON_ULT_1_EFFECT_3) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_PASSIVE, Event: types.TRIGGER_TURN}
}

func (skill CON_ULT_1_EFFECT_3) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_EFFECT,
		Target: owner.GetUUID(),
		Meta:   types.ActionEffect{Effect: types.EFFECT_HEAL, Value: utils.PercentOf(owner.GetStat(types.STAT_HP), 1)},
	})

	return nil
}

func (skill CON_ULT_1) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	owner.AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value:  CON_ULT_1_EFFECT_1{},
		Expire: 10,
		OnExpire: func(owner types.PlayerEntity, fight types.FightInstance) {
			for _, enemy := range fight.GetEnemiesFor(owner.GetUUID()) {
				fight.HandleAction(types.Action{
					Event:  types.ACTION_DMG,
					Target: enemy.GetUUID(),
					Source: owner.GetUUID(),
					Meta: types.ActionDamage{
						Damage: []types.Damage{
							{
								Type:  types.DMG_MAGICAL,
								Value: utils.PercentOf(owner.GetStat(types.STAT_AP), 75),
							},
						},
					},
				})

				fight.HandleAction(types.Action{
					Event:  types.ACTION_EFFECT,
					Target: enemy.GetUUID(),
					Source: owner.GetUUID(),
					Meta:   types.ActionEffect{Effect: types.EFFECT_STUN, Duration: 1},
				})
			}
		},
	})

	owner.AppendTempSkill(types.WithExpire[types.PlayerSkill]{Value: CON_ULT_1_EFFECT_2{}, Expire: 10})

	for _, enemy := range fightInstance.GetEnemiesFor(owner.GetUUID()) {
		fightInstance.HandleAction(types.Action{
			Event:  types.ACTION_EFFECT,
			Target: enemy.GetUUID(),
			Source: owner.GetUUID(),
			Meta: types.ActionEffect{
				Effect:   types.EFFECT_RESIST,
				Value:    -20,
				Duration: 1,
				Meta:     types.ActionEffectResist{IsPercent: true, DmgType: 4},
			},
		})
	}

	owner.AppendTempSkill(types.WithExpire[types.PlayerSkill]{Value: CON_ULT_1_EFFECT_3{}, Expire: 10})

	return nil
}

func (skill CON_ULT_1) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{}
}

func (skill CON_ULT_1) GetCD() int {
	return 10
}

func (skill CON_ULT_1) GetCooldown(upgrades int) int {
	return skill.GetCD()
}

func (skill CON_ULT_1) GetLevel() int {
	return 10
}

func (skill CON_ULT_1) GetDescription() string {
	return "Tworzy pulsujące pole energetyczne zadające obrażenia równe 1% maks. HP co rundę wszystkim wrogom. Pole trwa 10 tur.\nPrzeciwnicy zaatakowani podczas trwania umiejętności zostaną spowolnieni o 10 a obrażenia przez nich otrzymywane są zwiększone o 20% przy okazji nie mogą używać umiejętności.\nUżytkownik dostaje 10 SPD i regeneruje 1% maks zdrowia co turę podczas działania umiejętności\nPo zakończeniu zadaje wszystkim wrogom 75% AP i ogłusza ich na jedną turę."
}

func (skill CON_ULT_1) GetUpgradableDescription(upgrades int) string {
	return skill.GetDescription()
}

func (skill CON_ULT_1) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	return skill.Execute(owner, target, fightInstance, meta)
}

type CON_ULT_2 struct {
	ControlSkill
	NoEvents
	NoStats
	DefaultActiveTrigger
}

func (skill CON_ULT_2) GetCD() int {
	return 10
}

func (skill CON_ULT_2) GetCooldown(upgrades int) int {
	return skill.GetCD()
}

func (skill CON_ULT_2) GetCost() int {
	return 1
}

func (skill CON_ULT_2) GetUpgradableCost(upgrades int) int {
	return skill.GetCost()
}

func (skill CON_ULT_2) GetLevel() int {
	return 10
}

func (skill CON_ULT_2) GetName() string {
	return "Poziom 10 - kontrola"
}

func (skill CON_ULT_2) GetDescription() string {
	return "Przyzywa golema który uderza w ziemie zadając 100% AP i ogłuszając wszystkich wrogów na jedną turę.\nGolem nakłada na wszystkich sojuszników tarczę w wysokości 100% AP i trwa póki golem żyje lub zostanie zniszczona.\nGolem prowokuje wszystkich przeciwników. Przy uderzeniu wroga golem leczy się o 20% zadanych obrażeń.\nJego statystyki to:\nHP: HP gracza + 200% AP\nAD: AD gracza + 50% AP\nPancerz i odporność na magię: 200% użytkownika\nPo zakończeniu walki umiera."
}

func (skill CON_ULT_2) GetUpgradableDescription(upgrades int) string {
	return skill.GetDescription()
}

func (skill CON_ULT_2) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{}
}

func (skill CON_ULT_2) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	return skill.Execute(owner, target, fightInstance, meta)
}

var CON_ULT_2_ENTITY_TYPE = uuid.MustParse("00000000-ffff-ffff-0000-000000000003")

func (skill CON_ULT_2) Execute(owner types.PlayerEntity, _ types.Entity, fightInstance types.FightInstance, meta any) any {
	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_SUMMON,
		Source: owner.GetUUID(),
		Meta: types.ActionSummon{
			Flags:       types.SUMMON_FLAG_ATTACK,
			EntityType:  CON_ULT_2_ENTITY_TYPE,
			Entity: &mobs.SummonEntity{
				Owner:     owner.GetUUID(),
				UUID:      uuid.New(),
				Name:      "Golem " + owner.GetName(),
				CurrentHP: owner.GetStat(types.STAT_HP) + utils.PercentOf(owner.GetStat(types.STAT_AP), 200),
				Stats: map[types.Stat]int{
					types.STAT_HP:       owner.GetStat(types.STAT_HP) + utils.PercentOf(owner.GetStat(types.STAT_AP), 200),
					types.STAT_AD:       owner.GetStat(types.STAT_AD) + utils.PercentOf(owner.GetStat(types.STAT_AP), 50),
					types.STAT_DEF:      utils.PercentOf(owner.GetStat(types.STAT_AP), 200),
					types.STAT_MR:       utils.PercentOf(owner.GetStat(types.STAT_AP), 200),
					types.STAT_ATK_VAMP: 20,
				},
				TempSkill: make([]*types.WithExpire[types.PlayerSkill], 0),
				Effects:   make([]types.ActionEffect, 0),
				OnSummon: func(f types.FightInstance, summonEntity *mobs.SummonEntity) {
					f.HandleAction(types.Action{
						Event:  types.ACTION_EFFECT,
						Source: summonEntity.GetUUID(),
						Meta:   types.ActionEffect{Effect: types.EFFECT_TAUNT, Duration: -1},
					})
				},
			},
		},
	})

	return nil
}
