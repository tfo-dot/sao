package data

import (
	"sao/battle"
	"sao/types"
	"sao/utils"
	"sort"

	"github.com/google/uuid"
)

var Items = map[uuid.UUID]types.PlayerItem{
	ReimiBlessingUUID: {
		UUID:        ReimiBlessingUUID,
		Name:        "Błogosławieństwo Reimi",
		Description: "Przeleczenie daje tarczę.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_AD:       25,
			types.STAT_HP:       100,
			types.STAT_ATK_VAMP: 10,
		},
		Effects: []types.PlayerSkill{ReimiBlessingSkill{}},
	},
	GiantSlayerUUID: {
		UUID:        GiantSlayerUUID,
		Name:        "Pogromca gigantów",
		Description: "Zadaje dodatkowe obrażenia w zależności od pancerza przeciwnika.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_AD:     25,
			types.STAT_LETHAL: 10,
		},
		Effects: []types.PlayerSkill{GiantSlayerSkill{}},
	},
	GiantKillerUUID: {
		UUID:        GiantKillerUUID,
		Name:        "Zabójca gigantów",
		Description: "Zadaje dodatkowe obrażenia w zależności od zdrowia przeciwnika.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_AD:     25,
			types.STAT_LETHAL: 10,
		},
		Effects: []types.PlayerSkill{GiantKillerSkill{}},
	},
	MageKillerUUID: {
		UUID:        MageKillerUUID,
		Name:        "Zabójca magów",
		Description: "Cele osłonięte tarczą otrzymają większe obrażenia.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_AD:     25,
			types.STAT_LETHAL: 10,
		},
		Effects: []types.PlayerSkill{MageKillerSkill{}},
	},
	SandBladeUUID: {
		UUID:        SandBladeUUID,
		Name:        "Piaskowe ostrze",
		Description: "Zadawanie obrażeń zmniejsza leczenie wroga.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_AD:  30,
			types.STAT_SPD: 5,
		},
		Effects: []types.PlayerSkill{SandBladeSkill{}},
	},
	WaterBladeUUID: {
		UUID:        WaterBladeUUID,
		Name:        "Wodne ostrze",
		Description: "Zadawanie obrażeń leczy o brakujące zdrowie.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_AD:       25,
			types.STAT_ATK_VAMP: 10,
			types.STAT_HP:       50,
		},
		Effects: []types.PlayerSkill{WaterBladeSkill{}},
	},
	DefenseVisageUUID: {
		UUID:        DefenseVisageUUID,
		Name:        "Oblicze obrony",
		Description: "Dostajesz ATK w zależności od maks. HP.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_HP:  100,
			types.STAT_DEF: 15,
			types.STAT_MR:  15,
		},
		Effects: []types.PlayerSkill{DefenseVisageSkill{}},
	},
	AttackVisageUUID: {
		UUID:        AttackVisageUUID,
		Name:        "Oblicze ataku",
		Description: "Dostajesz HP w zależności od ATK.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_AD:  20,
			types.STAT_MR:  15,
			types.STAT_DEF: 15,
		},
		Effects: []types.PlayerSkill{AttackVisageSkill{}},
	},
	WarriorsLegacyUUID: {
		UUID:        WarriorsLegacyUUID,
		Name:        "Dziedzictwo wojownika",
		Description: "Zwiększa obrażenia w zależności od maks zdrowia.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_AD: 20,
			types.STAT_HP: 50,
		},
		Effects: []types.PlayerSkill{WarriorsLegacySkill{}},
	},
	SecondBreathUUID: {
		UUID:        SecondBreathUUID,
		Name:        "Drugie oddech",
		Description: "Zwiększa otrzymywane leczenie.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_HP:        200,
			types.STAT_DEF:       10,
			types.STAT_MR:        10,
			types.STAT_HEAL_SELF: 20,
		},
		Effects: []types.PlayerSkill{},
	},
	LilithsWrathUUID: {
		UUID:        LilithsWrathUUID,
		Name:        "Gniew Lilith",
		Description: "Co ture zadaje obrażenia w zależności od zdrowia użytkownika.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_HP:  200,
			types.STAT_DEF: 30,
		},
		Effects: []types.PlayerSkill{LilithsWrathSkill{}},
	},
	RyuLegacyUUID: {
		UUID:        RyuLegacyUUID,
		Name:        "Dziedzictwo Ryu",
		Description: "Zwiększa RES i DEF o 20%.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_HP:  150,
			types.STAT_DEF: 40,
			types.STAT_MR:  40,
		},
		Effects: []types.PlayerSkill{RyuLegacySkill{}},
	},
	DefenderBladeUUID: {
		UUID:        DefenderBladeUUID,
		Name:        "Ostrze obrońcy",
		Description: "Zwiększa ataki o twój RES i DEF.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_HP:  150,
			types.STAT_DEF: 20,
			types.STAT_MR:  30,
			types.STAT_AD:  20,
		},
		Effects: []types.PlayerSkill{DefenderBladeSkill{}},
	},
	GrudgeArmorUUID: {
		UUID:        GrudgeArmorUUID,
		Name:        "Pancerz zwady",
		Description: "Zadaje obrażenia wrogom, którzy cię uderzają i zmniejsza ich leczenie.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_HP:  150,
			types.STAT_DEF: 30,
		},
		Effects: []types.PlayerSkill{GrudgeArmorSkill{}},
	},
	AmplifyingCoatUUID: {
		UUID:        AmplifyingCoatUUID,
		Name:        "Płaszcz wzmacniający",
		Description: "Zwiększa maksymalne zdrowie o 20%.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_HP:  300,
			types.STAT_DEF: 10,
			types.STAT_MR:  10,
		},
		Effects: []types.PlayerSkill{AmplifyingCoatSkill{}},
	},
	ControllersBraceletUUID: {
		UUID:        ControllersBraceletUUID,
		Name:        "Bransoleta kontrolera",
		Description: "Nałożenie efektu CC leczy ciebie i sojusznika.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_AP:  20,
			types.STAT_AD:  10,
			types.STAT_SPD: 5,
		},
		Effects: []types.PlayerSkill{ControllersBraceletSkill{}},
	},
	CursedIceUUID: {
		UUID:        CursedIceUUID,
		Name:        "Przeklęty lód",
		Description: "Efekty spowolnienia są mocniejsze",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_AP: 20,
			types.STAT_AD: 20,
		},
		Effects: []types.PlayerSkill{CursedIceSkill{}},
	},
	ControllersRuneUUID: {
		UUID:        ControllersRuneUUID,
		Name:        "Runa kontrolera",
		Description: "Zabicie wroga objętego CC przywraca manę.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_AP: 20,
			types.STAT_AD: 20,
		},
		Effects: []types.PlayerSkill{ControllersRuneSkill{}},
	},
	ControllersNecklaceUUID: {
		UUID:        ControllersNecklaceUUID,
		Name:        "Naszyjnik kontrolera",
		Description: "Nałożenie efektu CC zwiększa twoją prędkość.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_AP:  20,
			types.STAT_AD:  10,
			types.STAT_SPD: 5,
		},
		Effects: []types.PlayerSkill{ControllersNecklaceSkill{}},
	},
	ControllersBladeUUID: {
		UUID:        ControllersBladeUUID,
		Name:        "Ostrze kontrolera",
		Description: "Atakowanie zmniejsza prędkość wrogów.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_SPD: 10,
			types.STAT_AD:  15,
		},
		Effects: []types.PlayerSkill{ControllersBladeSkill{}},
	},
	ControllersHatUUID: {
		UUID:        ControllersHatUUID,
		Name:        "Kapelusz kontrolera",
		Description: "Daje siłę adaptacyjną w zależności od many.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_MANA: 5,
		},
		Effects: []types.PlayerSkill{ControllersHatSkill{}},
	},
	ArdentCenserUUID: {
		UUID:        ArdentCenserUUID,
		Name:        "Ognisty trybularz",
		Description: "Leczenie i tarcze zwiększają obrażenia i prędkość sojusznika.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_HEAL_POWER: 10,
			types.STAT_AP:         30,
			types.STAT_HP:         50,
		},
		Effects: []types.PlayerSkill{ArdentCenserSkill{}},
	},
	SirensCallUUID: {
		UUID:        SirensCallUUID,
		Name:        "Syreni śpiew",
		Description: "Leczenie i tarcze przeskakują na sojusznika",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_HEAL_POWER: 10,
			types.STAT_AP:         40,
			types.STAT_HP:         50,
		},
		Effects: []types.PlayerSkill{SirensCallSkill{}},
	},
	FogsEmpowermentUUID: {
		UUID:        FogsEmpowermentUUID,
		Name:        "Mgliste wzmocnienie",
		Description: "Otrzymujesz AP w zależności od siły leczenia i tarcz.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_HEAL_POWER: 15,
			types.STAT_AP:         30,
		},
		Effects: []types.PlayerSkill{FogsEmpowermentSkill{}},
	},
	WindsEmpowermentUUID: {
		UUID:        WindsEmpowermentUUID,
		Name:        "Wietrzne wzmocnienie",
		Description: "Otrzymujesz SPD w zależności od siły leczenia i tarcz.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_HEAL_POWER: 10,
			types.STAT_AD:         15,
		},
		Effects: []types.PlayerSkill{WindsEmpowermentSkill{}},
	},
	KyokiBeltUUID: {
		UUID:        KyokiBeltUUID,
		Name:        "Pas Kyoki",
		Description: "Zwiększa obrażenia w zależności od losowego mnożnika.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      true,
		Stats: map[types.Stat]int{
			types.STAT_AD: 20,
			types.STAT_AP: 40,
		},
		Effects: []types.PlayerSkill{KyokiBeltSkill{}},
	},
	ShikiFlameUUID: {
		UUID:        ShikiFlameUUID,
		Name:        "Płomień Shiki",
		Description: "Trafienie zaklęciem zadaje dodatkowe obrażenia w zależności od zdrowia.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      true,
		Stats: map[types.Stat]int{
			types.STAT_AP:  100,
			types.STAT_SPD: 5,
		},
		Effects: []types.PlayerSkill{ShikiFlameSkill{}},
	},
	StormHarbingerUUID: {
		UUID:        StormHarbingerUUID,
		Name:        "Zwiastun burzy",
		Description: "Ataki zadają dodatkowe obrażenia w zależności od AP.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      true,
		Stats: map[types.Stat]int{
			types.STAT_AP: 50,
			types.STAT_AD: 10,
		},
		Effects: []types.PlayerSkill{StormHarbingerSkill{}},
	},
}

