package world

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sao/battle"
	"sao/battle/mobs"
	"sao/data"
	"sao/player"
	"sao/types"
	"sao/utils"
	"sao/world/calendar"
	"sao/world/location"
	"sao/world/npc"
	"sao/world/party"
	"sao/world/tournament"
	"sao/world/transaction"
	"sort"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/google/uuid"
)

type FloorMap map[string]location.Floor

type World struct {
	Players      map[uuid.UUID]*player.Player
	Transactions map[uuid.UUID]*transaction.Transaction
	NPCs         map[uuid.UUID]*npc.NPC
	Stores       map[uuid.UUID]*npc.NPCStore
	Floors       FloorMap
	Tournaments  map[uuid.UUID]*tournament.Tournament
	Fights       map[uuid.UUID]battle.Fight
	Entities     map[uuid.UUID]*battle.Entity
	Time         *calendar.Calendar
	Parties      map[uuid.UUID]*party.Party
	TestMode     bool
	DiscordToken string
	DChannel     chan types.DiscordMessageStruct
}

func CreateWorld(discordToken string, testMode bool) World {
	stockItem := npc.Stock{
		ItemType: types.ITEM_MATERIAL,
		ItemUUID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Price:    1,
		Quantity: 10,
		Limit:    10,
	}

	npcUuid := uuid.New()

	testStore := npc.NPCStore{
		RestockInterval: *calendar.StartCalendar(),
		LastRestock:     *calendar.StartCalendar(),
		Uuid:            uuid.New(),
		NPCUuid:         npcUuid,
		Name:            "Warzywniak babci stasi",
		Stock:           []*npc.Stock{&stockItem},
	}

	npcMap := map[uuid.UUID]*npc.NPC{
		npcUuid: {
			Name:     "Babcia stasia",
			Location: types.PlayerLocation{FloorName: "dev", LocationName: "Rynek"},
			Store:    &testStore,
		},
	}

	storeMap := map[uuid.UUID]*npc.NPCStore{
		testStore.Uuid: &testStore,
	}

	return World{
		make(map[uuid.UUID]*player.Player),
		make(map[uuid.UUID]*transaction.Transaction),
		npcMap,
		storeMap,
		location.GetFloors(testMode),
		make(map[uuid.UUID]*tournament.Tournament),
		make(map[uuid.UUID]battle.Fight),
		make(map[uuid.UUID]*battle.Entity),
		calendar.StartCalendar(),
		make(map[uuid.UUID]*party.Party),
		testMode,
		discordToken,
		make(chan types.DiscordMessageStruct, 10),
	}
}

func (w *World) RegisterNewPlayer(name, uid string) *player.Player {
	newPlayer := player.NewPlayer(name, uid)

	w.Players[newPlayer.GetUUID()] = &newPlayer

	return &newPlayer
}

func (w *World) MovePlayer(pUuid uuid.UUID, floorName, locationName, reason string) error {
	player := w.Players[pUuid]

	if floorData, exists := w.Floors[floorName]; !exists || (!floorData.Unlocked && player.Meta.Location.FloorName != floorName) {
		return fmt.Errorf("floor %v not found or locked", locationName)
	}

	if locationData := w.Floors[floorName].FindLocation(locationName); locationData == nil || (!locationData.Unlocked && player.Meta.Location.LocationName != locationName) {
		return fmt.Errorf("location %v not found or locked", locationName)
	}

	if player.Meta.FightInstance != nil {
		return errors.New("player is in fight")
	}

	//TODO If there is a reason, tell it to the player

	player.Meta.Location.FloorName = floorName
	player.Meta.Location.LocationName = locationName

	return nil
}

