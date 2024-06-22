package data

import (
	"sao/types"

	"github.com/google/uuid"
)

var Recipes = map[uuid.UUID]types.Recipe{
	uuid.MustParse("00000000-0000-0000-0000-000000000000"): {
		UUID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Name: "Błogosławieństwo Reimi",
		Ingredients: []types.WithCount[uuid.UUID]{
			{Item: BloodyEdgeUUID, Count: 1},
			{Item: GoodlyElementUUID, Count: 1},
		},
		Cost: 150,
		Product: types.ResultItem{
			UUID:  ReimiBlessing.UUID,
			Type:  types.ITEM_OTHER,
			Count: 1,
		},
	},
}

var ReimiBlessingRecipe = Recipes[uuid.MustParse("00000000-0000-0000-0000-000000000000")]
var ReimiBlessingRecipeUUID = uuid.MustParse("00000000-0000-0000-0000-000000000000")