var ReimiBlessingUUID = uuid.MustParse("00000000-0000-0000-0000-000000000000")
var ReimiBlessingSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000001")
var ReimiBlessingEffectUUID = uuid.MustParse("00000000-0000-0001-0001-000000000001")

var GiantSlayerUUID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
var GiantSlayerSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000001")

var GiantKillerUUID = uuid.MustParse("00000000-0000-0000-0000-000000000002")
var GiantKillerSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000002")

var MageKillerUUID = uuid.MustParse("00000000-0000-0000-0000-000000000003")
var MageKillerSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000003")

var SandBladeUUID = uuid.MustParse("00000000-0000-0000-0000-000000000004")
var SandBladeSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000004")

var WaterBladeUUID = uuid.MustParse("00000000-0000-0000-0000-000000000005")
var WaterBladeSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000005")

var DefenseVisageUUID = uuid.MustParse("00000000-0000-0000-0000-000000000006")
var DefenseVisageSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000006")
var DefenseVisageEffectUUID = uuid.MustParse("00000000-0000-0001-0001-000000000006")

var AttackVisageUUID = uuid.MustParse("00000000-0000-0000-0000-000000000007")
var AttackVisageSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000007")
var AttackVisageEffectUUID = uuid.MustParse("00000000-0000-0001-0001-000000000007")

var WarriorsLegacyUUID = uuid.MustParse("00000000-0000-0000-0000-000000000008")
var WarriorsLegacySkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000008")

