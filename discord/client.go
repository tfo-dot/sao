package discord

import (
	"context"
	"fmt"
	"sao/player"
	"sao/player/inventory"
	"sao/world"
	"sao/world/party"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
)

var World *world.World
var Client *bot.Client

func StartClient(token string) {
	client, err := disgo.New(token,
		bot.WithDefaultGateway(),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			e.Client().SetPresence(context.Background(), gateway.WithWatchingActivity("SAO"))
		}),
		bot.WithEventListenerFunc(commandListener),
		bot.WithEventListenerFunc(AutocompleteHandler),
		bot.WithEventListenerFunc(ComponentHandler),
	)

	if err != nil {
		panic(err)
	}

	Client = &client

	if _, err = (*Client).Rest().SetGuildCommands((*Client).ApplicationID(), snowflake.MustParse("1151589368373444690"), DISCORD_COMMANDS); err != nil {
		fmt.Println("error while registering commands")
	}

	if err = (*Client).OpenGateway(context.Background()); err != nil {
		panic(err)
	}

	go worldMessageListener()
}

func worldMessageListener() {
	for {
		msg, ok := <-World.DChannel

		if !ok {
			return
		}

		snowflake, err := snowflake.Parse(msg.ChannelID)

		if err != nil {
			panic(err)
		}

		_, err = (*Client).Rest().CreateMessage(snowflake, msg.MessageContent)

		if err != nil {
			panic(err)
		}
	}
}

