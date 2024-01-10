package npc

import "github.com/google/uuid"

type NPCMeta struct {
	canFight bool
}

type NPC struct {
	Name     string
	Location uuid.UUID
	Meta     NPCMeta
	Store    *NPCStore
}

func (n *NPC) CanFight() bool {
	return n.Meta.canFight
}

func (n NPC) CanTrade() bool {
	return n.Store != nil
}
