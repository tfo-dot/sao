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
	MSG_SUMMON_EXPIRED
	MSG_SUMMON_DIED
	MSG_ENTITY_DIED
)

type EventHandler struct {
	Target  uuid.UUID
	Handler func(source, target types.Entity, fightInstance types.FightInstance, meta any) any
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

type EntityMap map[uuid.UUID]*EntityEntry

type EntityEntry struct {
	Entity types.Entity
	Side   int
	Speed  int
	Turn   int
}
