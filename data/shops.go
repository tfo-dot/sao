package data

import (
	"encoding/json"
	"os"
	"sao/types"
	"strings"

	"github.com/google/uuid"
)

var Shops = GetShops()

func GetShops() map[uuid.UUID]*types.NPCStore {
	dirData, err := os.ReadDir(Config.GameDataLocation + "/locations/shops")

	if err != nil {
		panic(err)
	}

	shops := map[uuid.UUID]*types.NPCStore{}

	for _, file := range dirData {
		if file.IsDir() {
			continue
		}

		rawData, err := os.ReadFile(Config.GameDataLocation + "/locations/shops/" + file.Name())

		println("Parsing shop:", file.Name())

		if err != nil {
			panic(err)
		}

		var data map[string]any

		err = json.Unmarshal(rawData, &data)

		if err != nil {
			panic(err)
		}

		shopUUID := uuid.MustParse(data["Uuid"].(string))
		name := data["Name"].(string)
		location := strings.Split(data["Location"].(string), ",") //Convert to entity location

		var stocks = make([]*types.WithCount[uuid.UUID], 0)

		for _, item := range data["Stock"].([]any) {
			stock := item.(map[string]any)

			stocks = append(stocks, &types.WithCount[uuid.UUID]{
				Item:  uuid.MustParse(stock["iuuid"].(string)),
				Count: int(stock["Price"].(float64)),
			})
		}

		shops[shopUUID] = &types.NPCStore{
			Uuid:     shopUUID,
			Name:     name,
			Location: FloorMap[location[0]].FindLocation(location[1]),
			Stock:    stocks,
		}
	}

	return shops
}