var SecondBreathUUID = uuid.MustParse("00000000-0000-0000-0000-000000000009")

var LilithsWrathUUID = uuid.MustParse("00000000-0000-0000-0000-00000000000A")
var LilithsWrathSkillUUID = uuid.MustParse("00000000-0000-0001-0000-00000000000A")

var RyuLegacyUUID = uuid.MustParse("00000000-0000-0000-0000-00000000000B")
var RyuLegacySkillUUID = uuid.MustParse("00000000-0000-0001-0000-00000000000B")

var DefenderBladeUUID = uuid.MustParse("00000000-0000-0000-0000-00000000000C")
var DefenderBladeSkillUUID = uuid.MustParse("00000000-0000-0001-0000-00000000000C")

var GrudgeArmorUUID = uuid.MustParse("00000000-0000-0000-0000-00000000000D")
var GrudgeArmorSkillUUID = uuid.MustParse("00000000-0000-0001-0000-00000000000D")

var AmplifyingCoatUUID = uuid.MustParse("00000000-0000-0000-0000-00000000000E")
var AmplifyingCoatSkillUUID = uuid.MustParse("00000000-0000-0001-0000-00000000000E")

var ControllersBraceletUUID = uuid.MustParse("00000000-0000-0000-0000-00000000000F")
var ControllersBraceletSkillUUID = uuid.MustParse("00000000-0000-0001-0000-00000000000F")

var CursedIceUUID = uuid.MustParse("00000000-0000-0000-0000-000000000010")
var CursedIceSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000010")

var ControllersRuneUUID = uuid.MustParse("00000000-0000-0000-0000-000000000011")
var ControllersRuneSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000011")
var ControllersRuneEffectUUID = uuid.MustParse("00000000-0000-0001-0001-000000000011")

var ControllersNecklaceUUID = uuid.MustParse("00000000-0000-0000-0000-000000000012")
var ControllersNecklaceSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000012")
var ControllersNecklaceEffectUUID = uuid.MustParse("00000000-0000-0001-0001-000000000012")

var ControllersBladeUUID = uuid.MustParse("00000000-0000-0000-0000-000000000013")
var ControllersBladeSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000013")
var ControllersBladeEffectUUID = uuid.MustParse("00000000-0000-0001-0001-000000000013")

var ControllersHatUUID = uuid.MustParse("00000000-0000-0000-0000-000000000014")
var ControllersHatSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000014")

var ArdentCenserUUID = uuid.MustParse("00000000-0000-0000-0000-000000000015")
var ArdentCenserSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000015")
var ArdentCenserEffectUUID = uuid.MustParse("00000000-0000-0001-0001-000000000015")

var SirensCallUUID = uuid.MustParse("00000000-0000-0000-0000-000000000016")
var SirensCallSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000016")
var SirensCallEffectUUID = uuid.MustParse("00000000-0000-0001-0001-000000000016")

var FogsEmpowermentUUID = uuid.MustParse("00000000-0000-0000-0000-000000000017")
var FogsEmpowermentSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000017")

var WindsEmpowermentUUID = uuid.MustParse("00000000-0000-0000-0000-000000000018")
var WindsEmpowermentSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000018")

var LightingSupportUUID = uuid.MustParse("00000000-0000-0000-0000-000000000019")
var LightingSupportSkillUUID = uuid.MustParse("00000000-0000-0001-0000-000000000019")
var LightingSupportEffectUUID = uuid.MustParse("00000000-0000-0001-0001-000000000019")

var GodlySupportUUID = uuid.MustParse("00000000-0000-0000-0000-00000000001A")
var GodlySupportSkillUUID = uuid.MustParse("00000000-0000-0001-0000-00000000001A")

var KyokiBeltUUID = uuid.MustParse("00000000-0000-0000-0000-00000000001B")
var KyokiBeltSkillUUID = uuid.MustParse("00000000-0000-0001-0000-00000000001B")

var ShikiFlameUUID = uuid.MustParse("00000000-0000-0000-0000-00000000001C")
var ShikiFlameSkillUUID = uuid.MustParse("00000000-0000-0001-0000-00000000001C")

var StormHarbingerUUID = uuid.MustParse("00000000-0000-0000-0000-00000000001D")
var StormHarbingerSkillUUID = uuid.MustParse("00000000-0000-0001-0000-00000000001D")

type BasePassiveSkill struct{}

func (bps BasePassiveSkill) GetCD() int {
	return 0
}

func (bps BasePassiveSkill) GetCost() int {
	return 0
}

func (bps BasePassiveSkill) IsLevelSkill() bool {
	return false
}

func (bps BasePassiveSkill) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return nil
}

func (bps BasePassiveSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	return nil
}

type ReimiBlessingSkill struct{ BasePassiveSkill }

func (rbs ReimiBlessingSkill) GetName() string {
	return "Uświęcona tarcza"
}

func (rbs ReimiBlessingSkill) GetDescription() string {
	return "Przeleczenie daje tarczę."
}

func (rbs ReimiBlessingSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_HEAL_SELF,
			TargetType:    []types.TargetTag{types.TARGET_SELF},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (rbs ReimiBlessingSkill) GetUUID() uuid.UUID {
	return ReimiBlessingSkillUUID
}

func (rbs ReimiBlessingSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	ownerEntity := owner.(battle.Entity)

	oldEffect := ownerEntity.GetEffectByUUID(ReimiBlessingEffectUUID)

	maxShield := utils.PercentOf(ownerEntity.GetStat(types.STAT_HP), 25) + utils.PercentOf(ownerEntity.GetStat(types.STAT_AD), 25)

	if oldEffect != nil {
		ownerEntity.RemoveEffect(ReimiBlessingEffectUUID)
	} else {
		oldEffect = &battle.ActionEffect{
			Effect:   battle.EFFECT_SHIELD,
			Value:    0,
			Duration: -1,
			Uuid:     ReimiBlessingEffectUUID,
			Meta:     nil,
			Caster:   ownerEntity.GetUUID(),
			Source:   types.SOURCE_ITEM,
		}
	}

	if oldEffect.Value < 0 {
		oldEffect.Value = 0
	}

	oldEffect.Value += meta.(battle.ActionEffectHeal).Value

	if oldEffect.Value > maxShield {
		oldEffect.Value = maxShield
	}

	ownerEntity.ApplyEffect(*oldEffect)

	return nil
}

func (rbs ReimiBlessingSkill) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return map[types.CustomTrigger]func(owner interface{}){
		types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
			owner.(battle.Entity).ApplyEffect(battle.ActionEffect{
				Effect:   battle.EFFECT_SHIELD,
				Value:    0,
				Duration: -1,
				Uuid:     ReimiBlessingEffectUUID,
				Meta:     nil,
				Caster:   owner.(battle.Entity).GetUUID(),
				Source:   types.SOURCE_ITEM,
			})
		},
	}
}

