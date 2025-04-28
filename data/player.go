package data

import (
	"fmt"
	"os"
	saoParts "sao/parts"
	"sao/types"
	"strings"

	"github.com/tfo-dot/parts"
)

var PlayerDefaults PlayerDefaultStruct = GetPlayerDefaults()

type PlayerDefaultStruct struct {
	Stats    map[types.Stat]int `parts:"Level,ignoreEmpty"`
	Level    map[types.Stat]int `parts:"Level,ignoreEmpty"`
}

func GetPlayerDefaults() PlayerDefaultStruct {
	tempConfig := PlayerDefaultStruct{
		Stats: make(map[types.Stat]int),
		Level: make(map[types.Stat]int),
	}

	println("Loading player defaults:", Config.GameDataLocation+"/players/default.pts")

	code, err := os.ReadFile(Config.GameDataLocation + "/players/default.pts")

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

	{
		statMap, err := saoParts.FetchVal(vm, "LevelStats")

		if err != nil {
			panic(err)
		}

		for key, value := range statMap.(map[string]any) {
			keyRaw, err := vm.Enviroment.Resolve(fmt.Sprintf("STAT_%s", strings.TrimPrefix(key, "RT")))

			if err != nil {
				panic(err)
			}

			keyVal := types.Stat(keyRaw.Value.(int))

			tempConfig.Level[keyVal] = value.(int)
		}
	}

	{
		statMap, err := saoParts.FetchVal(vm, "StartingStats")

		if err != nil {
			panic(err)
		}

		for key, value := range statMap.(map[string]any) {
			keyRaw, err := vm.Enviroment.Resolve(fmt.Sprintf("STAT_%s", strings.TrimPrefix(key, "RT")))

			if err != nil {
				panic(err)
			}

			keyVal := types.Stat(keyRaw.Value.(int))

			tempConfig.Stats[keyVal] = value.(int)
		}
	}

	parts.ReadFromParts(vm, &tempConfig)

	return tempConfig
}
