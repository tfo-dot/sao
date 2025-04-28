package discord

import (
	"fmt"
	"sao/data"
	"sao/player"
	"sao/types"
	"sao/world/party"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/google/uuid"
	"slices"
)

func ModalSubmitHandler(event *events.ModalSubmitInteractionCreate) {
	if event.Data.CustomID != "shop/buy" {
		return
	}

	var componentCustomId string

	for _, comp := range event.Data.Components {
		componentCustomId = comp.ID()
	}

	segments := strings.Split(componentCustomId, "/")
	store := data.Shops[uuid.MustParse(segments[2])]
	itemIdx, _ := strconv.Atoi(segments[3])

	stringInput, _ := event.Data.TextInputComponent(componentCustomId)

	amount, _ := strconv.Atoi(stringInput.Value)
	stockItem := store.Stock[itemIdx]

	if amount <= 0 {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Nieprawidłowa ilość").Build())
		return
	}

	player := World.GetPlayer(event.User().ID.String())

	if player == nil {
		event.CreateMessage(noCharMessage)
		return
	}

	if player.Inventory.Gold < stockItem.Count*amount {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Za mało pieniędzy").Build())
		return
	}

	player.Inventory.Gold -= stockItem.Count * amount

	for range amount {
		item := data.Items[stockItem.Item]

		item.Count = amount

		player.Inventory.AddItem(&item)
	}

	event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Zakupiono").Build())
}

func HandleChoice(event *events.ComponentInteractionCreate, id string) {
	for idx, value := range Choices {
		if value.Id == id {
			value.Select(event)

			Choices = slices.Delete(Choices, idx, idx+1)
			return
		}
	}
}

func HandleSkillUpgrade(event *events.ComponentInteractionCreate) {
	data := strings.Split(event.ComponentInteraction.Data.CustomID(), "|")

	lvl := 0

	fmt.Sscanf(data[1], "%d", &lvl)

	upgrade := 0

	fmt.Sscanf(data[2], "%d", &upgrade)

	err := World.GetPlayer(event.Member().User.ID.String()).UpgradeSkill(lvl, upgrade)

	if err == nil {
		event.UpdateMessage(
			discord.
				NewMessageUpdateBuilder().
				SetContent("Odblokowano ulepszenie").
				ClearContainerComponents().
				ClearEmbeds().
				Build(),
		)

		return
	}

	msgContent := ""

	switch err.Error() {
	case "INVALID_UPGRADE":
		msgContent = "Nie prawidłowe ulepszenie"
	case "SKILL_NOT_FOUND":
		msgContent = "Nie znaleziono umiejętności"
	case "UPGRADE_ALREADY_UNLOCKED":
		msgContent = "Ulepszenie zostało już odblokowane"
	case "NO_ACTIONS_AVAILABLE":
		msgContent = "Wszystkie akcje zostały już wykorzystane"
	}

	event.CreateMessage(discord.MessageCreate{Content: msgContent})
	return
}

func HandleSkillUnlock(event *events.ComponentInteractionCreate) {
	data := strings.Split(event.ComponentInteraction.Data.CustomID(), "|")

	pathParsed, err := strconv.Atoi(data[1])

	if err != nil {
		event.CreateMessage(discord.MessageCreate{Content: "Nieznany błąd (odblokowanie umiejętności)"})
		return
	}

	path := types.SkillPath(pathParsed)

	lvl := 0

	fmt.Sscanf(data[2], "%d", &lvl)

	pl := World.GetPlayer(event.Member().User.ID.String())

	choice := 0
	fmt.Sscanf(data[3], "%d", &choice)

	err = pl.UnlockSkill(path, lvl, choice)

	if err == nil {
		event.UpdateMessage(discord.
			NewMessageUpdateBuilder().
			SetContent("Odblokowano umiejętność").
			ClearContainerComponents().
			ClearEmbeds().
			Build(),
		)

		return
	}

	msgContent := ""

	switch err.Error() {
	case "PLAYER_LVL_TOO_LOW":
		msgContent = "Nie masz wystarczającego poziomu"
	case "SKILL_ALREADY_UNLOCKED":
		msgContent = "Umiejętność jest już odblokowana"
	case "SKILL_NOT_FOUND":
		msgContent = "Nie znaleziono umiejętności"
	case "INVALID_CHOICE":
		msgContent = "Nie znaleziono umiejętności"
	case "NO_ACTIONS_AVAILABLE":
		msgContent = "Wszystkie akcje zostały już wykorzystane"
	}

	event.CreateMessage(discord.NewMessageCreateBuilder().SetContent(msgContent).Build())
	return
}

