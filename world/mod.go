package world

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sao/battle"
	"sao/battle/mobs"
	"sao/player"
	"sao/utils"
	"sao/world/calendar"
	"sao/world/location"
	ClientMessage "sao/world/messages/client"
	ServerMessage "sao/world/messages/server"
	"sao/world/npc"
	"strings"
	"time"

	"github.com/google/uuid"
)

type FloorMap map[string]location.Floor

type World struct {
	Players      map[uuid.UUID]*player.Player
	NPCs         map[uuid.UUID]npc.NPC
	Floors       FloorMap
	Fights       map[uuid.UUID]battle.Fight
	Entities     map[uuid.UUID]battle.Entity
	Time         *calendar.Calendar
	EventChannel chan []byte
	Conn         *net.Conn
}

func CreateWorld() World {
	floorMap := make(FloorMap, 1)

	floorMap["dev"] = location.Floor{
		Name: "dev",
		CID:  "1162450076438900958",
		Locations: []location.Location{
			{
				Name:     "Rynek",
				CID:      "1162450122249076756",
				CityPart: true,
				Enemies:  []mobs.MobType{},
			},
			{
				Name:     "Las",
				CID:      "1162450159251234876",
				CityPart: true,
				Enemies:  []mobs.MobType{mobs.MOB_TIGER},
			},
		},
	}

	return World{
		make(map[uuid.UUID]*player.Player),
		make(map[uuid.UUID]npc.NPC),
		floorMap,
		make(map[uuid.UUID]battle.Fight),
		make(map[uuid.UUID]battle.Entity),
		calendar.StartCalendar(),
		make(chan []byte, 10),
		nil,
	}
}

func (w *World) SetConnection(conn *net.Conn) {
	w.Conn = conn

	go w.ListenToSocket()
	go w.ListenToChannel()
}

func (w *World) ClearConnection() {
	w.Conn = nil
}

func (w *World) ListenToSocket() {
	bytes := make([]byte, 1024)

	for {
		num, err := (*w.Conn).Read(bytes)

		if err != nil {
			w.ClearConnection()
			break
		}

		w.HandlePacket(bytes[:num])
	}
}

func (w *World) HandlePacket(buf []byte) {
	switch ClientMessage.Message(buf[0]) {
	case ClientMessage.PLAYER_NEW:
		payload := ClientMessage.DecodePlayerNew(buf[1:])

		newPlayer := player.NewPlayer(payload.Gender, payload.Name, payload.Uid)

		w.Players[newPlayer.Meta.OwnUUID] = &newPlayer

		w.DispatchSpawn(newPlayer.Meta.OwnUUID)
	case ClientMessage.PLAYER_MOVE:
		payload := ClientMessage.DecodePlayerMove(buf[1:])

		var player player.Player

		for _, p := range w.Players {
			if p.Meta.UserID == payload.Uid {
				player = *p
				break
			}
		}

		w.MovePlayer(player.Meta.OwnUUID, w.Floors[player.Meta.Location.FloorName].FindLocation(payload.CID).Name)
	case ClientMessage.PLAYER_WANDER:
		payload := ClientMessage.DecodePlayerWander(buf[1:])

		var player *player.Player

		for _, p := range w.Players {
			if p.Meta.UserID == payload.Uid {
				player = p
				break
			}
		}

		//Dont care error jumpscare
		if player == nil {
			return
		}

		floor := w.Floors[player.Meta.Location.FloorName]
		location := floor.FindLocation(player.Meta.Location.LocationName)

		if location.CityPart {
			w.DispatchDebug("wander in city part not implemented")
		} else {
			w.PlayerEncounter(player.Meta.OwnUUID)
		}
	}
}

func (w *World) MovePlayer(pUuid uuid.UUID, locationName string) {
	player := w.Players[pUuid]

	player.Meta.Location.LocationName = locationName

	w.EventChannel <- ServerMessage.PlayerMoveSelfPacket(player.Meta.UserID, player.Meta.Location.FloorName, player.Meta.Location.LocationName)
}

