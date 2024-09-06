package world

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sao/battle"
	"sao/battle/mobs"
	"sao/config"
	"sao/data"
	"sao/player"
	"sao/types"
	"sao/utils"
	"sao/world/calendar"
	"sao/world/location"
	"sao/world/party"
	"sao/world/tournament"
	"sao/world/transaction"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/google/uuid"
)

type FloorMap map[string]location.Floor

type World struct {
	Players        map[uuid.UUID]*player.Player
	Transactions   map[uuid.UUID]*transaction.Transaction
	Stores         map[uuid.UUID]*types.NPCStore
	Floors         FloorMap
	Tournaments    map[uuid.UUID]*tournament.Tournament
	Fights         map[uuid.UUID]*battle.Fight
	Entities       map[uuid.UUID]*types.Entity
	Time           *calendar.Calendar
	Parties        map[uuid.UUID]*party.Party
	DiscordChannel chan types.DiscordEvent
	BufferChannel  chan types.DiscordMessageStruct
}

func (w *World) MessageHandler() {
	for {
		msg, ok := <-w.BufferChannel

		if !ok {
			panic("Buffer channel closed")
		}

		w.DiscordChannel <- types.DiscordSendMsg{Data: msg}
	}
}

func CreateWorld() World {
	return World{
		make(map[uuid.UUID]*player.Player),
		make(map[uuid.UUID]*transaction.Transaction),
		data.Shops,
		location.GetFloors(),
		make(map[uuid.UUID]*tournament.Tournament),
		make(map[uuid.UUID]*battle.Fight),
		make(map[uuid.UUID]*types.Entity),
		calendar.StartCalendar(),
		make(map[uuid.UUID]*party.Party),
		make(chan types.DiscordEvent, 10),
		make(chan types.DiscordMessageStruct, 10),
	}
}

func (w *World) MovePlayer(pUuid uuid.UUID, floorName, locationName, reason string) error {
	player := w.Players[pUuid]

	floorData, floorExists := w.Floors[floorName]

	if !floorExists || (!floorData.Unlocked && player.Meta.Location.Floor != floorName) {
		return errors.New("floor not found or locked")
	}

	if locationData := w.Floors[floorName].FindLocation(locationName); locationData == nil || !locationData.Unlocked {
		return errors.New("location not found or locked")
	}

	if player.Meta.FightInstance != nil {
		return errors.New("player is in fight")
	}

	hasPerPlayerUnlock := false
	for _, flag := range floorData.Flags {
		if flag == "per-player" {
			hasPerPlayerUnlock = true
			break
		}
	}

	if hasPerPlayerUnlock {
		playerUnlocked := false

		for _, floor := range player.Meta.UnlockedFloors {
			if floor == floorName {
				playerUnlocked = true
				break
			}
		}

		if !playerUnlocked {
			return errors.New("player has not unlocked this floor")
		}

	}

	player.Meta.Location.Floor = floorName
	player.Meta.Location.Location = locationName

	return nil
}

func (w *World) PlayerFight(pUuid uuid.UUID, threadId string, mentionAll bool, mobId string, mobCount int) {
	playerObj := w.Players[pUuid]

	floor := w.Floors[playerObj.Meta.Location.Floor]

	var location location.Location

	entityMap := make(battle.EntityMap)

	for _, loc := range floor.Locations {
		if loc.Name == playerObj.Meta.Location.Location {
			location = loc
		}
	}

	if playerObj.Meta.Party != nil {
		for _, member := range w.Parties[playerObj.Meta.Party.UUID].Players {
			entityMap[w.Players[member.PlayerUuid].GetUUID()] = battle.EntityEntry{
				Entity: w.Players[member.PlayerUuid],
				Side:   0,
			}
		}
	} else {
		entityMap[playerObj.GetUUID()] = battle.EntityEntry{Entity: playerObj, Side: 0}
	}

	for i := 0; i < mobCount; i++ {
		entity := mobs.Spawn(mobId)

		entityMap[entity.GetUUID()] = battle.EntityEntry{Entity: entity, Side: 1}
	}

	fight := battle.Fight{
		Entities:       entityMap,
		DiscordChannel: w.BufferChannel,
		Location:       &location,
		Meta: &battle.FightMeta{
			Tournament: nil,
			ThreadId:   threadId,
		},
	}

	fight.Init()

	fightUUID := w.RegisterFight(&fight)

	if mentionAll {
		mentionString := ""
		for _, entity := range fight.Entities {
			if entity.Entity.GetFlags()&types.ENTITY_AUTO == 0 {
				mentionString += fmt.Sprintf("<@%v>, ", entity.Entity.(*player.Player).Meta.UserID)
			}
		}

		mentionString = mentionString[:len(mentionString)-2]

		w.BufferChannel <- types.DiscordMessageStruct{
			ChannelID: threadId,
			MessageContent: discord.
				NewMessageCreateBuilder().
				SetContent(mentionString).
				Build(),
		}
	}

	playerObj.Meta.FightInstance = &fightUUID

	go w.ListenForFight(fightUUID)
}

