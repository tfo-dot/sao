package discord

import "github.com/disgoorg/disgo/discord"

var noCharMessage = discord.
	NewMessageCreateBuilder().
	SetContent("Nie znaleziono postaci...").
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

var transactionArleadyAccepted = discord.
	NewMessageCreateBuilder().
	SetContent("Transakcja już została zaakceptowana").
	SetEphemeral(true).
	Build()

var noTransaction = discord.
	NewMessageCreateBuilder().
	SetContent("Transakcja nie istnieje").
	SetEphemeral(true).
	Build()

var alreadyInParty = discord.
	NewMessageCreateBuilder().
	SetContent("Jesteś już w party").
	SetEphemeral(true).
	Build()

var unknownError = discord.
	NewMessageCreateBuilder().
	SetContent("Wystapił nieoczekiwany błąd").
	SetEphemeral(true).
	Build()

var messageUpdateClearComponents = discord.NewMessageUpdateBuilder().ClearContainerComponents().Build()