func (w *World) PlayerSearch(uuid uuid.UUID) {
	player := w.Players[uuid]

	floor := w.Floors[player.Meta.Location.FloorName]

	var location location.Location

	for _, loc := range floor.Locations {
		if loc.Name == player.Meta.Location.LocationName {
			location = loc
		}
	}

	//TODO Random (or dropped) items lying there

	if len(location.Enemies) == 0 {
		return
	}

	enemyEntry := utils.RandomElement(location.Enemies)
	enemyCount := utils.RandomNumber(enemyEntry.MinNum, enemyEntry.MaxNum)

	entityMap := make(battle.EntityMap)

	if player.Meta.Party != nil {
		partyData := w.Parties[*player.Meta.Party]
		partyMemberCount := len(partyData.Players)

		for _, member := range partyData.Players {
			effects := make([]battle.ActionEffect, 0)

			switch member.Role {
			case party.DPS:
				effects = append(effects, battle.ActionEffect{
					Effect:   battle.EFFECT_STAT_INC,
					Duration: -1,
					Source:   battle.SOURCE_PARTY,
					Meta: battle.ActionEffectStat{
						Stat:      types.STAT_ADAPTIVE,
						Value:     10 + (partyMemberCount-1)*5,
						IsPercent: true,
					},
				})
			case party.Tank:
				effects = append(effects, battle.ActionEffect{
					Effect:   battle.EFFECT_STAT_INC,
					Duration: -1,
					Source:   battle.SOURCE_PARTY,
					Meta: battle.ActionEffectStat{
						Stat:      types.STAT_DEF,
						Value:     25,
						IsPercent: false,
					},
				})

				effects = append(effects, battle.ActionEffect{
					Effect:   battle.EFFECT_STAT_INC,
					Duration: -1,
					Source:   battle.SOURCE_PARTY,
					Meta: battle.ActionEffectStat{
						Stat:      types.STAT_MR,
						Value:     25,
						IsPercent: false,
					},
				})

				effects = append(effects, battle.ActionEffect{
					Effect:   battle.EFFECT_STAT_INC,
					Duration: -1,
					Source:   battle.SOURCE_PARTY,
					Meta: battle.ActionEffectStat{
						Stat:      types.STAT_HP,
						Value:     (partyMemberCount - 1) * 5,
						IsPercent: true,
					},
				})

				effects = append(effects, battle.ActionEffect{
					Effect:   battle.EFFECT_STAT_INC,
					Duration: -1,
					Source:   battle.SOURCE_PARTY,
					Meta: battle.ActionEffectStat{
						Stat:      types.STAT_DEF,
						Value:     (partyMemberCount - 1) * 5,
						IsPercent: true,
					},
				})

				effects = append(effects, battle.ActionEffect{
					Effect:   battle.EFFECT_STAT_INC,
					Duration: -1,
					Source:   battle.SOURCE_PARTY,
					Meta: battle.ActionEffectStat{
						Stat:      types.STAT_MR,
						Value:     (partyMemberCount - 1) * 5,
						IsPercent: true,
					},
				})

				if partyMemberCount > 2 {
					effects = append(effects, battle.ActionEffect{
						Effect:   battle.EFFECT_TAUNT,
						Source:   battle.SOURCE_PARTY,
						Duration: -1,
						Meta:     nil,
					})
				}
			case party.Support:
				effects = append(effects, battle.ActionEffect{
					Effect:   battle.EFFECT_STAT_INC,
					Duration: -1,
					Source:   battle.SOURCE_PARTY,
					Meta: battle.ActionEffectStat{
						Stat:      types.STAT_HEAL_POWER,
						Value:     15 + (partyMemberCount-1)*5,
						IsPercent: true,
					},
				})
			}

			resolvedUser := w.Players[member.PlayerUuid]

			for _, effect := range effects {
				resolvedUser.ApplyEffect(effect)
			}

			entityMap[resolvedUser.GetUUID()] = battle.EntityEntry{Entity: resolvedUser, Side: 0}
		}
	} else {
		entityMap[player.GetUUID()] = battle.EntityEntry{Entity: player, Side: 0}
	}

	for i := 0; i < enemyCount; i++ {
		entity := mobs.Spawn(enemyEntry.Enemy)

		entityMap[entity.GetUUID()] = battle.EntityEntry{Entity: entity, Side: 1}
	}

	fight := battle.Fight{
		Entities:       entityMap,
		DiscordChannel: w.DChannel,
		StartTime:      w.Time.Copy(),
		Location:       &location,
	}

	fight.Init()

	fightUUID := w.RegisterFight(fight)

	player.Meta.FightInstance = &fightUUID

	go w.ListenForFight(fightUUID)
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
			if player.GetCurrentHP() == 0 {
				continue
			}

			//Missing mana
			if player.GetCurrentMana() < player.GetMaxMana() {
				player.Stats.CurrentMana += 1
			}

			//Can be healed
			if player.GetCurrentHP() < player.GetMaxHP() {
				healRatio := 50

				if w.PlayerInCity(pUuid) {
					healRatio = 25
				}

				player.Heal(player.GetMaxHP() / healRatio)
			}
		}
	}
}