func (w *World) PlayerSearch(pUuid uuid.UUID, threadId string, event *events.ApplicationCommandInteractionCreate) {
	player := w.Players[pUuid]
	floor := w.Floors[player.Meta.Location.Floor]

	if player.Meta.FightInstance != nil {
		return
	}

	var location location.Location

	for _, loc := range floor.Locations {
		if loc.Name == player.Meta.Location.Location {
			location = loc
		}
	}

	if len(location.Enemies) == 0 {
		return
	}

	isPrivate := false
	mentionAll := false
	canChoose := true

	for _, flag := range location.Flags {
		switch flag {
		case "private":
			isPrivate = true
		case "fixed":
			canChoose = false
		}
	}

	if isPrivate {
		if player.Meta.Party != nil {
			event.CreateMessage(
				discord.
					NewMessageCreateBuilder().
					SetContent("Nie możesz walczyć z przeciwnikami w prywatnych lokacjach będąc w party!").
					Build(),
			)
		}
	}

	if isPrivate {
		mentionAll = true
		client, err := disgo.New(config.Config.Token)

		if err != nil {
			panic(err)
		}

		thread, err := client.Rest().CreateThread(snowflake.MustParse(location.CID), discord.GuildPrivateThreadCreate{
			Name: "Walka " + player.GetName(),
		})

		if err != nil {
			panic(err)
		}

		threadId = thread.ID().String()
	}

	if canChoose {
		if len(location.Enemies) == 1 {
			if location.Enemies[0].MinNum == location.Enemies[0].MaxNum {
				go w.PlayerFight(pUuid, threadId, mentionAll, location.Enemies[0].Enemy, location.Enemies[0].MinNum)

				return
			}

			selectMenuUuid := uuid.New().String()

			options := make([]discord.StringSelectMenuOption, 0)

			for i := location.Enemies[0].MinNum; i <= location.Enemies[0].MaxNum; i++ {
				options = append(options, discord.NewStringSelectMenuOption(strconv.Itoa(i), strconv.Itoa(i)))
			}

			event.CreateMessage(
				discord.
					NewMessageCreateBuilder().
					AddActionRow(
						discord.
							NewStringSelectMenu("chc/"+selectMenuUuid, "Wybierz ilość przeciwników").
							WithMaxValues(1).
							AddOptions(options...),
					).
					Build(),
			)

			w.DiscordChannel <- types.DiscordChoiceMsg{
				Data: types.DiscordChoice{
					Id: selectMenuUuid,
					Select: func(cic *events.ComponentInteractionCreate) {
						choiceRaw := cic.StringSelectMenuInteractionData().Values[0]

						choice, err := strconv.Atoi(choiceRaw)

						if err != nil {
							return
						}

						cic.UpdateMessage(discord.
							NewMessageUpdateBuilder().
							ClearContainerComponents().
							SetContent("Wybrano " + choiceRaw + " przeciwników!").
							Build())

						go w.PlayerFight(pUuid, threadId, mentionAll, location.Enemies[0].Enemy, choice)
					},
				},
			}
		} else {
			selectMenuUuid := uuid.New().String()

			options := make([]discord.StringSelectMenuOption, 0)

			for idx, enemy := range location.Enemies {
				options = append(options, discord.NewStringSelectMenuOption(enemy.Enemy, strconv.Itoa(idx)))
			}

			event.CreateMessage(
				discord.
					NewMessageCreateBuilder().
					AddActionRow(
						discord.
							NewStringSelectMenu("chc/"+selectMenuUuid, "Wybierz typ przeciwnika").
							WithMaxValues(1).
							AddOptions(options...),
					).
					Build(),
			)

			w.DiscordChannel <- types.DiscordChoiceMsg{
				Data: types.DiscordChoice{
					Id: selectMenuUuid,
					Select: func(cic *events.ComponentInteractionCreate) {
						choiceRaw := cic.StringSelectMenuInteractionData().Values[0]

						choice, err := strconv.Atoi(choiceRaw)

						if err != nil {
							return
						}

						enemy := location.Enemies[choice]

						if enemy.MinNum == enemy.MaxNum {
							count := enemy.MinNum
							cic.UpdateMessage(discord.NewMessageUpdateBuilder().ClearContainerComponents().SetContent("Wybrano przeciwnika - " + mobs.Mobs[location.Enemies[choice].Enemy].Name).Build())

							go w.PlayerFight(pUuid, threadId, mentionAll, enemy.Enemy, count)
						} else {
							selectMenuUuid := uuid.New().String()

							options := make([]discord.StringSelectMenuOption, 0)

							for i := enemy.MinNum; i <= enemy.MaxNum; i++ {
								options = append(options, discord.NewStringSelectMenuOption(strconv.Itoa(i), strconv.Itoa(i)))
							}

							cic.UpdateMessage(
								discord.NewMessageUpdateBuilder().
									SetContent("Wybierz ilość przeciwników!").
									AddActionRow(
										discord.
											NewStringSelectMenu("chc/"+selectMenuUuid, "Wybierz ilość przeciwników").
											WithMaxValues(1).
											AddOptions(options...),
									).
									Build(),
							)

							w.DiscordChannel <- types.DiscordChoiceMsg{
								Data: types.DiscordChoice{
									Id: selectMenuUuid,
									Select: func(cic *events.ComponentInteractionCreate) {
										choiceRaw := cic.StringSelectMenuInteractionData().Values[0]

										choice, err := strconv.Atoi(choiceRaw)

										if err != nil {
											return
										}

										cic.UpdateMessage(discord.NewMessageUpdateBuilder().ClearContainerComponents().SetContent("Wybrano " + choiceRaw + " przeciwników!").Build())

										go w.PlayerFight(pUuid, threadId, mentionAll, enemy.Enemy, choice)
									},
								},
							}
						}
					},
				},
			}
		}

		return
	} else {
		enemy := utils.RandomElement(location.Enemies)
		enemyCount := utils.RandomNumber(enemy.MinNum, enemy.MaxNum)

		go w.PlayerFight(pUuid, threadId, mentionAll, enemy.Enemy, enemyCount)
	}
}

func (w *World) StartBackupClock() {
	for range time.Tick(15 * time.Minute) {
		w.CreateBackup()
	}
}

