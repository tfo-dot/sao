package discord

import (
	"fmt"
	"sao/battle"
	"sao/data"
	"sao/player"
	"sao/types"
	"sao/world/npc"
	"sao/world/party"
	"sao/world/transaction"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/google/uuid"
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
	store := World.Stores[uuid.MustParse(segments[2])]
	itemIdx, _ := strconv.Atoi(segments[3])

	stringInput, _ := event.Data.TextInputComponent(componentCustomId)

	amount, _ := strconv.Atoi(stringInput.Value)
	stockItem := store.Stock[itemIdx]

	if amount <= 0 {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Nieprawidłowa ilość").Build())
		return
	}

	event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Zakupiono").Build())

	if stockItem.Quantity < amount {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Za mało towaru").Build())
		return
	}

	player := World.GetPlayer(event.User().ID.String())

	if player == nil {
		event.CreateMessage(noCharMessage)
		return
	}

	if player.Inventory.Gold < stockItem.Price*amount {
		event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Za mało pieniędzy").Build())
		return
	}

	player.Inventory.Gold -= stockItem.Price * amount
	stockItem.Quantity -= amount

	for i := 0; i < amount; i++ {
		if stockItem.ItemType == types.ITEM_MATERIAL {
			item := data.Ingredients[stockItem.ItemUUID]

			item.Count = amount

			player.Inventory.AddIngredient(&item)
		} else {
			item := data.Items[stockItem.ItemUUID]

			item.Count = amount

			player.Inventory.AddItem(&item)
		}
	}

	event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Zakupiono").Build())
}

