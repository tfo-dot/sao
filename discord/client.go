package discord

import (
	"bytes"
	"context"
	"fmt"
	"sao/battle/mobs"
	"sao/data"
	"sao/player"
	"sao/player/inventory"
	"sao/types"
	"sao/utils"
	"sao/world"
	"sao/world/party"
	"sao/world/tournament"
	"slices"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/snowflake/v2"
	"github.com/google/uuid"
)

var World *world.World
var Client *bot.Client
var Choices = make([]types.DiscordChoice, 0)

func StartClient() {
	client, err := disgo.New(data.Config.Token,
		bot.WithEventListenerFunc(func(e *events.Ready) {
			println("Discord connection is here!")
			e.Client().SetPresence(context.Background(), gateway.WithWatchingActivity("SAO"))
		}),
		bot.WithEventListenerFunc(func(e *events.MessageCreate) {
			if e.Message.Content == "sao:dump" && e.Message.Author.ID.String() == data.Config.Owner {
				rawBackup := World.CreateBackup()

				e.Client().Rest().AddReaction(e.Message.ChannelID, e.Message.ID, data.Config.Emote)

				chanel, err := e.Client().Rest().CreateDMChannel(snowflake.MustParse(data.Config.Owner))

				if err != nil {
					return
				}

				_, err = e.Client().Rest().CreateMessage(
					chanel.ID(),
					discord.MessageCreate{Files: []*discord.File{{Name: "backup.json", Reader: bytes.NewReader(rawBackup)}}})

				if err != nil {
					return
				}
			}

			if e.Message.Content == "sao:reload" && e.Message.Author.ID.String() == data.Config.Owner {
				mobs.Mobs = mobs.GetMobs()
				data.FloorMap = data.GetFloors()
				data.PlayerDefaults = data.GetPlayerDefaults()
				data.Shops = data.GetShops()
				data.WorldConfig = data.GetWorldConfig()
				data.Items = data.GetItems()

				e.Client().Rest().AddReaction(e.Message.ChannelID, e.Message.ID, data.Config.Emote)
			}
		}),
		bot.WithEventListenerFunc(commandListener),
		bot.WithEventListenerFunc(AutocompleteHandler),
		bot.WithEventListenerFunc(ComponentHandler),
		bot.WithEventListenerFunc(ModalSubmitHandler),
		bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuildMessages, gateway.IntentMessageContent)),
	)

	if err != nil {
		panic(err)
	}

	Client = &client

	if _, err = (*Client).Rest().SetGuildCommands((*Client).ApplicationID(), snowflake.MustParse(data.Config.GuildID), DISCORD_COMMANDS); err != nil {
		fmt.Println("error while registering commands")
	}

	if err = (*Client).OpenGateway(context.Background()); err != nil {
		panic(err)
	}

	go worldMessageListener()
}