type GiantSlayerSkill struct{ BasePassiveSkill }

func (gss GiantSlayerSkill) GetName() string {
	return "Pogromca gigantów"
}

func (gss GiantSlayerSkill) GetDescription() string {
	return "Zadaje dodatkowe obrażenia w zależności od pancerza przeciwnika."
}

func (gss GiantSlayerSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_ATTACK_BEFORE,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (gss GiantSlayerSkill) GetUUID() uuid.UUID {
	return GiantSlayerSkillUUID
}

func (gss GiantSlayerSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	damageValue := utils.PercentOf((target.(battle.Entity).GetStat(types.STAT_DEF)), 10)

	return types.AttackTriggerMeta{Effects: []types.DamagePartial{
		{Value: damageValue, Type: types.DMG_PHYSICAL}},
	}
}

type GiantKillerSkill struct{ BasePassiveSkill }

func (gks GiantKillerSkill) GetName() string {
	return "Zabójca gigantów"
}

func (gks GiantKillerSkill) GetDescription() string {
	return "Zadaje dodatkowe obrażenia w zależności od zdrowia przeciwnika."
}

func (gks GiantKillerSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_ATTACK_BEFORE,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (gks GiantKillerSkill) GetUUID() uuid.UUID {
	return GiantKillerSkillUUID
}

func (gks GiantKillerSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	damageValue := utils.PercentOf((target.(battle.Entity).GetStat(types.STAT_HP)), 2)

	println((target.(battle.Entity).GetStat(types.STAT_HP)))
	println(damageValue)

	return types.AttackTriggerMeta{Effects: []types.DamagePartial{
		{Value: damageValue, Type: types.DMG_PHYSICAL}},
	}
}

type MageKillerSkill struct{ BasePassiveSkill }

func (mks MageKillerSkill) GetName() string {
	return "Zabójca magów"
}

func (mks MageKillerSkill) GetDescription() string {
	return "Atakowanie celi osłoniętych tarczą zwiększa obrażenia twojego ataku."
}

func (mks MageKillerSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_ATTACK_BEFORE,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_HAS_EFFECT},
			Meta:          map[string]interface{}{"effect": battle.EFFECT_SHIELD},
		},
	}
}

func (mks MageKillerSkill) GetUUID() uuid.UUID {
	return MageKillerSkillUUID
}

func (mks MageKillerSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	return types.AttackTriggerMeta{Effects: []types.DamagePartial{
		{Value: 10, Type: types.DMG_PHYSICAL, Percent: true}},
	}
}

type SandBladeSkill struct{ BasePassiveSkill }

func (sbs SandBladeSkill) GetName() string {
	return "Piaskowe ostrze"
}

func (sbs SandBladeSkill) GetDescription() string {
	return "Zadawanie obrażeń zmniejsza leczenie wroga."
}

func (sbs SandBladeSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_DAMAGE,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (sbs SandBladeSkill) GetUUID() uuid.UUID {
	return SandBladeSkillUUID
}

func (sbs SandBladeSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Source: owner.(battle.Entity).GetUUID(),
		Target: target.(battle.Entity).GetUUID(),
		Meta: battle.ActionEffect{
			Effect:   battle.EFFECT_HEAL_REDUCE,
			Value:    20,
			Duration: 1,
			Uuid:     uuid.New(),
			Meta:     nil,
			Caster:   owner.(battle.Entity).GetUUID(),
			Source:   types.SOURCE_ITEM,
		},
	})

	return nil
}

type WaterBladeSkill struct{ BasePassiveSkill }

func (wbs WaterBladeSkill) GetName() string {
	return "Wodne ostrze"
}

func (wbs WaterBladeSkill) GetDescription() string {
	return "Zadawanie obrażeń leczy o brakujące zdrowie."
}

func (wbs WaterBladeSkill) GetCD() int {
	return 10
}

func (wbs WaterBladeSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_ATTACK_HIT,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (wbs WaterBladeSkill) GetUUID() uuid.UUID {
	return WaterBladeSkillUUID
}

func (wbs WaterBladeSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	addPercentage := utils.PercentOf(owner.(battle.Entity).GetStat(types.STAT_AD), 1)

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Source: owner.(battle.Entity).GetUUID(),
		Target: owner.(battle.Entity).GetUUID(),
		Meta: battle.ActionEffectHeal{
			Value: utils.PercentOf(owner.(battle.Entity).GetStat(types.STAT_HP)-owner.(battle.Entity).GetCurrentHP(), 10+addPercentage),
		},
	})

	return nil
}

type DefenseVisageSkill struct{ BasePassiveSkill }

func (dvs DefenseVisageSkill) GetName() string {
	return "Oblicze obrony"
}

func (dvs DefenseVisageSkill) GetDescription() string {
	return "Dostajesz ATK w zależności od maks. HP."
}

func (dvs DefenseVisageSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType: types.TRIGGER_NONE,
		},
	}
}

func (dvs DefenseVisageSkill) GetUUID() uuid.UUID {
	return DefenseVisageSkillUUID
}

