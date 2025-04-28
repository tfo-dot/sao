package discord

import (
	"sao/world/tournament"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func AutocompleteHandler(event *events.AutocompleteInteractionCreate) {
	switch event.Data.CommandName {
	case "turniej":
		name := event.Data.String("nazwa")

		choices := make([]discord.AutocompleteChoice, 0)

		for _, tournamentObj := range World.Tournaments {
			if strings.HasPrefix(tournamentObj.Name, name) && tournamentObj.State == tournament.Waiting {
				choices = append(choices, discord.AutocompleteChoiceString{
					Name:  tournamentObj.Name,
					Value: tournamentObj.Name,
				})
			}
		}

		event.AutocompleteResult(choices)
	}
}