func commandListener(event *events.ApplicationCommandInteractionCreate) {
	data := event.SlashCommandInteractionData()

	user := event.User()
	member := event.Member()

	var playerChar *player.Player

	for _, pl := range World.Players {
		if pl.Meta.UserID == user.ID.String() {
			playerChar = pl
		}
	}

	if playerChar == nil && data.CommandName() != "create" {
		event.CreateMessage(noCharMessage)
		return
	}

	switch data.CommandName() {
	case "create":
		if !isAdmin(member) {
			event.CreateMessage(MessageContent("Nie masz uprawnień do tej komendy", true))
			return
		}

		charName := data.String("nazwa")
		charUser := data.User("gracz")

		World.RegisterNewPlayer(charName, charUser.ID.String())

		err := event.CreateMessage(MessageContent("Zarejestrowano postać "+charName, false))

		if err != nil {
			fmt.Println("error while sending message")
		}
		return
	case "ruch":
		locationName := data.String("nazwa")

		err := World.MovePlayer(playerChar.GetUUID(), playerChar.Meta.Location.FloorName, locationName, "")

		if err == nil {
			event.CreateMessage(MessageContent("Przeszedłeś do "+locationName, false))
		} else {
			fmt.Println(err)
			event.CreateMessage(MessageContent("Nie udało się przejść do "+locationName, true))
		}

		return
	case "tp":
		floorName := data.String("nazwa")

		currentLocation := World.Floors[playerChar.Meta.Location.FloorName].FindLocation(playerChar.Meta.Location.LocationName)

		if !currentLocation.TP {
			event.CreateMessage(
				MessageContent("Nie możesz się stąd teleportować (idź do miasta lub lokacji z tp)", true),
			)
			return
		}

		err := World.MovePlayer(playerChar.GetUUID(), floorName, World.Floors[floorName].Default, "")

		if err == nil {
			event.CreateMessage(MessageContent("Teleportowałeś się na"+floorName, false))
		}

		return
	case "info":
		if mentionedUser, exists := data.OptUser("gracz"); exists {
			user = mentionedUser

			playerChar = nil
			playerChar = World.GetPlayer(user.ID.String())

			if playerChar == nil {
				event.CreateMessage(
					discord.
						NewMessageCreateBuilder().
						SetContent("Użytkownik nie ma postaci").
						SetEphemeral(true).
						Build(),
				)
				return
			}
		}

		if playerChar == nil {
			event.CreateMessage(
				discord.
					NewMessageCreateBuilder().
					SetContent("Nie znaleziono postaci").
					SetEphemeral(true).
					Build(),
			)
			return
		}

		inFightText := "Nie"
		inPartyText := "Nie"

		inFight := playerChar.Meta.FightInstance != nil
		inParty := playerChar.Meta.Party != nil

		if inFight {
			inFightText = "Tak"
		}

		if inParty {
			inPartyText = "Tak"
		}

		event.CreateMessage(
			discord.NewMessageCreateBuilder().
				AddEmbeds(
					discord.NewEmbedBuilder().
						AddField("Nazwa", playerChar.GetName(), true).
						AddField("Gracz", fmt.Sprintf("<@%s>", playerChar.Meta.UserID), true).
						AddField("Lokacja", playerChar.Meta.Location.LocationName, true).
						AddField("Piętro", playerChar.Meta.Location.FloorName, true).
						AddField("HP", fmt.Sprintf("%d/%d", playerChar.GetCurrentHP(), playerChar.GetMaxHP()), true).
						AddField("Mana", fmt.Sprintf("%d/%d", playerChar.GetCurrentMana(), playerChar.GetMaxMana()), true).
						AddField("Atak", fmt.Sprintf("%d", playerChar.GetATK()), true).
						AddField("AP", fmt.Sprintf("%d", playerChar.GetAP()), true).
						AddField("DEF/RES", fmt.Sprintf("%d/%d", playerChar.GetDEF(), playerChar.GetMR()), true).
						AddField("Lvl", fmt.Sprintf("%d %d/%d", playerChar.XP.Level, playerChar.XP.Exp, (playerChar.XP.Level*100)+100), true).
						AddField("SPD/AGL", fmt.Sprintf("%d/%d", playerChar.GetSPD(), playerChar.GetAGL()), true).
						AddField("Walce?", inFightText, true).
						AddField("Party?", inPartyText, true).
						AddField("Gildia?", "Nie", true).
						Build(),
				).
				Build(),
		)
	case "skill":
		switch *data.SubCommandName {
		case "pokaż":
			if mentionedUser, exists := data.OptUser("gracz"); exists {
				user = mentionedUser

				playerChar = nil
				playerChar = World.GetPlayer(user.ID.String())

				if playerChar == nil {
					event.CreateMessage(
						discord.
							NewMessageCreateBuilder().
							SetContent("Użytkownik nie ma postaci").
							SetEphemeral(true).
							Build(),
					)
					return
				}
			}

			embed := discord.NewEmbedBuilder()

			if len(playerChar.Inventory.Skills) == 0 {
				embed.AddField("Umiejętności", "Brak", false)
			}

			for _, skill := range playerChar.Inventory.Skills {
				embed.AddField(skill.Name, skill.Description, false)
			}

			if len(playerChar.Inventory.LevelSkills) == 0 {
				embed.AddField("Umiejętności za lvl", "Brak", false)
			}

			for _, skill := range playerChar.Inventory.LevelSkills {
				embed.AddField(
					fmt.Sprintf("%s (LVL: %d)", skill.Name, skill.ForLevel),
					fmt.Sprintf("Ścieżka: %s\n\n%s", PathToString[skill.Path], skill.Description),
					false,
				)
			}

			if len(embed.Fields) >= 20 {
				event.CreateMessage(
					discord.
						NewMessageCreateBuilder().
						SetContent("Wykryto spaghetti od <@344048874656366592> aka jeszcze nie zrobione!...").
						Build(),
				)
			}

			event.CreateMessage(
				discord.
					NewMessageCreateBuilder().
					AddEmbeds(embed.Build()).
					Build(),
			)

			return
		case "odblokuj":
			lvl := data.Int("lvl")

			skillList := make([]inventory.PlayerSkill, 0)

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

			if _, exists := playerChar.Inventory.LevelSkills[lvl]; exists {
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
	case "plecak":
		switch *data.SubCommandName {
		case "pokaż":
			embed := discord.NewEmbedBuilder()

			if len(playerChar.Inventory.Items) == 0 {
				embed.AddField("Przedmioty", fmt.Sprintf("%d/%d", 0, playerChar.Inventory.Capacity), false)
			} else {
				count := 0

				for _, item := range playerChar.Inventory.Items {
					if !item.Hidden {
						count++
					}
				}

				embed.AddField("Przedmioty", fmt.Sprintf("%d/%d", count, playerChar.Inventory.Capacity), false)
			}

			for _, item := range playerChar.Inventory.Items {

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
	case "szukaj":
		go World.PlayerSearch(playerChar.GetUUID())

		event.CreateMessage(
			discord.
				NewMessageCreateBuilder().
				SetContent("Szukanie...").
				SetEphemeral(true).
				Build(),
		)
		return
	case "party":
		switch *data.SubCommandName {
		case "pokaż":
			if playerChar.Meta.Party == nil {
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

			party := World.Parties[*playerChar.Meta.Party]

			for _, member := range party.Players {
				memberObj := World.Players[member.PlayerUuid]

				partyMembersText += fmt.Sprintf("<@%s> - %s\n", memberObj.Meta.UserID, memberObj.GetName())
			}

			if partyMembersText[len(partyMembersText)-1] == '\n' {
				partyMembersText = partyMembersText[:len(partyMembersText)-1]
			}

			embed.AddField("Członkowie", partyMembersText, false)

			//TODO group by role

			// for _, entry := range party.Players {

			// 	if len(party.Players) == 0 {
			// 		embed.AddField(RoleToString[entry.Role], "Brak", false)
			// 		continue
			// 	}

			// 	playersText := ""

			// 	for _, player := range players {
			// 		playerObj := World.Players[player]

			// 		playersText += fmt.Sprintf("<@%s> - %s\n", playerObj.Meta.UserID, playerObj.GetName())
			// 	}

			// 	if playersText[len(playersText)-1] == '\n' {
			// 		playersText = playersText[:len(playersText)-1]
			// 	}

			// 	embed.AddField(RoleToString[role], playersText, false)
			// }

			event.CreateMessage(
				discord.
					NewMessageCreateBuilder().
					AddEmbeds(embed.Build()).
					Build(),
			)

			return
		case "zapros":
			if playerChar.Meta.Party == nil {
				event.CreateMessage(
					discord.
						NewMessageCreateBuilder().
						SetContent("Nie jesteś w party").
						SetEphemeral(true).
						Build(),
				)
				return
			}

			if len(World.Parties[*playerChar.Meta.Party].Players) >= 6 {
				event.CreateMessage(
					discord.
						NewMessageCreateBuilder().
						SetContent("Party jest pełne").
						SetEphemeral(true).
						Build(),
				)
				return
			}

			part := World.Parties[*playerChar.Meta.Party]

			for _, entry := range part.Players {
				if entry.PlayerUuid == playerChar.GetUUID() && entry.Role != party.Leader {
					event.CreateMessage(
						discord.
							NewMessageCreateBuilder().
							SetContent("Nie jesteś liderem").
							SetEphemeral(true).
							Build(),
					)
					return
				}
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
				SetContent(fmt.Sprintf("<@%s> (%s) zaprasza cię do party", user.ID.String(), playerChar.GetName())).
				AddActionRow(
					discord.NewPrimaryButton("Akceptuj", "party/res|"+(*playerChar.Meta.Party).String()),
					discord.NewDangerButton("Odrzuć", "party/rej|"+(*playerChar.Meta.Party).String()),
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
		case "wyrzuć":
			if playerChar.Meta.Party == nil {
				event.CreateMessage(
					discord.
						NewMessageCreateBuilder().
						SetContent("Nie jesteś w party").
						SetEphemeral(true).
						Build(),
				)
				return
			}

			part := World.Parties[*playerChar.Meta.Party]

			for _, entry := range part.Players {
				if entry.PlayerUuid == playerChar.GetUUID() && entry.Role != party.Leader {
					event.CreateMessage(
						discord.
							NewMessageCreateBuilder().
							SetContent("Nie jesteś liderem").
							SetEphemeral(true).
							Build(),
					)
					return
				}
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

					if *pl.Meta.Party != *playerChar.Meta.Party {
						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Gracz nie jest w twoim party").
								SetEphemeral(true).
								Build(),
						)
						return
					}

					for i, partyMember := range World.Parties[*playerChar.Meta.Party].Players {
						if partyMember.PlayerUuid == pl.GetUUID() {
							pl.Meta.Party = nil

							World.Parties[*playerChar.Meta.Party].Players = append(World.Parties[*playerChar.Meta.Party].Players[:i], World.Parties[*playerChar.Meta.Party].Players[i+1:]...)
							break
						}
					}

				}
			}

			return
		case "opuść":
			if playerChar.Meta.Party == nil {
				event.CreateMessage(
					discord.
						NewMessageCreateBuilder().
						SetContent("Nie jesteś w party").
						SetEphemeral(true).
						Build(),
				)
				return
			}

			part := World.Parties[*playerChar.Meta.Party]

			for _, entry := range part.Players {
				if entry.PlayerUuid == playerChar.GetUUID() && entry.Role != party.Leader {
					event.CreateMessage(
						discord.
							NewMessageCreateBuilder().
							SetContent("Nie możesz opuścić party będąc liderem").
							SetEphemeral(true).
							Build(),
					)
					return
				}
			}

			for i, partyMember := range World.Parties[*playerChar.Meta.Party].Players {
				if partyMember.PlayerUuid == playerChar.GetUUID() {
					World.Parties[*playerChar.Meta.Party].Players = append(World.Parties[*playerChar.Meta.Party].Players[:i], World.Parties[*playerChar.Meta.Party].Players[i+1:]...)
					break
				}
			}

			playerChar.Meta.Party = nil

			event.CreateMessage(
				discord.
					NewMessageCreateBuilder().
					SetContent("Opuściłeś party").
					SetEphemeral(true).
					Build(),
			)

			return
		case "zmień":
			if playerChar.Meta.Party == nil {
				event.CreateMessage(
					discord.
						NewMessageCreateBuilder().
						SetContent("Nie jesteś w party").
						SetEphemeral(true).
						Build(),
				)
				return
			}

			part := World.Parties[*playerChar.Meta.Party]

			for _, entry := range part.Players {
				if entry.PlayerUuid == playerChar.GetUUID() && entry.Role != party.Leader {
					event.CreateMessage(
						discord.
							NewMessageCreateBuilder().
							SetContent("Nie jesteś liderem").
							SetEphemeral(true).
							Build(),
					)
					return
				}
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

					if *pl.Meta.Party != *playerChar.Meta.Party {
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

					for i, partyMember := range World.Parties[*playerChar.Meta.Party].Players {
						if partyMember.PlayerUuid == pl.GetUUID() {
							switch role {
							case "Lider":
								World.Parties[*playerChar.Meta.Party].Players[i].Role = party.Leader
							case "DPS":
								World.Parties[*playerChar.Meta.Party].Players[i].Role = party.DPS
							case "Support":
								World.Parties[*playerChar.Meta.Party].Players[i].Role = party.Support
							case "Tank":
								World.Parties[*playerChar.Meta.Party].Players[i].Role = party.Tank
							}
							break
						}
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
		case "rozwiąż":
			if playerChar.Meta.Party == nil {
				event.CreateMessage(
					discord.
						NewMessageCreateBuilder().
						SetContent("Nie jesteś w party").
						SetEphemeral(true).
						Build(),
				)
				return
			}

			part := World.Parties[*playerChar.Meta.Party]

			for _, entry := range part.Players {
				if entry.PlayerUuid == playerChar.GetUUID() && entry.Role != party.Leader {
					event.CreateMessage(
						discord.
							NewMessageCreateBuilder().
							SetContent("Nie jesteś liderem").
							SetEphemeral(true).
							Build(),
					)
					return
				}
			}

			uuid := *playerChar.Meta.Party

			for _, partyMember := range World.Parties[uuid].Players {
				World.Players[partyMember.PlayerUuid].Meta.Party = nil
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
	case "ratuj":
		if playerChar.Meta.FightInstance != nil {
			event.CreateMessage(
				discord.
					NewMessageCreateBuilder().
					SetContent("Już jesteś w walce!").
					SetEphemeral(true).
					Build(),
			)
			return
		}

		options := make([]discord.StringSelectMenuOption, 0)

		cid := event.Channel().ID().String()

		for _, fight := range World.Fights {

			if fight.Location.CID != cid {
				continue
			} else {
				for entityUuid, entityEntry := range fight.Entities {
					if !entityEntry.Entity.IsAuto() {
						options = append(options, discord.NewStringSelectMenuOption(entityEntry.Entity.GetName(), entityUuid.String()))
					}
				}
			}

		}

		if len(options) == 0 {
			event.CreateMessage(
				discord.
					NewMessageCreateBuilder().
					SetContent("Nie ma nikogo do ratowania").
					SetEphemeral(true).
					Build(),
			)
			return
		}

		event.CreateMessage(
			discord.
				NewMessageCreateBuilder().
				SetContent("Wybierz kogo chcesz ratować").
				AddActionRow(
					discord.NewStringSelectMenu("f/save", "Wybierz", options...),
				).
				Build(),
		)
		return
	}
}
