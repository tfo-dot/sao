package data

import (
	"encoding/json"
	"os"
	"sao/config"
	"sao/types"
	"sao/utils"

	"github.com/google/uuid"
)

var Ingredients = GetIngredients()

func GetIngredients() map[uuid.UUID]types.Ingredient {
	dirData, err := os.ReadDir(config.Config.GameDataLocation + "/ingredients")

	if err != nil {
		panic(err)
	}

	var rawIngredients = make([]map[string]interface{}, 0)

	for _, file := range dirData {
		if file.IsDir() {
			continue
		}

		rawData, err := os.ReadFile(config.Config.GameDataLocation + "/ingredients/" + file.Name())

		if err != nil {
			panic(err)
		}

		var parsedJson interface{}

		err = json.Unmarshal(rawData, &parsedJson)

		if err != nil {
			panic(err)
		}

		if data, ok := parsedJson.(map[string]interface{}); ok {
			rawIngredients = append(rawIngredients, data)
		} else {
			for _, ingredient := range parsedJson.([]interface{}) {
				rawIngredients = append(rawIngredients, ingredient.(map[string]interface{}))
			}
		}
	}

	var ingredients = make(map[uuid.UUID]types.Ingredient)

	for _, ingredient := range rawIngredients {

		var UUID = uuid.MustParse(ingredient["UUID"].(string))
		var Name = ingredient["Name"].(string)

		var Stats = make(map[types.Stat]int)

		if _, ok := ingredient["Stats"]; !ok {
			for key, value := range ingredient["Stats"].(map[string]interface{}) {
				Stats[utils.StringToStat[key]] = int(value.(float64))
			}
		}

		ingredients[UUID] = types.Ingredient{
			UUID:  UUID,
			Name:  Name,
			Stats: Stats,
			Count: 1,
		}
	}

	return ingredients
}
