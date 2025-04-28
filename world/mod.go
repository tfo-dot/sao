package world

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sao/battle"
	"sao/battle/mobs"
	"sao/data"
	"sao/player"
	"sao/types"
	"sao/utils"
	"sao/world/party"
	"sao/world/tournament"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/google/uuid"
	"slices"
)

type World struct {
	Players        map[uuid.UUID]*player.Player
	Tournaments    map[uuid.UUID]*tournament.Tournament
	Fights         map[uuid.UUID]*battle.Fight
	Parties        map[uuid.UUID]*party.Party
	DiscordChannel chan types.DiscordEvent
}

func CreateWorld() World {
	return World{
		make(map[uuid.UUID]*player.Player),
		make(map[uuid.UUID]*tournament.Tournament),
		make(map[uuid.UUID]*battle.Fight),
		make(map[uuid.UUID]*party.Party),
		make(chan types.DiscordEvent, 10),
	}
}

func (w *World) SendMessage(channelId string, content discord.MessageCreate, dm bool) {
	w.DiscordChannel <- types.DiscordSendMsg{
		Data: types.DiscordMessageStruct{ChannelID: channelId, MessageContent: content, DM: dm},
	}
}

func (w *World) RequestChoice(choiceId string, choose func(cic *events.ComponentInteractionCreate)) {
	w.DiscordChannel <- types.DiscordChoiceMsg{
		Data: types.DiscordChoice{Id: choiceId, Select: choose},
	}
}

func (w *World) PlayerFight(pUuid uuid.UUID, location *types.Location, threadId string, mobId string, mobCount int) {
	playerObj := w.Players[pUuid]

	entityMap := make(battle.EntityMap)

	if playerObj.Meta.Party != nil {
		for _, member := range w.Parties[playerObj.Meta.Party.UUID].Players {
			mUuid := member.PlayerUuid

			entityMap[mUuid] = &battle.EntityEntry{Entity: w.Players[mUuid]}
		}
	} else {
		entityMap[pUuid] = &battle.EntityEntry{Entity: playerObj}
	}

	for range mobCount {
		entity := mobs.Spawn(mobId)

		entityMap[entity.GetUUID()] = &battle.EntityEntry{Entity: entity, Side: 1}
	}

	fight := battle.Fight{
		Entities:       entityMap,
		DiscordChannel: w.DiscordChannel,
		Location:       location,
		Meta:           &battle.FightMeta{ThreadId: threadId},
	}

	fight.Init()

	fightUUID := w.RegisterFight(&fight)

	mentionString := ""

	for _, entity := range fight.Entities {
		if entity.Entity.GetFlags()&types.ENTITY_AUTO == 0 {
			mentionString += fmt.Sprintf("<@%v>, ", entity.Entity.(*player.Player).Meta.UserID)
		}
	}

	w.SendMessage(fight.GetChannelId(), discord.MessageCreate{Content: mentionString[:len(mentionString)-2]}, false)

	playerObj.Meta.FightInstance = &fightUUID

	go w.ListenForFight(fightUUID)
}