func commandListener(event *events.ApplicationCommandInteractionCreate) {
	interactionData := event.SlashCommandInteractionData()

	user := event.User()
	member := event.Member()

	var playerChar *player.Player

	for _, pl := range World.Players {
		if pl.Meta.UserID == user.ID.String() {
			playerChar = pl
		}
	}

	if interactionData.CommandName() != "create" && interactionData.CommandName() != "turniej" && playerChar == nil {
		event.CreateMessage(noCharMessage)
		return
	}

	switch interactionData.CommandName() {
	case "create":
		if !isAdmin(member) {
			event.CreateMessage(MessageContent("Nie masz uprawnień do tej komendy", true))
			return
		}

		charName := interactionData.String("nazwa")
		charUser := interactionData.User("gracz")

		if World.GetPlayer(charUser.ID.String()) != nil {
			event.CreateMessage(MessageContent("Gracz ma już postać", true))
			return
		}

		newPlayer := player.NewPlayer(charName, charUser.ID.String())

		World.Players[newPlayer.GetUUID()] = &newPlayer

		err := event.CreateMessage(MessageContent("Zarejestrowano postać "+charName, false))

		if err != nil {
			fmt.Println("error while sending message")
			return
		}

		(*Client).Rest().AddMemberRole(*event.GuildID(), charUser.ID, snowflake.MustParse(data.Config.RoleID))
		return
	case "info":
		if mentionedUser, exists := interactionData.OptUser("gracz"); exists {
			user = mentionedUser

			playerChar = nil
			playerChar = World.GetPlayer(user.ID.String())

			if playerChar == nil {
				event.CreateMessage(MessageContent("Użytkownik nie ma postaci", true))
				return
			}
		}

		if playerChar == nil {
			event.CreateMessage(noCharMessage)
			return
		}

		inFightText := "Tak"

		if playerChar.Meta.FightInstance == nil {
			inFightText = "Nie"
		}

		inPartyText := "Tak"

		if playerChar.Meta.Party == nil {
			inPartyText = "Nie"
		}

		lvlText := fmt.Sprint(playerChar.XP.Level)

		if playerChar.XP.Level >= data.FloorMap.GetUnlockedFloorCount()*5 {
			lvlText += " MAX"
		} else {
			lvlText += fmt.Sprintf(" %d/%d", playerChar.XP.Exp, (playerChar.XP.Level*100)+100)
		}

		derivedStatsText := ""

		for _, stat := range playerChar.GetDerivedStats() {
			statValue := utils.PercentOf(playerChar.GetStat(stat.Base), stat.Percent)

			derivedStatsText += fmt.Sprintf("- %d %s (%d%% %s => %s)\n", statValue, types.StatToString[stat.Derived], stat.Percent, types.StatToString[stat.Base], types.StatToString[stat.Derived])
		}

		if derivedStatsText == "" {
			derivedStatsText = "Brak"
		}

		messageBuilder := discord.NewMessageCreateBuilder().
			AddEmbeds(
				discord.NewEmbedBuilder().
					AddField("Nazwa", playerChar.GetName(), true).
					AddField("Gracz", fmt.Sprintf("<@%s>", playerChar.Meta.UserID), true).
					AddField("HP", fmt.Sprintf("%d/%d", playerChar.Stats.HP, playerChar.GetStat(types.STAT_HP)), true).
					AddField("Mana", fmt.Sprintf("%d/%d", playerChar.GetCurrentMana(), playerChar.GetStat(types.STAT_MANA)), true).
					AddField("Atak", fmt.Sprintf("%d", playerChar.GetStat(types.STAT_AD)), true).
					AddField("AP", fmt.Sprintf("%d", playerChar.GetStat(types.STAT_AP)), true).
					AddField("DEF/RES", fmt.Sprintf("%d/%d", playerChar.GetStat(types.STAT_DEF), playerChar.GetStat(types.STAT_MR)), true).
					AddField("Lvl", lvlText, true).
					AddField("SPD/AGL", fmt.Sprintf("%d/%d", playerChar.GetStat(types.STAT_SPD), playerChar.GetStat(types.STAT_AGL)), true).
					AddField("W walce?", inFightText, true).
					AddField("W party?", inPartyText, true).
					AddField("Dynamiczne statystyki", derivedStatsText, true).
					Build(),
			)

		if playerChar.Stats.HP > 0 && playerChar.Stats.HP < playerChar.GetStat(types.STAT_HP) {
			messageBuilder.AddActionRow(discord.NewPrimaryButton("Czekaj do 100% HP", "utils|wait"))
		}

		event.CreateMessage(messageBuilder.Build())
	case "skill":
		switch *interactionData.SubCommandName {
		case "pokaż":
			if mentionedUser, exists := interactionData.OptUser("gracz"); exists {
				user = mentionedUser

				playerChar = nil
				playerChar = World.GetPlayer(user.ID.String())

				if playerChar == nil {
					event.CreateMessage(MessageContent("Użytkownik nie ma postaci", true))
					return
				}
			}

			embed := discord.NewEmbedBuilder()

			if len(playerChar.Inventory.LevelSkills) == 0 {
				embed.AddField("Umiejętności za lvl", "Brak", false)
			}

			tempArray := make([]LevelField, len(playerChar.Inventory.LevelSkills))

			notInLine := false

			for key, skillInfo := range playerChar.Inventory.LevelSkills {

				skillDescription := skillInfo.Skill.GetUpgradableDescription(skillInfo.Upgrades)

				skillName := skillInfo.Skill.GetName()

				if types.HasFlag(skillInfo.Skill.GetTrigger().Flags, types.FLAG_INSTANT_SKILL) {
					skillName = fmt.Sprintf("**%s**", skillName)
				}

				if types.HasFlag(skillInfo.Skill.GetTrigger().Flags, types.FLAG_IGNORE_CC) {
					skillName = fmt.Sprintf("__%s__", skillName)
				}

				if skillCd := skillInfo.Skill.GetCooldown(skillInfo.Upgrades); skillCd != 0 {
					skillName = fmt.Sprintf("%s CD:%d", skillName, skillCd)
				}

				tempArray = append(tempArray, LevelField{Level: key, Field: discord.EmbedField{
					Name:   skillName,
					Value:  skillDescription,
					Inline: &notInLine,
				}})
			}

			slices.SortFunc(tempArray, func(i, j LevelField) int {
				return i.Level - j.Level
			})

			for _, field := range tempArray {
				embed.AddField(field.Field.Name, field.Field.Value, false)
			}

			event.CreateMessage(MessageEmbed(embed.Build()))

			return
		case "odblokuj":
			lvl := interactionData.Int("lvl")

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

			if lvl%10 == 0 {
				path := playerChar.GetSkillPath()
				embed := discord.NewEmbedBuilder()

				buttons := make([]discord.InteractiveComponent, 0)

				for idx, skill := range inventory.AVAILABLE_SKILLS[path][lvl] {
					embed.AddField(
						fmt.Sprintf("[%d] %s", idx+1, skill.GetName()), skill.GetDescription(), false,
					)

					buttons = append(buttons, discord.NewPrimaryButton(
						fmt.Sprintf("[%d] %s", idx+1, skill.GetName()), fmt.Sprintf("su|%d|%d|%d", path, lvl, idx)),
					)
				}

				event.CreateMessage(discord.MessageCreate{
					Embeds:     []discord.Embed{embed.Build()},
					Components: []discord.ContainerComponent{discord.NewActionRow(buttons...)},
				})

				return
			}

			embed := discord.NewEmbedBuilder()

			buttons := make([]discord.InteractiveComponent, 0)

			for _, skillTree := range inventory.AVAILABLE_SKILLS {
				skills, exists := skillTree[lvl]

				if !exists {
					continue
				}

				for idx, skill := range skills {
					upgradeMsg := ""

					for _, upgrade := range skill.GetUpgrades() {
						upgradeMsg += fmt.Sprintf("\n- %s", upgrade.Description)
					}

					embed.AddField(
						skill.GetName(), fmt.Sprintf("%s\nUlepszenia:%s", skill.GetDescription(), upgradeMsg), false,
					)

					buttons = append(buttons, discord.NewPrimaryButton(
						skill.GetName(), fmt.Sprintf("su|%d|%d|%d", skill.GetPath(), skill.GetLevel(), idx)),
					)
				}
			}

			if len(buttons) == 0 {
				event.CreateMessage(MessageContent("Nie znaleziono umiejętności dla tego poziomu", true))
				return
			}

			containers := make([]discord.ContainerComponent, 0)
			containers = append(containers, discord.NewActionRow())

			tempButtons := make([]discord.InteractiveComponent, len(buttons))
			copy(tempButtons, buttons)

			for len(tempButtons) > 0 {
				if len(containers[len(containers)-1].(discord.ActionRowComponent).Components()) == 5 {
					containers = append(containers, discord.NewActionRow())
				}

				containers[len(containers)-1] = containers[len(containers)-1].(discord.ActionRowComponent).AddComponents(tempButtons[0])
				tempButtons = tempButtons[1:]
			}

			if len(containers[len(containers)-1].Components()) == 0 {
				containers = containers[:len(containers)-1]
			}

			if len(containers) > 5 {
				event.CreateMessage(MessageContent("Za dużo umiejętności na tym poziomie.", true))
				return
			}

			event.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{embed.Build()}, Components: containers,
			})

			return
		case "ulepsz":
			lvl := interactionData.Int("lvl")

			skillInfo, exists := playerChar.Inventory.LevelSkills[lvl]

			if !exists {
				event.CreateMessage(MessageContent("Nie masz takiej umiejętności", true))
				return
			}

			upgrades := skillInfo.Skill.GetUpgrades()

			if len(upgrades) == 0 {
				event.CreateMessage(MessageContent("Ta umiejętność nie ma ulepszeń", true))
				return
			}

			if playerChar.GetAvailableSkillActions() <= 0 {
				event.CreateMessage(MessageContent("Posiadasz za mało punktów akcji", true))
				return
			}

			availableUpgrades := make([]types.PlayerSkillUpgrade, 0)

			for idx, upgrade := range upgrades {
				if inventory.HasUpgrade(skillInfo.Upgrades, idx+1) {
					continue
				}

				availableUpgrades = append(availableUpgrades, upgrade)
			}

			if len(availableUpgrades) == 0 {
				event.CreateMessage(MessageContent("Nie masz dostępnych ulepszeń", true))
				return
			}

			buttons := make([]discord.InteractiveComponent, 0)
			embed := discord.NewEmbedBuilder()

			for idx, upgrade := range availableUpgrades {
				embed.AddField(fmt.Sprintf("Ulepszenie %v", idx+1), upgrade.Description, false)

				buttons = append(buttons, discord.NewPrimaryButton(
					fmt.Sprintf("Ulepsz %v", idx+1), fmt.Sprintf("sup|%d|%d", lvl, idx),
				))
			}

			event.CreateMessage(
				discord.NewMessageCreateBuilder().AddEmbeds(embed.Build()).AddActionRow(buttons...).Build(),
			)

		}
	case "plecak":
		switch *interactionData.SubCommandName {
		case "pokaż":
			embed := discord.NewEmbedBuilder()

			if len(playerChar.Inventory.Items) == 0 {
				embed.AddField("Przedmioty", fmt.Sprintf("%d/10", 0), false)
			} else {
				count := 0

				for _, item := range playerChar.Inventory.Items {
					if !item.Hidden {
						count++
					}
				}

				embed.AddField("Przedmioty", fmt.Sprintf("%d/10", count), false)
			}

			for _, item := range playerChar.Inventory.Items {
				if item.Hidden {
					continue
				}

				embed.AddField(item.Name, item.Description, false)
			}

			event.CreateMessage(MessageEmbed(
				embed.AddField("Złoto", fmt.Sprintf("%d", playerChar.Inventory.Gold), false).Build()),
			)

			return
		}
	case "szukaj":
		dChannel, error := (*Client).Rest().GetChannel(event.Channel().ID())

		if error != nil {
			event.CreateMessage(MessageContent("Nie można znaleźć kanału", true))
			return
		}

		var channelId string
		var threadId string

		if threadChannel, ok := dChannel.(discord.GuildThread); ok {
			threadChannel.ParentID()

			textChannel, error := (*Client).Rest().GetChannel(*threadChannel.ParentID())

			if error != nil {
				event.CreateMessage(MessageContent("Nie można znaleźć kanału", true))
				return
			}

			channelId = textChannel.ID().String()
			threadId = dChannel.ID().String()
		} else {
			channelId = dChannel.ID().String()
		}

		loc := data.FloorMap.FindLocation(func(l types.Location) bool { return l.CID == channelId })

		println(channelId)

		if loc == nil {
			event.CreateMessage(discord.
				NewMessageCreateBuilder().
				SetContent("Nie rozpoznaje tego kanału...").
				SetEphemeral(true).
				Build(),
			)
			return
		}

		if loc.CityPart {
			event.CreateMessage(discord.
				NewMessageCreateBuilder().
				SetContent("Nie ma czego tu szukać...").
				SetEphemeral(true).
				Build(),
			)
			return
		}

		World.PlayerSearch(playerChar.GetUUID(), threadId, event)

		event.CreateMessage(MessageContent("Szukanie...", true))

		go World.PlayerSearch(playerChar.GetUUID(), threadId, event)

		return
	case "party":
		if playerChar.Meta.Party == nil && *interactionData.SubCommandName != "zapros" {
			event.CreateMessage(
				discord.
					NewMessageCreateBuilder().
					SetContent("Nie jesteś w party").
					SetEphemeral(true).
					Build(),
			)
			return
		}

		switch *interactionData.SubCommandName {
		case "pokaż":
			partyObj := World.Parties[playerChar.Meta.Party.UUID]
			partyLeader := World.Players[partyObj.Leader]
			partyMembersText := ""

			for _, member := range partyObj.Players {
				memberObj := World.Players[member.PlayerUuid]

				partyMembersText += fmt.Sprintf("<@%s> - %s\n", memberObj.Meta.UserID, memberObj.GetName())
			}

			if partyMembersText[len(partyMembersText)-1] == '\n' {
				partyMembersText = partyMembersText[:len(partyMembersText)-1]
			}

			event.CreateMessage(
				MessageEmbed(
					discord.NewEmbedBuilder().
						AddField("Członkowie", partyMembersText, false).
						AddField("Lider", fmt.Sprintf("<@%s> - %s\n", partyLeader.Meta.UserID, partyLeader.GetName()), false).
						Build(),
				),
			)

			return
		case "zapros":
			if playerChar.Meta.Party == nil {
				World.RegisterParty(party.Party{
					Leader: playerChar.GetUUID(),
					Players: []*party.PartyEntry{
						{
							PlayerUuid: playerChar.GetUUID(),
							Role:       party.None,
						},
					},
				})
			}

			if len(World.Parties[playerChar.Meta.Party.UUID].Players) >= 6 {
				event.CreateMessage(MessageContent("Party jest pełne", true))
				return
			}

			part := World.Parties[playerChar.Meta.Party.UUID]

			if playerChar.GetUUID() != part.Leader {
				event.CreateMessage(MessageContent("Nie jesteś liderem", true))
				return
			}

			mentionedUser := interactionData.User("gracz")

			pl := World.GetPlayer(mentionedUser.ID.String())

			if pl.Meta.Party != nil {
				event.CreateMessage(MessageContent("Gracz jest już w party", true))
				return
			}

			ch, error := event.Client().Rest().CreateDMChannel(mentionedUser.ID)

			if error != nil {
				event.CreateMessage(MessageContent("Nie można wysłać wiadomości do gracza", true))
				return
			}

			chID := ch.ID()

			_, error = event.Client().Rest().CreateMessage(chID, discord.NewMessageCreateBuilder().
				SetContent(fmt.Sprintf("<@%s> (%s) zaprasza cię do party", user.ID.String(), playerChar.GetName())).
				AddActionRow(
					discord.NewPrimaryButton("Akceptuj", "party/res|"+(playerChar.Meta.Party.UUID).String()),
					discord.NewDangerButton("Odrzuć", "party/rej|"+(playerChar.Meta.Party.UUID).String()),
				).
				Build(),
			)

			if error != nil {
				event.CreateMessage(MessageContent("Nie można wysłać wiadomości do gracza", true))
				return
			}

			event.CreateMessage(MessageContent("Wysłano zaproszenie do party", true))
			return
		case "wyrzuć":
			part := World.Parties[playerChar.Meta.Party.UUID]

			if playerChar.GetUUID() != part.Leader {
				event.CreateMessage(MessageContent("Nie jesteś liderem", true))
				return
			}

			mentionedUser := interactionData.User("gracz")

			if mentionedUser.ID == user.ID {
				event.CreateMessage(MessageContent("Nie możesz wyrzucić samego siebie", true))
				return
			}

			pl := World.GetPlayer(mentionedUser.ID.String())

			if pl.Meta.Party == nil {
				event.CreateMessage(MessageContent("Gracz nie jest w party", true))
				return
			}

			if pl.Meta.Party.UUID != playerChar.Meta.Party.UUID {
				event.CreateMessage(MessageContent("Gracz nie jest w twoim party", true))
				return
			}

			for i, partyMember := range World.Parties[playerChar.Meta.Party.UUID].Players {
				if partyMember.PlayerUuid == pl.GetUUID() {
					pl.Meta.Party = nil

					World.Parties[playerChar.Meta.Party.UUID].Players = slices.Delete(World.Parties[playerChar.Meta.Party.UUID].Players, i, i+1)
				} else {
					World.Players[partyMember.PlayerUuid].Meta.Party.MembersCount--
				}
			}

			return
		case "opuść":
			part := World.Parties[playerChar.Meta.Party.UUID]

			if playerChar.GetUUID() == part.Leader {
				event.CreateMessage(MessageContent("Jesteś liderem...", true))
				return
			}

			//TODO fix this XDDD
			for i, partyMember := range World.Parties[playerChar.Meta.Party.UUID].Players {
				if partyMember.PlayerUuid == playerChar.GetUUID() {
					World.Parties[playerChar.Meta.Party.UUID].Players = slices.Delete(World.Parties[playerChar.Meta.Party.UUID].Players, i, i+1)
				} else {
					World.Players[partyMember.PlayerUuid].Meta.Party.MembersCount--
				}
			}

			playerChar.Meta.Party = nil

			event.CreateMessage(MessageContent("Opuściłeś party", true))
			return
		case "zmień":
			part := World.Parties[playerChar.Meta.Party.UUID]

			if playerChar.GetUUID() != part.Leader {
				event.CreateMessage(MessageContent("Nie jesteś liderem", true))
				return
			}

			mentionedUser := interactionData.User("gracz")

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

			pl := World.GetPlayer(mentionedUser.ID.String())

			if pl.Meta.Party == nil {
				event.CreateMessage(MessageContent("Gracz nie jest w party", true))
				return
			}

			if pl.Meta.Party.UUID != playerChar.Meta.Party.UUID {
				event.CreateMessage(MessageContent("Gracz nie jest w twoim party", true))
				return
			}

			role := interactionData.String("rola")

			for i, partyMember := range World.Parties[playerChar.Meta.Party.UUID].Players {
				if partyMember.PlayerUuid == pl.GetUUID() {
					switch role {
					case "Lider":
						World.Parties[playerChar.Meta.Party.UUID].Leader = pl.GetUUID()
					case "DPS":
						World.Parties[playerChar.Meta.Party.UUID].Players[i].Role = party.DPS
					case "Support":
						World.Parties[playerChar.Meta.Party.UUID].Players[i].Role = party.Support
					case "Tank":
						World.Parties[playerChar.Meta.Party.UUID].Players[i].Role = party.Tank
					}
					break
				}
			}

			event.CreateMessage(MessageContent("Zmieniono rolę", true))
			return
		case "rozwiąż":
			part := World.Parties[playerChar.Meta.Party.UUID]

			if playerChar.GetUUID() != part.Leader {
				event.CreateMessage(MessageContent("Nie jesteś liderem", true))
				return
			}

			uuid := playerChar.Meta.Party.UUID

			for _, partyMember := range World.Parties[uuid].Players {
				World.Players[partyMember.PlayerUuid].Meta.Party = nil
			}

			delete(World.Parties, uuid)

			event.CreateMessage(MessageContent("Rozwiązano party", true))
			return
		}
	case "sklep":
		switch *interactionData.SubCommandName {
		case "pokaż":

			channelId := event.Channel().ID()

			loc := data.FloorMap.FindLocation(func(l types.Location) bool { return l.CID == channelId.String() })

			if loc == nil {
				event.CreateMessage(MessageContent("Nie rozponaje tego kanału...", true))
				return
			}

			storesInLocation := make([]*types.NPCStore, 0)

			for _, store := range data.Shops {
				if loc == store.Location {
					storesInLocation = append(storesInLocation, store)
				}
			}

			if len(storesInLocation) == 0 {
				event.CreateMessage(MessageContent("Nie ma sklepów w tej lokalizacji", true))
				return
			}

			if len(storesInLocation) > 25 {
				event.CreateMessage(MessageContent("Znaleziono więcej niż 25 sklepów w danej lokacji", true))
				return
			}

			messageBuilder := discord.NewMessageCreateBuilder()

			for i := range (len(storesInLocation) / 5) + 1 {
				var stores []*types.NPCStore

				if len(storesInLocation) <= i*5+5 {
					stores = storesInLocation[i*5:]
				} else {
					stores = storesInLocation[i*5 : i*5+5]
				}

				buttonArray := make([]discord.InteractiveComponent, 0)

				for _, store := range stores {
					buttonArray = append(buttonArray, discord.NewPrimaryButton(store.Name, "shop/show/1/"+store.Uuid.String()))
				}

				if len(buttonArray) == 0 {
					continue
				}

				messageBuilder.AddActionRow(buttonArray...)
			}

			event.CreateMessage(messageBuilder.Build())
		}
	case "turniej":
		switch *interactionData.SubCommandName {
		case "stwórz":
			if !isAdmin(member) {
				event.CreateMessage(MessageContent("Nie masz uprawnień do tej komendy", true))
				return
			}

			maxCount, isMaxCountPresent := interactionData.OptInt("max")

			if !isMaxCountPresent {
				maxCount = -1
			}

			tournament := tournament.Tournament{
				Uuid:         uuid.New(),
				Name:         interactionData.String("nazwa"),
				MaxPlayers:   maxCount,
				Participants: make([]uuid.UUID, 0),
			}

			World.RegisterTournament(tournament)

			event.CreateMessage(MessageContent("Turniej stworzony", true))
			return

		case "rozpocznij":
			if !isAdmin(member) {
				event.CreateMessage(MessageContent("Nie masz uprawnień do tej komendy", true))
				return
			}

			tournamentName := interactionData.String("nazwa")

			var actualTournament *tournament.Tournament

			for _, t := range World.Tournaments {
				if t.Name == tournamentName {
					actualTournament = t
					break
				}
			}

			if actualTournament == nil {
				event.CreateMessage(MessageContent("Nie znaleziono turnieju", true))
				return
			}

			if err := World.StartTournament(actualTournament.Uuid); err != nil {
				event.CreateMessage(MessageContent("Napotkano błąd podczas rozpozczynania turnieju", true))
				return
			}

			event.CreateMessage(MessageContent("Rozpoczynam turniej", true))
			return
		}
	}
}
