package discord

import (
	"sao/data"
	"sao/world/tournament"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func AutocompleteHandler(event *events.AutocompleteInteractionCreate) {
	switch event.Data.CommandName {
	case "ruch":
		locationOption := event.Data.String("nazwa")

		userSnowflake := event.Member().User.ID

		var floorName string

		for _, pl := range World.Players {
			if pl.Meta.UserID == userSnowflake.String() {
				floorName = pl.Meta.Location.FloorName
			}
		}

		choices := make([]discord.AutocompleteChoice, 0)

		for _, location := range World.Floors[floorName].Locations {
			if !location.Unlocked {
				continue
			}

			if strings.HasPrefix(location.Name, locationOption) {
				choices = append(choices, discord.AutocompleteChoiceString{
					Name:  location.Name,
					Value: location.Name,
				})
			}
		}

		event.AutocompleteResult(choices)

	case "tp":
		locationOption := event.Data.String("nazwa")

		userSnowflake := event.Member().User.ID.String()

		var floorName string

		for _, pl := range World.Players {
			if pl.Meta.UserID == userSnowflake {
				floorName = pl.Meta.Location.FloorName
			}
		}

		choices := make([]discord.AutocompleteChoice, 0)

		for name, floorData := range World.Floors {

			if !floorData.Unlocked || name == floorName {
				continue
			}

			if strings.HasPrefix(name, locationOption) {
				choices = append(choices, discord.AutocompleteChoiceString{
					Name:  name,
					Value: name,
				})
			}
		}

		event.AutocompleteResult(choices)
	case "stwórz":
		itemOption := event.Data.String("nazwa")

		choices := make([]discord.AutocompleteChoice, 0)

		for _, item := range data.Recipes {
			if strings.HasPrefix(item.Name, itemOption) {
				choices = append(choices, discord.AutocompleteChoiceString{
					Name:  item.Name,
					Value: item.UUID.String(),
				})
			}
		}

		if len(choices) > 25 {
			event.AutocompleteResult(choices[:25])
		} else {
			event.AutocompleteResult(choices)
		}
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
