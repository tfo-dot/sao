package discord

import (
	"fmt"
	"sao/battle"
	"sao/player/inventory"
	"sao/utils"
	"sao/world/party"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/google/uuid"
)

func ComponentHandler(event *events.ComponentInteractionCreate) {
	customId := event.ComponentInteraction.Data.CustomID()

	if strings.HasPrefix(customId, "su") {
		data := strings.Split(customId, "|")

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

	if strings.HasPrefix(customId, "party") {
		if customId == "party/res" {
			userSnowflake := event.Member().User.ID.String()

			partyUuid := uuid.MustParse(customId[10:])

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
			event.UpdateMessage(messageUpdateClearComponents)

			event.CreateMessage(
				discord.
					NewMessageCreateBuilder().
					SetContent("Odrzucono zaproszenie").
					Build(),
			)
		}
	}

	if strings.HasPrefix(customId, "f") {
		action := customId[2:]
		meta := ""

		if strings.Contains(action, "+") {
			split := strings.Split(action, "+")

			action = split[0]
			meta = split[1]
		}

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

			event.UpdateMessage(messageUpdateClearComponents)

			fight.ActionChannel <- battle.Action{
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
		case "defend":
			event.UpdateMessage(messageUpdateClearComponents)

			fight.ActionChannel <- battle.Action{
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
				if player.CanUseLvlSkill(skill) {
					options = append(options, discord.NewStringSelectMenuOption(skill.Name, fmt.Sprint(skill.ForLevel)))
				}
			}

			event.UpdateMessage(
				discord.
					NewMessageUpdateBuilder().
					ClearContainerComponents().
					AddActionRow(
						discord.
							NewStringSelectMenu(
								"f/skill/usage",
								"Wybierz umiejętność",
							).
							WithMaxValues(1).
							AddOptions(options...),
					).
					Build(),
			)
		case "skill/usage":
			rawSkillUuid := event.ComponentInteraction.StringSelectMenuInteractionData().Values[0]

			for _, skill := range player.Inventory.LevelSkills {
				if fmt.Sprint(skill.ForLevel) == rawSkillUuid {
					event.UpdateMessage(
						discord.
							NewMessageUpdateBuilder().
							ClearContainerComponents().
							Build(),
					)

					if skill.Trigger.Event.TargetCount == -1 {
						fight.ActionChannel <- battle.Action{
							Event:  battle.ACTION_SKILL,
							Source: player.GetUUID(),
							Target: uuid.Nil,
							Meta: battle.ActionSkillMeta{
								IsForLevel: true,
								Lvl:        skill.ForLevel,
								SkillUuid:  uuid.Nil,
								Targets: utils.Map(
									fight.GetEnemiesFor(player.GetUUID()),
									func(entity battle.Entity) uuid.UUID { return entity.GetUUID() },
								),
							},
						}

						event.CreateMessage(
							discord.
								NewMessageCreateBuilder().
								SetContent("Użyto umiejętności!").
								SetEphemeral(true).
								Build(),
						)
					} else {
						options := utils.Map(
							fight.GetEnemiesFor(player.GetUUID()),
							func(entity battle.Entity) discord.StringSelectMenuOption {
								return discord.NewStringSelectMenuOption(entity.GetName(), entity.GetUUID().String())
							},
						)

						event.UpdateMessage(
							discord.
								NewMessageUpdateBuilder().
								ClearContainerComponents().
								AddActionRow(
									discord.
										NewStringSelectMenu(
											"f/skill/usage/target+"+rawSkillUuid,
											"Wybierz cel umiejętności",
											options...,
										).
										WithMaxValues(skill.Trigger.Event.TargetCount),
								).
								Build(),
						)
					}

					return
				}
			}
		case "skill/usage/target":
			rawSkillUuid := meta

			for _, skill := range player.Inventory.LevelSkills {
				if fmt.Sprint(skill.ForLevel) == rawSkillUuid {
					event.UpdateMessage(
						discord.
							NewMessageUpdateBuilder().
							ClearContainerComponents().
							Build(),
					)

					fight.ActionChannel <- battle.Action{
						Event:  battle.ACTION_SKILL,
						Source: player.GetUUID(),
						Target: uuid.Nil,
						Meta: battle.ActionSkillMeta{
							IsForLevel: true,
							Lvl:        skill.ForLevel,
							SkillUuid:  uuid.Nil,
							Targets: utils.Map(
								event.ComponentInteraction.StringSelectMenuInteractionData().Values,
								func(s string) uuid.UUID { return uuid.MustParse(s) },
							),
						},
					}

					event.CreateMessage(
						discord.
							NewMessageCreateBuilder().
							SetContent("Użyto umiejętności!").
							SetEphemeral(true).
							Build(),
					)

					return
				}
			}
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

						fight.ActionChannel <- battle.Action{
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

			fight.ActionChannel <- battle.Action{
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
			store := World.Stores[uuid.MustParse(segments[2])]

			message := discord.NewMessageCreateBuilder()

			embed := discord.NewEmbedBuilder()

			embed.SetTitle("Sklep: " + store.Name)

			embed.SetDescription(store.LastRestock.String())

			message.AddEmbeds(embed.Build())

			event.CreateMessage(message.Build())
		}
	}
}