func (dvs DefenseVisageSkill) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return map[types.CustomTrigger]func(owner interface{}){
		types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
			owner.(battle.PlayerEntity).AppendDerivedStat(types.DerivedStat{
				Base:    types.STAT_HP,
				Derived: types.STAT_AD,
				Percent: 5,
				Source:  dvs.GetUUID(),
			})
		},
	}
}

type AttackVisageSkill struct{ BasePassiveSkill }

func (avs AttackVisageSkill) GetName() string {
	return "Oblicze ataku"
}

func (avs AttackVisageSkill) GetDescription() string {
	return "Dostajesz HP w zależności od ATK."
}

func (avs AttackVisageSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType: types.TRIGGER_NONE,
		},
	}
}

func (avs AttackVisageSkill) GetUUID() uuid.UUID {
	return AttackVisageSkillUUID
}

func (avs AttackVisageSkill) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return map[types.CustomTrigger]func(owner interface{}){
		types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
			owner.(battle.PlayerEntity).AppendDerivedStat(types.DerivedStat{
				Base:    types.STAT_HP,
				Derived: types.STAT_AD,
				Percent: 10,
				Source:  avs.GetUUID(),
			})
		},
	}
}

type WarriorsLegacySkill struct{ BasePassiveSkill }

func (wls WarriorsLegacySkill) GetName() string {
	return "Dziedzictwo wojownika"
}

func (wls WarriorsLegacySkill) GetDescription() string {
	return "Zwiększa obrażenia w zależności od maks zdrowia."
}

func (wls WarriorsLegacySkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_ATTACK_BEFORE,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (wls WarriorsLegacySkill) GetUUID() uuid.UUID {
	return WarriorsLegacySkillUUID
}

func (wls WarriorsLegacySkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	dmgPercent := utils.PercentOf(owner.(battle.Entity).GetStat(types.STAT_HP_PLUS), 1)

	return types.AttackTriggerMeta{
		Effects: []types.DamagePartial{
			{
				Value:   dmgPercent,
				Type:    types.DMG_PHYSICAL,
				Percent: true,
			},
		},
	}
}

type LilithsWrathSkill struct{ BasePassiveSkill }

func (lws LilithsWrathSkill) GetName() string {
	return "Gniew Lilith"
}

func (lws LilithsWrathSkill) GetDescription() string {
	return "Co ture zadaje obrażenia w zależności od zdrowia użytkownika."
}

func (lws LilithsWrathSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_TURN,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (lws LilithsWrathSkill) GetUUID() uuid.UUID {
	return LilithsWrathSkillUUID
}

func (lws LilithsWrathSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_DMG,
		Source: owner.(battle.Entity).GetUUID(),
		Target: target.(battle.Entity).GetUUID(),
		Meta: []battle.Damage{{
			Value:    utils.PercentOf(owner.(battle.Entity).GetStat(types.STAT_HP), 5),
			Type:     types.DMG_PHYSICAL,
			CanDodge: false,
		}},
	})

	return nil
}

type RyuLegacySkill struct{ BasePassiveSkill }

func (rls RyuLegacySkill) GetName() string {
	return "Dziedzictwo Ryu"
}

func (rls RyuLegacySkill) GetDescription() string {
	return "Zwiększa RES i DEF o 20%."
}

func (rls RyuLegacySkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType: types.TRIGGER_NONE,
		},
	}
}

func (rls RyuLegacySkill) GetUUID() uuid.UUID {
	return RyuLegacySkillUUID
}

func (rls RyuLegacySkill) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return map[types.CustomTrigger]func(owner interface{}){
		types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
			owner.(battle.PlayerEntity).AppendDerivedStat(types.DerivedStat{
				Base:    types.STAT_DEF,
				Derived: types.STAT_DEF,
				Percent: 20,
				Source:  rls.GetUUID(),
			})

			owner.(battle.PlayerEntity).AppendDerivedStat(types.DerivedStat{
				Base:    types.STAT_MR,
				Derived: types.STAT_MR,
				Percent: 20,
				Source:  rls.GetUUID(),
			})
		},
	}
}

type DefenderBladeSkill struct{ BasePassiveSkill }

func (dbs DefenderBladeSkill) GetName() string {
	return "Ostrze obrońcy"
}

func (dbs DefenderBladeSkill) GetDescription() string {
	return "Zwiększa ataki o twój RES i DEF."
}

func (dbs DefenderBladeSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_ATTACK_BEFORE,
			TargetType:    []types.TargetTag{types.TARGET_SELF},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (dbs DefenderBladeSkill) GetUUID() uuid.UUID {
	return DefenderBladeSkillUUID
}

func (dbs DefenderBladeSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	defStat := owner.(battle.Entity).GetStat(types.STAT_DEF)
	mrStat := owner.(battle.Entity).GetStat(types.STAT_MR)

	return types.AttackTriggerMeta{Effects: []types.DamagePartial{
		{
			Value: utils.PercentOf(defStat, 2) + utils.PercentOf(mrStat, 3),
			Type:  types.DMG_PHYSICAL,
		},
	}}
}

type GrudgeArmorSkill struct{ BasePassiveSkill }

func (gas GrudgeArmorSkill) GetName() string {
	return "Pancerz zwady"
}

func (gas GrudgeArmorSkill) GetDescription() string {
	return "Zadaje obrażenia wrogom, którzy cię uderzają i zmniejsza ich leczenie."
}

func (gas GrudgeArmorSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_ATTACK_GOT_HIT,
			TargetType:    []types.TargetTag{types.TARGET_SELF},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (gas GrudgeArmorSkill) GetUUID() uuid.UUID {
	return GrudgeArmorSkillUUID
}

func (gas GrudgeArmorSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_DMG,
		Source: owner.(battle.Entity).GetUUID(),
		Target: target.(battle.Entity).GetUUID(),
		Meta: battle.ActionDamage{
			Damage: []battle.Damage{
				{
					Value:    utils.PercentOf(owner.(battle.Entity).GetStat(types.STAT_DEF), 10),
					Type:     types.DMG_TRUE,
					CanDodge: false,
				},
			},
		},
	})

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Source: owner.(battle.Entity).GetUUID(),
		Target: target.(battle.Entity).GetUUID(),
		Meta: battle.ActionEffect{
			Effect:   battle.EFFECT_HEAL_REDUCE,
			Value:    20,
			Duration: 1,
			Uuid:     uuid.New(),
			Meta:     nil,
			Caster:   owner.(battle.Entity).GetUUID(),
			Source:   types.SOURCE_ITEM,
		},
	})

	return nil
}

