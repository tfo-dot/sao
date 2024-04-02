package npc

import (
	"sao/types"
)

type NPC struct {
	Name     string
	Location types.PlayerLocation
	Store    *NPCStore
}

func (n NPC) CanTrade() bool {
	return n.Store != nil
}
