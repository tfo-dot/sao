package inventory

import (
	"fmt"
	"sao/types"
	"sao/utils"

	"github.com/disgoorg/disgo/discord"
	"github.com/google/uuid"
)

type SpecialSkill struct{}

func (skill SpecialSkill) GetPath() types.SkillPath {
	return types.PathSpecial
}

func (skill SpecialSkill) GetUUID() uuid.UUID {
	return uuid.Nil
}

func (skill SpecialSkill) IsLevelSkill() bool {
	return true
}

func (skill SpecialSkill) CanUse(owner types.PlayerEntity, fightInstance types.FightInstance) bool {
	return true
}

func (skill SpecialSkill) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	return nil
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

func (skill SPC_LVL_1) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	baseIncrease := 10
	baseDuration := 1

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 2) {
		baseIncrease = 12
	}

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 3) {
		baseDuration++
	}

	randomStat := utils.RandomElement(
		[]types.Stat{types.STAT_DEF, types.STAT_MR, types.STAT_SPD, types.STAT_AD, types.STAT_AP},
	)

	fightInstance.SendMessage(
		fightInstance.GetChannelId(),
		discord.NewMessageCreateBuilder().SetContentf("Zwiększono statystykę %s o %d%% na %d tur", types.StatToString[randomStat], baseIncrease, baseDuration).Build(),
		false,
	)

	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_EFFECT,
		Target: target.GetUUID(),
		Source: owner.GetUUID(),
		Meta: types.ActionEffect{
			Effect:   types.EFFECT_STAT_INC,
			Value:    baseIncrease,
			Duration: baseDuration,
			Meta:     types.ActionEffectStat{Stat: randomStat, IsPercent: true},
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
		{Id: "Cooldown", Description: "Zmniejsza czas odnowienia o 1 turę"},
		{Id: "Percent", Description: "Zwiększa wartość procentową do 12%"},
		{Id: "Duration", Description: "Zwiększa czas trwania o 1 turę"},
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
	stats := map[types.Stat]int{}

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

func (skill SPC_LVL_2) GetDerivedStats(upgrades int) []types.DerivedStat {
	return []types.DerivedStat{}
}

func (skill SPC_LVL_2) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{Id: "Skill", Description: "Kradzież życia działa na umiejętności"},
		{Id: "Increase", Description: "Zwiększa otrzymywaną statystykę do 10"},
		{Id: "ShieldInc", Description: "Moc leczenia i tarcz (na sobie) zwiększona o 20%"},
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

func (skill SPC_LVL_3) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta any) any {
	baseDmg := 25
	baseHeal := 20

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 2) {
		baseDmg = 35 + utils.PercentOf(owner.GetStat(types.STAT_AP), 10)
		baseDmg += utils.PercentOf(owner.GetStat(types.STAT_AD), 10)
	}

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 3) {
		baseHeal = 25
	}

	healValue := utils.PercentOf(baseDmg, baseHeal)

	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_DMG,
		Target: target.GetUUID(),
		Source: owner.GetUUID(),
		Meta:   types.ActionDamage{Damage: []types.Damage{{Type: types.DMG_TRUE, Value: baseDmg}}},
	})

	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_EFFECT,
		Source: owner.GetUUID(),
		Meta:   types.ActionEffect{Effect: types.EFFECT_HEAL, Value: healValue},
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
		{Id: "Cooldown", Description: "Zmniejsza czas odnowienia o 1 turę"},
		{Id: "Damage", Description: "Obrażenia zwiększone o 30 + 10%AP + 10%ATK"},
		{Id: "Heal", Description: "Przelicznik leczenie zwiększony do 25%"},
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

func (skill SPC_LVL_4) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta any) any {
	tempSkill := target.(types.PlayerEntity).GetLvlSkill(meta.(types.SkillChoice).Choice)

	owner.AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value:      tempSkill,
		Expire:     1,
		AfterUsage: HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 3),
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
		{Id: "Cooldown", Description: "Zmniejsza czas odnowienia o 1 turę"},
		{Id: "Cost", Description: "Zmniejsza koszt o 1"},
		{Id: "Duration", Description: "Umiejętność wygasa do końca walki"},
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
		{Id: "Skill", Description: "Zwiększa także obrażenia umiejętności"},
		{Id: "Duration", Description: "Efekt utrzymuje się przez całą turę"},
		{Id: "DmgReduction", Description: "Podczas trwania zmniejsza obrażenia o 10%"},
	}
}

