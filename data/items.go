package data

import (
	"os"
	"sao/config"
	"sao/types"
	"sao/utils"

	"github.com/Shopify/go-lua"
	"github.com/google/uuid"
)

var Items = GetItems()

func GetItems() map[uuid.UUID]types.PlayerItem {
	dirData, err := os.ReadDir(config.Config.GameDataLocation + "/items")

	if err != nil {
		panic(err)
	}

	items := map[uuid.UUID]types.PlayerItem{}

	for _, file := range dirData {
		if file.IsDir() {
			continue
		}

		state := lua.NewState()

		lua.OpenLibraries(state)

		state.NewTable()

		state.SetGlobal("Effects")

		state.NewTable()

		state.SetGlobal("ReservedUIDs")

		state.NewTable()

		state.PushGoFunction(func(state *lua.State) int {
			value := lua.CheckInteger(state, 1)
			percent := lua.CheckInteger(state, 2)

			state.PushInteger(utils.PercentOf(value, percent))

			return 1
		})

		state.SetField(-2, "percentOf")

		state.SetGlobal("utils")

		err := lua.DoFile(state, config.Config.GameDataLocation+"/items/"+file.Name())

		if err != nil {
			panic(err)
		}

		item := types.PlayerItem{
			UUID:        uuid.MustParse(GetLuaString(state, "UUID")),
			Name:        GetLuaString(state, "Name"),
			Description: GetLuaString(state, "Description"),
			TakesSlot:   GetLuaBool(state, "TakesSlot"),
			Stacks:      GetLuaBool(state, "Stacks"),
			Consume:     GetLuaBool(state, "Consume"),
			Count:       int(GetLuaFloat(state, "Count")),
			MaxCount:    int(GetLuaFloat(state, "MaxCount")),
			Hidden:      GetLuaBool(state, "Hidden"),
			Stats:       map[types.Stat]int{},
			Effects:     []types.PlayerSkill{},
		}

		state.Global("Stats")

		state.PushNil()

		for state.Next(-2) {
			key, _ := state.ToString(-2)
			value, _ := state.ToNumber(-1)

			item.Stats[utils.StringToStat[key]] = int(value)

			state.Pop(1)
		}

		state.Pop(1)

		state.Global("Effects")

		state.RawGetInt(-1, 0)

		state.PushNil()

		for state.Next(-2) {
			key, _ := state.ToString(-2)

			println(key)

			state.Pop(1)
		}

		items[item.UUID] = item
	}

	return items
}

func GetLuaString(state *lua.State, str string) string {
	state.Global(str)

	value, ok := state.ToString(-1)

	if !ok {
		panic("Cannot convert to string")
	}

	state.Pop(1)
	return value
}

func GetLuaBool(state *lua.State, str string) bool {
	state.Global(str)

	value := state.ToBoolean(-1)

	state.Pop(1)
	return value
}

func GetLuaFloat(state *lua.State, str string) float64 {
	state.Global(str)

	value, ok := state.ToNumber(-1)

	if !ok {
		panic("Cannot convert to float")
	}

	state.Pop(1)
	return value
}
