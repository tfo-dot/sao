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
	STAT_HEAL_POWER
	STAT_ADAPTIVE
)

var StatToString = map[Stat]string{
	STAT_HP:         "HP",
	STAT_SPD:        "SPD",
	STAT_AGL:        "AGL",
	STAT_AD:         "AD",
	STAT_DEF:        "DEF",
	STAT_MR:         "RES",
	STAT_MANA:       "MANA",
	STAT_AP:         "AP",
	STAT_HEAL_POWER: "Moc leczenia",
	STAT_ADAPTIVE:   "SIła adaptacyjna",
}

type AdaptiveType int

const (
	ADAPTIVE_ATK AdaptiveType = iota
	ADAPTIVE_AP
)
