package inventory

import (
	"fmt"
	"sao/battle"
	"sao/types"
	"sao/utils"

	"github.com/disgoorg/disgo/discord"
	"github.com/google/uuid"
)

type SpecialSkill struct{}

func (skill SpecialSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	return nil
}

func (skill SpecialSkill) GetPath() types.SkillPath {
	return types.PathSpecial
}

func (skill SpecialSkill) GetUUID() uuid.UUID {
	return uuid.Nil
}

func (skill SpecialSkill) IsLevelSkill() bool {
	return true
}

func (skill SpecialSkill) CanUse(owner interface{}, fightInstance interface{}, upgrades int) bool {
	return true
}

type SPC_LVL_1 struct {
	SpecialSkill
	DefaultCost
	NoEvents
	NoStats
}

func (skill SPC_LVL_1) GetName() string {
	return "Poziom 1 - Specjalista"
}

func (skill SPC_LVL_1) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE, Flags: types.FLAG_INSTANT_SKILL}
}

func (skill SPC_LVL_1) GetUpgradableTrigger(upgrades int) types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE, Flags: types.FLAG_INSTANT_SKILL}
}

func (skill SPC_LVL_1) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	baseIncrease := 10
	baseDuration := 1

	if HasUpgrade(upgrades, 2) {
		baseIncrease = 12
	}

	if HasUpgrade(upgrades, 3) {
		baseDuration++
	}

	randomStat := utils.RandomElement(
		[]types.Stat{types.STAT_DEF, types.STAT_MR, types.STAT_SPD, types.STAT_AD, types.STAT_AP},
	)

	{
		channelId := fightInstance.(*battle.Fight).Location.CID

		if fightInstance.(*battle.Fight).Tournament != nil {
			channelId = fightInstance.(*battle.Fight).Tournament.Location
		}

		fightInstance.(*battle.Fight).DiscordChannel <- types.DiscordMessageStruct{
			ChannelID:      channelId,
			MessageContent: discord.NewMessageCreateBuilder().SetContentf("Zwiększono statystykę %s o %d%% na %d tur", types.StatToString[randomStat], baseIncrease, baseDuration).Build(),
			DM:             false,
		}
	}

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Target: target.(battle.Entity).GetUUID(),
		Source: owner.(battle.PlayerEntity).GetUUID(),
		Meta: battle.ActionEffect{
			Effect:   battle.EFFECT_STAT_INC,
			Value:    baseIncrease,
			Duration: baseDuration,
			Caster:   owner.(battle.PlayerEntity).GetUUID(),
			Meta: battle.ActionEffectStat{
				Stat:      randomStat,
				Value:     baseIncrease,
				IsPercent: true,
			},
		},
	})

	return nil
}

func (skill SPC_LVL_1) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill SPC_LVL_1) GetCooldown(upgrades int) int {
	baseCD := skill.GetCD()

	if HasUpgrade(upgrades, 1) {
		return baseCD - 1
	}

	return baseCD
}

func (skill SPC_LVL_1) GetDescription() string {
	return "Zwiększa losową statystykę (DEF, RES, SPD, ATK, AP) o 10% na jedną turę"
}

func (skill SPC_LVL_1) GetLevel() int {
	return 1
}

func (skill SPC_LVL_1) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "Cooldown",
			Description: "Zmniejsza czas odnowienia o 1 turę",
		},
		{
			Id:          "Percent",
			Description: "Zwiększa wartość procentową do 12%",
		},
		{
			Id:          "Duration",
			Description: "Zwiększa czas trwania o 1 turę",
		},
	}
}

func (skill SPC_LVL_1) GetUpgradableDescription(upgrades int) string {
	percent := 10

	if HasUpgrade(upgrades, 2) {
		percent = 12
	}

	duration := 1

	if HasUpgrade(upgrades, 3) {
		duration = 2
	}

	return fmt.Sprintf("Zwiększa losową statystykę (DEF, RES, SPD, ATK, AP) o %d%% na %d tur.", percent, duration)
}

