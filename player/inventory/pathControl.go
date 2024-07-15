package inventory

import (
	"fmt"
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

type CON_LVL_1 struct {
	ControlSkill
	DefaultCost
	DefaultActiveTrigger
	NoEvents
	NoStats
}

func (skill CON_LVL_1) GetName() string {
	return "Poziom 1 - Kontrola"
}

func (skill CON_LVL_1) GetUpgrades() []PlayerSkillUpgrade {
	return []PlayerSkillUpgrade{
		{
			Name:        "Ulepszenie 1",
			Id:          "Cooldown",
			Events:      nil,
			Description: "Zmniejsza czas odnowienia o 1 turę",
		},
		{
			Name:        "Ulepszenie 2",
			Id:          "Speed",
			Events:      nil,
			Description: "Zwiększa prędkość użytkownika o 10 na 1 turę",
		},
		{
			Name:        "Ulepszenie 3",
			Id:          "Slow",
			Events:      nil,
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

	if HasUpgrade(upgrades, 3) {
		//TODO add slow effect after stun ends
		fmt.Println("SLOW")
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

func (skill CON_LVL_1) GetLevel() int {
	return 1
}

var CON_LVL_2_UUID = uuid.MustParse("00000000-0001-0000-0000-000000000001")

type CON_LVL_2 struct {
	ControlSkill
	NoExecute
	NoEvents
	NoTrigger
}

func (skill CON_LVL_2) GetName() string {
	return "Poziom 2 - Kontrola"
}

func (skill CON_LVL_2) GetDescription() string {
	return "Dostajesz 5 SPD i AGL"
}

func (skill CON_LVL_2) GetLevel() int {
	return 2
}

type CON_LVL_EFFECT_1 struct {
	NoLevel
	NoCooldown
	NoEvents
	NoCost
}

func (skill CON_LVL_EFFECT_1) Execute(owner, target, fightInstance, meta interface{}) interface{} {
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

	return nil
}

func (skill CON_LVL_EFFECT_1) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_ATTACK_HIT,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
			Meta:          nil,
		},
	}
}

func (skill CON_LVL_EFFECT_1) GetName() string {
	return "Kontrola 2 - Ulepszenie 1"
}

func (skill CON_LVL_EFFECT_1) GetDescription() string {
	return "Po trafieniu spowalnia przeciwnika o 10 SPD"
}

type CON_LVL_EFFECT_3 struct {
	NoLevel
	NoCooldown
	NoCost
	NoEvents
}

func (skill CON_LVL_EFFECT_3) Execute(owner, target, fightInstance, meta interface{}) interface{} {
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

	return nil
}

func (skill CON_LVL_EFFECT_3) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_ATTACK_HIT,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
			Meta:          nil,
		},
	}
}

func (skill CON_LVL_EFFECT_3) GetName() string {
	return "Kontrola 2 - Ulepszenie 3"
}

func (skill CON_LVL_EFFECT_3) GetDescription() string {
	return "Po trafieniu spowalnia przeciwnika o 10 SPD"
}