func (w *World) StartClock() {
	go w.StartBackupClock()

	for range time.Tick(1 * time.Minute) {
		w.Time.Tick()

		for pUuid, player := range w.Players {
			//Not in fight
			if player.Meta.FightInstance != nil {
				continue
			}

			//Dead
			//TODO only resurrect if not in hardcore mode
			if player.GetCurrentHP() <= 0 {
				location := w.Floors[player.Meta.Location.Floor].FindLocation(player.Meta.Location.Location)

				player.Stats.HP = player.GetStat(types.STAT_HP)

				w.BufferChannel <- types.DiscordMessageStruct{
					ChannelID: location.CID,
					MessageContent: discord.
						NewMessageCreateBuilder().
						AddEmbeds(
							discord.NewEmbedBuilder().SetTitle("Wskrzeszenie!").SetDescriptionf("%s zostaje wskrzeszony...", player.GetName()).Build(),
						).
						Build(),
				}

				continue
			}

			//Missing mana
			if player.GetCurrentMana() < player.GetStat(types.STAT_MANA) {
				player.Stats.CurrentMana += 1
			}

			//Can be healed
			if player.GetCurrentHP() < player.GetStat(types.STAT_HP) {
				healRatio := 50

				{
					player := w.Players[pUuid]

					floor := w.Floors[player.Meta.Location.Floor]
					location := floor.FindLocation(player.Meta.Location.Location)

					if location.CityPart {
						healRatio = 25
					}
				}

				player.Heal(player.GetStat(types.STAT_HP) / healRatio)

				if player.Meta.WaitToHeal && player.GetCurrentHP() == player.GetStat(types.STAT_HP) {
					player.Meta.WaitToHeal = false

					client, err := disgo.New(config.Config.Token)

					if err != nil {
						w.BufferChannel <- types.DiscordMessageStruct{
							ChannelID: config.Config.LogChannelID,
							MessageContent: discord.NewMessageCreateBuilder().
								SetContent("Nie można wysłać wiadomości do gracza (HP Heal)").
								Build(),
						}
					}

					ch, err := client.Rest().CreateDMChannel(snowflake.MustParse(player.Meta.UserID))

					if err != nil {
						w.BufferChannel <- types.DiscordMessageStruct{
							ChannelID: config.Config.LogChannelID,
							MessageContent: discord.NewMessageCreateBuilder().
								SetContent("Nie można wysłać wiadomości do gracza (HP Heal)").
								Build(),
						}
					}

					_, err = client.Rest().CreateMessage(ch.ID(), discord.NewMessageCreateBuilder().
						SetContent("Twoja postać ma już 100% HP, baw się dobrze!").
						Build(),
					)

					if err != nil {
						w.BufferChannel <- types.DiscordMessageStruct{
							ChannelID: config.Config.LogChannelID,
							MessageContent: discord.NewMessageCreateBuilder().
								SetContent("Nie można wysłać wiadomości do gracza (HP Heal)").
								Build(),
						}
					}
				}

				if player.Stats.HP >= player.GetStat(types.STAT_HP) && player.Meta.WaitToHeal {
					player.Meta.WaitToHeal = false
				}
			}
		}
	}
}

// Fight stuff
func (w *World) RegisterFight(fight *battle.Fight) uuid.UUID {
	uuid := uuid.New()

	w.Fights[uuid] = fight

	for _, entity := range fight.Entities {
		if entity.Entity.GetFlags()&types.ENTITY_AUTO == 0 {
			w.Entities[entity.Entity.GetUUID()] = &entity.Entity
		}
	}

	return uuid
}

