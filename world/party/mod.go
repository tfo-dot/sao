package party

import "github.com/google/uuid"

type Party struct {
	Players []uuid.UUID
	//Some hacker array since we know there won't be more than 1 leader XD
	Roles *map[PartyRole][]uuid.UUID
}

type PartyRole int

const (
	Leader PartyRole = iota
	DPS
	Support
	Tank
)