type AmplifyingCoatSkill struct{ BasePassiveSkill }

func (acs AmplifyingCoatSkill) GetName() string {
	return "Płaszcz wzmacniający"
}

func (acs AmplifyingCoatSkill) GetDescription() string {
	return "Zwiększa zdrowie o 20%."
}

func (acs AmplifyingCoatSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_NONE,
			TargetType:    []types.TargetTag{types.TARGET_SELF},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (acs AmplifyingCoatSkill) GetUUID() uuid.UUID {
	return AmplifyingCoatSkillUUID
}

func (acs AmplifyingCoatSkill) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return map[types.CustomTrigger]func(owner interface{}){
		types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
			owner.(battle.PlayerEntity).AppendDerivedStat(types.DerivedStat{
				Base:    types.STAT_HP,
				Derived: types.STAT_HP,
				Percent: 20,
				Source:  acs.GetUUID(),
			})
		},
	}
}

type ControllersBraceletSkill struct{ BasePassiveSkill }

func (cbs ControllersBraceletSkill) GetName() string {
	return "Bransoleta kontrolera"
}

func (cbs ControllersBraceletSkill) GetDescription() string {
	return "Leczy po nałożeniu efektu CC."
}

func (cbs ControllersBraceletSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_APPLY_CROWD_CONTROL,
			TargetType:    []types.TargetTag{types.TARGET_ALLY},
			TargetDetails: []types.TargetDetails{types.DETAIL_LOW_HP},
		},
	}
}

func (cbs ControllersBraceletSkill) GetUUID() uuid.UUID {
	return ControllersBraceletSkillUUID
}

func (cbs ControllersBraceletSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	healValue := utils.PercentOf(owner.(battle.Entity).GetStat(types.STAT_AD), 15) + utils.PercentOf(owner.(battle.Entity).GetStat(types.STAT_AP), 15)

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Source: owner.(battle.Entity).GetUUID(),
		Target: target.(battle.Entity).GetUUID(),
		Meta:   battle.ActionEffectHeal{Value: healValue},
	})

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Source: owner.(battle.Entity).GetUUID(),
		Target: owner.(battle.Entity).GetUUID(),
		Meta:   battle.ActionEffectHeal{Value: healValue},
	})

	return nil
}

type CursedIceSkill struct{ BasePassiveSkill }

func (cis CursedIceSkill) GetName() string {
	return "Spaczony lód"
}

func (cis CursedIceSkill) GetDescription() string {
	return "Spowolnienia działają o 20% mocniej."
}

func (cis CursedIceSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_APPLY_CROWD_CONTROL,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (cis CursedIceSkill) GetUUID() uuid.UUID {
	return CursedIceSkillUUID
}

func (cis CursedIceSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	if meta.(battle.ActionEffect).Effect == battle.EFFECT_STAT_DEC {
		if meta.(battle.ActionEffect).Meta.(battle.ActionEffectStat).Stat == types.STAT_SPD {
			return types.EffectTriggerMeta{
				Effects: []types.IncreasePartial{
					{
						Value:   20,
						Percent: true,
					},
				}}
		}
	}

	return nil
}

type ControllersRuneSkill struct{ BasePassiveSkill }

func (crs ControllersRuneSkill) GetName() string {
	return "Runa kontrolera"
}

func (crs ControllersRuneSkill) GetDescription() string {
	return "Zabicie wroga objętego CC przywraca ci punkt many."
}

func (crs ControllersRuneSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_EXECUTE,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_HAS_EFFECT},
		},
	}
}

func (crs ControllersRuneSkill) GetUUID() uuid.UUID {
	return ControllersRuneSkillUUID
}

func (crs ControllersRuneSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Source: owner.(battle.Entity).GetUUID(),
		Target: owner.(battle.Entity).GetUUID(),
		Meta: battle.ActionEffect{
			Effect:   battle.EFFECT_MANA_RESTORE,
			Value:    1,
			Duration: 0,
			Uuid:     ControllersRuneEffectUUID,
			Meta:     nil,
			Caster:   owner.(battle.Entity).GetUUID(),
			Source:   types.SOURCE_ITEM,
		},
	})

	return nil
}

type ControllersNecklaceSkill struct{ BasePassiveSkill }

func (cns ControllersNecklaceSkill) GetName() string {
	return "Naszyjnik kontrolera"
}

func (cns ControllersNecklaceSkill) GetDescription() string {
	return "Nałożenie efektu CC przyśpiesza cię."
}

func (cns ControllersNecklaceSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_APPLY_CROWD_CONTROL,
			TargetType:    []types.TargetTag{types.TARGET_SELF},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (cns ControllersNecklaceSkill) GetUUID() uuid.UUID {
	return ControllersNecklaceSkillUUID
}

func (cns ControllersNecklaceSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Source: owner.(battle.Entity).GetUUID(),
		Target: owner.(battle.Entity).GetUUID(),
		Meta: battle.ActionEffect{
			Effect:   battle.EFFECT_STAT_INC,
			Value:    0,
			Duration: 1,
			Uuid:     ControllersNecklaceEffectUUID,
			Meta: battle.ActionEffectStat{
				Stat:      types.STAT_SPD,
				Value:     10,
				IsPercent: false,
			},
			Caster: owner.(battle.Entity).GetUUID(),
			Source: types.SOURCE_ITEM,
		},
	})

	return nil
}

type ControllersBladeSkill struct{ BasePassiveSkill }

func (cbs ControllersBladeSkill) GetName() string {
	return "Ostrze kontrolera"
}

func (cbs ControllersBladeSkill) GetDescription() string {
	return "Zadawanie obrażeń spowalnia wroga."
}

