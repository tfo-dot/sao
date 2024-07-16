package types

import "github.com/disgoorg/disgo/discord"

type DiscordMessageStruct struct {
	ChannelID      string
	MessageContent discord.MessageCreate
	DM             bool
}

type Stat int

const (
	STAT_NONE Stat = iota
	STAT_HP
	STAT_HP_PLUS
	STAT_SPD
	STAT_AGL
	STAT_AD
	STAT_DEF
	STAT_MR
	STAT_MANA
	STAT_MANA_PLUS
	STAT_AP
	STAT_HEAL_SELF
	STAT_HEAL_POWER
	STAT_LETHAL
	STAT_LETHAL_PERCENT
	STAT_MAGIC_PEN
	STAT_MAGIC_PEN_PERCENT
	STAT_ADAPTIVE
	STAT_ADAPTIVE_PERCENT
	STAT_OMNI_VAMP
	STAT_ATK_VAMP
)

var StatToString = map[Stat]string{
	STAT_HP:                "HP",
	STAT_SPD:               "SPD",
	STAT_AGL:               "AGL",
	STAT_AD:                "AD",
	STAT_DEF:               "DEF",
	STAT_MR:                "RES",
	STAT_MANA:              "MANA",
	STAT_AP:                "AP",
	STAT_HEAL_SELF:         "Otrzymane leczenie",
	STAT_HEAL_POWER:        "Moc leczenia",
	STAT_LETHAL:            "Przebicie pancerza",
	STAT_LETHAL_PERCENT:    "Przebicie procentowe pancerza",
	STAT_MAGIC_PEN:         "Przebicie odporności na magię",
	STAT_MAGIC_PEN_PERCENT: "Przebicie procentowe odporności na magię",
	STAT_ADAPTIVE:          "SIła adaptacyjna",
}

type AdaptiveAttackType int

const (
	ADAPTIVE_ATK AdaptiveAttackType = iota
	ADAPTIVE_AP
)

type AdaptiveDefenseType int

const (
	ADAPTIVE_DEF AdaptiveDefenseType = iota
	ADAPTIVE_RES
)

type EffectSource int

const (
	SOURCE_ND EffectSource = iota
	SOURCE_PARTY
	SOURCE_LOCATION
	SOURCE_ITEM
)

type EntityFlag int

const (
	ENTITY_AUTO EntityFlag = 1 << iota
	ENTITY_SUMMON
)

var PathToString = map[SkillPath]string{
	PathControl:   "Kontrola",
	PathDamage:    "Obrażenia",
	PathEndurance: "Wytrzymałość",
	PathSpecial:   "Specjalista",
}

type DamageType int

const (
	DMG_PHYSICAL DamageType = iota
	DMG_MAGICAL
	DMG_TRUE
)
