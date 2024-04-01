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
