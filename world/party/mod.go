package party

import "github.com/google/uuid"

type Party struct {
	Players []*PartyEntry
}

type PartyEntry struct {
	PlayerUuid uuid.UUID
	Role       PartyRole
}

type PartyRole int

const (
	Leader PartyRole = iota
	DPS
	Support
	Tank
	None
)