func (cbs ControllersBladeSkill) GetCD() int {
	return 3
}

func (cbs ControllersBladeSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_ATTACK_HIT,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
		Cooldown: &types.CooldownMeta{
			PassEvent: types.TRIGGER_ATTACK_HIT,
		},
	}
}

func (cbs ControllersBladeSkill) GetUUID() uuid.UUID {
	return ControllersBladeSkillUUID
}

func (cbs ControllersBladeSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Source: owner.(battle.Entity).GetUUID(),
		Target: target.(battle.Entity).GetUUID(),
		Meta: battle.ActionEffect{
			Effect:   battle.EFFECT_STAT_INC,
			Value:    0,
			Duration: 1,
			Uuid:     ControllersBladeEffectUUID,
			Meta: battle.ActionEffectStat{
				Stat:      types.STAT_SPD,
				Value:     -10,
				IsPercent: false,
			},
			Caster: owner.(battle.Entity).GetUUID(),
			Source: types.SOURCE_ITEM,
		},
	})

	return nil
}

type ControllersHatSkill struct{ BasePassiveSkill }

func (chs ControllersHatSkill) GetName() string {
	return "Kapelusz kontrolera"
}

func (chs ControllersHatSkill) GetDescription() string {
	return "Daje siłe adaptacyjną w zależności od dodatkowej many"
}

func (chs ControllersHatSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType: types.TRIGGER_NONE,
		},
	}
}

func (chs ControllersHatSkill) GetUUID() uuid.UUID {
	return ControllersHatSkillUUID
}

func (chs ControllersHatSkill) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return map[types.CustomTrigger]func(owner interface{}){
		types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
			owner.(battle.PlayerEntity).AppendDerivedStat(types.DerivedStat{
				Base:    types.STAT_MANA_PLUS,
				Derived: types.STAT_ADAPTIVE,
				Percent: 100,
				Source:  chs.GetUUID(),
			})
		},
	}
}

type ArdentCenserSkill struct{ BasePassiveSkill }

func (acs ArdentCenserSkill) GetName() string {
	return "Ognisty trybularz"
}

func (acs ArdentCenserSkill) GetDescription() string {
	return "Leczenie i tarcze przyśpieszają sojuszników i zwiększają ich obrażenia."
}

func (acs ArdentCenserSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_HEAL_OTHER,
			TargetType:    []types.TargetTag{types.TARGET_ALLY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (acs ArdentCenserSkill) GetUUID() uuid.UUID {
	return ArdentCenserSkillUUID
}

func (acs ArdentCenserSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	target.(battle.PlayerEntity).AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value: types.BaseAttackIncreaseSkill{
			Calculate: func(meta types.AttackTriggerMeta) types.AttackTriggerMeta {
				return types.AttackTriggerMeta{
					Effects: []types.DamagePartial{
						{
							Value:   utils.PercentOf(owner.(battle.Entity).GetStat(types.STAT_AP), 25),
							Percent: false,
							Type:    1,
						},
					},
				}
			},
		},
	})

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Source: owner.(battle.Entity).GetUUID(),
		Target: target.(battle.Entity).GetUUID(),
		Meta: battle.ActionEffect{
			Effect:   battle.EFFECT_STAT_INC,
			Value:    0,
			Duration: 1,
			Uuid:     ArdentCenserEffectUUID,
			Meta: battle.ActionEffectStat{
				Stat:      types.STAT_SPD,
				Value:     10,
				IsPercent: false,
			},
			Caster: owner.(battle.Entity).GetUUID(),
			Source: types.SOURCE_ITEM,
		},
	})

	return nil
}

type SirensCallSkill struct{ BasePassiveSkill }

func (scs SirensCallSkill) GetName() string {
	return "Syreni śpiew"
}

func (scs SirensCallSkill) GetDescription() string {
	return "Część leczenia i tarcz przeskakuje na innego sojusznika."
}

func (scs SirensCallSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_HEAL_OTHER,
			TargetType:    []types.TargetTag{types.TARGET_ALLY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (scs SirensCallSkill) GetUUID() uuid.UUID {
	return SirensCallSkillUUID
}

func (scs SirensCallSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	healValue := utils.PercentOf(meta.(battle.ActionEffectHeal).Value, 10)
	healTarget := target.(battle.Entity)

	if target.(battle.Entity).GetUUID() == owner.(battle.Entity).GetUUID() {
		validTargets := fightInstance.(*battle.Fight).GetAlliesFor(owner.(battle.Entity).GetUUID())

		if len(validTargets) == 0 {
			return nil
		}

		sortInit := battle.EntitySort{
			Entities: validTargets,
			Order:    []types.TargetDetails{types.DETAIL_LOW_HP},
			Meta:     nil,
		}

		sort.Sort(sortInit)

		healTarget = sortInit.Entities[0]
	}

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Source: owner.(battle.Entity).GetUUID(),
		Target: healTarget.GetUUID(),
		Meta: battle.ActionEffectHeal{
			Value: healValue,
		},
	})

	return nil
}

type FogsEmpowermentSkill struct{ BasePassiveSkill }

func (fes FogsEmpowermentSkill) GetName() string {
	return "Mgliste wzmocnienie"
}

func (fes FogsEmpowermentSkill) GetDescription() string {
	return "Otrzymujesz AP w zależności od siły leczenia i tarcz."
}

func (fes FogsEmpowermentSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_NONE,
			TargetType:    []types.TargetTag{types.TARGET_SELF},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (fes FogsEmpowermentSkill) GetUUID() uuid.UUID {
	return FogsEmpowermentSkillUUID
}

func (fes FogsEmpowermentSkill) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return map[types.CustomTrigger]func(owner interface{}){
		types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
			owner.(battle.PlayerEntity).AppendDerivedStat(types.DerivedStat{
				Base:    types.STAT_HEAL_POWER,
				Derived: types.STAT_AP,
				Percent: 1000,
				Source:  fes.GetUUID(),
			})
		},
	}
}

type WindsEmpowermentSkill struct{ BasePassiveSkill }

