package client

import (
	"sao/player"
	"sao/utils"
)

type Message byte

const (
	PLAYER_NEW Message = iota
	PLAYER_MOVE
	PLAYER_WANDER
)

type PayloadPlayerNew struct {
	Gender player.PlayerGender
	Uid    string
	Name   string
}

func DecodePlayerNew(buf []byte) PayloadPlayerNew {
	offset, uid := utils.ReadStringWithOffset(1, buf)
	_, name := utils.ReadStringWithOffset(offset, buf)

	return PayloadPlayerNew{Gender: player.PlayerGender(buf[0]), Uid: uid, Name: name}
}

type PayloadPlayerMove struct {
	Uid string
	CID string
}

func DecodePlayerMove(buf []byte) PayloadPlayerMove {
	offset, uid := utils.ReadStringWithOffset(0, buf)
	_, cid := utils.ReadStringWithOffset(offset, buf)

	return PayloadPlayerMove{Uid: uid, CID: cid}
}

type PayloadPlayerWander struct {
	Uid string
}

func DecodePlayerWander(buf []byte) PayloadPlayerWander {
	_, uid := utils.ReadStringWithOffset(0, buf)

	return PayloadPlayerWander{Uid: uid}
}
