package server

import (
	"sao/utils"

	"github.com/google/uuid"
)

type Message byte

const (
	PLAYER_SPAWN Message = iota
	PLAYER_MOVE_OTHER
	PLAYER_MOVE_SELF
	PLAYER_DEATH
	FIGHT_START
	FIGHT_PACKET
	FIGHT_END
	DEBUG Message = 255
)

type FightMessage byte

const (
	FIGHT_ACTION_NEEDED FightMessage = iota
	FIGHT_ENTITY_DIED
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

func FightStartPacket(fUuid uuid.UUID) []byte {
	rawData := make([]byte, 1+36)

	rawData[0] = byte(FIGHT_START)
	copy(rawData[1:], fUuid.String()[:])

	return rawData
}

func ActionNeededPacket(fUuid uuid.UUID, playerUuid uuid.UUID) []byte {
	rawData := make([]byte, 2+36+36)

	rawData[0] = byte(FIGHT_PACKET)
	rawData[1] = byte(FIGHT_ACTION_NEEDED)

	copy(rawData[2:], fUuid.String()[:])
	copy(rawData[2+36:], playerUuid.String()[:])

	return rawData
}
