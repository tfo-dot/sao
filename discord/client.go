package discord

import (
	"context"
	"fmt"
	"sao/player/inventory"
	"sao/world"
	"sao/world/party"
	"strings"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/google/uuid"
)

var PathToString = map[inventory.SkillPath]string{
	inventory.PathControl:   "Kontrola",
	inventory.PathDamage:    "Obrażenia",
	inventory.PathEndurance: "Wytrzymałość",
	inventory.PathMobility:  "Mobilność",
}

var RoleToString = map[party.PartyRole]string{
	party.Leader:  "Lider",
	party.DPS:     "DPS",
	party.Support: "Support",
	party.Tank:    "Tank",
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
		discord.SlashCommandCreate{
			Name:        "szukaj",
			Description: "Szukaj zajęcia",
		},
		discord.SlashCommandCreate{
			Name:        "party",
			Description: "Zarządzaj party",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionSubCommand{
					Name:        "pokaż",
					Description: "Pokaż party",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "zapros",
					Description: "Zaproś do party",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionUser{
							Name:        "gracz",
							Description: "Gracz",
							Required:    true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "wyrzuć",
					Description: "Wyrzuć z party",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionUser{
							Name:        "gracz",
							Description: "Gracz",
							Required:    true,
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "opuść",
					Description: "Opuść party",
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "zmień",
					Description: "Zmień rolę",
					Options: []discord.ApplicationCommandOption{
						discord.ApplicationCommandOptionUser{
							Name:        "gracz",
							Description: "Gracz",
							Required:    true,
						},
						discord.ApplicationCommandOptionString{
							Name:        "rola",
							Description: "Rola",
							Required:    true,
							Choices: []discord.ApplicationCommandOptionChoiceString{
								{
									Name:  "Lider",
									Value: "Lider",
								},
								{
									Name:  "DPS",
									Value: "DPS",
								},
								{
									Name:  "Support",
									Value: "Support",
								},
								{
									Name:  "Tank",
									Value: "Tank",
								},
							},
						},
					},
				},
				discord.ApplicationCommandOptionSubCommand{
					Name:        "rozwiąż",
					Description: "Rozwiąż party",
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

	if strings.HasPrefix(event.ComponentInteraction.Data.CustomID(), "party") {
		if event.ComponentInteraction.Data.CustomID() == "party/res" {
			userSnowflake := event.Member().User.ID.String()

			partyUuid := uuid.MustParse(event.ComponentInteraction.Data.CustomID()[10:])

			for _, pl := range World.Players {
				if pl.Meta.UserID == userSnowflake {

					if pl.Meta.Party != nil {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Jesteś już w party").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					World.Parties[partyUuid].Players = append(World.Parties[partyUuid].Players, pl.GetUUID())

					pl.Meta.Party = &partyUuid

					event.CreateMessage(
						discord.
							NewMessageCreateBuilder().
							SetContent("Dołączono do party").
							SetEphemeral(true).
							Build(),
					)

					return
				}
			}
		}

		if event.ComponentInteraction.Data.CustomID() == "party/rej" {
			event.CreateMessage(
				discord.
					NewMessageCreateBuilder().
					SetContent("Odrzucono zaproszenie").
					SetEphemeral(true).
					Build(),
			)

			//TODO msg delete (or update)
		}
	}
}

//TODO add checks if player has character XD

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

				inParty := pl.Meta.Party != nil

				inPartyText := "Nie"

				if inParty {
					inPartyText = "Tak"
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
								AddField("Party?", inPartyText, true).
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

					if _, exists := pl.Inventory.LevelSkills[lvl]; exists {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Odblokowano już umiejętność na tym poziomie").
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
	case "szukaj":
		user := event.User()

		for _, pl := range World.Players {
			if pl.Meta.UserID == user.ID.String() {
				World.PlayerEncounter(pl.GetUUID())

				//TODO implement rest
			}
		}
	case "party":
		switch *data.SubCommandName {
		case "pokaż":
			user := event.User()

			for _, pl := range World.Players {
				if pl.Meta.UserID == user.ID.String() {

					if pl.Meta.Party == nil {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Nie jesteś w party").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					embed := discord.NewEmbedBuilder()

					partyMembersText := ""

					party := World.Parties[*pl.Meta.Party]

					for _, member := range party.Players {
						memberObj := World.Players[member]

						partyMembersText += fmt.Sprintf("<@%s> - %s\n", memberObj.Meta.UserID, memberObj.GetName())
					}

					if partyMembersText[len(partyMembersText)-1] == '\n' {
						partyMembersText = partyMembersText[:len(partyMembersText)-1]
					}

					embed.AddField("Członkowie", partyMembersText, false)

					for role, players := range *party.Roles {

						if len(players) == 0 {
							embed.AddField(RoleToString[role], "Brak", false)
							continue
						}

						playersText := ""

						for _, player := range players {
							playerObj := World.Players[player]

							playersText += fmt.Sprintf("<@%s> - %s\n", playerObj.Meta.UserID, playerObj.GetName())
						}

						if playersText[len(playersText)-1] == '\n' {
							playersText = playersText[:len(playersText)-1]
						}

						embed.AddField(RoleToString[role], playersText, false)
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
		case "zapros":
			user := event.User()

			for _, pl := range World.Players {
				if pl.Meta.UserID == user.ID.String() {
					if pl.Meta.Party == nil {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Nie jesteś w party").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					if len(World.Parties[*pl.Meta.Party].Players) == 6 {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Party jest pełne").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					part := World.Parties[*pl.Meta.Party]

					if (*part.Roles)[party.Leader][0] != pl.GetUUID() {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Nie jesteś liderem").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					mentionedUser := data.User("gracz")

					for _, pl := range World.Players {
						if pl.Meta.UserID == mentionedUser.ID.String() {
							if pl.Meta.Party != nil {
								event.CreateMessage(
									discord.
										NewMessageCreateBuilder().
										SetContent("Gracz jest już w party").
										SetEphemeral(true).
										Build(),
								)
								return
							}
						}
					}

					ch, error := event.Client().Rest().CreateDMChannel(mentionedUser.ID)

					if error != nil {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Nie można wysłać wiadomości do gracza").
								SetEphemeral(true).
								Build(),
						)
					}

					chID := ch.ID()

					_, error = event.Client().Rest().CreateMessage(chID, discord.NewMessageCreateBuilder().
						SetContent(fmt.Sprintf("<@%s> (%s) zaprasza cię do party", user.ID.String(), pl.GetName())).
						AddActionRow(
							discord.NewPrimaryButton("Akceptuj", "party/res|"+(*pl.Meta.Party).String()),
							discord.NewDangerButton("Odrzuć", "party/rej|"+(*pl.Meta.Party).String()),
						).
						Build(),
					)

					if error != nil {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Nie można wysłać wiadomości do gracza").
								SetEphemeral(true).
								Build(),
						)
					} else {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Wysłano zaproszenie").
								SetEphemeral(true).
								Build(),
						)
					}

					return
				}
			}

		case "wyrzuć":
			user := event.User()

			for _, plyer := range World.Players {
				if plyer.Meta.UserID == user.ID.String() {

					if plyer.Meta.Party == nil {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Nie jesteś w party").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					part := World.Parties[*plyer.Meta.Party]

					if (*part.Roles)[party.Leader][0] != plyer.GetUUID() {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Nie jesteś liderem").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					mentionedUser := data.User("gracz")

					if mentionedUser.ID == user.ID {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Nie możesz wyrzucić samego siebie").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					for _, pl := range World.Players {
						if pl.Meta.UserID == mentionedUser.ID.String() {
							if pl.Meta.Party == nil {
								event.CreateMessage(
									discord.
										NewMessageCreateBuilder().
										SetContent("Gracz nie jest w party").
										SetEphemeral(true).
										Build(),
								)
								return
							}

							if *pl.Meta.Party != *plyer.Meta.Party {
								event.CreateMessage(
									discord.
										NewMessageCreateBuilder().
										SetContent("Gracz nie jest w twoim party").
										SetEphemeral(true).
										Build(),
								)
								return
							}

							for i, partyMember := range World.Parties[*plyer.Meta.Party].Players {
								if partyMember == pl.GetUUID() {
									pl.Meta.Party = nil

									World.Parties[*plyer.Meta.Party].Players = append(World.Parties[*plyer.Meta.Party].Players[:i], World.Parties[*plyer.Meta.Party].Players[i+1:]...)
									break
								}
							}

						}
					}

					return
				}
			}

		case "opuść":
			user := event.User()

			for _, pl := range World.Players {
				if pl.Meta.UserID == user.ID.String() {

					if pl.Meta.Party == nil {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Nie jesteś w party").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					part := World.Parties[*pl.Meta.Party]

					if (*part.Roles)[party.Leader][0] == pl.GetUUID() {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Nie możesz opuścić party będąc liderem").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					for i, partyMember := range World.Parties[*pl.Meta.Party].Players {
						if partyMember == pl.GetUUID() {
							World.Parties[*pl.Meta.Party].Players = append(World.Parties[*pl.Meta.Party].Players[:i], World.Parties[*pl.Meta.Party].Players[i+1:]...)
							break
						}
					}

					pl.Meta.Party = nil

					event.CreateMessage(
						discord.
							NewMessageCreateBuilder().
							SetContent("Opuściłeś party").
							SetEphemeral(true).
							Build(),
					)

					return
				}
			}

		case "zmień":
			user := event.User()

			for _, playr := range World.Players {
				if playr.Meta.UserID == user.ID.String() {

					if playr.Meta.Party == nil {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Nie jesteś w party").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					part := World.Parties[*playr.Meta.Party]

					if (*part.Roles)[party.Leader][0] != playr.GetUUID() {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Nie jesteś liderem").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					mentionedUser := data.User("gracz")

					if mentionedUser.ID == user.ID {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Nie możesz zmienić roli samego siebie").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					for _, pl := range World.Players {
						if pl.Meta.UserID == mentionedUser.ID.String() {
							if pl.Meta.Party == nil {
								event.CreateMessage(
									discord.
										NewMessageCreateBuilder().
										SetContent("Gracz nie jest w party").
										SetEphemeral(true).
										Build(),
								)
								return
							}

							if *pl.Meta.Party != *playr.Meta.Party {
								event.CreateMessage(
									discord.
										NewMessageCreateBuilder().
										SetContent("Gracz nie jest w twoim party").
										SetEphemeral(true).
										Build(),
								)
								return
							}

							role := data.String("rola")

							switch role {
							case "Lider":
								(*part.Roles)[party.Leader][0] = pl.GetUUID()
							case "DPS":
								(*part.Roles)[party.DPS] = append((*part.Roles)[party.DPS], pl.GetUUID())
							case "Support":
								(*part.Roles)[party.Support] = append((*part.Roles)[party.Support], pl.GetUUID())
							case "Tank":
								(*part.Roles)[party.Tank] = append((*part.Roles)[party.Tank], pl.GetUUID())
							}

							event.CreateMessage(
								discord.
									NewMessageCreateBuilder().
									SetContent("Zmieniono rolę").
									SetEphemeral(true).
									Build(),
							)
						}
					}
				}
			}

		case "rozwiąż":
			user := event.User()

			for _, pl := range World.Players {
				if pl.Meta.UserID == user.ID.String() {
					if pl.Meta.Party == nil {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Nie jesteś w party").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					part := World.Parties[*pl.Meta.Party]

					if (*part.Roles)[party.Leader][0] != pl.GetUUID() {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Nie jesteś liderem").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					uuid := *pl.Meta.Party

					for _, partyMember := range World.Parties[uuid].Players {
						World.Players[partyMember].Meta.Party = nil
					}

					delete(World.Parties, uuid)

					event.CreateMessage(
						discord.
							NewMessageCreateBuilder().
							SetContent("Rozwiązano party").
							SetEphemeral(true).
							Build(),
					)

					return
				}
			}
		}
	}
}
