package data

import (
	"sao/types"

	"github.com/google/uuid"
)

var Ingredients = map[uuid.UUID]types.Ingredient{
	BloodyEdgeUUID: {
		UUID: BloodyEdgeUUID,
		Name: "Krwiste ostrze",
		Stats: map[types.Stat]int{
			types.STAT_LIFESTEAL: 10,
			types.STAT_AD:        10,
		},
		Count: 1,
	},
	GoodlyElementUUID: {
		UUID:  GoodlyElementUUID,
		Name:  "Niebiański pierwiastek",
		Stats: nil,
		Count: 1,
	},
}

var BloodyEdge = Ingredients[uuid.MustParse("00000000-0000-0000-0000-000000000000")]
var BloodyEdgeUUID = uuid.MustParse("00000000-0000-0000-0000-000000000000")

var GoodlyElement = Ingredients[uuid.MustParse("00000000-0000-0000-0000-000000000001")]
var GoodlyElementUUID = uuid.MustParse("00000000-0000-0000-0000-000000000001")
