package data

import (
	"encoding/json"
	"os"
	"sao/config"
	"sao/types"

	"github.com/google/uuid"
)

var Recipes = GetRecipes()

func GetRecipes() map[uuid.UUID]types.Recipe {
	dirData, err := os.ReadDir(config.Config.GameDataLocation + "/recipes")

	if err != nil {
		panic(err)
	}

	var rawRecipes = make([]map[string]interface{}, 0)

	for _, file := range dirData {
		if file.IsDir() {
			continue
		}

		println("Loading recipe: " + file.Name())

		rawData, err := os.ReadFile(config.Config.GameDataLocation + "/recipes/" + file.Name())

		if err != nil {
			panic(err)
		}

		var parsedJson interface{}

		err = json.Unmarshal(rawData, &parsedJson)

		if err != nil {
			panic(err)
		}

		if data, ok := parsedJson.(map[string]interface{}); ok {
			rawRecipes = append(rawRecipes, data)
		} else {
			for _, ingredient := range parsedJson.([]interface{}) {
				rawRecipes = append(rawRecipes, ingredient.(map[string]interface{}))
			}
		}
	}

	var recipes = make(map[uuid.UUID]types.Recipe)

	for _, recipe := range rawRecipes {
		var UUID = uuid.MustParse(recipe["UUID"].(string))
		var Name = recipe["Name"].(string)

		var Ingredients = make([]types.WithCount[uuid.UUID], 0)

		for _, value := range recipe["Ingredients"].([]interface{}) {
			casted := value.(map[string]interface{})

			Ingredients = append(Ingredients, types.WithCount[uuid.UUID]{
				Item:  uuid.MustParse(casted["Item"].(string)),
				Count: int(casted["Count"].(float64)),
			})
		}

		var Cost = int(recipe["Cost"].(float64))

		var Product = types.ResultItem{
			UUID:  uuid.MustParse(recipe["Product"].(map[string]interface{})["UUID"].(string)),
			Type:  StringToType[recipe["Product"].(map[string]interface{})["Type"].(string)],
			Count: int(recipe["Product"].(map[string]interface{})["Count"].(float64)),
		}

		recipes[UUID] = types.Recipe{
			UUID:        UUID,
			Name:        Name,
			Ingredients: Ingredients,
			Cost:        Cost,
			Product:     Product,
		}
	}

	return recipes
}

var StringToType = map[string]types.ItemType{
	"Other":    types.ITEM_OTHER,
	"Material": types.ITEM_MATERIAL,
}