func (w *World) ListenForFight(fightUuid uuid.UUID) {
	fight, ok := w.Fights[fightUuid]

	if !ok {
		return
	}

	go fight.Run()

	channelId := fight.GetChannelId()

	for {
		eventData, ok := <-fight.ExternalChannel

		if !ok {
			w.DeregisterFight(fightUuid)
			break
		}

		switch eventData.GetEvent() {
		case battle.MSG_FIGHT_END:

			if eventData.GetData().(bool) {
				w.BufferChannel <- types.DiscordMessageStruct{
					ChannelID: channelId,
					MessageContent: discord.
						NewMessageCreateBuilder().
						AddEmbeds(
							discord.
								NewEmbedBuilder().
								SetTitle("Koniec walki!").
								SetDescriptionf("Wszyscy gracze uciekli z walki!").
								Build(),
						).
						Build(),
				}

				w.DeregisterFight(fightUuid)

				return
			}

			wonSideIDX := fight.SidesLeft()[0]
			wonEntities := fight.FromSide(wonSideIDX)

			allAuto := true

			for _, entity := range wonEntities {
				if entity.GetFlags()&types.ENTITY_AUTO != types.ENTITY_AUTO {
					allAuto = false
					break
				}
			}

			if allAuto {
				w.DeregisterFight(fightUuid)

				wonSideText := ""

				for _, entity := range wonEntities {
					wonSideText += fmt.Sprintf("%v\n", entity.GetName())
				}

				wonSideText = wonSideText[:len(wonSideText)-1]

				w.BufferChannel <- types.DiscordMessageStruct{
					ChannelID: channelId,
					MessageContent: discord.
						NewMessageCreateBuilder().
						AddEmbeds(
							discord.
								NewEmbedBuilder().
								SetTitle("Koniec walki!").
								SetDescriptionf("Wygrali:\n" + wonSideText).
								Build(),
						).
						Build(),
				}

				break
			}

			enemies := make([]types.Entity, 0)

			xpMap := make(map[uuid.UUID]int)
			goldMap := make(map[uuid.UUID]int)

			for _, entity := range fight.Entities {
				if entity.Side == wonSideIDX {
					continue
				}

				enemies = append(enemies, entity.Entity)
			}

			for _, entity := range enemies {
				if entity.GetFlags()&types.ENTITY_AUTO == 0 {
					continue
				}

				if entity.HasOnDefeat() {
					for _, wonEntity := range wonEntities {
						if wonEntity.GetFlags()&types.ENTITY_AUTO != 0 {
							continue
						}

						entity.(types.DefeatableEntity).OnDefeat(wonEntity.(types.PlayerEntity))
					}
				}
			}

			if fight.Meta.Tournament == nil {
				overallXp := 0
				overallGold := 0
				lootedItems := make([]types.Loot, 0)

				for _, entity := range fight.Entities {
					if entity.Side == wonSideIDX {
						continue
					}

					lootList := entity.Entity.GetLoot()

					for _, loot := range lootList {
						switch loot.Type {
						case types.LOOT_EXP:
							overallXp += loot.Count
						case types.LOOT_GOLD:
							overallGold += loot.Count
						case types.LOOT_ITEM:
							lootedItems = append(lootedItems, loot)
						}
					}
				}

				partyInfo := wonEntities[0].(*player.Player).Meta.Party

				unlockedFloors := w.GetUnlockedFloorCount()

				for _, loot := range fight.AdditionalLoot {
					for _, entity := range wonEntities {
						if entity.GetUUID() == loot.Target {
							player := w.Players[entity.GetUUID()]

							switch loot.Value.Type {
							case types.LOOT_EXP:
								player.AddEXP(unlockedFloors, loot.Value.Count)
								if _, ok := xpMap[entity.GetUUID()]; !ok {
									xpMap[entity.GetUUID()] = loot.Value.Count
								} else {
									xpMap[entity.GetUUID()] += loot.Value.Count
								}
							case types.LOOT_GOLD:
								player.AddEXP(unlockedFloors, loot.Value.Count)
								if _, ok := xpMap[entity.GetUUID()]; !ok {
									goldMap[entity.GetUUID()] = loot.Value.Count
								} else {
									goldMap[entity.GetUUID()] += loot.Value.Count
								}
							case types.LOOT_ITEM:
								itemUuid := loot.Value.Meta.Uuid

								if loot.Value.Meta.Type == types.ITEM_OTHER {
									itemObj := data.Items[itemUuid]

									itemObj.Count = loot.Value.Count

									player.Inventory.Items = append(player.Inventory.Items, &itemObj)
								} else {
									ingredient := data.Ingredients[itemUuid]

									ingredient.Count = loot.Value.Count

									player.Inventory.AddIngredient(&ingredient)
								}
							}
						}
					}
				}

				if partyInfo != nil {
					partyData := w.Parties[partyInfo.UUID]
					partyLeader := w.Players[partyData.Leader]

					for _, member := range partyData.Players {
						player := w.Players[member.PlayerUuid]

						player.AddEXP(unlockedFloors, overallXp/len(partyData.Players))

						if _, ok := xpMap[member.PlayerUuid]; !ok {
							xpMap[member.PlayerUuid] = overallXp / len(partyData.Players)
						} else {
							xpMap[member.PlayerUuid] += overallXp / len(partyData.Players)
						}

						player.AddGold(overallGold / len(partyData.Players))

						if _, ok := goldMap[member.PlayerUuid]; !ok {
							goldMap[member.PlayerUuid] = overallGold / len(partyData.Players)
						} else {
							goldMap[member.PlayerUuid] += overallGold / len(partyData.Players)
						}
					}

					for _, loot := range lootedItems {
						itemUuid := loot.Meta.Uuid

						if loot.Meta.Type == types.ITEM_OTHER {
							itemObj := data.Items[itemUuid]

							itemObj.Count = loot.Count

							partyLeader.Inventory.Items = append(partyLeader.Inventory.Items, &itemObj)
						} else {
							ingredient := data.Ingredients[itemUuid]

							ingredient.Count = loot.Count

							partyLeader.Inventory.AddIngredient(&ingredient)
						}
					}
				} else {
					for _, entity := range wonEntities {
						entityUuid := entity.GetUUID()

						if entity.GetFlags()&types.ENTITY_AUTO != 0 {
							continue
						}

						player := w.Players[entityUuid]

						player.AddEXP(unlockedFloors, overallXp)

						if _, ok := xpMap[entityUuid]; !ok {
							xpMap[entityUuid] = overallXp
						} else {
							xpMap[entityUuid] += overallXp
						}

						player.AddGold(overallGold)

						if _, ok := goldMap[entityUuid]; !ok {
							goldMap[entityUuid] = overallGold
						} else {
							goldMap[entityUuid] += overallGold
						}

						for _, loot := range lootedItems {
							itemUuid := loot.Meta.Uuid

							if loot.Meta.Type == types.ITEM_OTHER {
								itemObj := data.Items[itemUuid]

								itemObj.Count = loot.Count

								player.Inventory.Items = append(player.Inventory.Items, &itemObj)
							} else {
								ingredient := data.Ingredients[itemUuid]

								ingredient.Count = loot.Count

								player.Inventory.AddIngredient(&ingredient)
							}
						}
					}
				}
			}

			wonSideText := ""

			for _, entity := range wonEntities {
				wonSideText += fmt.Sprintf("%v", entity.GetName())

				if entity.GetFlags()&types.ENTITY_AUTO == 0 {
					wonSideText += fmt.Sprintf(" (<@%v>)", entity.(*player.Player).Meta.UserID)
				}

				wonSideText += fmt.Sprintf(" (%v/%v HP)", entity.GetCurrentHP(), entity.GetStat(types.STAT_HP))

				wonSideText += "\n"
			}

			wonSideText = wonSideText[:len(wonSideText)-1]

			lootSummaryText := ""

			for _, entity := range wonEntities {
				if entity.GetFlags()&types.ENTITY_AUTO != 0 {
					continue
				}

				xpGotten, exists := xpMap[entity.GetUUID()]

				if !exists {
					xpGotten = 0
				}

				goldGotten, exists := goldMap[entity.GetUUID()]

				if !exists {
					goldGotten = 0
				}

				lootSummaryText += fmt.Sprintf("%v - XP: %d, Złoto: %d\n", entity.GetName(), xpGotten, goldGotten)
			}

			lootSummaryText = lootSummaryText[:len(lootSummaryText)-1]

			w.BufferChannel <- types.DiscordMessageStruct{
				ChannelID: channelId,
				MessageContent: discord.
					NewMessageCreateBuilder().
					AddEmbeds(
						discord.
							NewEmbedBuilder().
							SetTitle("Koniec walki!").
							SetDescriptionf("Wygrali:\n"+wonSideText).
							Build(),
						discord.NewEmbedBuilder().SetTitle("Podsumowanie").SetDescription(lootSummaryText).Build(),
					).
					Build(),
			}

			if fight.Meta.Tournament != nil {
				w.Tournaments[fight.Meta.Tournament.Tournament].ExternalChannel <- tournament.MatchFinishedData{Winner: wonEntities[0].GetUUID()}
			}
		case battle.MSG_FIGHT_START:
			oneSide := fight.FromSide(0)
			otherSide := fight.FromSide(1)

			oneSideText := ""

			for _, entity := range oneSide {
				oneSideText += fmt.Sprintf("%v", entity.GetName())
				if entity.GetFlags()&types.ENTITY_AUTO == 0 {
					oneSideText += fmt.Sprintf(" (<@%v>)", entity.(*player.Player).Meta.UserID)
				}

				oneSideText += "\n"
			}

			oneSideText = oneSideText[:len(oneSideText)-1]

			otherSideText := ""

			for _, entity := range otherSide {
				otherSideText += fmt.Sprintf("%v", entity.GetName())

				if entity.GetFlags()&types.ENTITY_AUTO == 0 {
					otherSideText += fmt.Sprintf(" (<@%v>)", entity.(*player.Player).Meta.UserID)
				}

				otherSideText += "\n"
			}

			otherSideText = otherSideText[:len(otherSideText)-1]

			w.BufferChannel <- types.DiscordMessageStruct{
				ChannelID: channelId,
				MessageContent: discord.NewMessageCreateBuilder().
					SetContent("Walka się rozpoczyna!").
					AddEmbeds(discord.NewEmbedBuilder().
						SetTitle("Walka").
						AddField("Po jednej!", oneSideText, false).
						AddField("Po drugiej!", otherSideText, false).
						Build()).
					Build(),
			}
		case battle.MSG_ACTION_NEEDED:
			entityUuid := eventData.GetData().(uuid.UUID)

			player := w.Players[entityUuid]

			canUseActionWhileCC := false

			for _, item := range player.Inventory.Items {
				for _, effect := range item.Effects {
					effectTrigger := effect.GetTrigger()

					if effectTrigger.Flags&types.FLAG_IGNORE_CC != 0 {
						canUseActionWhileCC = true
						break
					}
				}
			}

			if !canUseActionWhileCC {
				for _, skill := range player.Inventory.LevelSkills {
					effectTrigger := skill.GetUpgradableTrigger(player.Inventory.LevelSkillsUpgrades[skill.GetLevel()])

					if effectTrigger.Flags&types.FLAG_IGNORE_CC != 0 {
						canUseActionWhileCC = true
						break
					}
				}
			}

			if !canUseActionWhileCC {
				if player.Meta.Fury != nil {
					for _, skill := range player.Meta.Fury.GetSkills() {
						effectTrigger := skill.GetTrigger()

						if effectTrigger.Flags&types.FLAG_IGNORE_CC != 0 {
							canUseActionWhileCC = true
							break
						}
					}
				}
			}

			if player.GetEffectByType(types.EFFECT_TAUNTED) != nil {
				if !canUseActionWhileCC {
					effect := player.GetEffectByType(types.EFFECT_TAUNTED)

					fight.PlayerActions <- types.Action{
						Event:  types.ACTION_ATTACK,
						Source: player.GetUUID(),
						Target: effect.Meta.(uuid.UUID),
					}

					w.BufferChannel <- types.DiscordMessageStruct{
						ChannelID: channelId,
						MessageContent: discord.NewMessageCreateBuilder().
							SetContentf("<@%v> jest zmuszony do ataku! Pomijamy turę!", player.Meta.UserID).
							Build(),
					}

					continue
				}
			}

			attackButton := discord.NewPrimaryButton("Atak", "f/attack")

			if !player.CanAttack() {
				attackButton = attackButton.AsDisabled()
			}

			defendButton := discord.NewPrimaryButton("Obrona", "f/defend")

			if !player.CanDefend() {
				defendButton = defendButton.AsDisabled()
			}

			skillButton := discord.NewPrimaryButton("Skill", "f/skill")

			filteredSkillsCount := 0

			if player.Meta.Fury != nil {
				for _, skill := range player.Meta.Fury.GetSkills() {
					player.CanUseSkill(skill)
				}
			}

			for _, skill := range player.Inventory.LevelSkills {
				if player.CanUseSkill(skill) && skill.CanUse(player, fight) {
					filteredSkillsCount++
				}
			}

			if filteredSkillsCount == 0 {
				skillButton = skillButton.AsDisabled()
			}

			itemButton := discord.NewPrimaryButton("Przedmiot", "f/item")

			filteredItemsCount := 0

			for _, item := range player.Inventory.Items {
				if item.Consume && item.Count > 0 && !item.Hidden {
					filteredItemsCount++
				}
			}

			if filteredItemsCount == 0 {
				itemButton = itemButton.AsDisabled()
			}

			escapeButton := discord.NewDangerButton("Ucieczka", "f/escape")

			w.BufferChannel <- types.DiscordMessageStruct{
				ChannelID: channelId,
				MessageContent: discord.NewMessageCreateBuilder().
					AddEmbeds(discord.NewEmbedBuilder().
						SetTitle("Czas na turę!").
						SetDescriptionf("Kolej <@%s>!", player.Meta.UserID).
						SetAuthorName(player.Name).
						Build(),
					).
					AddActionRow(
						attackButton,
						defendButton,
						skillButton,
						itemButton,
						escapeButton,
					).
					Build(),
			}
		case battle.MSG_SUMMON_EXPIRED:
			entityUuid := eventData.GetData().(uuid.UUID)

			temp := fight.Entities[entityUuid].Entity

			delete(fight.Entities, entityUuid)

			w.BufferChannel <- types.DiscordMessageStruct{
				ChannelID: channelId,
				MessageContent: discord.
					NewMessageCreateBuilder().
					AddEmbeds(
						discord.
							NewEmbedBuilder().
							SetTitle("Przyzwany stwór uciekł!").
							SetDescriptionf("%s uciekł z pola walki!", temp.GetName()).
							Build(),
					).
					Build(),
			}
		default:
			panic("Unhandled event")
		}

		//Fallback for when the channel is closed
		if fight.IsFinished() {
			w.DeregisterFight(fightUuid)
			break
		}
	}
}

