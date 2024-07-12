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

func (p *Party) Serialize() map[string]interface{} {
	members := make([]map[string]interface{}, 0)

	for _, player := range p.Players {
		members = append(members, map[string]interface{}{
			"player": player.PlayerUuid.String(),
			"role":   player.Role,
		})
	}

	return map[string]interface{}{
		"players": members,
		"leader":  p.Leader.String(),
	}
}

func Deserialize(data map[string]interface{}) *Party {
	party := &Party{
		Leader: uuid.MustParse(data["leader"].(string)),
	}

	for _, player := range data["players"].([]map[string]interface{}) {
		party.Players = append(party.Players, &PartyEntry{
			PlayerUuid: uuid.MustParse(player["player"].(string)),
			Role:       PartyRole(player["role"].(float64)),
		})
	}

	return party
}