func (skill CON_LVL_2) GetUpgrades() []PlayerSkillUpgrade {
	itemEvents1 := map[types.CustomTrigger]func(owner battle.PlayerEntity){
		types.CUSTOM_TRIGGER_UNLOCK: func(owner battle.PlayerEntity) {
			owner.AddItem(&types.PlayerItem{
				UUID:        CON_LVL_2_UUID,
				Name:        "Kontrola 2 - Ulepszenie 1",
				Description: "=-=-=",
				TakesSlot:   false,
				Stacks:      false,
				Consume:     false,
				Count:       1,
				MaxCount:    1,
				Hidden:      true,
				Stats:       nil,
				Effects:     []types.PlayerSkill{CON_LVL_EFFECT_1{}},
			})
		},
	}

	itemEvents2 := map[types.CustomTrigger]func(owner battle.PlayerEntity){
		types.CUSTOM_TRIGGER_UNLOCK: func(owner battle.PlayerEntity) {
			owner.AddItem(&types.PlayerItem{
				UUID:        CON_LVL_2_UUID,
				Name:        "Kontrola 2 - Ulepszenie 2",
				Description: "=-=-=",
				TakesSlot:   false,
				Stacks:      false,
				Consume:     false,
				Count:       1,
				MaxCount:    1,
				Hidden:      true,
				Stats: map[types.Stat]int{
					types.STAT_SPD: 5,
					types.STAT_AGL: 5,
				},
			})
		},
	}

	itemEvents3 := map[types.CustomTrigger]func(owner battle.PlayerEntity){
		types.CUSTOM_TRIGGER_UNLOCK: func(owner battle.PlayerEntity) {
			owner.AddItem(&types.PlayerItem{
				UUID:        CON_LVL_2_UUID,
				Name:        "Kontrola 2 - Ulepszenie 3",
				Description: "=-=-=",
				TakesSlot:   false,
				Stacks:      false,
				Consume:     false,
				Count:       1,
				MaxCount:    1,
				Hidden:      true,
				Stats:       nil,
				Effects:     []types.PlayerSkill{CON_LVL_EFFECT_3{}},
			})
		},
	}

	return []PlayerSkillUpgrade{
		{
			Name:        "Ulepszenie 1",
			Id:          "OnHit",
			Events:      &itemEvents1,
			Description: "Po trafieniu spowalnia przeciwnika o 10 SPD",
		},
		{
			Name:        "Ulepszenie 2",
			Id:          "Increase",
			Events:      &itemEvents2,
			Description: "Zwiększa wartości dwukrotnie",
		},
		{
			Name:        "Ulepszenie 3",
			Id:          "OnHit",
			Events:      &itemEvents3,
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

func (skill CON_LVL_3) GetUpgrades() []PlayerSkillUpgrade {
	return []PlayerSkillUpgrade{
		{
			Name:        "Ulepszenie 1",
			Id:          "Ally",
			Events:      nil,
			Description: "Może być użyte na sojuszniku",
		},
		{
			Name:        "Ulepszenie 2",
			Id:          "Cooldown",
			Events:      nil,
			Description: "Zmniejsza czas odnowienia o 1 turę",
		},
		{
			Name:        "Ulepszenie 3",
			Id:          "Cost",
			Events:      nil,
			Description: "Nie kosztuje many",
		},
	}
}

// TODO add ally target
func (skill CON_LVL_3) GetTrigger() types.Trigger {
	return types.Trigger{
		Type:  types.TRIGGER_ACTIVE,
		Event: nil,
		//Can use while stunned
		Flags: 1 << 1,
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
						Type:  battle.DMG_PHYSICAL,
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
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_ATTACK_HIT,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
			Meta:          nil,
			OptionalEvent: types.TRIGGER_ATTACK_MISS,
		},
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

func (skill CON_LVL_4) GetUpgrades() []PlayerSkillUpgrade {
	return []PlayerSkillUpgrade{
		{
			Name:        "Ulepszenie 1",
			Id:          "Ripple",
			Events:      nil,
			Description: "Zadaje 20% obrażeń losowemu przeciwnikowi",
		},
		{
			Name:        "Ulepszenie 2",
			Id:          "CanMiss",
			Events:      nil,
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

func (skill CON_LVL_5) GetName() string {
	return "Poziom 5 - Kontrola"
}

func (skill CON_LVL_5) GetDescription() string {
	//DEF, MR, SPD, AD, AP, HP
	return "Tworzy klona który ma 25% statystyk użytkownika, może tylko atakować i bronić"
}

func (skill CON_LVL_5) GetLevel() int {
	return 5
}

func (skill CON_LVL_5) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	//TODO Upgrades
	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_SUMMON,
		Target: owner.(battle.Entity).GetUUID(),
		Source: owner.(battle.PlayerEntity).GetUUID(),
		Meta: battle.ActionSummon{
			Flags:       battle.SUMMON_FLAG_EXPIRE,
			ExpireTimer: 5,
			Entity: &mobs.MobEntity{
				Id:        "Klon",
				MaxHP:     90,
				HP:        90,
				SPD:       40,
				ATK:       25,
				Effects:   make(mobs.EffectList, 0),
				UUID:      uuid.New(),
				Props:     make(map[string]interface{}, 0),
				Loot:      []battle.Loot{{Type: battle.LOOT_EXP, Count: 55}},
				TempSkill: make([]types.WithExpire[types.PlayerSkill], 0),
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

func (skill CON_LVL_5) GetUpgrades() []PlayerSkillUpgrade {
	return []PlayerSkillUpgrade{
		{
			Name:        "Ulepszenie 1",
			Id:          "Actions",
			Events:      nil,
			Description: "20% szans że użyje umiejętności",
		},
		{
			Name:        "Ulepszenie 2",
			Id:          "Stats",
			Events:      nil,
			Description: "Klon ma 50% statystyk",
		},
		{
			Name:        "Ulepszenie 3",
			Id:          "MaxCount",
			Events:      nil,
			Description: "Możesz mieć 2 klony. Po przyzwaniu klon prowokuje wszystkich przeciwników",
		},
	}
}