type SPC_LVL_5_EFFECT struct {
	NoCooldown
	NoCost
	NoStats
	NoLevel
	NoEvents
}

func (skill SPC_LVL_5_EFFECT) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta any) any {
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

func (skill SPC_LVL_5) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	spdReduction := owner.GetStat(types.STAT_SPD) - owner.GetDefaultStat(types.STAT_SPD)

	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_EFFECT,
		Source: owner.GetUUID(),
		Meta: types.ActionEffect{
			Effect:   types.EFFECT_STAT_DEC,
			Value:    spdReduction,
			Duration: 1,
			Meta:     types.ActionEffectStat{Stat: types.STAT_SPD},
		},
	})

	fightInstance.SendMessage(
		fightInstance.GetChannelId(),
		discord.NewMessageCreateBuilder().SetContentf("Zwiększenie obrażeń wynosi %d", spdReduction).Build(),
		false,
	)

	owner.AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value:      SPC_LVL_5_EFFECT{},
		Expire:     1,
		AfterUsage: HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 2),
	})

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 3) {
		fightInstance.HandleAction(types.Action{
			Event:  types.ACTION_EFFECT,
			Target: owner.GetUUID(),
			Source: owner.GetUUID(),
			Meta: types.ActionEffect{
				Effect:   types.EFFECT_RESIST,
				Value:    10,
				Duration: 1,
				Meta:     types.ActionEffectResist{IsPercent: true, DmgType: 4},
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

type SPC_LVL_6 struct {
	SpecialSkill
	NoStats
	NoEvents
	DefaultCost
}

func (skill SPC_LVL_6) GetName() string {
	return "Poziom 6 - Specjalista"
}

func (skill SPC_LVL_6) GetDescription() string {
	return "Zmniejsza leczenie i tarcze przeciwnika o 10%"
}

func (skill SPC_LVL_6) GetLevel() int {
	return 6
}

func (skill SPC_LVL_6) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{
		{Id: "Duration", Description: "Efekt trwa przez dwie tury"},
		{Id: "EffectIncrease", Description: "Zwiększa redukcje do 20%"},
		{Id: "Damage", Description: "Zadaje obrażenia równe 1% maks HP celu"},
	}
}

func (skill SPC_LVL_6) GetUpgradableDescription(upgrades int) string {
	duration := "jedną turę"

	if HasUpgrade(upgrades, 1) {
		duration = "dwie tury"
	}

	effectIncrease := 10
	if HasUpgrade(upgrades, 2) {
		effectIncrease = 20
	}

	damage := ""
	if HasUpgrade(upgrades, 3) {
		damage = "\nZadaje obrażenia równe 1% maks HP celu"
	}

	return fmt.Sprintf("Zmniejsza leczenie i tarcze przeciwnika o %d%% na %s.%s", effectIncrease, duration, damage)
}

func (skill SPC_LVL_6) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	effectValue := 10

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 2) {
		effectValue = 20
	}

	duration := 1

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 1) {
		duration++
	}

	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_EFFECT,
		Target: target.GetUUID(),
		Source: owner.GetUUID(),
		Meta: types.ActionEffect{
			Effect:   types.EFFECT_STAT_DEC,
			Value:    effectValue,
			Duration: 1,
			Meta:     types.ActionEffectStat{Stat: types.STAT_HEAL_SELF},
		},
	})

	if HasUpgrade(owner.GetUpgrades(skill.GetLevel()), 3) {
		fightInstance.HandleAction(types.Action{
			Event:  types.ACTION_DMG,
			Target: target.GetUUID(),
			Source: owner.GetUUID(),
			Meta: types.ActionDamage{
				Damage: []types.Damage{{Type: types.DMG_TRUE, Value: target.GetStat(types.STAT_HP) / 100}},
			},
		})
	}

	return nil
}

func (skill SPC_LVL_6) GetCD() int {
	return BaseCooldowns[skill.GetLevel()]
}

func (skill SPC_LVL_6) GetCooldown(upgrades int) int {
	return skill.GetCD()
}

func (skill SPC_LVL_6) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_ACTIVE, Target: &types.TargetTrigger{Target: types.TARGET_ENEMY, MaxTargets: 1}}
}

func (skill SPC_LVL_6) GetUpgradableTrigger(upgrades int) types.Trigger {
	return skill.GetTrigger()
}