func (w *World) PlayerInCity(uuid uuid.UUID) bool {
	player := w.Players[uuid]

	floor := w.Floors[player.Meta.Location.FloorName]
	location := floor.FindLocation(player.Meta.Location.LocationName)

	return location.CityPart
}

// Fight stuff
func (w *World) RegisterFight(fight battle.Fight) uuid.UUID {
	uuid := uuid.New()

	w.Fights[uuid] = fight

	for _, entity := range fight.Entities {
		if entity.Entity.IsAuto() {
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

	channelId := fight.Location.CID

	if fight.Tournament != nil {
		channelId = w.Tournaments[fight.Tournament.Tournament].Channel
	}

	for {
		eventData, ok := <-fight.ExternalChannel

		if !ok {
			w.DeregisterFight(fightUuid)
			break
		}

		switch eventData.GetEvent() {
		case battle.MSG_FIGHT_END:
			wonSideIDX := fight.Entities.SidesLeft()[0]
			wonEntities := fight.Entities.FromSide(wonSideIDX)
			if fight.Tournament == nil {
				overallXp := 0
				overallGold := 0
				lootedItems := make([]battle.Loot, 0)

				for _, entity := range fight.Entities {
					if entity.Side == wonSideIDX {
						continue
					}

					lootList := entity.Entity.GetLoot()

					for _, loot := range lootList {
						switch loot.Type {
						case battle.LOOT_EXP:
							overallXp += (*loot.Meta)["value"].(int)
						case battle.LOOT_GOLD:
							overallGold += (*loot.Meta)["value"].(int)
						case battle.LOOT_ITEM:
							lootedItems = append(lootedItems, loot)
						}
					}
				}

				partyUuid := wonEntities[0].(battle.PlayerEntity).GetParty()

				if partyUuid != nil {
					partyData := w.Parties[*partyUuid]
					partyLeader := w.Players[partyData.Leader]

					for _, member := range partyData.Players {
						player := w.Players[member.PlayerUuid]

						player.AddEXP(overallXp / len(partyData.Players))
						player.AddGold(overallGold / len(partyData.Players))
					}

					for _, loot := range lootedItems {
						itemUuid := (*loot.Meta)["uuid"].(uuid.UUID)

						if (*loot.Meta)["type"].(types.ItemType) == types.ITEM_OTHER {
							itemObj := data.Items[itemUuid]

							itemObj.Count = (*loot.Meta)["count"].(int)

							partyLeader.Inventory.Items = append(partyLeader.Inventory.Items, &itemObj)
						} else {
							ingredient := data.Ingredients[itemUuid]

							ingredient.Count = (*loot.Meta)["count"].(int)

							partyLeader.Inventory.AddIngredient(&ingredient)
						}
					}
				} else {
					for _, entity := range wonEntities {
						entityUuid := entity.GetUUID()

						if entity.IsAuto() {
							continue
						}

						player := w.Players[entityUuid]

						player.AddEXP(overallXp)
						player.AddGold(overallGold)

						for _, loot := range lootedItems {
							itemUuid := (*loot.Meta)["uuid"].(uuid.UUID)

							if (*loot.Meta)["type"].(types.ItemType) == types.ITEM_OTHER {
								itemObj := data.Items[itemUuid]

								itemObj.Count = (*loot.Meta)["count"].(int)

								player.Inventory.Items = append(player.Inventory.Items, &itemObj)
							} else {
								ingredient := data.Ingredients[itemUuid]

								ingredient.Count = (*loot.Meta)["count"].(int)

								player.Inventory.AddIngredient(&ingredient)
							}
						}
					}
				}
			}

			wonSideText := ""

			for _, entity := range wonEntities {
				wonSideText += fmt.Sprintf("%v", entity.GetName())

				if !entity.IsAuto() {
					wonSideText += fmt.Sprintf(" (<@%v>)", entity.(battle.PlayerEntity).GetUID())
				}

				wonSideText += "\n"
			}

			wonSideText = wonSideText[:len(wonSideText)-1]

			w.DChannel <- types.DiscordMessageStruct{
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

			if fight.Tournament != nil {
				w.Tournaments[fight.Tournament.Tournament].ExternalChannel <- tournament.MatchFinishedData{Winner: wonEntities[0].GetUUID()}
			}
		case battle.MSG_FIGHT_START:
			oneSide := fight.Entities.FromSide(0)
			otherSide := fight.Entities.FromSide(1)

			oneSideText := ""

			for _, entity := range oneSide {
				oneSideText += fmt.Sprintf("%v", entity.GetName())
				if !entity.IsAuto() {
					oneSideText += fmt.Sprintf(" (<@%v>)", entity.(battle.PlayerEntity).GetUID())
				}

				oneSideText += "\n"
			}

			oneSideText = oneSideText[:len(oneSideText)-1]

			otherSideText := ""

			for _, entity := range otherSide {
				otherSideText += fmt.Sprintf("%v", entity.GetName())

				if !entity.IsAuto() {
					otherSideText += fmt.Sprintf(" (<@%v>)", entity.(battle.PlayerEntity).GetUID())
				}

				otherSideText += "\n"
			}

			otherSideText = otherSideText[:len(otherSideText)-1]

			w.DChannel <- types.DiscordMessageStruct{
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

			if player.GetEffect(battle.EFFECT_TAUNTED) != nil {
				effect := player.GetEffect(battle.EFFECT_TAUNTED)

				fight.PlayerActions <- battle.Action{
					Event:  battle.ACTION_ATTACK,
					Source: player.GetUUID(),
					Target: effect.Meta.(uuid.UUID),
				}

				w.DChannel <- types.DiscordMessageStruct{
					ChannelID: channelId,
					MessageContent: discord.NewMessageCreateBuilder().
						SetContentf("<@%v> jest zmuszony do ataku! Pomijamy turę!", player.GetUID()).
						Build(),
				}

				continue
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

			for _, skill := range player.GetAllSkills() {
				if player.CanUseSkill(skill) {
					filteredSkillsCount++
				}
			}

			for _, skill := range player.Inventory.LevelSkills {
				if player.CanUseLvlSkill(skill) {
					filteredSkillsCount++
				}
			}

			if filteredSkillsCount == 0 {
				skillButton = skillButton.AsDisabled()
			}

			itemButton := discord.NewPrimaryButton("Przedmiot", "f/item")

			filteredItemsCount := 0

			for _, item := range player.GetAllItems() {
				if item.Consume && item.Count > 0 && !item.Hidden {
					filteredItemsCount++
				}
			}

			if filteredItemsCount == 0 {
				itemButton = itemButton.AsDisabled()
			}

			escapeButton := discord.NewDangerButton("Ucieczka", "f/escape")

			w.DChannel <- types.DiscordMessageStruct{
				ChannelID: channelId,
				MessageContent: discord.NewMessageCreateBuilder().
					AddEmbeds(discord.NewEmbedBuilder().
						SetTitle("Czas na turę!").
						SetDescriptionf("Kolej <@%s>!", player.GetUID()).
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
		if !entity.Entity.IsAuto() {
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

	w.DChannel <- types.DiscordMessageStruct{
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

	client, err := disgo.New(w.DiscordToken)

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

		w.DChannel <- types.DiscordMessageStruct{
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

			w.DChannel <- types.DiscordMessageStruct{
				ChannelID: tournamentObj.Channel,
				MessageContent: discord.NewMessageCreateBuilder().
					SetContentf("Szczęśliwy gracz to <@%v>\nNie musisz walczyć w tej rundzie i możesz spokojnie oglądać!", w.Players[tournamentObj.Participants[luckyPlayer]].GetUID()).
					Build(),
			}

			participants = append(participants[:luckyPlayer], participants[luckyPlayer+1:]...)
		}

		//Hackery shuffle
		for i := range participants {
			j := rand.Intn(i + 1)
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

			w.DChannel <- types.DiscordMessageStruct{
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

					w.DChannel <- types.DiscordMessageStruct{
						ChannelID: tournamentObj.Channel,
						MessageContent: discord.NewMessageCreateBuilder().
							SetContentf("Turniej zakończony! Wygrał %v (<@%v>)", player.GetName(), player.GetUID()).
							Build(),
					}
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
		DiscordChannel: w.DChannel,
		StartTime:      w.Time.Copy(),
		Location:       &fightingLocation,
		Tournament: &battle.TournamentData{
			Tournament: tUuid,
			Location:   w.Tournaments[tUuid].Channel,
		},
	}

	fight.Init()

	fightUUID := w.RegisterFight(fight)

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

	w.DChannel <- types.DiscordMessageStruct{
		ChannelID:      w.Players[transactionObj.LeftSide.Who].GetUID(),
		MessageContent: discord.NewMessageCreateBuilder().SetContent("Transakcja rozpoczęta!").Build(),
		DM:             true,
	}

	w.DChannel <- types.DiscordMessageStruct{
		ChannelID:      w.Players[transactionObj.RightSide.Who].GetUID(),
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
		storeData = append(storeData, map[string]interface{}{
			"uuid":            store.Uuid,
			"name":            store.Name,
			"restockInterval": store.RestockInterval.Serialize(),
			"lastRestock":     store.LastRestock.Serialize(),
			"stock":           store.Stock,
		})
	}

	tournamentData := make([]map[string]interface{}, 0)

	for _, tournament := range w.Tournaments {
		tournamentData = append(tournamentData, tournament.Serialize())
	}

	return map[string]interface{}{
		"players":     playerData,
		"stores":      storeData,
		"time":        w.Time.Serialize(),
		"test":        w.TestMode,
		"tournaments": tournamentData,
	}
}

func (w *World) CreateBackup() {
	backupPath := "backup"

	if w.TestMode {
		backupPath = "test" + backupPath
	}

	_, err := os.Stat(backupPath)

	if os.IsNotExist(err) {
		os.Mkdir(backupPath, os.ModePerm)
	}

	newBackupPath := fmt.Sprintf("%s/%v.json", backupPath, time.Now().Format("2006-01-02_15-04-05"))

	backupFile, err := os.Create(newBackupPath)

	if err != nil {
		panic(err)
	}

	defer backupFile.Close()

	jsonFile, err := json.Marshal(w.Serialize())

	if err != nil {
		panic(err)
	}

	backupFile.Write(jsonFile)

	w.DChannel <- types.DiscordMessageStruct{
		ChannelID: "1151922672595390588",
		MessageContent: discord.NewMessageCreateBuilder().
			SetContent("Backup zrobiony!").
			Build(),
	}
}

func (w *World) LoadBackup() {
	backupPath := "backup"

	if w.TestMode {
		backupPath = "test" + backupPath
	}

	_, err := os.Stat(backupPath)

	if os.IsNotExist(err) {
		return
	}

	backupData := make(map[string]interface{})

	allBackups, err := os.ReadDir(backupPath)

	if err != nil {
		//HONESTLY I DON'T KNOW WHAT COULD HAPPEN HERE
		panic(err)
	}

	sort.Slice(allBackups, func(i, j int) bool {
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

	if len(allBackups) == 0 {
		return
	}

	backupContent, err := os.ReadFile("./" + backupPath + "/" + allBackups[0].Name())

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

	//TODO Load stores

	for _, tData := range backupData["tournaments"].([]interface{}) {
		tempData := tData.(map[string]interface{})
		parsedData := tournament.Deserialize(tempData)

		w.Tournaments[parsedData.Uuid] = &parsedData
	}
}