func (w *World) ForceMovePlayer(pUuid uuid.UUID, floorName string, locationName, reason string) {
	player := w.Players[pUuid]

	player.Meta.Location.FloorName = floorName
	player.Meta.Location.LocationName = locationName

	w.EventChannel <- ServerMessage.PlayerMoveOtherPacket(player.Meta.UserID, player.Meta.Location.FloorName, player.Meta.Location.LocationName, reason)
}

func (w *World) DispatchSpawn(uuid uuid.UUID) {
	player := w.Players[uuid]

	w.EventChannel <- ServerMessage.PlayerSpawnPacket(player.Meta.UserID, player.Meta.Location.FloorName, player.Meta.Location.LocationName)
}

func (w *World) ListenToChannel() {
	for {
		if w.Conn == nil {
			//Don't handle message until connection is there
			break
		}

		msg, ok := <-w.EventChannel

		if !ok {
			fmt.Printf("Channel closed?")
			w.ClearConnection()
			break
		}

		_, err := (*w.Conn).Write(msg)

		if err != nil {
			fmt.Printf("Error (2): %s", err.Error())
			w.ClearConnection()
			break
		}
	}
}

func (w *World) DispatchDebug(msg string) {
	w.EventChannel <- ServerMessage.DebugPacket(msg)
}

func (w *World) PlayerEncounter(uuid uuid.UUID) {
	player := w.Players[uuid]

	floor := w.Floors[player.Meta.Location.FloorName]
	location := floor.FindLocation(player.Meta.Location.LocationName)

	randomEnemy := utils.RandomElement[mobs.MobType](location.Enemies)

	enemies := mobs.MobEncounter(randomEnemy)

	entityMap := make(battle.EntityMap)

	entityMap[player.GetUUID()] = battle.EntityEntry{Entity: player, Side: 0}

	for _, enemy := range enemies {
		entityMap[enemy.GetUUID()] = battle.EntityEntry{Entity: enemy, Side: 1}
	}

	fight := battle.Fight{Entities: entityMap}

	w.EventChannel <- ServerMessage.FightStartPacket(entityMap)

	fight.Init(w.Time.Copy())

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
					healRatio = 100
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
			w.Entities[entity.Entity.GetUUID()] = entity.Entity
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
		case battle.MSG_ACTION_NEEDED:
			err := uuid.Nil.UnmarshalBinary(payload[1:17])

			if err != nil {
				panic(err)
			}

			rawFUid, _ := fightUuid.MarshalBinary()

			msg := make([]byte, 1+16+16)

			msg[0] = byte(10)
			copy(msg[1:17], payload[1:17])
			copy(msg[17:33], rawFUid)

			w.EventChannel <- msg
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
				} else {
					println("NP character won (loot distribution skipped)")
				}
			}

		default:
			fmt.Printf("Unknown message %d\n", payload[0])
			panic("Not implemented (unknown message)")
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

// Http stuff
func (w *World) HTTPGetTime(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)

	mapData := make(map[string]interface{})

	mapData["day"] = w.Time.Day
	mapData["month"] = w.Time.Month
	mapData["year"] = w.Time.Year
	mapData["hour"] = w.Time.Time.Hour
	mapData["tick"] = w.Time.Time.Tick

	data, err := json.Marshal(mapData)

	if err != nil {
		panic(err)
	}

	res.Write(data)
}

func (w *World) HTTPGetEntity(res http.ResponseWriter, req *http.Request) {
	entityUuid := strings.Split(req.URL.Path, "/")[2]
	uuid, err := uuid.Parse(entityUuid)

	if err != nil {
		res.WriteHeader(400)
		return
	}

	entity, ok := w.Entities[uuid]

	if ok {
		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(200)
		res.Write(SerializeEntity(entity))
	} else {
		entity, ok := w.Players[uuid]

		if ok {
			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(200)
			res.Write(SerializeEntity(entity))
		} else {
			res.WriteHeader(404)
		}
	}
}

func (w *World) HTTPGetPlayer(res http.ResponseWriter, req *http.Request) {
	userID := strings.Split(req.URL.Path, "/")[3]

	for _, player := range w.Players {
		if player.Meta.UserID == userID {
			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(200)
			res.Write(SerializePlayer(*player))
			return
		}
	}

	res.WriteHeader(404)
}

