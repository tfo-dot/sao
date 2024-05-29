package tournament

import "github.com/google/uuid"

type Tournament struct {
	Uuid uuid.UUID
	Name string
	Type TournamentType
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

func (t *Tournament) Serialize() map[string]interface{} {

	tStages := make([]map[string]interface{}, 0)

	for _, stage := range t.Stages {
		tStages = append(tStages, stage.Serialize())
	}

	return map[string]interface{}{
		"uuid":         t.Uuid,
		"name":         t.Name,
		"type":         t.Type,
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
		Type:         TournamentType(rawData["type"].(float64)),
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