type SPC_ULT_1 struct {
	SpecialSkill
	NoStats
	NoEvents
	DefaultActiveTrigger
}

func (skill SPC_ULT_1) GetName() string {
	return "Poziom 10 - Specjalista"
}

func (skill SPC_ULT_1) GetDescription() string {
	return "Jesteś niemożliwy do trafienia przez 10 tur. Spowalniasz do bazowej wartości SPD, wszyscy sojusznicy otrzymują połowe dodatkowej prędkości. Gdy jesteś niemożliwy do trafienia twoje ataki zadają 50% mniej obrażeń ale za to przywracają 1many i dają ci 10 SPD i 5 DGD do końca walki."
}

func (skill SPC_ULT_1) GetLevel() int {
	return 10
}

func (skill SPC_ULT_1) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{}
}

func (skill SPC_ULT_1) GetUpgradableDescription(upgrades int) string {
	return skill.GetDescription()
}

func (skill SPC_ULT_1) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	return skill.Execute(owner, target, fightInstance, meta)
}

func (skill SPC_ULT_1) GetCD() int {
	return 10
}

func (skill SPC_ULT_1) GetCooldown(upgrades int) int {
	return skill.GetCD()
}

func (skill SPC_ULT_1) GetCost() int {
	return 1
}

func (skill SPC_ULT_1) GetUpgradableCost(upgrades int) int {
	return skill.GetCost()
}

type SPC_ULT_1_EFFECT_1 struct {
	NoCooldown
	NoCost
	NoStats
	NoLevel
	NoEvents
}

func (skill SPC_ULT_1_EFFECT_1) GetName() string {
	return "Poziom 10 - Specjalista - Efekt 1"
}

func (skill SPC_ULT_1_EFFECT_1) GetDescription() string {
	return ""
}

func (skill SPC_ULT_1_EFFECT_1) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_PASSIVE, Event: types.TRIGGER_ATTACK_BEFORE}
}

func (skill SPC_ULT_1_EFFECT_1) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	owner.RestoreMana(1)

	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_EFFECT,
		Target: owner.GetUUID(),
		Source: owner.GetUUID(),
		Meta: types.ActionEffect{
			Effect:   types.EFFECT_STAT_INC,
			Value:    10,
			Duration: -1,
			Meta:     types.ActionEffectStat{Stat: types.STAT_SPD},
		},
	})

	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_EFFECT,
		Target: owner.GetUUID(),
		Source: owner.GetUUID(),
		Meta: types.ActionEffect{
			Effect:   types.EFFECT_STAT_INC,
			Value:    5,
			Duration: -1,
			Meta:     types.ActionEffectStat{Stat: types.STAT_AGL},
		},
	})

	return types.AttackTriggerMeta{Effects: []types.DamagePartial{{Percent: true, Value: -50, Type: types.DMG_PHYSICAL}}}
}

type SPC_ULT_1_EFFECT_2 struct {
	NoCooldown
	NoCost
	NoStats
	NoLevel
	NoEvents
}

func (skill SPC_ULT_1_EFFECT_2) GetName() string {
	return "Poziom 10 - Specjalista - Efekt 2"
}

func (skill SPC_ULT_1_EFFECT_2) GetDescription() string {
	return ""
}

func (skill SPC_ULT_1_EFFECT_2) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_PASSIVE, Event: types.TRIGGER_DAMAGE_GOT_HIT}
}

func (skill SPC_ULT_1_EFFECT_2) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	return types.AttackTriggerMeta{
		Effects:    []types.DamagePartial{},
		ShouldMiss: true,
	}
}

func (skill SPC_ULT_1) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	valueIncrease := owner.GetStat(types.STAT_SPD) - owner.GetDefaultStat(types.STAT_SPD)

	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_EFFECT,
		Target: owner.GetUUID(),
		Source: owner.GetUUID(),
		Meta: types.ActionEffect{
			Effect:   types.EFFECT_STAT_DEC,
			Value:    valueIncrease,
			Duration: 10,
			Meta:     types.ActionEffectStat{Stat: types.STAT_SPD},
		},
	})

	for _, entity := range fightInstance.GetAlliesFor(owner.GetUUID()) {
		fightInstance.HandleAction(types.Action{
			Event:  types.ACTION_EFFECT,
			Target: entity.GetUUID(),
			Source: owner.GetUUID(),
			Meta: types.ActionEffect{
				Effect:   types.EFFECT_STAT_INC,
				Value:    valueIncrease,
				Duration: 10,
				Meta:     types.ActionEffectStat{Stat: types.STAT_SPD},
			},
		})
	}

	owner.AppendTempSkill(types.WithExpire[types.PlayerSkill]{Value: SPC_ULT_1_EFFECT_1{}, Expire: 10})
	owner.AppendTempSkill(types.WithExpire[types.PlayerSkill]{Value: SPC_ULT_1_EFFECT_2{}, Expire: 10})

	return nil
}