func (w *World) HTTPGetPlayerActions(res http.ResponseWriter, req *http.Request) {
	userID := strings.Split(req.URL.Path, "/")[3]

	for _, player := range w.Players {
		if player.Meta.UserID == userID {
			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(200)
			res.Write(SerializePlayerActions(player.GetAvailableActions()))
			return
		}
	}

	res.WriteHeader(404)
}

func (w *World) HTTPGetFight(res http.ResponseWriter, req *http.Request) {
	fightID := strings.Split(req.URL.Path, "/")[2]

	parsedFightID, err := uuid.Parse(fightID)

	if err != nil {
		res.WriteHeader(400)
		return
	}

	fight, foundFight := w.Fights[parsedFightID]

	if !foundFight {
		res.WriteHeader(404)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)
	res.Write(SerializeFight(fight))
}

func (w *World) HTTPGetStore(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(200)

	mapData := make(map[string]interface{})

	mapData["items"] = []map[string]interface{}{
		{
			"name": "Health Potion",
			"uuid": "1162450076438900958",
			"cost": map[string]interface{}{
				"gold": 10,
			},
		},
	}

	data, err := json.Marshal(mapData)

	if err != nil {
		panic(err)
	}

	res.Write(data)
}

func (w *World) HTTPGetPlayerStore(res http.ResponseWriter, req *http.Request) {
	userID := strings.Split(req.URL.Path, "/")[3]

	for _, player := range w.Players {
		if player.Meta.UserID == userID {
			res.Header().Set("Content-Type", "application/json")
			res.WriteHeader(200)
			// res.Write(SerializePlayerStore(player))
			return
		}
	}

	res.WriteHeader(404)
}

func SerializeEntity(entity battle.Entity) []byte {
	mapData := make(map[string]interface{})

	mapData["isPlayer"] = !entity.IsAuto()
	mapData["name"] = entity.GetName()
	mapData["hp"] = entity.GetCurrentHP()
	mapData["maxHp"] = entity.GetMaxHP()
	mapData["spd"] = entity.GetSPD()
	mapData["atk"] = entity.GetATK()

	data, err := json.Marshal(mapData)

	if err != nil {
		panic(err)
	}

	return data
}

func SerializePlayer(entity player.Player) []byte {
	mapData := make(map[string]interface{})

	mapData["name"] = entity.Name
	mapData["uuid"] = entity.Meta.OwnUUID.String()
	mapData["location"] = map[string]interface{}{
		"floor":    entity.Meta.Location.FloorName,
		"location": entity.Meta.Location.LocationName,
	}

	if entity.Meta.FightInstance != nil {
		mapData["fight"] = entity.Meta.FightInstance.String()
	} else {
		mapData["fight"] = nil
	}

	data, err := json.Marshal(mapData)

	if err != nil {
		panic(err)
	}

	return data
}

func SerializePlayerActions(actions []battle.ActionPartial) []byte {
	actionList := make([]map[string]interface{}, len(actions))

	for i, action := range actions {
		tempAction := make(map[string]interface{})

		tempAction["event"] = action.Event
		if action.Meta != nil {
			tempAction["meta"] = action.Meta.String()
		}

		actionList[i] = tempAction
	}

	data, err := json.Marshal(actionList)

	if err != nil {
		panic(err)
	}

	return data
}

func SerializeFight(fight battle.Fight) []byte {
	mapData := make(map[string]interface{})

	entityList := make([]map[string]interface{}, len(fight.Entities))

	for _, entity := range fight.Entities {
		tempEntity := make(map[string]interface{})

		tempEntity["uuid"] = entity.Entity.GetUUID().String()
		tempEntity["side"] = entity.Side
		tempEntity["isPlayer"] = !entity.Entity.IsAuto()

		entityList = append(entityList, tempEntity)
	}

	mapData["entities"] = entityList
	mapData["finished"] = fight.IsFinished()

	data, err := json.Marshal(mapData)

	if err != nil {
		panic(err)
	}

	return data
}
