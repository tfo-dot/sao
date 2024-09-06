package battle

import (
	"sao/types"

	"github.com/google/uuid"
)

type FightMessage byte

const (
	MSG_ACTION_NEEDED FightMessage = iota
	MSG_FIGHT_START
	MSG_FIGHT_END
	MSG_ENTITY_RESCUE
	MSG_SUMMON_EXPIRED
	MSG_ENTITY_DIED
)

const SPEED_GAUGE = 100

type EventHandler struct {
	Target  uuid.UUID
	Handler func(source, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{}
	Trigger types.SkillTrigger
}

type FightMeta struct {
	ThreadId   string
	Tournament *TournamentData
}

type TournamentData struct {
	Tournament uuid.UUID
	Location   string
}

type EntityMap map[uuid.UUID]EntityEntry

type EntityEntry struct {
	Entity types.Entity
	Side   int
}
