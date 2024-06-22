package types

import "github.com/disgoorg/disgo/discord"

type DiscordMessageStruct struct {
	ChannelID      string
	MessageContent discord.MessageCreate
	DM             bool
}

type Stat int

const (
	STAT_HP Stat = iota
	STAT_SPD
	STAT_AGL
	STAT_AD
	STAT_DEF
	STAT_MR
	STAT_MANA
	STAT_AP
	STAT_HEAL_SELF
	STAT_HEAL_POWER
	STAT_LETHAL
	STAT_LETHAL_PERCENT
	STAT_MAGIC_PEN
	STAT_MAGIC_PEN_PERCENT
	STAT_ADAPTIVE
	STAT_VAMP
	STAT_LIFESTEAL
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

type AdaptiveType int

const (
	ADAPTIVE_ATK AdaptiveType = iota
	ADAPTIVE_AP
)

type AdaptiveRes int

const (
	ADAPTIVE_DEF AdaptiveRes = iota
	ADAPTIVE_RES
)
