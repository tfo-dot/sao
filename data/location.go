package data

import (
	"os"
	saoParts "sao/parts"
	"sao/types"
	"strings"

	"github.com/tfo-dot/parts"
)

type Floors map[string]types.Floor

var FloorMap = GetFloors()

func GetFloors() Floors {
	dirData, err := os.ReadDir(Config.GameDataLocation + "/locations/floors")

	if err != nil {
		panic(err)
	}

	var floors = make(map[string]types.Floor)

	for _, file := range dirData {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".pts") {
			continue
		}

		var floorInfo types.Floor

		println("Loading floor data:", Config.GameDataLocation+"/locations/floors/"+file.Name())

		code, err := os.ReadFile(Config.GameDataLocation + "/locations/floors/" + file.Name())

		if err != nil {
			panic(err)
		}

		vm, err := parts.GetVMWithSource(string(code))

		if err != nil {
			panic(err)
		}

		saoParts.AddConsts(vm)
		saoParts.AddFunctions(vm)

		err = vm.Run()

		if err != nil {
			panic(err)
		}

		parts.ReadFromParts(vm, &floorInfo)

		res, err := saoParts.FetchVal(vm, "Locations")

		if err != nil {
			panic(err)
		}

		for _, loc := range res.([]any) {
			locData := loc.(map[string]any)

			mobs := make([]types.EnemyMeta, 0)

			if val, has := locData["RTEnemies"]; has {
				for _, mob := range val.([]any) {
					mobData := mob.(map[string]any)

					mobs = append(mobs, types.EnemyMeta{
						MinNum: mobData["RTMinNum"].(int),
						MaxNum: mobData["RTMaxNum"].(int),
						Enemy:  mobData["RTEnemy"].(string),
					})
				}
			}

			floorInfo.Locations = append(floorInfo.Locations, types.Location{
				Name:     locData["RTName"].(string),
				CID:      locData["RTCID"].(string),
				CityPart: locData["RTCityPart"].(bool),
				TP:       locData["RTTP"].(bool),
				Unlocked: locData["RTUnlocked"].(bool),
				Enemies:  mobs,
			})
		}

		floors[floorInfo.Name] = floorInfo
	}

	return floors
}

func (f Floors) FindLocation(check func(types.Location) bool) *types.Location {
	for _, flor := range f {
		for _, loc := range flor.Locations {
			if check(loc) {
				return &loc
			}
		}
	}

	return nil
}

func (f Floors) GetUnlockedFloorCount() int {
	unlockedFloors := 0

	for _, floor := range f {
		if floor.Unlocked && floor.CountsAsUnlocked {
			unlockedFloors++
		}
	}

	return unlockedFloors
}