func (w *World) PlayerSearch(pUuid uuid.UUID, threadId string, event *events.ApplicationCommandInteractionCreate) {
	player := w.Players[pUuid]

	cid := event.Channel().ID()

	location := data.FloorMap.FindLocation(func(l types.Location) bool { return l.CID == cid.String() })

	if player.Meta.FightInstance != nil || location == nil || len(location.Enemies) == 0 {
		return
	}

	if len(location.Enemies) == 1 {
		enemy := location.Enemies[0]

		if enemy.MinNum == enemy.MaxNum {
			go w.PlayerFight(pUuid, location, threadId, enemy.Enemy, enemy.MinNum)

			return
		}

		selectMenuUuid := uuid.New().String()

		options := make([]discord.StringSelectMenuOption, 0)

		for i := enemy.MinNum; i <= enemy.MaxNum; i++ {
			options = append(options, discord.NewStringSelectMenuOption(strconv.Itoa(i), strconv.Itoa(i)))
		}

		event.CreateMessage(discord.
			NewMessageCreateBuilder().
			AddActionRow(discord.
				NewStringSelectMenu("chc/"+selectMenuUuid, "Wybierz ilość przeciwników").
				WithMaxValues(1).
				AddOptions(options...),
			).
			Build(),
		)

		w.RequestChoice(selectMenuUuid, func(cic *events.ComponentInteractionCreate) {
			choiceRaw := cic.StringSelectMenuInteractionData().Values[0]

			choice, _ := strconv.Atoi(choiceRaw)

			cic.UpdateMessage(discord.
				NewMessageUpdateBuilder().
				ClearContainerComponents().
				SetContent("Wybrano " + choiceRaw + " przeciwników!").
				Build(),
			)

			go w.PlayerFight(pUuid, location, threadId, enemy.Enemy, choice)
		})

		return
	}

	selectMenuUuid := uuid.New().String()

	options := make([]discord.StringSelectMenuOption, 0)

	for idx, enemy := range location.Enemies {
		options = append(options, discord.NewStringSelectMenuOption(enemy.Enemy, strconv.Itoa(idx)))
	}

	event.CreateMessage(discord.
		NewMessageCreateBuilder().
		AddActionRow(discord.
			NewStringSelectMenu("chc/"+selectMenuUuid, "Wybierz typ przeciwnika").
			WithMaxValues(1).
			AddOptions(options...),
		).
		Build(),
	)

	w.RequestChoice(selectMenuUuid, func(cic *events.ComponentInteractionCreate) {
		choiceRaw := cic.StringSelectMenuInteractionData().Values[0]

		choice, _ := strconv.Atoi(choiceRaw)

		enemy := location.Enemies[choice]

		if enemy.MinNum == enemy.MaxNum {
			cic.UpdateMessage(discord.
				NewMessageUpdateBuilder().
				ClearContainerComponents().
				SetContent("Wybrano przeciwnika - " + mobs.Mobs[location.Enemies[choice].Enemy].Name).
				Build(),
			)

			go w.PlayerFight(pUuid, location, threadId, enemy.Enemy, enemy.MinNum)

			return
		}

		selectMenuUuid := uuid.New().String()

		enemyCount := enemy.MaxNum - enemy.MinNum

		options := make([]discord.StringSelectMenuOption, enemyCount)

		for i := range enemyCount + 1 {
			options[i] = discord.NewStringSelectMenuOption(
				strconv.Itoa(i+enemy.MinNum), strconv.Itoa(i+enemy.MinNum),
			)
		}

		cic.UpdateMessage(discord.
			NewMessageUpdateBuilder().
			SetContent("Wybierz ilość przeciwników!").
			AddActionRow(discord.
				NewStringSelectMenu("chc/"+selectMenuUuid, "Wybierz ilość przeciwników").
				WithMaxValues(1).
				AddOptions(options...),
			).
			Build(),
		)

		w.RequestChoice(selectMenuUuid, func(cic *events.ComponentInteractionCreate) {
			choiceRaw := cic.StringSelectMenuInteractionData().Values[0]

			choice, _ := strconv.Atoi(choiceRaw)

			cic.UpdateMessage(discord.
				NewMessageUpdateBuilder().
				ClearContainerComponents().
				SetContent("Wybrano " + choiceRaw + " przeciwników!").
				Build(),
			)

			go w.PlayerFight(pUuid, location, threadId, enemy.Enemy, choice)
		})
	})
}

func (w *World) TickPlayer(p *player.Player) {
	if p.Meta.FightInstance != nil {
		return
	}

	if !data.WorldConfig.Hardcore && p.GetCurrentHP() <= 0 {
		p.Stats.HP = p.GetStat(types.STAT_HP)
		p.Stats.CurrentMana = p.GetStat(types.STAT_MANA)

		w.SendMessage(
			data.Config.LogChannelID,
			discord.MessageCreate{
				Embeds: []discord.Embed{{
					Title:       "Wskrzeszenie!",
					Description: fmt.Sprintf("%s zostaje wskrzeszony...", p.GetName()),
				}},
			},
			false)

		return
	}

	if p.GetCurrentMana() < p.GetStat(types.STAT_MANA) {
		p.Stats.CurrentMana += 1
	}

	if p.GetCurrentHP() < p.GetStat(types.STAT_HP) {
		p.Heal(p.GetStat(types.STAT_HP) / 25)

		if p.Meta.WaitToHeal && p.GetCurrentHP() >= p.GetStat(types.STAT_HP) {
			p.Meta.WaitToHeal = false

			client, err := disgo.New(data.Config.Token)

			if err != nil {
				w.SendMessage(
					data.Config.LogChannelID,
					discord.MessageCreate{Content: "Nie można wysłać wiadomości do gracza (HP Heal)"},
					false,
				)
			}

			ch, err := client.Rest().CreateDMChannel(snowflake.MustParse(p.Meta.UserID))

			if err != nil {
				w.SendMessage(
					data.Config.LogChannelID,
					discord.NewMessageCreateBuilder().
						SetContent("Nie można wysłać wiadomości do gracza (HP Heal)").
						Build(),
					false,
				)
			}

			_, err = client.Rest().CreateMessage(
				ch.ID(), discord.MessageCreate{Content: "Twoja postać ma już 100% HP, baw się dobrze!"},
			)

			if err != nil {
				w.SendMessage(
					data.Config.LogChannelID,
					discord.MessageCreate{Content: "Nie można wysłać wiadomości do gracza (HP Heal)"},
					false,
				)
			}
		}
	}
}

