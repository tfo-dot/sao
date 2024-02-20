package discord

import "github.com/disgoorg/disgo/discord"

var noCharMessage = discord.
	NewMessageCreateBuilder().
	SetContent("Nie znaleziono twojej postaci...").
	SetEphemeral(true).
	Build()

var fightNotYoursMessage = discord.
	NewMessageCreateBuilder().
	SetContent("To nie twoja walka...").
	SetEphemeral(true).
	Build()

var fightAlreadyEndedMessage = discord.
	NewMessageCreateBuilder().
	SetContent("Walka zakończona").
	SetEphemeral(true).
	Build()

var notYourTurnMessage = discord.
	NewMessageCreateBuilder().
	SetContent("Nie twoja tura!").
	SetEphemeral(true).
	Build()

var messageUpdateClearComponents = discord.NewMessageUpdateBuilder().ClearContainerComponents().Build()