func (w *World) DeregisterFight(uuid uuid.UUID) {
	tmp := w.Fights[uuid]

	for _, entity := range tmp.Entities {
		if entity.Entity.GetFlags()&types.ENTITY_AUTO == 0 {
			entity.Entity.(*player.Player).Meta.FightInstance = nil
		} else {
			delete(w.Entities, entity.Entity.GetUUID())
		}
	}

	delete(w.Fights, uuid)
}

func (w *World) GetPlayer(uid string) *player.Player {
	for _, pl := range w.Players {
		if pl.Meta.UserID == uid {
			return pl
		}
	}

	return nil
}

func (w *World) RegisterTournament(tournament tournament.Tournament) {
	w.Tournaments[tournament.Uuid] = &tournament

	var playerText string = ""
	if tournament.MaxPlayers == -1 {
		playerText = "Nieograniczona"
	} else {
		playerText = fmt.Sprintf("%v/%v", len(tournament.Participants), tournament.MaxPlayers)
	}

	w.BufferChannel <- types.DiscordMessageStruct{
		ChannelID: "1225150345009827841",
		MessageContent: discord.NewMessageCreateBuilder().
			AddEmbeds(discord.NewEmbedBuilder().
				SetTitle("Nowy turniej!").
				SetDescriptionf("Zapisy na turniej `%v` otwarte", tournament.Name).
				SetFooterText("Ilość miejsc: " + playerText).
				Build()).
			AddActionRow(
				discord.NewPrimaryButton("Dołącz", "t/join/"+tournament.Uuid.String()),
			).
			Build(),
	}
}

