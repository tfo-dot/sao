package location

import (
	"encoding/json"
	"os"
	"sao/config"
)

type Location struct {
	Name     string
	CID      string
	CityPart bool
	Effects  []LocationEffect
	TP       bool
	Enemies  []EnemyMeta
	Unlocked bool
	Flags    []string
}

type EnemyMeta struct {
	MinNum int
	MaxNum int
	Enemy  string
}

type Floor struct {
	Name             string
	CID              string
	Default          string
	Locations        []Location
	Effects          []LocationEffect
	Flags            []string
	Unlocked         bool
	CountsAsUnlocked bool
}

type LocationEffect struct {
	Effect int
	Value  int
	Meta   *map[string]interface{}
}

func (f Floor) FindLocation(str string) *Location {
	for _, loc := range f.Locations {
		if loc.CID == str || loc.Name == str {
			return &loc
		}
	}

	return nil
}

var Floors = GetFloors()

func GetFloors() map[string]Floor {
	dirData, err := os.ReadDir(config.Config.GameDataLocation + "/locations/floors")

	if err != nil {
		panic(err)
	}

	var rawFloors = make([]map[string]interface{}, 0)

	for _, file := range dirData {
		if file.IsDir() {
			continue
		}

		rawData, err := os.ReadFile(config.Config.GameDataLocation + "/locations/floors/" + file.Name())

		if err != nil {
			panic(err)
		}

		var parsedJson interface{}

		err = json.Unmarshal(rawData, &parsedJson)

		if err != nil {
			panic(err)
		}

		if data, ok := parsedJson.(map[string]interface{}); ok {
			rawFloors = append(rawFloors, data)
		}
	}

	var floors = make(map[string]Floor)

	for _, floor := range rawFloors {
		Name := floor["Name"].(string)
		CID := floor["CID"].(string)
		Default := floor["Default"].(string)
		Unlocked := floor["Unlocked"].(bool)
		CountsAsUnlocked := floor["CountsAsUnlocked"].(bool)

		floorFlags := floor["Flags"].([]interface{})

		var flags = make([]string, len(floorFlags))

		for i, flag := range floorFlags {
			flags[i] = flag.(string)
		}

		var Locations = make([]Location, 0)

		for _, loc := range floor["Locations"].([]interface{}) {
			lName := loc.(map[string]interface{})["Name"].(string)
			lCID := loc.(map[string]interface{})["CID"].(string)
			lCityPart := loc.(map[string]interface{})["CityPart"].(bool)
			lTP := loc.(map[string]interface{})["TP"].(bool)
			lUnlocked := loc.(map[string]interface{})["Unlocked"].(bool)
			lFlags := loc.(map[string]interface{})["Flags"].([]interface{})

			var Effects = make([]LocationEffect, 0)

			for _, eff := range loc.(map[string]interface{})["Effects"].([]interface{}) {
				e := eff.(map[string]interface{})
				Effects = append(Effects, LocationEffect{
					Effect: int(e["Effect"].(float64)),
					Value:  int(e["Value"].(float64)),
					Meta:   nil,
				})
			}

			var Enemies = make([]EnemyMeta, 0)

			for _, en := range loc.(map[string]interface{})["Enemies"].([]interface{}) {
				e := en.(map[string]interface{})
				Enemies = append(Enemies, EnemyMeta{
					MinNum: int(e["MinNum"].(float64)),
					MaxNum: int(e["MaxNum"].(float64)),
					Enemy:  e["Enemy"].(string),
				})
			}

			var locationFlags = make([]string, len(lFlags))

			for i, flag := range lFlags {
				locationFlags[i] = flag.(string)
			}

			Locations = append(Locations, Location{
				Name:     lName,
				CID:      lCID,
				CityPart: lCityPart,
				Effects:  Effects,
				TP:       lTP,
				Enemies:  Enemies,
				Unlocked: lUnlocked,
				Flags:    locationFlags,
			})
		}

		var Effects = make([]LocationEffect, 0)

		for _, eff := range floor["Effects"].([]interface{}) {
			e := eff.(map[string]interface{})
			Effects = append(Effects, LocationEffect{
				Effect: int(e["Effect"].(float64)),
				Value:  int(e["Value"].(float64)),
				Meta:   nil,
			})
		}

		floors[Name] = Floor{
			Name:             Name,
			CID:              CID,
			Default:          Default,
			Locations:        Locations,
			Effects:          Effects,
			Unlocked:         Unlocked,
			CountsAsUnlocked: CountsAsUnlocked,
			Flags:            flags,
		}

	}
	return floors
}