type SPC_LVL_2 struct {
	SpecialSkill
	NoExecute
	NoEvents
	NoTrigger
}

func (skill SPC_LVL_2) GetName() string {
	return "Poziom 2 - Specjalista"
}

func (skill SPC_LVL_2) GetDescription() string {
	return "Otrzymujesz 5 kradzieży życia"
}

func (skill SPC_LVL_2) GetLevel() int {
	return 2
}

func (skill SPC_LVL_2) GetStats(upgrades int) map[types.Stat]int {
	stats := map[types.Stat]int{
		types.STAT_ATK_VAMP: 5,
	}

	vampValue := 5
	vampType := types.STAT_ATK_VAMP

	if HasUpgrade(upgrades, 1) {
		vampType = types.STAT_OMNI_VAMP
	}

	if HasUpgrade(upgrades, 2) {
		vampValue = 10
	}

	if HasUpgrade(upgrades, 3) {
		stats[types.STAT_HEAL_SELF] = 20
	}

	stats[vampType] = vampValue

	return stats
}

func (skill SPC_LVL_2) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "Skill",
			Description: "Kradzież życia działa na umiejętności",
		},
		{
			Id:          "Increase",
			Description: "Zwiększa otrzymywaną statystykę do 10",
		},
		{
			Id:          "ShieldInc",
			Description: "Moc leczenia i tarcz (na sobie) zwiększona o 20%",
		},
	}
}

func (skill SPC_LVL_2) GetUpgradableDescription(upgrades int) string {
	vampValue := 5
	vampType := "kradzieży życia"

	if HasUpgrade(upgrades, 1) {
		vampType = "wampiryzmu"
	}

	if HasUpgrade(upgrades, 2) {
		vampValue += 5
	}

	additionalEffect := ""

	if HasUpgrade(upgrades, 3) {
		additionalEffect = "\nLeczenie i tarcze (na sobie) zwiększone o 20%."
	}

	return fmt.Sprintf("Otrzymujesz %d %s.%s", vampValue, vampType, additionalEffect)
}

type SPC_LVL_3 struct {
	SpecialSkill
	DefaultCost
	NoEvents
	NoStats
}

func (skill SPC_LVL_3) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE, Target: &types.TargetTrigger{Target: types.TARGET_ENEMY, MaxTargets: 1}}
}

func (skill SPC_LVL_3) GetUpgradableTrigger(upgrades int) types.Trigger {
	return skill.GetTrigger()
}

func (skill SPC_LVL_3) GetName() string {
	return "Poziom 3 - Specjalista"
}

func (skill SPC_LVL_3) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	baseDmg := 25
	baseHeal := 20

	if HasUpgrade(upgrades, 2) {
		baseDmg = 35
		baseDmg += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AP), 10)
		baseDmg += utils.PercentOf(owner.(battle.PlayerEntity).GetStat(types.STAT_AD), 10)
	}

	if HasUpgrade(upgrades, 3) {
		baseHeal = 25
	}

	healValue := utils.PercentOf(baseDmg, baseHeal)

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_DMG,
		Target: target.(battle.Entity).GetUUID(),
		Source: owner.(battle.PlayerEntity).GetUUID(),
		Meta: battle.ActionDamage{
			Damage: []battle.Damage{
				{
					Type:  types.DMG_TRUE,
					Value: baseDmg,
				},
			},
		},
	})

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Target: owner.(battle.PlayerEntity).GetUUID(),
		Source: owner.(battle.PlayerEntity).GetUUID(),
		Meta: battle.ActionEffect{
			Effect:   battle.EFFECT_HEAL_SELF,
			Value:    healValue,
			Duration: 0,
			Caster:   owner.(battle.PlayerEntity).GetUUID(),
		},
	})

	return nil
}

func (skill SPC_LVL_3) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill SPC_LVL_3) GetCooldown(upgrades int) int {
	baseCD := skill.GetCD()

	if HasUpgrade(upgrades, 1) {
		return baseCD - 1
	}

	return baseCD
}

