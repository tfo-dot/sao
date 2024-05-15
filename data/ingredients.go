package data

import (
	"sao/types"

	"github.com/google/uuid"
)

var Ingredients = map[uuid.UUID]types.Ingredient{
	uuid.MustParse("00000000-0000-0000-0000-000000000000"): {UUID: uuid.MustParse("00000000-0000-0000-0000-000000000000"), Name: "Tlen", Count: 0},
}