type SPC_ULT_2 struct {
	SpecialSkill
	NoStats
	NoEvents
	DefaultActiveTrigger
}

func (skill SPC_ULT_2) GetName() string {
	return "Poziom 10 - Specjalista"
}

func (skill SPC_ULT_2) GetDescription() string {
	return "Usuwasz negatywne efekty sojuszników. Leczysz wszystkich sojuszników o 25% ich maksymalnego zdrowia, otrzymują tarczę w wysokości 20% twojego HP. Za każdego sojusznika otrzymujesz 10 SPD do końca walki. Efekt pasywny: za każde 10 dodatkowego SPD twoje obrażenia są zwiększone o 5%"
}

func (skill SPC_ULT_2) GetLevel() int {
	return 10
}

func (skill SPC_ULT_2) GetUpgrades() []types.PlayerSkillUpgrade {
	return []types.PlayerSkillUpgrade{}
}

func (skill SPC_ULT_2) GetUpgradableDescription(upgrades int) string {
	return skill.GetDescription()
}

func (skill SPC_ULT_2) UpgradableExecute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	return skill.Execute(owner, target, fightInstance, meta)
}

type SPC_ULT_2_EFFECT struct {
	NoCooldown
	NoCost
	NoStats
	NoLevel
	NoEvents
}

func (skill SPC_ULT_2_EFFECT) GetName() string {
	return "Poziom 10 - Specjalista - Efekt"
}

func (skill SPC_ULT_2_EFFECT) GetDescription() string {
	return ""
}

func (skill SPC_ULT_2_EFFECT) GetTrigger() types.Trigger {
	return types.Trigger{Type: types.TRIGGER_PASSIVE, Event: types.TRIGGER_ATTACK_BEFORE}
}

func (skill SPC_ULT_2_EFFECT) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	statDiff := owner.GetStat(types.STAT_SPD) - owner.GetDefaultStat(types.STAT_SPD)

	percentChange := 0

	for statDiff > 10 {
		percentChange++
		statDiff -= 10
	}

	return types.AttackTriggerMeta{
		Effects: []types.DamagePartial{{Percent: true, Value: percentChange * 5, Type: types.DMG_PHYSICAL}},
	}
}

func (skill SPC_ULT_2) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	allies := fightInstance.GetAlliesFor(owner.GetUUID())

	for _, entity := range allies {
		entity.Cleanse()

		fightInstance.HandleAction(types.Action{
			Event:  types.ACTION_EFFECT,
			Source: owner.GetUUID(),
			Meta: types.ActionEffect{
				Effect: types.EFFECT_HEAL,
				Value:  utils.PercentOf(entity.GetStat(types.STAT_HP), 25),
			},
		})

		fightInstance.HandleAction(types.Action{
			Event:  types.ACTION_EFFECT,
			Source: owner.GetUUID(),
			Meta: types.ActionEffect{
				Effect: types.EFFECT_SHIELD,
				Value:  utils.PercentOf(owner.GetStat(types.STAT_HP), 20),
			},
		})
	}

	fightInstance.HandleAction(types.Action{
		Event:  types.ACTION_EFFECT,
		Source: owner.GetUUID(),
		Meta: types.ActionEffect{
			Effect:   types.EFFECT_STAT_INC,
			Value:    10 * len(allies),
			Duration: -1,
			Meta:     types.ActionEffectStat{Stat: types.STAT_SPD},
		},
	})

	owner.AppendTempSkill(types.WithExpire[types.PlayerSkill]{Value: SPC_ULT_2_EFFECT{}, Expire: -1})

	return nil
}

func (skill SPC_ULT_2) GetCD() int {
	return 10
}

func (skill SPC_ULT_2) GetCooldown(upgrades int) int {
	return skill.GetCD()
}

func (skill SPC_ULT_2) GetCost() int {
	return 1
}

func (skill SPC_ULT_2) GetUpgradableCost(upgrades int) int {
	return skill.GetCost()
}