func ComponentHandler(event *events.ComponentInteractionCreate) {
	customId := event.ComponentInteraction.Data.CustomID()

	if customId == "utils|wait" {
		event.CreateMessage(
			discord.
				NewMessageCreateBuilder().
				SetContent("Wyślę ci prywatną wiadomość gdy twoja postać będzie miała 100% HP").
				SetEphemeral(true).
				Build(),
		)

		player := World.GetPlayer(event.User().ID.String())

		player.Meta.WaitToHeal = true

		return
	}

	if strings.HasPrefix(customId, "chc/") {
		customId, _ = strings.CutPrefix(customId, "chc/")

		for idx, value := range Choices {
			if value.Id == customId {
				value.Select(event)

				Choices = append(Choices[:idx], Choices[idx+1:]...)
				return
			}
		}

		return
	}

	if strings.HasPrefix(customId, "sup") {
		data := strings.Split(customId, "|")

		lvl := 0

		fmt.Sscanf(data[1], "%d", &lvl)

		upgrade := 0

		fmt.Sscanf(data[2], "%d", &upgrade)

		userSnowflake := event.Member().User.ID.String()

		for _, pl := range World.Players {
			if pl.Meta.UserID == userSnowflake {

				res := pl.UpgradeSkill(lvl, upgrade)

				if res == nil {
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

				switch res.Error() {
				case "INVALID_UPGRADE":
					msgContent = "Nie prawidłowe ulepszenie"
				case "SKILL_NOT_FOUND":
					msgContent = "Nie znaleziono umiejętności"
				case "UPGRADE_ALREADY_UNLOCKED":
					msgContent = "Ulepszenie zostało już odblokowane"
				case "NO_ACTIONS_AVAILABLE":
					msgContent = "Wszystkie akcje zostały już wykorzystane"
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

	if strings.HasPrefix(customId, "su") {
		data := strings.Split(customId, "|")

		path := types.SkillPath(0)

		switch data[1] {
		case "0":
			path = types.PathControl
		case "1":
			path = types.PathEndurance
		case "2":
			path = types.PathDamage
		case "3":
			path = types.PathSpecial
		}

		lvl := 0

		fmt.Sscanf(data[2], "%d", &lvl)

		userSnowflake := event.Member().User.ID.String()

		for _, pl := range World.Players {
			if pl.Meta.UserID == userSnowflake {
				if len(data) == 4 {
					choice := 0
					fmt.Sscanf(data[3], "%d", &choice)

					res := pl.UnlockSkill(path, lvl, choice)

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
					case "INVALID_CHOICE":
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

				res := pl.UnlockSkill(path, lvl, 0)

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
				case "NO_ACTIONS_AVAILABLE":
					msgContent = "Wszystkie akcje zostały już wykorzystane"
				}

				if res.Error() == "MULTIPLE_OPTIONS" {
					options := make([]discord.StringSelectMenuOption, 0)

					selectMenuUuid := uuid.New().String()

					Choices = append(Choices, types.DiscordChoice{
						Id: selectMenuUuid,
						Select: func(event *events.ComponentInteractionCreate) {
							choice := event.StringSelectMenuInteractionData().Values[0]

							parsed := 0
							fmt.Sscanf(choice, "%d", &parsed)

							res := pl.UnlockSkill(path, lvl, parsed)

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
							case "INVALID_CHOICE":
								msgContent = "Nie znaleziono umiejętności"
							case "NO_ACTIONS_AVAILABLE":
								msgContent = "Wszystkie akcje zostały już wykorzystane"
							}

							event.CreateMessage(
								discord.
									NewMessageCreateBuilder().
									SetContent(msgContent).
									Build(),
							)
						},
					})

					event.UpdateMessage(
						discord.
							NewMessageUpdateBuilder().
							ClearContainerComponents().
							AddActionRow(
								discord.
									NewStringSelectMenu("chc/"+selectMenuUuid, "Wybierz umiejętność").
									WithMaxValues(1).
									AddOptions(options...),
							).
							Build(),
					)
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

	if strings.HasPrefix(customId, "party") {
		if strings.HasPrefix(customId, "party/res") {
			userSnowflake := event.User().ID.String()

			partyUuid := uuid.MustParse(strings.Split(customId, "|")[1])

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

					World.Parties[partyUuid].Players = append(World.Parties[partyUuid].Players, &party.PartyEntry{
						PlayerUuid: pl.GetUUID(),
						Role:       party.None,
					})

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

		if customId == "party/rej" {
			event.UpdateMessage(
				discord.
					NewMessageUpdateBuilder().
					ClearContainerComponents().
					SetContent("Odrzucono zaproszenie").
					Build(),
			)
		}
	}

	if strings.HasPrefix(customId, "f") {
		action := customId[2:]

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

		if !ok && action != "save" {
			event.UpdateMessage(messageUpdateClearComponents)

			(*Client).Rest().CreateMessage(event.Channel().ID(),
				discord.
					NewMessageCreateBuilder().
					SetContent("Wykryto spaghetti od <@344048874656366592>...").
					Build(),
			)
		}

		if fight.IsFinished() && action != "save" {
			event.CreateMessage(fightAlreadyEndedMessage)
			return
		}

		authorName := event.Message.Embeds[0].Author.Name

		if player.GetName() != authorName && action != "save" {
			event.CreateMessage(notYourTurnMessage)
			return
		}

		switch action {
		case "attack":
			playerEnemies := fight.GetEnemiesFor(player.GetUUID())

			if len(playerEnemies) == 0 {
				event.CreateMessage(
					discord.
						NewMessageCreateBuilder().
						SetContent("Brak przeciwników").
						SetEphemeral(true).
						Build(),
				)

				return
			}

			if len(playerEnemies) == 1 {
				fight.PlayerActions <- battle.Action{
					Event:  battle.ACTION_ATTACK,
					Source: player.GetUUID(),
					Target: playerEnemies[0].GetUUID(),
				}

				event.CreateMessage(
					discord.
						NewMessageCreateBuilder().
						SetContent("Zaatakowano!").
						SetEphemeral(true).
						Build(),
				)

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

					fight.PlayerActions <- battle.Action{
						Event:  battle.ACTION_ATTACK,
						Source: player.GetUUID(),
						Target: parsedUuids[0],
					}

					event.UpdateMessage(messageUpdateClearComponents)
				},
			})

			selectMenu := discord.NewStringSelectMenu("chc/"+selectMenuUuid, "Wybierz wrogów").AddOptions(options...).WithMaxValues(1)

			event.UpdateMessage(
				discord.
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

			fight.PlayerActions <- battle.Action{
				Event:  battle.ACTION_DEFEND,
				Source: player.GetUUID(),
				Target: player.GetUUID(),
			}

			event.CreateMessage(
				discord.
					NewMessageCreateBuilder().
					SetContent("Przygotowujesz się na atak!").
					SetEphemeral(true).
					Build(),
			)

			return
		case "skill":
			options := make([]discord.StringSelectMenuOption, 0)

			for _, skill := range player.Inventory.LevelSkills {
				if player.CanUseSkill(skill) {
					if !skill.CanUse(&player, &fight, player.Inventory.LevelSkillsUpgrades[skill.GetLevel()]) {
						continue
					}
					options = append(options, discord.NewStringSelectMenuOption(skill.GetName(), fmt.Sprintf("l|%d", skill.GetLevel())))
				}
			}

			if player.Meta.Fury != nil {
				for _, skill := range player.Meta.Fury.GetSkills() {
					if player.CanUseSkill(skill) {
						options = append(options, discord.NewStringSelectMenuOption(skill.GetName(), "f|"+skill.GetUUID().String()))
					}
				}
			}

			selectMenuUuid := uuid.New().String()

			Choices = append(Choices, types.DiscordChoice{
				Id: selectMenuUuid,
				Select: func(event *events.ComponentInteractionCreate) {
					selected := event.StringSelectMenuInteractionData().Values[0]

					switch selected[0] {
					case 'l':
						value, err := strconv.Atoi(strings.Split(selected, "|")[1])

						if err != nil {
							event.CreateMessage(
								discord.
									NewMessageCreateBuilder().
									SetContent("Wystapił nieoczekiwany błąd").
									SetEphemeral(true).
									Build(),
							)
						}

						skillObj, skillExists := player.Inventory.LevelSkills[value]

						if !skillExists {
							event.CreateMessage(
								discord.
									NewMessageCreateBuilder().
									SetContent("Wystapił nieoczekiwany błąd").
									SetEphemeral(true).
									Build(),
							)
						}

						skillTrigger := skillObj.GetUpgradableTrigger(player.Inventory.LevelSkillsUpgrades[value])

						if skillTrigger.Target == nil || skillTrigger.Target.Target == types.TARGET_SELF {
							skipTurn := skillTrigger.Flags&types.FLAG_INSTANT_SKILL == 0

							fight.PlayerActions <- battle.Action{
								Event:       battle.ACTION_SKILL,
								Source:      player.GetUUID(),
								Target:      player.GetUUID(),
								ConsumeTurn: &skipTurn,
								Meta: battle.ActionSkillMeta{
									IsForLevel: true,
									Lvl:        skillObj.GetLevel(),
								},
							}

							event.UpdateMessage(messageUpdateClearComponents)

							return
						}

						if skillTrigger.Target.Target == types.TARGET_ENEMY {
							playerEnemies := fight.GetEnemiesFor(player.GetUUID())

							if len(playerEnemies) == 0 {
								event.CreateMessage(
									discord.
										NewMessageCreateBuilder().
										SetContent("Brak przeciwników").
										SetEphemeral(true).
										Build(),
								)
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

									fight.PlayerActions <- battle.Action{
										Event:       battle.ACTION_SKILL,
										Source:      player.GetUUID(),
										Target:      player.GetUUID(),
										ConsumeTurn: &skipTurn,
										Meta: battle.ActionSkillMeta{
											IsForLevel: true,
											Lvl:        skillObj.GetLevel(),
											Targets:    parsedUuids,
										},
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
								discord.
									NewMessageUpdateBuilder().
									ClearContainerComponents().
									AddActionRow(selectMenu).
									Build(),
							)

							return
						}

						if skillTrigger.Target.Target == types.TARGET_ALLY {
							playerAllies := fight.GetEnemiesFor(player.GetUUID())

							if len(playerAllies) == 0 {
								event.CreateMessage(
									discord.
										NewMessageCreateBuilder().
										SetContent("Brak sojuszników").
										SetEphemeral(true).
										Build(),
								)
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

									fight.PlayerActions <- battle.Action{
										Event:       battle.ACTION_SKILL,
										Source:      player.GetUUID(),
										Target:      player.GetUUID(),
										ConsumeTurn: &skipTurn,
										Meta: battle.ActionSkillMeta{
											IsForLevel: true,
											Lvl:        skillObj.GetLevel(),
											Targets:    parsedUuids,
										},
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

						skillTargets := make([]battle.Entity, 0)

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

								fight.PlayerActions <- battle.Action{
									Event:       battle.ACTION_SKILL,
									Source:      player.GetUUID(),
									Target:      player.GetUUID(),
									ConsumeTurn: &skipTurn,
									Meta: battle.ActionSkillMeta{
										IsForLevel: true,
										Lvl:        skillObj.GetLevel(),
										Targets:    parsedUuids,
									},
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
							discord.
								NewMessageUpdateBuilder().
								ClearContainerComponents().
								AddActionRow(selectMenu).
								Build(),
						)

						return

					case 'f':
						value := uuid.MustParse(strings.Split(selected, "|")[1])

						if player.Meta.Fury != nil {
							event.CreateMessage(
								discord.
									NewMessageCreateBuilder().
									SetContent("Wystapił nieoczekiwany błąd").
									SetEphemeral(true).
									Build(),
							)
						}

						var skillObj types.PlayerSkill

						for _, skill := range player.Meta.Fury.GetSkills() {
							if skill.GetUUID() == value {
								skillObj = skill
							}
						}

						skillTrigger := skillObj.GetTrigger()

						if skillTrigger.Target == nil || skillTrigger.Target.Target == types.TARGET_SELF {
							skipTurn := skillTrigger.Flags&types.FLAG_INSTANT_SKILL == 0

							fight.PlayerActions <- battle.Action{
								Event:       battle.ACTION_SKILL,
								Source:      player.GetUUID(),
								Target:      player.GetUUID(),
								ConsumeTurn: &skipTurn,
								Meta: battle.ActionSkillMeta{
									IsForLevel: false,
									SkillUuid:  skillObj.GetUUID(),
								},
							}

							event.UpdateMessage(messageUpdateClearComponents)

							return
						}

						if skillTrigger.Target.Target == types.TARGET_ENEMY {
							playerEnemies := fight.GetEnemiesFor(player.GetUUID())

							if len(playerEnemies) == 0 {
								event.CreateMessage(
									discord.
										NewMessageCreateBuilder().
										SetContent("Brak przeciwników").
										SetEphemeral(true).
										Build(),
								)
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

									fight.PlayerActions <- battle.Action{
										Event:       battle.ACTION_SKILL,
										Source:      player.GetUUID(),
										Target:      player.GetUUID(),
										ConsumeTurn: &skipTurn,
										Meta: battle.ActionSkillMeta{
											IsForLevel: false,
											SkillUuid:  skillObj.GetUUID(),
											Targets:    parsedUuids,
										},
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
								discord.
									NewMessageUpdateBuilder().
									ClearContainerComponents().
									AddActionRow(selectMenu).
									Build(),
							)

							return
						}

						if skillTrigger.Target.Target == types.TARGET_ALLY {
							playerAllies := fight.GetEnemiesFor(player.GetUUID())

							if len(playerAllies) == 0 {
								event.CreateMessage(
									discord.
										NewMessageCreateBuilder().
										SetContent("Brak sojuszników").
										SetEphemeral(true).
										Build(),
								)
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

									fight.PlayerActions <- battle.Action{
										Event:       battle.ACTION_SKILL,
										Source:      player.GetUUID(),
										Target:      player.GetUUID(),
										ConsumeTurn: &skipTurn,
										Meta: battle.ActionSkillMeta{
											IsForLevel: false,
											SkillUuid:  skillObj.GetUUID(),
											Targets:    parsedUuids,
										},
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

						skillTargets := make([]battle.Entity, 0)

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

								fight.PlayerActions <- battle.Action{
									Event:       battle.ACTION_SKILL,
									Source:      player.GetUUID(),
									Target:      player.GetUUID(),
									ConsumeTurn: &skipTurn,
									Meta: battle.ActionSkillMeta{
										IsForLevel: false,
										SkillUuid:  skillObj.GetUUID(),
										Targets:    parsedUuids,
									},
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
							discord.
								NewMessageUpdateBuilder().
								ClearContainerComponents().
								AddActionRow(selectMenu).
								Build(),
						)

						return

					}
				},
			})

			event.UpdateMessage(
				discord.
					NewMessageUpdateBuilder().
					ClearContainerComponents().
					AddActionRow(
						discord.
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
						event.UpdateMessage(
							discord.
								NewMessageUpdateBuilder().
								ClearContainerComponents().
								Build(),
						)

						fight.PlayerActions <- battle.Action{
							Event:  battle.ACTION_ITEM,
							Source: player.GetUUID(),
							Target: player.GetUUID(),
							Meta: battle.ActionItemMeta{
								Item:    item.UUID,
								Targets: []uuid.UUID{},
							},
						}

						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Użyto przedmiotu").
								SetEphemeral(true).
								Build(),
						)

						return
					}

					event.CreateMessage(
						discord.
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

			fight.PlayerActions <- battle.Action{
				Event:  battle.ACTION_RUN,
				Source: player.GetUUID(),
				Target: player.GetUUID(),
			}
		case "save":
			event.UpdateMessage(
				discord.
					NewMessageUpdateBuilder().
					SetContent("Wykryto spaghetti od <@344048874656366592> aka jeszcze nie zrobione!...").
					Build(),
			)

			return
		}
	}

	if strings.HasPrefix(customId, "shop") {
		segments := strings.Split(customId, "/")
		action := segments[1]

		switch action {
		case "show":
			page, _ := strconv.Atoi(segments[2])
			pageStart := (page - 1) * 5
			pageEnd := page * 5
			store := World.Stores[uuid.MustParse(segments[3])]

			message := discord.NewMessageCreateBuilder()

			embed := discord.NewEmbedBuilder()

			embed.SetTitle("Sklep: " + store.Name)

			var pageStock []*npc.Stock

			if pageEnd > len(store.Stock) {
				pageStock = store.Stock[pageStart:]
			} else {
				pageStock = store.Stock[pageStart:pageEnd]
			}

			productButtons := make([]discord.InteractiveComponent, 0)

			for itemIdx, stock := range pageStock {
				var itemName string

				if stock.ItemType == types.ITEM_MATERIAL {
					itemName = data.Ingredients[stock.ItemUUID].Name
				} else {
					itemName = data.Items[stock.ItemUUID].Name
				}

				productButton := discord.NewPrimaryButton(itemName, "shop/buy/"+segments[3]+"/"+fmt.Sprint(itemIdx))

				if stock.Quantity <= 0 {
					productButton = productButton.AsDisabled()
				}

				embed.AddField(itemName, fmt.Sprintf("Cena: %d\nIlość: %d/%d", stock.Price, stock.Quantity, stock.Limit), true)
				productButtons = append(productButtons, productButton)
			}

			embed.SetFooterTextf("Ostatnia dostawa %s, strona %d/%d", store.LastRestock.String(), 1, 1)

			message.AddEmbeds(embed.Build())
			message.AddActionRow(productButtons...)

			prevPageButton := discord.NewPrimaryButton("Poprzednia strona", "shop/show/"+fmt.Sprint(page-1)+"/"+segments[3])

			if page-1 < 1 {
				prevPageButton = prevPageButton.AsDisabled()
			}

			nextPageButton := discord.NewPrimaryButton("Następna strona", "shop/show/"+fmt.Sprint(page+1)+"/"+segments[3])

			if pageEnd > len(store.Stock) {
				nextPageButton = nextPageButton.AsDisabled()
			}

			message.AddActionRow(prevPageButton, nextPageButton)

			event.CreateMessage(message.Build())
		case "buy":
			store := World.Stores[uuid.MustParse(segments[2])]
			itemIdx, _ := strconv.Atoi(segments[3])

			stockItem := store.Stock[itemIdx]

			if stockItem.Quantity <= 0 {
				event.CreateMessage(discord.NewMessageCreateBuilder().SetContent("Brak towaru").Build())
				return
			}

			modal := discord.NewModalCreateBuilder()

			modal.SetTitle("Kupno")
			modal.SetCustomID("shop/buy")
			modal.AddActionRow(
				discord.NewShortTextInput("shop/buy/"+segments[2]+"/"+fmt.Sprint(itemIdx), "Ilość"),
			)

			event.Modal(modal.Build())
		}
	}

	if strings.HasPrefix(customId, "t") {
		segments := strings.Split(customId, "/")
		action := segments[1]

		switch action {
		case "join":
			tournamentUuid := uuid.MustParse(segments[2])

			var playerChar *player.Player

			for _, pl := range World.Players {
				if pl.Meta.UserID == event.Member().User.ID.String() {
					playerChar = pl
					break
				}
			}

			if playerChar == nil {
				event.CreateMessage(noCharMessage)
				return
			}

			for _, user := range World.Tournaments[tournamentUuid].Participants {
				if user == playerChar.GetUUID() {
					event.CreateMessage(
						discord.
							NewMessageCreateBuilder().
							SetContent("Jesteś już zapisany").
							SetEphemeral(true).
							Build(),
					)
				}
			}

			World.Tournaments[tournamentUuid].Participants = append(World.Tournaments[tournamentUuid].Participants, playerChar.GetUUID())

			var playerText string = ""
			if World.Tournaments[tournamentUuid].MaxPlayers == -1 {
				playerText = fmt.Sprintf("Nieograniczona (%v graczy)", len(World.Tournaments[tournamentUuid].Participants))
			} else {
				playerText = fmt.Sprintf("%v/%v", len(World.Tournaments[tournamentUuid].Participants), World.Tournaments[tournamentUuid].MaxPlayers)
			}

			event.UpdateMessage(
				discord.NewMessageUpdateBuilder().
					SetEmbeds(discord.NewEmbedBuilder().
						SetTitle("Nowy turniej!").
						SetDescriptionf("Zapisy na turniej `%v` otwarte!", World.Tournaments[tournamentUuid].Name).
						SetFooterText("Ilość miejsc: " + playerText).
						Build()).
					Build(),
			)
		}
	}

	if strings.HasPrefix(customId, "trade") {
		segments := strings.Split(customId, "|")

		tradeUuid := uuid.MustParse(segments[1])

		if segments[0] == "trade/res" {
			trade := World.Transactions[tradeUuid]

			if trade == nil {
				event.CreateMessage(
					discord.
						NewMessageCreateBuilder().
						SetContent("Transakcja nie istnieje").
						SetEphemeral(true).
						Build(),
				)
				return
			}

			if trade.State != transaction.TransactionPending {
				event.CreateMessage(
					discord.
						NewMessageCreateBuilder().
						SetContent("Transakcja już została zaakceptowana").
						SetEphemeral(true).
						Build(),
				)
				return
			}

			trade.State = transaction.TransactionProgress

			World.InitTrade(tradeUuid)

			event.UpdateMessage(
				discord.
					NewMessageUpdateBuilder().
					ClearContainerComponents().
					ClearEmbeds().
					SetContent("Odrzucono transakcję").
					Build(),
			)
		} else {
			trade := World.Transactions[tradeUuid]

			if trade == nil {
				event.CreateMessage(
					discord.
						NewMessageCreateBuilder().
						SetContent("Transakcja nie istnieje").
						SetEphemeral(true).
						Build(),
				)
				return
			}

			if trade.State != transaction.TransactionPending {
				event.CreateMessage(
					discord.
						NewMessageCreateBuilder().
						SetContent("Transakcja już została zaakceptowana").
						SetEphemeral(true).
						Build(),
				)
				return
			}

			World.RejectTrade(tradeUuid)

			event.UpdateMessage(
				discord.
					NewMessageUpdateBuilder().
					ClearContainerComponents().
					ClearEmbeds().
					SetContent("Odrzucono transakcję").
					Build(),
			)
		}
	}
}
