package tournament

import "github.com/google/uuid"

type TournamentEvent int

const (
	MatchFinished TournamentEvent = iota
)

type TournamentEventData interface {
	GetEvent() TournamentEvent
	GetData() interface{}
}

type MatchFinishedData struct {
	Winner uuid.UUID
}

func (mfd MatchFinishedData) GetEvent() TournamentEvent {
	return MatchFinished
}

func (mfd MatchFinishedData) GetData() interface{} {
	return mfd.Winner
}
