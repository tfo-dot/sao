package data

import (
	"encoding/json"
	"os"
	"sao/config"
	"sao/types"
	"strings"

	"github.com/google/uuid"
)

var Shops = GetShops()

func GetShops() map[uuid.UUID]*types.NPCStore {
	dirData, err := os.ReadDir(config.Config.GameDataLocation + "/locations/shops")

	if err != nil {
		panic(err)
	}

	rawShops := []map[string]interface{}{}

	for _, file := range dirData {
		if file.IsDir() {
			continue
		}

		rawData, err := os.ReadFile(config.Config.GameDataLocation + "/locations/shops/" + file.Name())

		println("Parsing shop:", file.Name())

		if err != nil {
			panic(err)
		}

		var parsedJson interface{}

		err = json.Unmarshal(rawData, &parsedJson)

		if err != nil {
			panic(err)
		}

		if data, ok := parsedJson.(map[string]interface{}); ok {
			rawShops = append(rawShops, data)
		}
	}

	shops := map[uuid.UUID]*types.NPCStore{}

	for _, shop := range rawShops {
		shopUUID := uuid.MustParse(shop["Uuid"].(string))
		name := shop["Name"].(string)
		location := shop["Location"].(string) //Convert to entity location

		var stocks = make([]*types.Stock, 0)

		for _, item := range shop["Stock"].([]interface{}) {
			stock := item.(map[string]interface{})

			Item := types.ItemType(stock["Item"].(float64))
			Price := int(stock["Price"].(float64))

			stocks = append(stocks, &types.Stock{
				ItemType: Item,
				ItemUUID: uuid.MustParse(stock["iuuid"].(string)),
				Price:    Price,
			})
		}

		shops[shopUUID] = &types.NPCStore{
			Uuid: shopUUID,
			Name: name,
			Location: types.EntityLocation{
				Location: strings.Split(location, ",")[1],
				Floor:    strings.Split(location, ",")[0],
			},
			Stock: stocks,
		}
	}

	return shops
}