func (w *World) JoinTournament(uuid uuid.UUID, player *player.Player) error {
	tournamentObj := w.Tournaments[uuid]

	if tournamentObj == nil {
		return errors.New("tournament not found")
	}

	if tournamentObj.State != tournament.Waiting {
		return errors.New("tournament running")
	}

	if tournamentObj.MaxPlayers != -1 && len(tournamentObj.Participants) >= tournamentObj.MaxPlayers {
		return errors.New("tournament is full")
	}

	tournamentObj.Participants = append(tournamentObj.Participants, player.GetUUID())

	if tournamentObj.MaxPlayers != -1 && len(tournamentObj.Participants) == tournamentObj.MaxPlayers {
		w.StartTournament(uuid)
	}

	return nil
}

func (w *World) StartTournament(tUuid uuid.UUID) error {
	tournamentObj := w.Tournaments[tUuid]

	if tournamentObj == nil {
		return errors.New("tournament not found")
	}

	if tournamentObj.State != tournament.Waiting {
		return errors.New("tournament running")
	}

	if tournamentObj.MaxPlayers != -1 && len(tournamentObj.Participants) < 2 {
		return errors.New("not enough players")
	}

	tournamentObj.State = tournament.Running

	var fightingLocation location.Location

	for _, floor := range w.Floors {
		for _, location := range floor.Locations {
			for _, effect := range location.Flags {
				if effect == "arena" {
					fightingLocation = location
				}
			}
		}
	}

	client, err := disgo.New(config.Config.Token)

	if err != nil {
		panic(err)
	}

	msg, err := client.Rest().CreateMessage(snowflake.MustParse(fightingLocation.CID), discord.NewMessageCreateBuilder().SetContent("Turniej rozpoczęty!").Build())

	if err != nil {
		panic(err)
	}

	thread, err := client.Rest().CreateThreadFromMessage(snowflake.MustParse(fightingLocation.CID), msg.ID, discord.ThreadCreateFromMessage{
		Name: "Turniej",
	})

	if err != nil {
		panic(err)
	}

	tournamentObj.ExternalChannel = make(chan tournament.TournamentEventData)
	tournamentObj.Channel = thread.ID().String()

	switch tournamentObj.Type {
	case tournament.SingleElimination:

		w.BufferChannel <- types.DiscordMessageStruct{
			ChannelID: tournamentObj.Channel,
			MessageContent: discord.NewMessageCreateBuilder().
				SetContent("Losowanie czas zacząć!").
				Build(),
		}

		matches := make([]*tournament.TournamentMatch, 0)

		var participants = tournamentObj.Participants

		if len(tournamentObj.Participants)%2 != 0 {
			luckyPlayer := utils.RandomNumber(0, len(tournamentObj.Participants)-1)

			matches = append(matches, &tournament.TournamentMatch{
				Players: []uuid.UUID{tournamentObj.Participants[luckyPlayer]},
				Winner:  &tournamentObj.Participants[luckyPlayer],
				State:   tournament.FinishedMatch,
			})

			w.BufferChannel <- types.DiscordMessageStruct{
				ChannelID: tournamentObj.Channel,
				MessageContent: discord.NewMessageCreateBuilder().
					SetContentf("Szczęśliwy gracz to <@%v>\nNie musisz walczyć w tej rundzie i możesz spokojnie oglądać!", w.Players[tournamentObj.Participants[luckyPlayer]].Meta.UserID).
					Build(),
			}

			participants = append(participants[:luckyPlayer], participants[luckyPlayer+1:]...)
		}

		//Hackery shuffle
		for i := range participants {
			j := utils.RandomNumber(0, len(participants)-1)

			participants[i], participants[j] = participants[j], participants[i]
		}

		for i := 0; i < len(participants); i += 2 {
			matches = append(matches, &tournament.TournamentMatch{
				Players: []uuid.UUID{participants[i], participants[i+1]},
				Winner:  nil,
				State:   tournament.BeforeMatch,
			})

			participant0 := w.Players[participants[i]]
			participant1 := w.Players[participants[i+1]]

			w.BufferChannel <- types.DiscordMessageStruct{
				ChannelID: tournamentObj.Channel,
				MessageContent: discord.NewMessageCreateBuilder().
					SetContentf("Los pociągnięty!\nMecz #%v: %v vs %v", (i/2)+1, participant0.GetName(), participant1.GetName()).
					Build(),
			}
		}

		tournamentObj.Stages = append(tournamentObj.Stages, &tournament.TournamentStage{
			Matches: matches,
		})
	}

	for idx, match := range tournamentObj.Stages[0].Matches {
		if match.State == tournament.BeforeMatch {
			w.StartMatch(tUuid, idx)
			break
		}
	}

	go w.ListenForTournament(tUuid)

	return nil
}

