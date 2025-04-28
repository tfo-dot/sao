package party

import "github.com/google/uuid"

type Party struct {
	Players []*PartyEntry
	Leader  uuid.UUID
}

type PartyEntry struct {
	PlayerUuid uuid.UUID
	Role       PartyRole
}

type PartyRole int

const (
	DPS PartyRole = iota
	Support
	Tank
	None
)

func (p *Party) Serialize() map[string]any {
	members := make([]map[string]any, 0)

	for _, player := range p.Players {
		members = append(members, map[string]any{"player": player.PlayerUuid.String(), "role": player.Role})
	}

	return map[string]any{"players": members, "leader": p.Leader.String()}
}

func Deserialize(data map[string]any) *Party {
	party := &Party{
		Leader: uuid.MustParse(data["leader"].(string)),
	}

	for _, player := range data["players"].([]any) {

		plr := player.(map[string]any)

		party.Players = append(party.Players, &PartyEntry{
			PlayerUuid: uuid.MustParse(plr["player"].(string)),
			Role:       PartyRole(plr["role"].(float64)),
		})
	}

	return party
}
