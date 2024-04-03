package npc

import (
	"sao/types"
	"sao/world/calendar"

	"github.com/google/uuid"
)

type NPC struct {
	Name     string
	Location types.PlayerLocation
	Store    *NPCStore
}

type NPCStore struct {
	Uuid            uuid.UUID
	Name            string
	RestockInterval calendar.Calendar
	LastRestock     calendar.Calendar
	Stock           []*Stock
}

type Stock struct {
	ItemType types.ItemType
	ItemUUID uuid.UUID
	Price    int
	Quantity int
	Limit    int
}
