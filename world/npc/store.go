package npc

import "sao/world/calendar"

type NPCStore struct {
	RestockInterval calendar.Calendar
	LastRestock     calendar.Calendar
}