func (skill SPC_LVL_3) GetDescription() string {
	return "Zadaje 25 obrażeń nieuchronnych i leczy o 20% zadanych obrażeń"
}

func (skill SPC_LVL_3) GetLevel() int {
	return 3
}

func (skill SPC_LVL_3) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "Cooldown",
			Description: "Zmniejsza czas odnowienia o 1 turę",
		},
		{
			Id:          "Damage",
			Description: "Obrażenia zwiększone o 30 + 10%AP + 10%ATK",
		},
		{
			Id:          "Heal",
			Description: "Przelicznik leczenie zwiększony do 25%",
		},
	}
}

func (skill SPC_LVL_3) GetUpgradableDescription(upgrades int) string {
	dmgValue := "25"
	healValue := 20

	if HasUpgrade(upgrades, 2) {
		dmgValue = "55 + 10%AP + 10%ATK"
	}

	if HasUpgrade(upgrades, 3) {
		healValue = 25
	}

	return fmt.Sprintf("Zadaje %s obrażeń nieuchronnych i leczy o %d%% zadanych obrażeń", dmgValue, healValue)
}

type SPC_LVL_4 struct {
	SpecialSkill
	NoEvents
	NoStats
	DefaultActiveTrigger
}

func (skill SPC_LVL_4) GetName() string {
	return "Poziom 4 - Specjalista"
}

func (skill SPC_LVL_4) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	tempSkill := target.(battle.PlayerEntity).GetLvlSkill(meta.(types.SkillChoice).Choice)

	owner.(battle.PlayerEntity).AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value:      tempSkill,
		Expire:     1,
		AfterUsage: HasUpgrade(upgrades, 3),
	})

	return nil
}

func (skill SPC_LVL_4) GetCD() int {
	return BaseCooldowns[skill.GetLevel()] + 1
}

func (skill SPC_LVL_4) GetCooldown(upgrades int) int {
	baseCD := skill.GetCD()

	if HasUpgrade(upgrades, 1) {
		return baseCD - 1
	}

	return baseCD
}

func (skill SPC_LVL_4) GetDescription() string {
	return "Pożycza umiejętność sojusznika"
}

func (skill SPC_LVL_4) GetLevel() int {
	return 4
}

func (skill SPC_LVL_4) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "Cooldown",
			Description: "Zmniejsza czas odnowienia o 1 turę",
		},
		{
			Id:          "Cost",
			Description: "Zmniejsza koszt o 1",
		},
		{
			Id:          "Duration",
			Description: "Umiejętność wygasa do końca walki",
		},
	}
}

func (skill SPC_LVL_4) GetUpgradableCost(upgrades int) int {
	if HasUpgrade(upgrades, 2) {
		return 1
	}

	return 2
}

func (skill SPC_LVL_4) GetCost() int {
	return 2
}

func (skill SPC_LVL_4) GetUpgradableDescription(upgrades int) string {
	duration := "na jedną turę"

	if HasUpgrade(upgrades, 3) {
		duration = "do końca walki"
	}

	return fmt.Sprintf("Pożycza umiejętność sojusznika %s.", duration)
}

type SPC_LVL_5 struct {
	SpecialSkill
	DefaultCost
	NoEvents
	NoStats
}

func (skill SPC_LVL_5) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE, Flags: types.FLAG_INSTANT_SKILL}
}

func (skill SPC_LVL_5) GetUpgradableTrigger(upgrades int) types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE, Flags: types.FLAG_INSTANT_SKILL}
}

func (skill SPC_LVL_5) GetName() string {
	return "Poziom 5 - Specjalista"
}

func (skill SPC_LVL_5) GetDescription() string {
	return "Zmniejsza SPD do początkowej wartości, zwiększa kolejny atak o procent zabranej statystyki"
}

func (skill SPC_LVL_5) GetLevel() int {
	return 5
}

