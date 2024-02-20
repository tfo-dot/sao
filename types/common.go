package types

import "github.com/disgoorg/disgo/discord"

type DiscordMessageStruct struct {
	ChannelID      string
	MessageContent discord.MessageCreate
}