func (w *World) ListenForTournament(tUuid uuid.UUID) {
	tournamentObj := w.Tournaments[tUuid]

	if tournamentObj == nil {
		return
	}

	for {
		data, ok := <-tournamentObj.ExternalChannel

		if !ok {
			break
		}

		switch data.GetEvent() {
		case tournament.MatchFinished:
			matchData := data.GetData().(uuid.UUID)

			w.FinishMatch(tUuid, matchData)

			currentStage := tournamentObj.Stages[len(tournamentObj.Stages)-1]

			allFinished := true

			for _, match := range currentStage.Matches {
				if match.State != tournament.FinishedMatch {
					allFinished = false
					break
				}
			}

			if allFinished {
				w.NextStage(tUuid)

				if tournamentObj.State == tournament.Finished {
					currentStage := tournamentObj.Stages[len(tournamentObj.Stages)-1]

					matchWinner := currentStage.Matches[0].Winner

					player := w.Players[*matchWinner]

					w.BufferChannel <- types.DiscordMessageStruct{
						ChannelID: tournamentObj.Channel,
						MessageContent: discord.NewMessageCreateBuilder().
							SetContentf("Turniej zakończony! Wygrał %v (<@%v>)", player.GetName(), player.Meta.UserID).
							Build(),
					}

					w.FinishTournament(tUuid)

					return
				}

				currentStage := tournamentObj.Stages[len(tournamentObj.Stages)-1]

				for idx, match := range currentStage.Matches {
					if match.State == tournament.BeforeMatch {
						w.StartMatch(tUuid, idx)
						break
					}
				}
			} else {
				for idx, match := range currentStage.Matches {
					if match.State == tournament.BeforeMatch {
						w.StartMatch(tUuid, idx)
						break
					}
				}
			}
		}
	}
}

func (w *World) StartMatch(tUuid uuid.UUID, matchIdx int) {
	tournamentObj := w.Tournaments[tUuid]

	if tournamentObj == nil {
		return
	}

	if tournamentObj.State != tournament.Running {
		return
	}

	stage := tournamentObj.Stages[len(tournamentObj.Stages)-1]

	if matchIdx >= len(stage.Matches) {
		return
	}

	match := stage.Matches[matchIdx]

	if match.State != tournament.BeforeMatch {
		return
	}

	match.State = tournament.RunningMatch

	player0 := w.Players[match.Players[0]]
	player1 := w.Players[match.Players[1]]

	entityMap := make(battle.EntityMap)

	entityMap[player0.GetUUID()] = battle.EntityEntry{Entity: player0, Side: 0}
	entityMap[player1.GetUUID()] = battle.EntityEntry{Entity: player1, Side: 1}

	var fightingLocation location.Location

	for _, floor := range w.Floors {
		for _, location := range floor.Locations {
			for _, effect := range location.Flags {
				if effect == "arena" {
					fightingLocation = location
				}
			}
		}
	}

	fight := battle.Fight{
		Entities:       entityMap,
		DiscordChannel: w.BufferChannel,
		Location:       &fightingLocation,
		Meta: &battle.FightMeta{
			ThreadId: "",
			Tournament: &battle.TournamentData{
				Tournament: tUuid,
				Location:   w.Tournaments[tUuid].Channel,
			},
		},
	}

	fight.Init()

	fightUUID := w.RegisterFight(&fight)

	player0.Meta.FightInstance = &fightUUID
	player1.Meta.FightInstance = &fightUUID

	go w.ListenForFight(fightUUID)
}

func (w *World) FinishMatch(tUuid uuid.UUID, winner uuid.UUID) {
	tournamentObj := w.Tournaments[tUuid]

	if tournamentObj == nil {
		return
	}

	if tournamentObj.State != tournament.Running {
		return
	}

	stage := tournamentObj.Stages[len(tournamentObj.Stages)-1]

	for _, match := range stage.Matches {
		for _, player := range match.Players {
			if player == winner {
				match.Winner = &winner

				match.State = tournament.FinishedMatch
			}
		}
	}

	if len(stage.Matches) == 1 {
		tournamentObj.State = tournament.Finished
	}
}

func (w *World) NextStage(tUuid uuid.UUID) {
	tournamentObj := w.Tournaments[tUuid]

	if tournamentObj == nil {
		return
	}

	if tournamentObj.State != tournament.Running {
		return
	}

	stage := tournamentObj.Stages[len(tournamentObj.Stages)-1]

	//check if all matches finished
	for _, match := range stage.Matches {
		if match.State != tournament.FinishedMatch {
			return
		}
	}

	//check if there is only one match
	if len(stage.Matches) == 1 {
		return
	}

	//create new stage
	newMatches := make([]*tournament.TournamentMatch, 0)

	for i := 0; i < len(stage.Matches); i += 2 {
		newMatches = append(newMatches, &tournament.TournamentMatch{
			Players: []uuid.UUID{*stage.Matches[i].Winner, *stage.Matches[i+1].Winner},
			Winner:  nil,
			State:   tournament.BeforeMatch,
		})
	}

	tournamentObj.Stages = append(tournamentObj.Stages, &tournament.TournamentStage{
		Matches: newMatches,
	})
}

func (w *World) FinishTournament(tUuid uuid.UUID) {
	tournamentObj := w.Tournaments[tUuid]

	if tournamentObj == nil {
		return
	}

	if tournamentObj.State != tournament.Finished {
		return
	}

	delete(w.Tournaments, tUuid)
}

func (w *World) CreatePendingTransaction(left, right uuid.UUID) *transaction.Transaction {
	tUuid := uuid.New()

	w.Transactions[tUuid] = &transaction.Transaction{
		Uuid:      tUuid,
		LeftSide:  &transaction.TransactionSide{Who: left},
		RightSide: &transaction.TransactionSide{Who: right},
		State:     transaction.TransactionPending,
	}

	return w.Transactions[tUuid]
}

func (w *World) InitTrade(tUuid uuid.UUID) {
	transactionObj := w.Transactions[tUuid]

	if transactionObj == nil {
		return
	}

	if transactionObj.State != transaction.TransactionPending {
		return
	}

	transactionObj.State = transaction.TransactionProgress

	w.BufferChannel <- types.DiscordMessageStruct{
		ChannelID:      w.Players[transactionObj.LeftSide.Who].Meta.UserID,
		MessageContent: discord.NewMessageCreateBuilder().SetContent("Transakcja rozpoczęta!").Build(),
		DM:             true,
	}

	w.BufferChannel <- types.DiscordMessageStruct{
		ChannelID:      w.Players[transactionObj.RightSide.Who].Meta.UserID,
		MessageContent: discord.NewMessageCreateBuilder().SetContent("Transakcja rozpoczęta!").Build(),
		DM:             true,
	}
}