func (skill SPC_LVL_5) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{
			Id:          "Skill",
			Description: "Zwiększa także obrażenia umiejętności",
		},
		{
			Id:          "Duration",
			Description: "Efekt utrzymuje się przez całą turę",
		},
		{
			Id:          "DmgReduction",
			Description: "Podczas trwania zmniejsza obrażenia o 10%",
		},
	}
}

type SPC_LVL_5_EFFECT struct {
	NoCooldown
	NoCost
	NoStats
	NoLevel
	NoEvents
}

func (skill SPC_LVL_5_EFFECT) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	tempMeta := meta.(types.AttackTriggerMeta)

	for _, effect := range tempMeta.Effects {
		effect.Value = utils.PercentOf(effect.Value, 20)
	}

	return tempMeta
}

func (skill SPC_LVL_5_EFFECT) GetName() string {
	return "Poziom 5 - Specjalista - Efekt"
}

func (skill SPC_LVL_5_EFFECT) GetDescription() string {
	return "Poziom 5 - Specjalista - Efekt"
}

func (skill SPC_LVL_5_EFFECT) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_PASSIVE, Event: types.TRIGGER_ATTACK_BEFORE}
}

func (skill SPC_LVL_5) UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{} {
	spdReduction := owner.(battle.PlayerEntity).GetStat(types.STAT_SPD) - owner.(battle.PlayerEntity).GetDefaultStat(types.STAT_SPD)

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Target: owner.(battle.PlayerEntity).GetUUID(),
		Source: owner.(battle.PlayerEntity).GetUUID(),
		Meta: battle.ActionEffect{
			Effect:   battle.EFFECT_STAT_DEC,
			Value:    spdReduction,
			Duration: 1,
			Caster:   owner.(battle.PlayerEntity).GetUUID(),
			Meta: battle.ActionEffectStat{
				Stat:      types.STAT_SPD,
				Value:     spdReduction,
				IsPercent: false,
			},
		},
	})

	fightInstance.(*battle.Fight).DiscordChannel <- types.DiscordMessageStruct{
		ChannelID:      fightInstance.(*battle.Fight).Location.CID,
		MessageContent: discord.NewMessageCreateBuilder().SetContentf("Zwiększenie obrażeń wynosi %d", spdReduction).Build(),
		DM:             false,
	}

	owner.(battle.PlayerEntity).AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value:      SPC_LVL_5_EFFECT{},
		Expire:     1,
		AfterUsage: HasUpgrade(upgrades, 2),
	})

	if HasUpgrade(upgrades, 3) {
		fightInstance.(*battle.Fight).HandleAction(battle.Action{
			Event:  battle.ACTION_EFFECT,
			Target: owner.(battle.PlayerEntity).GetUUID(),
			Source: owner.(battle.PlayerEntity).GetUUID(),
			Meta: battle.ActionEffect{
				Effect:   battle.EFFECT_RESIST,
				Value:    10,
				Duration: 1,
				Caster:   owner.(battle.PlayerEntity).GetUUID(),
				Meta: battle.ActionEffectResist{
					Value:     10,
					IsPercent: true,
				},
			},
		})
	}

	return nil
}

func (skill SPC_LVL_5) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill SPC_LVL_5) GetCooldown(upgrades int) int {
	return skill.GetCD()
}

func (skill SPC_LVL_5) GetUpgradableDescription(upgrades int) string {
	trigger := ""

	if HasUpgrade(upgrades, 1) {
		trigger += " lub umiejętność"
	}

	effectDuration := "kolejny"
	if HasUpgrade(upgrades, 2) {
		effectDuration = ""
	}

	resist := ""
	if HasUpgrade(upgrades, 3) {
		resist = "\nPodczas trwania zmniejsza obrażenia o 10%"
	}

	return fmt.Sprintf("Zmniejsza SPD do początkowej wartości na jedną turę. Zwiększa%s atak%s o zabraną wartość.%s", effectDuration, trigger, resist)
}