func ComponentHandler(event *events.ComponentInteractionCreate) {
	customId := event.ComponentInteraction.Data.CustomID()

	if customId == "utils|wait" {
		event.CreateMessage(
			MessageContent("Wyślę ci prywatną wiadomość gdy twoja postać będzie miała 100% HP", true),
		)

		World.GetPlayer(event.User().ID.String()).Meta.WaitToHeal = true

		return
	}

	if strings.HasPrefix(customId, "chc/") {
		HandleChoice(event, strings.TrimPrefix(customId, "chc/"))
		return
	}

	if strings.HasPrefix(customId, "sup") {
		HandleSkillUpgrade(event)
		return
	}

	if strings.HasPrefix(customId, "su") {
		HandleSkillUnlock(event)
		return
	}

	if strings.HasPrefix(customId, "party") {
		if strings.HasPrefix(customId, "party/res") {
			pl := World.GetPlayer(event.User().ID.String())
			partyUuid := uuid.MustParse(strings.Split(customId, "|")[1])

			if pl.Meta.Party != nil {
				event.CreateMessage(alreadyInParty)
				return
			}

			World.Parties[partyUuid].Players = append(World.Parties[partyUuid].Players, &party.PartyEntry{
				PlayerUuid: pl.GetUUID(),
				Role:       party.None,
			})

			for _, player := range World.Parties[partyUuid].Players {
				World.Players[player.PlayerUuid].Meta.Party.MembersCount = len(World.Parties[partyUuid].Players)
			}

			pl.Meta.Party = &player.PartialParty{
				UUID:         partyUuid,
				Role:         party.None,
				MembersCount: len(World.Parties[partyUuid].Players),
			}

			event.CreateMessage(MessageContent("Dołączono do party", true))
			return
		}

		if customId == "party/rej" {
			event.UpdateMessage(discord.
				NewMessageUpdateBuilder().
				ClearContainerComponents().
				SetContent("Odrzucono zaproszenie").
				Build(),
			)
		}
	}

	if strings.HasPrefix(customId, "f") {
		action := strings.Split(customId, "/")[1]
		userIdTurn := strings.Split(customId, "/")[2]

		player := World.GetPlayer(event.User().ID.String())

		if player == nil {
			event.CreateMessage(noCharMessage)
			return
		}

		if player.Meta.FightInstance == nil {
			event.CreateMessage(fightNotYoursMessage)
			return
		}

		fight, ok := World.Fights[*player.Meta.FightInstance]

		if !ok {
			event.CreateMessage(MessageContent("Walka nie istnieje...", true))
		}

		if fight.IsFinished() {
			event.CreateMessage(fightAlreadyEndedMessage)
			return
		}

		if player.GetUID() != userIdTurn {
			event.CreateMessage(notYourTurnMessage)
			return
		}

		switch action {
		case "attack":
			playerEnemies := fight.GetEnemiesFor(player.GetUUID())

			if len(playerEnemies) == 0 {
				event.CreateMessage(MessageContent("Brak przeciwników", true))
				return
			}

			if len(playerEnemies) == 1 {
				fight.PlayerActions <- types.Action{
					Event:  types.ACTION_ATTACK,
					Source: player.GetUUID(),
					Target: playerEnemies[0].GetUUID(),
				}

				event.CreateMessage(MessageContent("Zaatakowano", true))
				return
			}

			selectMenuUuid := uuid.New().String()

			options := make([]discord.StringSelectMenuOption, 0)

			for _, enemy := range playerEnemies {
				options = append(options, discord.NewStringSelectMenuOption(enemy.GetName(), enemy.GetUUID().String()))
			}

			Choices = append(Choices, types.DiscordChoice{
				Id: selectMenuUuid,
				Select: func(event *events.ComponentInteractionCreate) {
					selected := event.StringSelectMenuInteractionData().Values

					parsedUuids := make([]uuid.UUID, len(selected))

					for idx, rawUuid := range selected {
						parsedUuids[idx] = uuid.MustParse(rawUuid)
					}

					fight.PlayerActions <- types.Action{
						Event:  types.ACTION_ATTACK,
						Source: player.GetUUID(),
						Target: parsedUuids[0],
					}

					event.UpdateMessage(messageUpdateClearComponents)
				},
			})

			selectMenu := discord.NewStringSelectMenu("chc/"+selectMenuUuid, "Wybierz wrogów").AddOptions(options...).WithMaxValues(1)

			event.UpdateMessage(discord.
				NewMessageUpdateBuilder().
				SetContent("Wybierz cel").
				ClearEmbeds().
				ClearContainerComponents().
				AddActionRow(selectMenu).
				Build(),
			)

			return
		case "defend":
			event.UpdateMessage(messageUpdateClearComponents)

			fight.PlayerActions <- types.Action{Event: types.ACTION_DEFEND, Source: player.GetUUID()}

			event.Acknowledge()
			return
		case "skill":
			options := make([]discord.StringSelectMenuOption, 0)

			for _, skill := range player.Inventory.LevelSkills {
				if player.CanUseSkill(skill.Skill) {
					if !skill.Skill.CanUse(player, fight) {
						continue
					}
					options = append(options, discord.NewStringSelectMenuOption(skill.Skill.GetName(), fmt.Sprintf("l|%d", skill.Skill.GetLevel())))
				}
			}

			selectMenuUuid := uuid.New().String()

			Choices = append(Choices, types.DiscordChoice{
				Id: selectMenuUuid,
				Select: func(event *events.ComponentInteractionCreate) {
					selected := event.StringSelectMenuInteractionData().Values[0]

					value, err := strconv.Atoi(strings.Split(selected, "|")[1])

					if err != nil {
						event.CreateMessage(unknownError)
						return
					}

					skillObj, skillExists := player.Inventory.LevelSkills[value]

					if !skillExists {
						event.CreateMessage(unknownError)
						return
					}

					skillTrigger := skillObj.Skill.GetUpgradableTrigger(skillObj.Upgrades)

					if skillTrigger.Target == nil || skillTrigger.Target.Target == types.TARGET_SELF {
						skipTurn := skillTrigger.Flags&types.FLAG_INSTANT_SKILL == 0

						fight.PlayerActions <- types.Action{
							Event:       types.ACTION_SKILL,
							Source:      player.GetUUID(),
							ConsumeTurn: &skipTurn,
							Meta:        types.ActionSkillMeta{IsForLevel: true, Lvl: value},
						}

						event.UpdateMessage(messageUpdateClearComponents)
						return
					}

					if skillTrigger.Target.Target == types.TARGET_ENEMY {
						playerEnemies := fight.GetEnemiesFor(player.GetUUID())

						if len(playerEnemies) == 0 {
							event.CreateMessage(MessageContent("Brak przeciwników", true))
							return
						}

						selectMenuUuidDeep := uuid.New().String()

						options := make([]discord.StringSelectMenuOption, 0)

						for _, enemy := range playerEnemies {
							options = append(options, discord.NewStringSelectMenuOption(enemy.GetName(), enemy.GetUUID().String()))
						}

						Choices = append(Choices, types.DiscordChoice{
							Id: selectMenuUuidDeep,
							Select: func(event *events.ComponentInteractionCreate) {

								selected := event.StringSelectMenuInteractionData().Values

								parsedUuids := make([]uuid.UUID, len(selected))

								for idx, rawUuid := range selected {
									parsedUuids[idx] = uuid.MustParse(rawUuid)
								}

								skipTurn := skillTrigger.Flags&types.FLAG_INSTANT_SKILL == 0

								fight.PlayerActions <- types.Action{
									Event:       types.ACTION_SKILL,
									Source:      player.GetUUID(),
									Target:      player.GetUUID(),
									ConsumeTurn: &skipTurn,
									Meta:        types.ActionSkillMeta{IsForLevel: true, Lvl: value, Targets: parsedUuids},
								}

								event.UpdateMessage(messageUpdateClearComponents)
							},
						})

						selectMenu := discord.NewStringSelectMenu("chc/"+selectMenuUuidDeep, "Wybierz wrogów").AddOptions(options...)

						if skillTrigger.Target.MaxTargets >= 0 {
							selectMenu = selectMenu.WithMaxValues(skillTrigger.Target.MaxTargets)
						} else {
							selectMenu = selectMenu.WithMaxValues(len(options))
						}

						event.UpdateMessage(
							discord.NewMessageUpdateBuilder().ClearContainerComponents().AddActionRow(selectMenu).Build(),
						)

						return
					}

					if skillTrigger.Target.Target == types.TARGET_ALLY {
						playerAllies := fight.GetEnemiesFor(player.GetUUID())

						if len(playerAllies) == 0 {
							event.CreateMessage(MessageContent("Brak sojuszników", true))
							return
						}

						selectMenuUuidDeep := uuid.New().String()

						options := make([]discord.StringSelectMenuOption, 0)

						for _, ally := range playerAllies {
							options = append(options, discord.NewStringSelectMenuOption(ally.GetName(), ally.GetUUID().String()))
						}

						Choices = append(Choices, types.DiscordChoice{
							Id: selectMenuUuidDeep,
							Select: func(event *events.ComponentInteractionCreate) {
								selected := event.StringSelectMenuInteractionData().Values

								parsedUuids := make([]uuid.UUID, len(selected))

								for idx, rawUuid := range selected {
									parsedUuids[idx] = uuid.MustParse(rawUuid)
								}

								skipTurn := skillTrigger.Flags&types.FLAG_INSTANT_SKILL == 0

								fight.PlayerActions <- types.Action{
									Event:       types.ACTION_SKILL,
									Source:      player.GetUUID(),
									Target:      player.GetUUID(),
									ConsumeTurn: &skipTurn,
									Meta:        types.ActionSkillMeta{IsForLevel: true, Lvl: value, Targets: parsedUuids},
								}

								event.UpdateMessage(messageUpdateClearComponents)
							},
						})

						selectMenu := discord.NewStringSelectMenu("chc/"+selectMenuUuidDeep, "Wybierz sojuszników").AddOptions(options...)

						if skillTrigger.Target.MaxTargets >= 0 {
							selectMenu = selectMenu.WithMaxValues(skillTrigger.Target.MaxTargets)
						} else {
							selectMenu = selectMenu.WithMaxValues(len(options))
						}

						event.UpdateMessage(
							discord.
								NewMessageUpdateBuilder().
								ClearContainerComponents().
								AddActionRow(selectMenu).
								Build(),
						)

						return
					}

					skillTargets := make([]types.Entity, 0)

					if skillTrigger.Target.Target&types.TARGET_SELF != 0 {
						skillTargets = append(skillTargets, player)
					}

					if skillTrigger.Target.Target&types.TARGET_ENEMY != 0 {
						skillTargets = append(skillTargets, fight.GetEnemiesFor(player.GetUUID())...)
					}

					if skillTrigger.Target.Target&types.TARGET_ALLY != 0 {
						skillTargets = append(skillTargets, fight.GetAlliesFor(player.GetUUID())...)
					}

					selectMenuUuidDeep := uuid.New().String()

					options := make([]discord.StringSelectMenuOption, 0)

					for _, target := range skillTargets {
						options = append(options, discord.NewStringSelectMenuOption(target.GetName(), target.GetUUID().String()))
					}

					Choices = append(Choices, types.DiscordChoice{
						Id: selectMenuUuidDeep,
						Select: func(event *events.ComponentInteractionCreate) {

							selected := event.StringSelectMenuInteractionData().Values

							parsedUuids := make([]uuid.UUID, len(selected))

							for idx, rawUuid := range selected {
								parsedUuids[idx] = uuid.MustParse(rawUuid)
							}

							skipTurn := skillTrigger.Flags&types.FLAG_INSTANT_SKILL == 0

							fight.PlayerActions <- types.Action{
								Event:       types.ACTION_SKILL,
								Source:      player.GetUUID(),
								ConsumeTurn: &skipTurn,
								Meta:        types.ActionSkillMeta{IsForLevel: true, Lvl: value, Targets: parsedUuids},
							}

							event.UpdateMessage(messageUpdateClearComponents)
						},
					})

					selectMenu := discord.NewStringSelectMenu("chc/"+selectMenuUuidDeep, "Wybierz cele").AddOptions(options...)

					if skillTrigger.Target.MaxTargets >= 0 {
						selectMenu = selectMenu.WithMaxValues(skillTrigger.Target.MaxTargets)
					} else {
						selectMenu = selectMenu.WithMaxValues(len(options))
					}

					event.UpdateMessage(
						discord.NewMessageUpdateBuilder().ClearContainerComponents().AddActionRow(selectMenu).Build(),
					)

					return
				},
			})

			event.UpdateMessage(discord.
				NewMessageUpdateBuilder().
				ClearContainerComponents().
				AddActionRow(discord.
					NewStringSelectMenu("chc/"+selectMenuUuid, "Wybierz umiejętność").
					WithMaxValues(1).
					AddOptions(options...),
				).
				Build(),
			)
		case "item":
			options := make([]discord.StringSelectMenuOption, 0)

			for _, item := range player.Inventory.Items {
				if item.Hidden {
					continue
				}

				if !item.Consume && item.Count <= 0 {
					continue
				}

				options = append(options, discord.NewStringSelectMenuOption(item.Name, item.UUID.String()))
			}

			event.UpdateMessage(
				discord.
					NewMessageUpdateBuilder().
					ClearContainerComponents().
					AddActionRow(
						discord.
							NewStringSelectMenu(
								"f/item/usage",
								"Wybierz przedmiot",
							).
							WithMaxValues(1).
							AddOptions(options...),
					).
					Build(),
			)
		case "item/usage":
			rawItemUuid := event.ComponentInteraction.StringSelectMenuInteractionData().Values[0]

			itemUuid := uuid.MustParse(rawItemUuid)

			for _, item := range player.Inventory.Items {
				if item.UUID == itemUuid {
					if item.Consume {
						fight.PlayerActions <- types.Action{
							Event:  types.ACTION_ITEM,
							Source: player.GetUUID(),
							Meta:   types.ActionItemMeta{Item: item.UUID, Targets: []uuid.UUID{}},
						}

						event.CreateMessage(discord.
							NewMessageCreateBuilder().
							SetContent("Użyto przedmiotu").
							SetEphemeral(true).
							Build(),
						)

						return
					}

					event.CreateMessage(discord.
						NewMessageCreateBuilder().
						SetContent("Nie można użyć tego przedmiotu").
						SetEphemeral(true).
						Build(),
					)

					return
				}
			}
		case "escape":
			event.UpdateMessage(messageUpdateClearComponents)

			fight.PlayerActions <- types.Action{
				Event:  types.ACTION_RUN,
				Source: player.GetUUID(),
				Target: player.GetUUID(),
			}
		}
	}

	if strings.HasPrefix(customId, "shop") {
		segments := strings.Split(customId, "/")

		switch segments[1] {
		case "show":
			page, _ := strconv.Atoi(segments[2])
			pageStart := (page - 1) * 5
			pageEnd := page * 5
			store := data.Shops[uuid.MustParse(segments[3])]

			message := discord.NewMessageCreateBuilder()

			embed := discord.NewEmbedBuilder()

			embed.SetTitle("Sklep: " + store.Name)

			var pageStock []*types.WithCount[uuid.UUID]

			if pageEnd > len(store.Stock) {
				pageStock = store.Stock[pageStart:]
			} else {
				pageStock = store.Stock[pageStart:pageEnd]
			}

			productButtons := make([]discord.InteractiveComponent, 0)

			for itemIdx, stock := range pageStock {
				itemName := data.Items[stock.Item].Name

				productButton := discord.NewPrimaryButton(itemName, "shop/buy/"+segments[3]+"/"+fmt.Sprint(itemIdx))

				embed.AddField(itemName, fmt.Sprintf("Cena: %d", stock.Count), true)
				productButtons = append(productButtons, productButton)
			}

			embed.SetFooterTextf("Strona %d/%d", page, (len(store.Stock)/5)+1)

			message.AddEmbeds(embed.Build())
			message.AddActionRow(productButtons...)

			prevPageButton := discord.NewPrimaryButton("Poprzednia strona", "shop/show/"+fmt.Sprint(page-1)+"/"+segments[3])

			if page-1 < 1 {
				prevPageButton = prevPageButton.AsDisabled()
			}

			nextPageButton := discord.NewPrimaryButton("Następna strona", "shop/show/"+fmt.Sprint(page+1)+"/"+segments[3])

			if pageEnd-1 > len(store.Stock) {
				nextPageButton = nextPageButton.AsDisabled()
			}

			message.AddActionRow(prevPageButton, nextPageButton)

			err := event.CreateMessage(message.Build())

			panic(err)
		case "buy":
			itemIdx, _ := strconv.Atoi(segments[3])

			event.Modal(discord.
				NewModalCreateBuilder().
				SetTitle("Kupno").
				SetCustomID("shop/buy").
				AddActionRow(
					discord.NewShortTextInput("shop/buy/"+segments[2]+"/"+fmt.Sprint(itemIdx), "Ilość"),
				).Build(),
			)
		}
	}

	if strings.HasPrefix(customId, "t") {
		segments := strings.Split(customId, "/")

		switch segments[1] {
		case "join":
			tUuid := uuid.MustParse(segments[2])

			pl := World.GetPlayer(event.Member().User.ID.String())

			if pl == nil {
				event.CreateMessage(noCharMessage)
				return
			}

			if slices.Contains(World.Tournaments[tUuid].Participants, pl.GetUUID()) {
				event.CreateMessage(MessageContent("Jesteś już zapisany", true))
				return
			}

			World.Tournaments[tUuid].Participants = append(World.Tournaments[tUuid].Participants, pl.GetUUID())

			var playerText string = ""

			if World.Tournaments[tUuid].MaxPlayers == -1 {
				playerText = fmt.Sprintf("Nieograniczona (%v graczy)", len(World.Tournaments[tUuid].Participants))
			} else {
				playerText = fmt.Sprintf("%v/%v", len(World.Tournaments[tUuid].Participants), World.Tournaments[tUuid].MaxPlayers)
			}

			event.UpdateMessage(discord.NewMessageUpdateBuilder().
				SetEmbeds(discord.NewEmbedBuilder().
					SetTitle("Nowy turniej!").
					SetDescriptionf("Zapisy na turniej `%v` otwarte!", World.Tournaments[tUuid].Name).
					SetFooterText("Ilość miejsc: " + playerText).
					Build()).
				Build(),
			)
		}
	}
}
