package discord

import (
	"context"
	"fmt"
	"sao/player/inventory"
	"sao/world"
	"strings"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

var PathToString = map[inventory.SkillPath]string{
	inventory.PathControl:   "Kontrola",
	inventory.PathDamage:    "Obrażenia",
	inventory.PathEndurance: "Wytrzymałość",
	inventory.PathMobility:  "Mobilność",
}

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
			},
		},
		discord.SlashCommandCreate{
			Name:        "info",
			Description: "Informacje o postaci",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "gracz",
					Description: "Gracz",
					Required:    false,
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "skill",
			Description: "Zarządzaj umiejętnościami",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "pokaż",
					Description: "Pokaż umiejętności",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionUser{
							Name:        "gracz",
							Description: "Gracz",
							Required:    false,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "odblokuj",
					Description: "Odblokuj umiejętność",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionInt{
							Name:        "lvl",
							Description: "Umiejętność którego chcesz odblokować",
							Required:    true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "ulepsz",
					Description: "Ulepsz umiejętność",
				},
			},
		},
		discord.SlashCommandCreate{
			Name:        "plecak",
			Description: "Zarządzaj ekwipunkiem",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "pokaż",
					Description: "Pokaż ekwipunek",
				},
			},
		},
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
		bot.WithEventListenerFunc(messageComponentHandler),
	)

	if err != nil {
		panic(err)
	}

	if _, err = client.Rest().SetGuildCommands(client.ApplicationID(), snowflake.MustParse("1151589368373444690"), cmds); err != nil {
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

		userSnowflake := event.Member().User.ID.String()

		var floorName string

		for _, pl := range World.Players {
			if pl.Meta.UserID == userSnowflake {
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

func messageComponentHandler(event *events.ComponentInteractionCreate) {
	if strings.HasPrefix(event.ComponentInteraction.Data.CustomID(), "su") {
		data := strings.Split(event.ComponentInteraction.Data.CustomID(), "|")

		path := inventory.SkillPath(0)

		switch data[1] {
		case "0":
			path = inventory.PathControl
		case "1":
			path = inventory.PathDamage
		case "2":
			path = inventory.PathEndurance
		case "3":
			path = inventory.PathMobility
		}

		lvl := 0

		fmt.Sscanf(data[2], "%d", &lvl)

		userSnowflake := event.Member().User.ID.String()

		for _, pl := range World.Players {
			if pl.Meta.UserID == userSnowflake {
				res := pl.Inventory.UnlockSkill(path, lvl, pl.XP.Level, pl)

				if res == nil {
					event.UpdateMessage(
						discord.
							NewMessageUpdateBuilder().
							SetContent("Odblokowano umiejętność").
							ClearContainerComponents().
							ClearEmbeds().
							Build(),
					)

					return
				}

				msgContent := ""

				switch res.Error() {
				case "PLAYER_LVL_TOO_LOW":
					msgContent = "Nie masz wystarczającego poziomu"
				case "SKILL_ALREADY_UNLOCKED":
					msgContent = "Umiejętność jest już odblokowana"
				case "SKILL_NOT_FOUND":
					msgContent = "Nie znaleziono umiejętności"
				}

				event.CreateMessage(
					discord.
						NewMessageCreateBuilder().
						SetContent(msgContent).
						Build(),
				)
				return
			}
		}
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
			if snowflake.MustParse(pl.Meta.UserID) == event.Member().User.ID {
				World.MovePlayer(uuid, pl.Meta.Location.FloorName, locationName, "")
				return
			}
		}
	case "tp":
		floorName := data.String("nazwa")

		for uuid, pl := range World.Players {
			if snowflake.MustParse(pl.Meta.UserID) == event.Member().User.ID {
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
	case "info":

		user := event.User()

		if mentionedUser, exists := data.OptUser("gracz"); exists {
			user = mentionedUser
		}

		for _, pl := range World.Players {
			if pl.Meta.UserID == user.ID.String() {

				inFight := pl.Meta.FightInstance != nil

				inFightText := "Nie"

				if inFight {
					inFightText = "Tak"
				}

				event.CreateMessage(
					discord.NewMessageCreateBuilder().
						AddEmbeds(
							discord.NewEmbedBuilder().
								AddField("Nazwa", pl.GetName(), true).
								AddField("Gracz", fmt.Sprintf("<@%s>", pl.Meta.UserID), true).
								AddField("Lokacja", pl.Meta.Location.LocationName, true).
								AddField("Piętro", pl.Meta.Location.FloorName, true).
								AddField("HP", fmt.Sprintf("%d/%d", pl.GetCurrentHP(), pl.GetMaxHP()), true).
								AddField("Mana", fmt.Sprintf("%d/%d", pl.GetCurrentMana(), pl.GetMaxMana()), true).
								AddField("Atak", fmt.Sprintf("%d", pl.GetATK()), true).
								AddField("AP", fmt.Sprintf("%d", pl.GetAP()), true).
								AddField("DEF/RES", fmt.Sprintf("%d/%d", pl.GetDEF(), pl.GetMR()), true).
								AddField("Lvl", fmt.Sprintf("%d %d/%d", pl.XP.Level, pl.XP.Exp, (pl.XP.Level*100)+100), true).
								AddField("SPD/AGL", fmt.Sprintf("%d/%d", pl.GetSPD(), pl.GetAGL()), true).
								AddField("Walce?", inFightText, true).
								AddField("Party?", "Nie", true).
								AddField("Gildia?", "Nie", true).
								Build(),
						).
						Build(),
				)
				return
			}
		}

		event.CreateMessage(
			discord.
				NewMessageCreateBuilder().
				SetContent("Nie znaleziono postaci").
				SetEphemeral(true).
				Build(),
		)
	case "skill":
		switch *data.SubCommandName {
		case "pokaż":
			user := event.User()

			if mentionedUser, exists := data.OptUser("gracz"); exists {
				user = mentionedUser
			}

			for _, pl := range World.Players {
				if pl.Meta.UserID == user.ID.String() {
					embed := discord.NewEmbedBuilder()

					if len(pl.Inventory.Skills) == 0 {
						embed.AddField("Skille", "Brak", false)
					}

					//TODO add inv.skills

					if len(pl.Inventory.LevelSkills) == 0 {
						embed.AddField("Skille za lvl", "Brak", false)
					}

					for _, skill := range pl.Inventory.LevelSkills {
						embed.AddField(
							fmt.Sprintf("%s (LVL: %d)", skill.Name, skill.ForLevel),
							fmt.Sprintf("Ścieżka: %s\n\n%s", PathToString[skill.Path], skill.Description),
							false,
						)
					}

					event.CreateMessage(
						discord.
							NewMessageCreateBuilder().
							AddEmbeds(embed.Build()).
							Build(),
					)

					return
				}
			}
		case "odblokuj":
			user := event.User()

			for _, pl := range World.Players {
				if pl.Meta.UserID == user.ID.String() {

					lvl := data.Int("lvl")

					skillList := make([]inventory.PSkill, 0)

					for _, skill := range inventory.AVAILABLE_SKILLS {
						if skill.ForLevel == lvl {
							skillList = append(skillList, skill)
						}
					}

					if len(skillList) == 0 {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Nie znaleziono umiejętności dla tego poziomu").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					embed := discord.NewEmbedBuilder()

					buttons := make([]discord.InteractiveComponent, 0)

					for _, skill := range skillList {

						upgradeMsg := ""

						for _, upgrade := range skill.Upgrades {
							upgradeMsg += fmt.Sprintf("\n- %s - %s", upgrade.Name, upgrade.Description)
						}

						embed.AddField(
							skill.Name,
							fmt.Sprintf("Ścieżka: %s\n\n%s\nUlepszenia:%s", PathToString[skill.Path], skill.Description, upgradeMsg),
							false,
						)

						buttons = append(buttons, discord.NewPrimaryButton(
							skill.Name,
							fmt.Sprintf("su|%d|%d", skill.Path, skill.ForLevel),
						))
					}

					event.CreateMessage(
						discord.
							NewMessageCreateBuilder().
							AddEmbeds(embed.Build()).
							AddActionRow(buttons...).
							Build(),
					)

					return
				}
			}
		}

	case "plecak":
		switch *data.SubCommandName {
		case "pokaż":
			user := event.User()

			for _, pl := range World.Players {
				if pl.Meta.UserID == user.ID.String() {
					embed := discord.NewEmbedBuilder()

					if len(pl.Inventory.Items) == 0 {
						embed.AddField("Przedmioty", fmt.Sprintf("%d/%d", 0, pl.Inventory.Capacity), false)
					} else {
						count := 0

						for _, item := range pl.Inventory.Items {
							if !item.Hidden {
								count++
							}
						}

						embed.AddField("Przedmioty", fmt.Sprintf("%d/%d", count, pl.Inventory.Capacity), false)
					}

					for _, item := range pl.Inventory.Items {

						if item.Hidden {
							continue
						}

						embed.AddField(item.Name, item.Description, false)
					}

					event.CreateMessage(
						discord.
							NewMessageCreateBuilder().
							AddEmbeds(embed.Build()).
							Build(),
					)

					return
				}
			}
		}
	}
}