func (w *World) StartClock() {
	counter := 0

	for range time.Tick(1 * time.Minute) {
		for _, player := range w.Players {
			w.TickPlayer(player)
		}

		counter++

		if counter >= 15 {
			counter = 0

			w.CreateBackup()
		}
	}
}

func (w *World) RegisterFight(fight *battle.Fight) uuid.UUID {
	uuid := uuid.New()

	w.Fights[uuid] = fight

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
				w.SendMessage(
					channelId,
					discord.MessageCreate{
						Embeds: []discord.Embed{discord.
							NewEmbedBuilder().
							SetTitle("Koniec walki!").
							SetDescriptionf("Wszyscy gracze uciekli z walki!").
							Build(),
						},
					},
					false,
				)

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

				w.SendMessage(
					channelId,
					discord.MessageCreate{Embeds: []discord.Embed{discord.
						NewEmbedBuilder().
						SetTitle("Koniec walki!").
						SetDescriptionf("Wygrali:\n" + wonSideText[:len(wonSideText)-1]).
						Build(),
					}},
					false,
				)

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

			if fight.Meta.Tournament == nil {
				overallXp := 0
				overallGold := 0

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
						}
					}
				}

				var partyInfo *player.PartialParty

				for _, entity := range wonEntities {
					if entity.GetFlags()&types.ENTITY_AUTO != 0 {
						continue
					}

					partyInfo = entity.(*player.Player).Meta.Party
					break
				}

				if partyInfo != nil {
					partyData := w.Parties[partyInfo.UUID]

					for _, member := range partyData.Players {
						player := w.Players[member.PlayerUuid]

						player.AddEXP(overallXp / len(partyData.Players))

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
				} else {
					for _, entity := range wonEntities {
						entityUuid := entity.GetUUID()

						if entity.GetFlags()&types.ENTITY_AUTO != 0 {
							continue
						}

						player := w.Players[entityUuid]

						player.AddEXP(overallXp)

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
					}
				}
			}

			wonSideText := ""

			for _, entity := range wonEntities {
				wonSideText += fmt.Sprintf("%v", entity.GetName())

				if entity.GetFlags()&types.ENTITY_AUTO == 0 {
					wonSideText += fmt.Sprintf(" (<@%v>)", entity.(*player.Player).Meta.UserID)
				}

				wonSideText += fmt.Sprintf(" (%v/%v HP)\n", entity.GetCurrentHP(), entity.GetStat(types.STAT_HP))
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

			w.SendMessage(
				channelId,
				discord.MessageCreate{Embeds: []discord.Embed{discord.
					NewEmbedBuilder().
					SetTitle("Koniec walki!").
					SetDescriptionf("Wygrali:\n" + wonSideText).
					Build(),
					discord.NewEmbedBuilder().SetTitle("Podsumowanie").SetDescription(lootSummaryText).Build(),
				}},
				false,
			)

			if fight.Meta.Tournament != nil {
				w.Tournaments[fight.Meta.Tournament.Tournament].ExternalChannel <- tournament.MatchFinishedData{
					Winner: wonEntities[0].GetUUID(),
				}
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

			w.SendMessage(
				channelId,
				discord.NewMessageCreateBuilder().
					SetContent("Walka się rozpoczyna!").
					AddEmbeds(discord.NewEmbedBuilder().
						SetTitle("Walka").
						AddField("Po jednej!", oneSideText, false).
						AddField("Po drugiej!", otherSideText, false).
						Build()).
					Build(),
				false,
			)
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
					effectTrigger := skill.Skill.GetUpgradableTrigger(skill.Upgrades)

					if effectTrigger.Flags&types.FLAG_IGNORE_CC != 0 {
						canUseActionWhileCC = true
						break
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

					w.SendMessage(
						channelId,
						discord.MessageCreate{
							Content: fmt.Sprintf("<@%v> jest zmuszony do ataku! Pomijamy turę!", player.Meta.UserID),
						},
						false,
					)

					continue
				}
			}

			attackButton := discord.NewPrimaryButton("Atak", fmt.Sprintf("f/attack/%s", player.Meta.UserID))

			if !player.CanAttack() {
				attackButton = attackButton.AsDisabled()
			}

			defendButton := discord.NewPrimaryButton("Obrona", fmt.Sprintf("f/defend/%s", player.Meta.UserID))

			if !player.CanDefend() {
				defendButton = defendButton.AsDisabled()
			}

			skillButton := discord.NewPrimaryButton("Skill", fmt.Sprintf("f/skill/%s", player.Meta.UserID))

			filteredSkillsCount := 0

			for _, skill := range player.Inventory.LevelSkills {
				if player.CanUseSkill(skill.Skill) && skill.Skill.CanUse(player, fight) {
					filteredSkillsCount++
				}
			}

			if filteredSkillsCount == 0 {
				skillButton = skillButton.AsDisabled()
			}

			itemButton := discord.NewPrimaryButton("Przedmiot", fmt.Sprintf("f/item/%s", player.Meta.UserID))

			filteredItemsCount := 0

			for _, item := range player.Inventory.Items {
				if item.Consume && item.Count > 0 && !item.Hidden {
					filteredItemsCount++
				}
			}

			if filteredItemsCount == 0 {
				itemButton = itemButton.AsDisabled()
			}

			escapeButton := discord.NewDangerButton("Ucieczka", fmt.Sprintf("f/escape/%s", player.Meta.UserID))

			w.SendMessage(
				channelId,
				discord.NewMessageCreateBuilder().
					AddEmbeds(discord.NewEmbedBuilder().
						SetTitle("Czas na turę!").
						SetDescriptionf("Kolej <@%s>!", player.Meta.UserID).
						SetAuthorName(player.Name).
						Build(),
					).
					AddActionRow(attackButton, defendButton, skillButton, itemButton, escapeButton).
					Build(),
				false,
			)

		case battle.MSG_SUMMON_EXPIRED:
			entityUuid := eventData.GetData().(uuid.UUID)

			temp := fight.Entities[entityUuid].Entity

			delete(fight.Entities, entityUuid)

			w.SendMessage(
				channelId,
				discord.
					NewMessageCreateBuilder().
					AddEmbeds(
						discord.
							NewEmbedBuilder().
							SetTitle("Przyzwany stwór uciekł!").
							SetDescriptionf("%s uciekł z pola walki!", temp.GetName()).
							Build(),
					).
					Build(),
				false,
			)

		case battle.MSG_SUMMON_DIED:
			entityUuid := eventData.GetData().(uuid.UUID)

			temp := fight.Entities[entityUuid].Entity

			delete(fight.Entities, entityUuid)

			w.SendMessage(channelId,
				discord.
					NewMessageCreateBuilder().
					AddEmbeds(
						discord.
							NewEmbedBuilder().
							SetTitle("Przyzwany stwór umarł!").
							SetDescriptionf("Niech Jack ma %s w opiece!", temp.GetName()).
							Build(),
					).
					Build(),
				false,
			)

		case battle.MSG_ENTITY_DIED:
			entityUuid := eventData.GetData().(uuid.UUID)

			temp := fight.Entities[entityUuid].Entity

			w.SendMessage(
				channelId,
				discord.NewMessageCreateBuilder().
					AddEmbeds(
						discord.
							NewEmbedBuilder().
							SetTitle("Ktoś pożegnał się z życiem...").
							SetDescriptionf("Niech Jack ma %s w opiece!", temp.GetName()).
							Build(),
					).
					Build(),
				false,
			)
		default:
			panic("Unhandled event")
		}

		if len(fight.ExternalChannel) == 0 && fight.IsFinished() {
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

	w.SendMessage(
		"1225150345009827841",
		discord.NewMessageCreateBuilder().
			AddEmbeds(discord.NewEmbedBuilder().
				SetTitle("Nowy turniej!").
				SetDescriptionf("Zapisy na turniej `%v` otwarte", tournament.Name).
				SetFooterText("Ilość miejsc: "+playerText).
				Build()).
			AddActionRow(
				discord.NewPrimaryButton("Dołącz", "t/join/"+tournament.Uuid.String()),
			).
			Build(),
		false,
	)
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

	var fightingLocation types.Location

	for _, floor := range data.FloorMap {
		for _, location := range floor.Locations {
			for _, effect := range location.Flags {
				if effect == "arena" {
					fightingLocation = location
				}
			}
		}
	}

	client, err := disgo.New(data.Config.Token)

	if err != nil {
		panic(err)
	}

	msg, err := client.Rest().CreateMessage(
		snowflake.MustParse(fightingLocation.CID),
		discord.MessageCreate{Content: "Turniej rozpoczęty!"},
	)

	if err != nil {
		panic(err)
	}

	thread, err := client.Rest().CreateThreadFromMessage(
		snowflake.MustParse(fightingLocation.CID),
		msg.ID,
		discord.ThreadCreateFromMessage{Name: "Turniej"},
	)

	if err != nil {
		panic(err)
	}

	tournamentObj.ExternalChannel = make(chan tournament.TournamentEventData)
	tournamentObj.Channel = thread.ID().String()

	w.SendMessage(
		tournamentObj.Channel,
		discord.NewMessageCreateBuilder().
			SetContent("Losowanie czas zacząć!").
			Build(),
		false,
	)

	matches := make([]*tournament.TournamentMatch, 0)

	var participants = tournamentObj.Participants

	if len(tournamentObj.Participants)%2 != 0 {
		luckyPlayer := utils.RandomNumber(0, len(tournamentObj.Participants)-1)

		matches = append(matches, &tournament.TournamentMatch{
			Players: []uuid.UUID{tournamentObj.Participants[luckyPlayer]},
			Winner:  &tournamentObj.Participants[luckyPlayer],
			State:   tournament.FinishedMatch,
		})

		w.SendMessage(
			tournamentObj.Channel,
			discord.NewMessageCreateBuilder().
				SetContentf(
					"Szczęśliwy gracz to <@%v>\nNie musisz walczyć w tej rundzie i możesz spokojnie oglądać!",
					w.Players[tournamentObj.Participants[luckyPlayer]].Meta.UserID,
				).
				Build(),
			false,
		)

		participants = slices.Delete(participants, luckyPlayer, luckyPlayer+1)
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

		player0 := w.Players[participants[i]]
		player1 := w.Players[participants[i+1]]

		w.SendMessage(
			tournamentObj.Channel,
			discord.MessageCreate{
				Content: fmt.Sprintf(
					"Los pociągnięty!\nMecz #%v: %v vs %v", (i/2)+1, player0.GetName(), player1.GetName(),
				),
			},
			false,
		)
	}

	tournamentObj.Stages = append(tournamentObj.Stages, &tournament.TournamentStage{
		Matches: matches,
	})

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

					w.SendMessage(
						tournamentObj.Channel,
						discord.NewMessageCreateBuilder().
							SetContentf("Turniej zakończony! Wygrał %v (<@%v>)", player.GetName(), player.Meta.UserID).
							Build(),
						false,
					)

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

	entityMap[player0.GetUUID()] = &battle.EntityEntry{Entity: player0, Side: 0}
	entityMap[player1.GetUUID()] = &battle.EntityEntry{Entity: player1, Side: 1}

	fightingLocation := data.FloorMap.FindLocation(func(loc types.Location) bool {
		return slices.Contains(loc.Flags, "arena")
	})

	fight := battle.Fight{
		Entities:       entityMap,
		DiscordChannel: w.DiscordChannel,
		Location:       fightingLocation,
		Meta: &battle.FightMeta{
			ThreadId:   "",
			Tournament: &battle.TournamentData{Tournament: tUuid, Location: w.Tournaments[tUuid].Channel},
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

	tournamentObj.FinishMatch(winner)
}

func (w *World) NextStage(tUuid uuid.UUID) {
	tournamentObj := w.Tournaments[tUuid]

	if tournamentObj == nil {
		return
	}

	tournamentObj.NextStage()
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

func (w *World) Serialize() map[string]any {
	playerData := make(map[uuid.UUID]map[string]any)

	for _, player := range w.Players {
		playerData[player.GetUUID()] = player.Serialize()
	}

	tournamentData := make([]map[string]any, 0)

	for _, tournament := range w.Tournaments {
		tournamentData = append(tournamentData, tournament.Serialize())
	}

	partyData := make(map[uuid.UUID]map[string]any)

	for key, party := range w.Parties {
		partyData[key] = party.Serialize()
	}

	return map[string]any{
		"players":     playerData,
		"parties":     partyData,
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
	_, err := os.Stat(data.Config.BackupLocation)

	if os.IsNotExist(err) {
		os.Mkdir(data.Config.BackupLocation, os.ModePerm)
	}

	newBackupPath := fmt.Sprintf("%s/%v.json", data.Config.BackupLocation, time.Now().Format("2006-01-02_15-04-05"))

	backupFile, err := os.Create(newBackupPath)

	if err != nil {
		panic(err)
	}

	defer backupFile.Close()

	rawData := w.DumpBackup()

	backupFile.Write(rawData)

	w.SendMessage(
		data.Config.LogChannelID,
		discord.MessageCreate{Content: "Backup zrobiony!"},
		false,
	)

	return rawData
}

func (w *World) LoadBackup() error {
	content, err := w.FindNewestBackup()

	if err != nil {
		return err
	}

	if len(content) == 0 {
		return nil
	}

	err = w.LoadBackupData(content)

	if err != nil {
		return err
	}

	return nil
}

func (w *World) FindNewestBackup() ([]byte, error) {
	_, err := os.Stat(data.Config.BackupLocation)

	if os.IsNotExist(err) {
		return []byte{}, nil
	}

	allBackups, err := os.ReadDir(data.Config.BackupLocation)

	if err != nil {
		return []byte{}, err
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

		leftTime, err := time.Parse("2006-01-02_15-04-05", leftName[:len(leftName)-5])

		if err != nil {
			panic(err)
		}

		rightTime, err := time.Parse("2006-01-02_15-04-05", rightName[:len(rightName)-5])

		if err != nil {
			panic(err)
		}

		return leftTime.After(rightTime)
	})

	if len(filteredBackups) == 0 {
		fmt.Println("No backups found")

		return []byte{}, nil
	} else {
		fmt.Println("Loading backup", filteredBackups[0].Name())
	}

	backupContent, err := os.ReadFile("./" + data.Config.BackupLocation + "/" + filteredBackups[0].Name())

	if err != nil {
		return []byte{}, err
	}

	return backupContent, nil
}

func (w *World) LoadBackupData(data []byte) error {
	backupData := make(map[string]any)

	err := json.Unmarshal(data, &backupData)

	if err != nil {
		return err
	}

	w.Players = make(map[uuid.UUID]*player.Player)

	for _, playerData := range backupData["players"].(map[string]any) {
		player := player.Deserialize(playerData.(map[string]any))

		w.Players[player.GetUUID()] = player
	}

	w.Parties = make(map[uuid.UUID]*party.Party)

	for key, partyData := range backupData["parties"].(map[string]any) {
		deserializedParty := party.Deserialize(partyData.(map[string]any))

		for _, member := range deserializedParty.Players {
			w.Players[member.PlayerUuid].Meta.Party = &player.PartialParty{
				UUID:         uuid.MustParse(key),
				Role:         member.Role,
				MembersCount: len(deserializedParty.Players),
			}
		}

		w.Parties[uuid.MustParse(key)] = deserializedParty
	}

	for _, tData := range backupData["tournaments"].([]any) {
		parsedData := tournament.Deserialize(tData.(map[string]any))

		w.Tournaments[parsedData.Uuid] = &parsedData
	}

	return nil
}

func (w *World) RegisterParty(party party.Party) {
	partyUuid := uuid.New()

	for _, member := range party.Players {
		w.Players[member.PlayerUuid].Meta.Party = &player.PartialParty{
			UUID: partyUuid, Role: member.Role, MembersCount: len(party.Players),
		}
	}

	w.Parties[partyUuid] = &party
}
