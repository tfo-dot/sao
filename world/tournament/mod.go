package tournament

import (
	"github.com/google/uuid"
)

type Tournament struct {
	Uuid uuid.UUID
	Name string
	//-1 for unlimited
	MaxPlayers      int
	Channel         string
	Participants    []uuid.UUID
	State           TournamentState
	Stages          []*TournamentStage
	ExternalChannel chan TournamentEventData
}

type TournamentStage struct {
	Matches []*TournamentMatch
	IDX     int
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

type TournamentState int

const (
	Waiting TournamentState = iota
	Running
	Finished
)

func (t *Tournament) NextStage() {
	if t.State != Running {
		return
	}

	stage := t.Stages[len(t.Stages)-1]

	for _, match := range stage.Matches {
		if match.State != FinishedMatch {
			return
		}
	}

	if len(stage.Matches) == 1 {
		return
	}

	newMatches := make([]*TournamentMatch, 0)

	for i := 0; i < len(stage.Matches); i += 2 {
		newMatches = append(newMatches, &TournamentMatch{
			Players: []uuid.UUID{*stage.Matches[i].Winner, *stage.Matches[i+1].Winner},
			Winner:  nil,
			State:   BeforeMatch,
		})
	}

	t.Stages = append(t.Stages, &TournamentStage{Matches: newMatches})
}

func (t *Tournament) FinishMatch(winner uuid.UUID) {
	if t.State != Running {
		return
	}

	stage := t.Stages[len(t.Stages)-1]

	for _, match := range stage.Matches {
		for _, player := range match.Players {
			if player == winner {
				match.Winner = &winner

				match.State = FinishedMatch
			}
		}
	}

	if len(stage.Matches) == 1 {
		t.State = Finished
	}
}

func (t *Tournament) Serialize() map[string]interface{} {

	tStages := make([]map[string]interface{}, 0)

	for _, stage := range t.Stages {
		tStages = append(tStages, stage.Serialize())
	}

	return map[string]interface{}{
		"uuid":         t.Uuid,
		"name":         t.Name,
		"max_players":  t.MaxPlayers,
		"participants": t.Participants,
		"state":        t.State,
		"stages":       tStages,
	}
}

func (ts *TournamentStage) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"matches": ts.Matches,
		"idx":     ts.IDX,
	}
}

func Deserialize(rawData map[string]interface{}) Tournament {
	rawParticipants := rawData["participants"].([]interface{})

	participants := make([]uuid.UUID, 0)

	for _, participant := range rawParticipants {
		participants = append(participants, uuid.MustParse(participant.(string)))
	}

	t := Tournament{
		Uuid:         uuid.MustParse(rawData["uuid"].(string)),
		Name:         rawData["name"].(string),
		MaxPlayers:   int(rawData["max_players"].(float64)),
		Participants: participants,
		State:        TournamentState(rawData["state"].(float64)),
	}

	tStages := rawData["stages"].([]interface{})

	for _, stage := range tStages {
		t.Stages = append(t.Stages, DeserializeStage(stage.(map[string]interface{})))
	}

	return t
}

func DeserializeStage(rawData map[string]interface{}) *TournamentStage {
	ts := TournamentStage{
		IDX:     rawData["idx"].(int),
		Matches: make([]*TournamentMatch, 0),
	}

	matches := rawData["matches"].([]map[string]interface{})

	for _, match := range matches {
		parsedPlayers := make([]uuid.UUID, 0)

		for _, player := range match["players"].([]string) {
			parsedPlayers = append(parsedPlayers, uuid.MustParse(player))
		}

		var winner *uuid.UUID

		if match["winner"] != nil {
			tempUuid := uuid.MustParse(match["winner"].(string))
			winner = &tempUuid
		} else {
			winner = nil
		}

		ts.Matches = append(ts.Matches, &TournamentMatch{
			Players: parsedPlayers,
			Winner:  winner,
			State:   match["state"].(MatchState),
		})
	}

	return &ts
}