func (w *World) RejectTrade(tUuid uuid.UUID) {
	transactionObj := w.Transactions[tUuid]

	if transactionObj == nil {
		return
	}

	if transactionObj.State != transaction.TransactionPending {
		return
	}

	delete(w.Transactions, tUuid)
}

func (w *World) Serialize() map[string]interface{} {
	playerData := make(map[uuid.UUID]map[string]interface{})

	for _, player := range w.Players {
		playerData[player.GetUUID()] = player.Serialize()
	}

	storeData := make([]map[string]interface{}, 0)

	for _, store := range w.Stores {

		stocks := make([]map[string]interface{}, 0)

		for _, stock := range store.Stock {
			stocks = append(stocks, map[string]interface{}{
				"type":  stock.ItemType,
				"uuid":  stock.ItemUUID,
				"price": stock.Price,
			})
		}

		storeData = append(storeData, map[string]interface{}{
			"uuid":  store.Uuid,
			"name":  store.Name,
			"stock": stocks,
		})
	}

	tournamentData := make([]map[string]interface{}, 0)

	for _, tournament := range w.Tournaments {
		tournamentData = append(tournamentData, tournament.Serialize())
	}

	partyData := make(map[uuid.UUID]map[string]interface{})

	for key, party := range w.Parties {
		partyData[key] = party.Serialize()
	}

	return map[string]interface{}{
		"players":     playerData,
		"parties":     partyData,
		"stores":      storeData,
		"time":        w.Time.Serialize(),
		"tournaments": tournamentData,
	}
}

func (w *World) DumpBackup() []byte {
	jsonFile, err := json.Marshal(w.Serialize())

	if err != nil {
		panic(err)
	}

	return jsonFile
}

func (w *World) CreateBackup() []byte {
	_, err := os.Stat(config.Config.BackupLocation)

	if os.IsNotExist(err) {
		os.Mkdir(config.Config.BackupLocation, os.ModePerm)
	}

	newBackupPath := fmt.Sprintf("%s/%v.json", config.Config.BackupLocation, time.Now().Format("2006-01-02_15-04-05"))

	backupFile, err := os.Create(newBackupPath)

	if err != nil {
		panic(err)
	}

	defer backupFile.Close()

	rawData := w.DumpBackup()

	backupFile.Write(rawData)

	w.BufferChannel <- types.DiscordMessageStruct{
		ChannelID: config.Config.LogChannelID,
		MessageContent: discord.NewMessageCreateBuilder().
			SetContent("Backup zrobiony!").
			Build(),
	}

	return rawData
}

func (w *World) LoadBackup() {
	_, err := os.Stat(config.Config.BackupLocation)

	if os.IsNotExist(err) {
		return
	}

	backupData := make(map[string]interface{})

	allBackups, err := os.ReadDir(config.Config.BackupLocation)

	if err != nil {
		panic(err)
	}

	filteredBackups := make([]os.DirEntry, 0)

	for _, backup := range allBackups {
		if strings.HasSuffix(backup.Name(), "json") {
			filteredBackups = append(filteredBackups, backup)
		}
	}

	sort.Slice(filteredBackups, func(i, j int) bool {
		leftName := allBackups[i].Name()
		rightName := allBackups[j].Name()

		leftName = leftName[:len(leftName)-5]
		rightName = rightName[:len(rightName)-5]

		leftTime, err := time.Parse("2006-01-02_15-04-05", leftName)

		if err != nil {
			panic(err)
		}

		rightTime, err := time.Parse("2006-01-02_15-04-05", rightName)

		if err != nil {
			panic(err)
		}

		return leftTime.After(rightTime)
	})

	if len(filteredBackups) == 0 {
		fmt.Println("No backups found")

		return
	} else {
		fmt.Println("Loading backup", filteredBackups[0].Name())
	}

	backupContent, err := os.ReadFile("./" + config.Config.BackupLocation + "/" + filteredBackups[0].Name())

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(backupContent, &backupData)

	if err != nil {
		panic(err)
	}

	w.Time = calendar.Deserialize(backupData["time"].(map[string]interface{}))

	w.Players = make(map[uuid.UUID]*player.Player)

	for _, playerData := range backupData["players"].(map[string]interface{}) {
		player := player.Deserialize(playerData.(map[string]interface{}))

		w.Players[player.GetUUID()] = player
	}

	w.Parties = make(map[uuid.UUID]*party.Party)

	for key, partyData := range backupData["parties"].(map[string]interface{}) {
		deserializedParty := party.Deserialize(partyData.(map[string]interface{}))

		for _, member := range deserializedParty.Players {
			w.Players[member.PlayerUuid].Meta.Party = &player.PartialParty{
				UUID:         uuid.MustParse(key),
				Role:         member.Role,
				MembersCount: len(deserializedParty.Players),
			}
		}

		w.Parties[uuid.MustParse(key)] = deserializedParty
	}

	for _, tData := range backupData["tournaments"].([]interface{}) {
		tempData := tData.(map[string]interface{})
		parsedData := tournament.Deserialize(tempData)

		w.Tournaments[parsedData.Uuid] = &parsedData
	}
}

func (w *World) GetUnlockedFloorCount() int {
	unlockedFloors := 0

	for _, floor := range w.Floors {
		if floor.Unlocked && floor.CountsAsUnlocked {
			unlockedFloors++
		}
	}

	return unlockedFloors
}

func (w *World) RegisterParty(party party.Party) {
	partyUuid := uuid.New()

	for _, member := range party.Players {
		playerData := w.Players[member.PlayerUuid]

		playerData.Meta.Party = &player.PartialParty{
			UUID:         partyUuid,
			Role:         member.Role,
			MembersCount: len(party.Players),
		}
	}

	w.Parties[partyUuid] = &party
}
