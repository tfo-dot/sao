package world

import (
	"encoding/json"
	"fmt"
	"os"
	"sao/battle"
	"sao/battle/mobs"
	"sao/player"
	"sao/types"
	"sao/utils"
	"sao/world/calendar"
	"sao/world/location"
	"sao/world/npc"
	"sao/world/party"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/google/uuid"
)

type FloorMap map[string]location.Floor

type World struct {
	Players  map[uuid.UUID]*player.Player
	NPCs     map[uuid.UUID]npc.NPC
	Floors   FloorMap
	Fights   map[uuid.UUID]battle.Fight
	Entities map[uuid.UUID]*battle.Entity
	Time     *calendar.Calendar
	Parties  map[uuid.UUID]*party.Party
	TestMode bool
	DChannel chan types.DiscordMessageStruct
}

func CreateWorld(testMode bool) World {
	floorMap := make(FloorMap, 1)

	pathToRes := "./data/release/world/floors/"

	if testMode {
		pathToRes = "./data/test/world/floors/"
	}

	files, err := os.ReadDir(pathToRes)

	if err != nil {
		panic(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		f, err := os.Open(pathToRes + file.Name())

		if err != nil {
			panic(err)
		}

		var floor location.Floor

		err = json.NewDecoder(f).Decode(&floor)

		if err != nil {
			panic(err)
		}

		floorMap[floor.Name] = floor

		fmt.Printf("Loaded floor %v\n", floor)
	}

	return World{
		make(map[uuid.UUID]*player.Player),
		make(map[uuid.UUID]npc.NPC),
		floorMap,
		make(map[uuid.UUID]battle.Fight),
		make(map[uuid.UUID]*battle.Entity),
		calendar.StartCalendar(),
		make(map[uuid.UUID]*party.Party),
		testMode,
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

	if len(reason) == 0 {
		fmt.Println("No reason for move, it was by player wish")
	}

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

	if len(location.Enemies) == 0 {
		return
	}

	enemyEntry := utils.RandomElement(location.Enemies)
	enemyCount := utils.RandomNumber(enemyEntry.MinNum, enemyEntry.MaxNum)

	entityMap := make(battle.EntityMap)

	if player.Meta.Party != nil {
		partyData := w.Parties[*player.Meta.Party]
		partyMemberCount := len(partyData.Players)

		for role, members := range *partyData.Roles {
			switch role {
			case party.Leader:
				for _, member := range members {
					resolvedMember := w.Players[member]

					if _, exists := entityMap[member]; exists {
						continue
					}

					entityMap[member] = battle.EntityEntry{Entity: resolvedMember, Side: 0}
				}
			case party.DPS:
				for _, member := range members {
					resolvedMember := w.Players[member]

					if _, exists := entityMap[member]; exists {
						continue
					}

					resolvedMember.ApplyEffect(battle.ActionEffect{
						Effect:   battle.EFFECT_STAT_INC,
						Duration: -1,
						Meta: battle.ActionEffectStat{
							Stat:      battle.STAT_ADAPTIVE,
							Value:     10 + (partyMemberCount-1)*5,
							IsPercent: true,
						},
					})

					entityMap[member] = battle.EntityEntry{Entity: resolvedMember, Side: 0}
				}
			case party.Tank:
				for _, member := range members {
					resolvedMember := w.Players[member]

					if _, exists := entityMap[member]; exists {
						continue
					}

					resolvedMember.ApplyEffect(battle.ActionEffect{
						Effect:   battle.EFFECT_STAT_INC,
						Duration: -1,
						Meta: battle.ActionEffectStat{
							Stat:      battle.STAT_DEF,
							Value:     25,
							IsPercent: false,
						},
					})

					resolvedMember.ApplyEffect(battle.ActionEffect{
						Effect:   battle.EFFECT_STAT_INC,
						Duration: -1,
						Meta: battle.ActionEffectStat{
							Stat:      battle.STAT_MR,
							Value:     25,
							IsPercent: false,
						},
					})

					resolvedMember.ApplyEffect(battle.ActionEffect{
						Effect:   battle.EFFECT_STAT_INC,
						Duration: -1,
						Meta: battle.ActionEffectStat{
							Stat:      battle.STAT_HP,
							Value:     (partyMemberCount - 1) * 5,
							IsPercent: true,
						},
					})

					resolvedMember.ApplyEffect(battle.ActionEffect{
						Effect:   battle.EFFECT_STAT_INC,
						Duration: -1,
						Meta: battle.ActionEffectStat{
							Stat:      battle.STAT_DEF,
							Value:     (partyMemberCount - 1) * 5,
							IsPercent: true,
						},
					})

					resolvedMember.ApplyEffect(battle.ActionEffect{
						Effect:   battle.EFFECT_STAT_INC,
						Duration: -1,
						Meta: battle.ActionEffectStat{
							Stat:      battle.STAT_MR,
							Value:     (partyMemberCount - 1) * 5,
							IsPercent: true,
						},
					})

					entityMap[member] = battle.EntityEntry{Entity: resolvedMember, Side: 0}
				}
			case party.Support:
				for _, member := range members {
					resolvedMember := w.Players[member]

					if _, exists := entityMap[member]; exists {
						continue
					}

					resolvedMember.ApplyEffect(battle.ActionEffect{
						Effect:   battle.EFFECT_STAT_INC,
						Duration: -1,
						Meta: battle.ActionEffectStat{
							Stat:      battle.STAT_HEAL_POWER,
							Value:     15 + (partyMemberCount-1)*5,
							IsPercent: true,
						},
					})

					entityMap[member] = battle.EntityEntry{Entity: resolvedMember, Side: 0}
				}
			}
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

func (w *World) StartClock() {
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

	for {
		payload, ok := <-fight.ExternalChannel

		if !ok {
			w.DeregisterFight(fightUuid)
			break
		}

		msgType := payload[0]

		switch battle.FightMessage(msgType) {
		case battle.MSG_FIGHT_END:
			wonSideIDX := fight.Entities.SidesLeft()[0]

			allLoot := make([]battle.Loot, 0)

			for _, entity := range fight.Entities {
				if entity.Side == wonSideIDX {
					continue
				}

				allLoot = append(allLoot, entity.Entity.GetLoot()...)
			}

			wonEntities := fight.Entities.FromSide(wonSideIDX)

			for _, entity := range wonEntities {
				if !entity.IsAuto() {
					entity.(battle.PlayerEntity).ReceiveMultipleLoot(allLoot)
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
				ChannelID: fight.Location.CID,
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
				ChannelID: fight.Location.CID,
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
			entityUuid, err := uuid.FromBytes(payload[1:17])
			player := w.Players[entityUuid]

			if err != nil {
				panic(err)
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
				ChannelID: fight.Location.CID,
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
