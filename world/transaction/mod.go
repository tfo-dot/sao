package transaction

import (
	"sao/types"

	"github.com/google/uuid"
)

type Transaction struct {
	Uuid      uuid.UUID
	LeftSide  *TransactionSide
	RightSide *TransactionSide
	State     TransactionState
}

type TransactionSide struct {
	Who   uuid.UUID
	With  []TransactionEntry
	Gold  int
	State TransactionSideState
}

type TransactionEntry struct {
	Item     uuid.UUID
	ItemType types.ItemType
	Amount   int
}

type TransactionState int

const (
	TransactionPending TransactionState = iota
	TransactionProgress
)

type TransactionSideState int

const (
	TransactionSideAccept TransactionSideState = iota
	TransactionSideDecline
)
