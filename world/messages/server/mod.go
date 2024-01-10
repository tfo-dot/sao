package server

import (
	"sao/battle"
	"sao/player"
	"sao/utils"
)

type Message byte

const (
	PLAYER_SPAWN Message = iota
	PLAYER_MOVE_OTHER
	PLAYER_MOVE_SELF
	PLAYER_DEATH
	FIGHT_START
	FIGHT_END
	DEBUG Message = 255
)

func PlayerSpawnPacket(uid, floorName, locationName string) []byte {
	eventData := make([]byte, 1+len(uid)+1+len(floorName)+1+len(locationName)+1)

	eventData[0] = byte(PLAYER_SPAWN)

	offset := utils.WriteStringWithOffset(eventData, 1, uid)
	offset = utils.WriteStringWithOffset(eventData, offset, floorName)
	utils.WriteStringWithOffset(eventData, offset, locationName)

	return eventData
}

func DebugPacket(msg string) []byte {
	msgLen := len(msg)

	if msgLen > 255 {
		panic("Not implemented (debug message too long)")
	}

	rawData := make([]byte, 2+msgLen)
	rawData[0] = byte(DEBUG)
	utils.WriteStringWithOffset(rawData, 1, msg)

	return rawData
}

func PlayerMoveOtherPacket(uid, floorName, locationName, reason string) []byte {
	eventData := make([]byte, 1+len(uid)+1+len(floorName)+1+len(locationName)+1+len(reason)+1)

	eventData[0] = byte(PLAYER_MOVE_OTHER)

	offset := utils.WriteStringWithOffset(eventData, 1, uid)
	offset = utils.WriteStringWithOffset(eventData, offset, floorName)
	offset = utils.WriteStringWithOffset(eventData, offset, locationName)
	utils.WriteStringWithOffset(eventData, offset, reason)

	return eventData
}

func PlayerMoveSelfPacket(uid, floorName, locationName string) []byte {
	eventData := make([]byte, 1+len(uid)+1+len(floorName)+1+len(locationName)+1)

	eventData[0] = byte(PLAYER_MOVE_SELF)

	offset := utils.WriteStringWithOffset(eventData, 1, uid)
	offset = utils.WriteStringWithOffset(eventData, offset, floorName)
	utils.WriteStringWithOffset(eventData, offset, locationName)

	return eventData
}

func FightStartPacket(entities battle.EntityMap) []byte {
	rawData := make([]byte, 2+len(entities)*(1+36))

	rawData[0] = byte(FIGHT_START)
	rawData[1] = byte(len(entities))

	offset := 2

	for _, entityEntry := range entities {
		rawData[offset] = byte(entityEntry.Side)
		offset++
		if !entityEntry.Entity.IsAuto() {
			offset = utils.WriteStringWithOffset(rawData, offset, entityEntry.Entity.(*player.Player).Meta.UserID)
		} else {
			offset = utils.WriteStringWithOffset(rawData, offset, entityEntry.Entity.GetUUID().String())
		}
	}

	return rawData
}
