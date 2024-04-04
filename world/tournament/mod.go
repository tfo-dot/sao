package tournament

import "github.com/google/uuid"

type Tournament struct {
	Uuid uuid.UUID
	Name string
	Type TournamentType
	//-1 for unlimited
	MaxPlayers   int
	Participants []uuid.UUID
	State        TournamentState
	Stages       []*TournamentStage
}

type TournamentStage struct {
	Matches []*TournamentMatch
}

type TournamentMatch struct {
	Players []uuid.UUID
	Winner  *uuid.UUID
	State   MatchState
}

type MatchState int

const (
	BeforeMatch MatchState = iota
	RunningMatch
	FinishedMatch
)

type TournamentType int

const (
	SingleElimination TournamentType = iota
)

type TournamentState int

const (
	Waiting TournamentState = iota
	Running
	Finished
)

//TODO hopefully someday finish this