func (wes WindsEmpowermentSkill) GetName() string {
	return "Wietrzne wzmocnienie"
}

func (wes WindsEmpowermentSkill) GetDescription() string {
	return "Otrzymujesz SPD w zależności od siły leczenia i tarcz. Oraz leczysz przy ataku"
}

func (wes WindsEmpowermentSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_ATTACK_HIT,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (wes WindsEmpowermentSkill) GetUUID() uuid.UUID {
	return WindsEmpowermentSkillUUID
}

func (wes WindsEmpowermentSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	healValue := utils.PercentOf(owner.(battle.Entity).GetStat(types.STAT_AD), 10)

	validTargets := fightInstance.(*battle.Fight).GetAlliesFor(owner.(battle.Entity).GetUUID())

	if len(validTargets) == 0 {
		return nil
	}

	sortInit := battle.EntitySort{
		Entities: validTargets,
		Order:    []types.TargetDetails{types.DETAIL_LOW_HP},
		Meta:     nil,
	}

	sort.Sort(sortInit)

	healTarget := sortInit.Entities[0]

	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Source: owner.(battle.Entity).GetUUID(),
		Target: healTarget.GetUUID(),
		Meta:   battle.ActionEffectHeal{Value: healValue},
	})

	return nil
}

func (wes WindsEmpowermentSkill) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return map[types.CustomTrigger]func(owner interface{}){
		types.CUSTOM_TRIGGER_UNLOCK: func(owner interface{}) {
			owner.(battle.PlayerEntity).AppendDerivedStat(types.DerivedStat{
				Base:    types.STAT_HEAL_POWER,
				Derived: types.STAT_SPD,
				Percent: 100,
				Source:  wes.GetUUID(),
			})
		},
	}
}

type LightingSupportSkill struct{ BasePassiveSkill }

func (lss LightingSupportSkill) GetName() string {
	return "Błyskawiczne wsparcie"
}

func (lss LightingSupportSkill) GetDescription() string {
	return "Po nałożeniu CC kolejne obrażenia sojusznika są zwiększone."
}

func (lss LightingSupportSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_APPLY_CROWD_CONTROL,
			TargetType:    []types.TargetTag{types.TARGET_ALLY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (lss LightingSupportSkill) GetUUID() uuid.UUID {
	return LightingSupportSkillUUID
}

func (lss LightingSupportSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	target.(battle.PlayerEntity).AppendTempSkill(types.WithExpire[types.PlayerSkill]{
		Value: types.BaseAttackIncreaseSkill{
			Calculate: func(meta types.AttackTriggerMeta) types.AttackTriggerMeta {
				return types.AttackTriggerMeta{
					Effects: []types.DamagePartial{
						{
							Value:   utils.PercentOf(owner.(battle.Entity).GetStat(types.STAT_AP), 25),
							Percent: false,
							Type:    1,
						},
					},
				}
			},
		},
	})

	return nil
}

type GodlySupportSkill struct{ BasePassiveSkill }

func (gss GodlySupportSkill) GetName() string {
	return "Niebiańskie wsparcie"
}

func (gss GodlySupportSkill) GetDescription() string {
	return "Po użyciu super umiejętności leczysz wszystkich."
}

func (gss GodlySupportSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_CAST_ULT,
			TargetType:    []types.TargetTag{types.TARGET_ALLY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (gss GodlySupportSkill) GetUUID() uuid.UUID {
	return GodlySupportSkillUUID
}

func (gss GodlySupportSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	fightInstance.(*battle.Fight).HandleAction(battle.Action{
		Event:  battle.ACTION_EFFECT,
		Source: owner.(battle.Entity).GetUUID(),
		Target: target.(battle.Entity).GetUUID(),
		Meta:   battle.ActionEffectHeal{Value: utils.PercentOf(owner.(battle.Entity).GetStat(types.STAT_AP), 2)},
	})

	return nil
}

type KyokiBeltSkill struct{ BasePassiveSkill }

func (kbs KyokiBeltSkill) GetName() string {
	return "Pasek Kyoki"
}

func (kbs KyokiBeltSkill) GetDescription() string {
	return "Obrażenia magiczne są zwiększone przez losowy mnożnik (0.8-1.8)."
}

func (kbs KyokiBeltSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_DAMAGE_BEFORE,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (kbs KyokiBeltSkill) GetUUID() uuid.UUID {
	return KyokiBeltSkillUUID
}

func (kbs KyokiBeltSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	return types.DamageTriggerMeta{
		Effects: []types.DamagePartial{
			{
				Value:   utils.RandomNumber(0, 100) - 20,
				Type:    types.DMG_MAGICAL,
				Percent: true,
			},
		},
	}
}

type ShikiFlameSkill struct{ BasePassiveSkill }

func (sfs ShikiFlameSkill) GetName() string {
	return "Płomień Shiki"
}

func (sfs ShikiFlameSkill) GetDescription() string {
	return "Zaklęcia zadają obrażenia w zależności od zdrowia wroga"
}

func (sfs ShikiFlameSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_DAMAGE_BEFORE,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (sfs ShikiFlameSkill) GetUUID() uuid.UUID {
	return ShikiFlameSkillUUID
}

type StormHarbingerSkill struct{ BasePassiveSkill }

func (shs StormHarbingerSkill) GetName() string {
	return "Zwiastun burzy"
}

func (shs StormHarbingerSkill) GetDescription() string {
	return "Zwiększa obrażenia przy ataku w zależności od AP."
}

func (shs StormHarbingerSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_ATTACK_BEFORE,
			TargetType:    []types.TargetTag{types.TARGET_ENEMY},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
		},
	}
}

func (shs StormHarbingerSkill) GetUUID() uuid.UUID {
	return StormHarbingerSkillUUID
}

func (shs StormHarbingerSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	return types.AttackTriggerMeta{
		Effects: []types.DamagePartial{
			{
				Value:   utils.PercentOf(owner.(battle.Entity).GetStat(types.STAT_AP), 20),
				Type:    types.DMG_MAGICAL,
				Percent: false,
			},
		},
	}
}
