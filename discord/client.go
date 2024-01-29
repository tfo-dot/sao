package discord

import (
	"context"
	"fmt"
	"sao/world"
	"strings"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

var (
	cmds = []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "create",
			Description: "Stwórz postać",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "nazwa",
					Description: "Nazwa postaci",
					Required:    true,
				},
				discord.ApplicationCommandOptionUser{
					Name:        "gracz",
					Description: "Gracz",
					Required:    true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "ruch",
			Description: "Przenieś się do innej lokacji",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Autocomplete: true,
					Name:         "nazwa",
					Description:  "Nazwa lokacji",
					Required:     true,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "tp",
			Description: "Teleportuj się na inne piętro",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:         "nazwa",
					Description:  "Nazwa piętra",
					Required:     true,
					Autocomplete: true,
				},
			}},
	}
)

var World *world.World

func StartClient(token string) {
	client, err := disgo.New(token,
		bot.WithDefaultGateway(),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			e.Client().SetPresence(context.Background(), gateway.WithWatchingActivity("SAO"))
		}),
		bot.WithEventListenerFunc(commandListener),
		bot.WithEventListenerFunc(autocompleteHandler),
	)

	if err != nil {
		panic(err)
	}

	gid, err := snowflake.Parse("1151589368373444690")

	if err != nil {
		panic(err)
	}

	if _, err = client.Rest().SetGuildCommands(client.ApplicationID(), gid, cmds); err != nil {
		fmt.Println("error while registering commands")
	}

	if err = client.OpenGateway(context.Background()); err != nil {
		panic(err)
	}
}

func autocompleteHandler(event *events.AutocompleteInteractionCreate) {
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

		userSnowflake := event.Member().User.ID

		var floorName string

		for _, pl := range World.Players {
			if pl.Meta.UserID == userSnowflake.String() {
				floorName = pl.Meta.Location.FloorName
			}
		}

		choices := make([]discord.AutocompleteChoice, 0)

		for name := range World.Floors {
			if name == floorName {
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
	}
}

func commandListener(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()
	switch data.CommandName() {
	case "create":
		if !event.Member().Permissions.Has(discord.PermissionAdministrator) {
			event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Nie masz uprawnień do tej komendy").SetEphemeral(true).Build())

			return
		}

		charName := data.String("nazwa")
		charUser := data.User("gracz")

		World.RegisterNewPlayer(charName, charUser.ID.String())

		err := event.CreateMessage(discord.NewMessageCreateBuilder().
			SetContent("Zarejestrowano postać " + charName).
			Build(),
		)

		if err != nil {
			fmt.Println("error while sending message")
			return
		}
	case "ruch":
		locationName := data.String("nazwa")

		for uuid, pl := range World.Players {
			snowflake, _ := snowflake.Parse(pl.Meta.UserID)

			if snowflake == event.Member().User.ID {
				World.MovePlayer(uuid, pl.Meta.Location.FloorName, locationName, "")
				return
			}
		}
	case "tp":
		floorName := data.String("nazwa")

		for uuid, pl := range World.Players {
			snowflake, _ := snowflake.Parse(pl.Meta.UserID)

			if snowflake == event.Member().User.ID {
				currentLocation := World.Floors[pl.Meta.Location.FloorName].FindLocation(pl.Meta.Location.LocationName)

				if !currentLocation.TP {
					event.CreateMessage(
						discord.
							NewMessageCreateBuilder().
							SetContent("Nie możesz się stąd teleportować (idź do miasta lub lokacji z tp)").
							SetEphemeral(true).
							Build(),
					)
					return
				}

				World.MovePlayer(uuid, floorName, World.Floors[floorName].Default, "")
				return
			}
		}
	}
}
