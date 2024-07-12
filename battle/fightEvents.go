package battle

import "github.com/google/uuid"

type FightEvent interface {
	GetEvent() FightMessage
	GetData() any
}

type FightStartMsg struct{}

func (fsm FightStartMsg) GetEvent() FightMessage {
	return MSG_FIGHT_START
}

func (fsm FightStartMsg) GetData() any {
	return nil
}

type FightEndMsg struct{}

func (fsm FightEndMsg) GetEvent() FightMessage {
	return MSG_FIGHT_END
}

func (fsm FightEndMsg) GetData() any {
	return nil
}

type FightActionNeededMsg struct {
	Entity uuid.UUID
}

func (fsm FightActionNeededMsg) GetEvent() FightMessage {
	return MSG_ACTION_NEEDED
}

func (fsm FightActionNeededMsg) GetData() any {
	return fsm.Entity
}
