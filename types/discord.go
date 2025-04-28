package types

import "github.com/disgoorg/disgo/discord"

type DiscordMessageStruct struct {
	ChannelID      string
	MessageContent discord.MessageCreate
	DM             bool
}

type DiscordEvent interface {
	GetEvent() DiscordMessage
	GetData() any
}

type DiscordMessage int

const (
	MSG_SEND DiscordMessage = iota
	MSG_CHOICE
)

type DiscordSendMsg struct {
	Data DiscordMessageStruct
}

func (fsm DiscordSendMsg) GetEvent() DiscordMessage {
	return MSG_SEND
}

func (fsm DiscordSendMsg) GetData() any {
	return fsm.Data
}

type DiscordChoiceMsg struct {
	Data DiscordChoice
}

func (fsm DiscordChoiceMsg) GetEvent() DiscordMessage {
	return MSG_CHOICE
}

func (fsm DiscordChoiceMsg) GetData() any {
	return fsm.Data
}