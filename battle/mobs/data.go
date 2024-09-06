package mobs

import (
	"os"
	"sao/battle"
	"sao/config"
	saoLua "sao/lua"
	"sao/types"
	"sao/utils"

	"github.com/Shopify/go-lua"
	"github.com/google/uuid"
)

var Mobs map[string]MobEntity = GetMobs()

func GetMobs() map[string]MobEntity {
	dirData, err := os.ReadDir(config.Config.GameDataLocation + "/mobs")

	if err != nil {
		panic(err)
	}

	mobs := map[string]MobEntity{}

	for _, file := range dirData {
		if file.IsDir() {
			continue
		}

		state := lua.NewState()

		lua.OpenLibraries(state)

		saoLua.AddFightFunctions(state)
		saoLua.AddEntityFunctions(state)
		saoLua.AddPlayerFunctions(state)
		saoLua.AddStatTypes(state)

		println("Loading mob: " + file.Name())

		err := lua.DoFile(state, config.Config.GameDataLocation+"/mobs/"+file.Name())

		if err != nil {
			panic(err)
		}

		MobId := utils.GetLuaString(state, "Id")
		MobName := utils.GetLuaString(state, "Name")
		MobHP := utils.GetLuaInt(state, "HP")
		MobATK := utils.GetLuaInt(state, "ATK")
		MobSPD := utils.GetLuaInt(state, "SPD")

		loot := make([]types.Loot, 0)

		state.Global("Loot")

		tab, err := utils.GetTableAsArray(state)

		if err != nil {
			panic(err)
		}

		for _, lootItem := range tab {
			lootItem := lootItem.(map[string]interface{})

			loot = append(loot, types.Loot{
				Type:  types.LootType(lootItem["Type"].(float64)),
				Count: int(lootItem["Count"].(float64)),
			})
		}

		var onDefeat func(types.PlayerEntity)

		state.Global("OnDefeat")

		if state.IsFunction(-1) {
			state.Pop(1)

			onDefeat = func(player types.PlayerEntity) {
				state.Global("OnDefeat")

				state.PushUserData(player)

				state.Call(1, 0)
			}
		} else {
			state.Pop(1)
		}

		var onAction func(*MobEntity, *battle.Fight) []types.Action

		state.Global("Action")

		if state.IsFunction(-1) {
			state.Pop(1)

			onAction = func(mob *MobEntity, fightInstance *battle.Fight) []types.Action {
				state.Global("Action")

				state.PushUserData(mob)
				state.PushUserData(fightInstance)

				state.Call(2, 1)

				temp, err := utils.GetTableAsArray(state)

				if err != nil {
					panic(err)
				}

				actions := make([]types.Action, len(temp))

				for idx, action := range temp {
					action := action.(map[string]interface{})

					actions[idx] = saoLua.ParseActionReturn(action, state)
				}

				return actions
			}
		} else {
			state.Pop(1)
		}

		mobs[MobId] = MobEntity{
			Id:           MobId,
			HP:           MobHP,
			Effects:      make([]types.ActionEffect, 0),
			UUID:         uuid.New(),
			Name:         MobName,
			Props:        make(map[string]interface{}),
			Loot:         loot,
			TempSkill:    make([]*types.WithExpire[types.PlayerSkill], 0),
			OnDefeatFunc: onDefeat,
			ActionFunc:   onAction,
			Stats: map[types.Stat]int{
				types.STAT_AD:  MobATK,
				types.STAT_SPD: MobSPD,
				types.STAT_HP:  MobHP,
			},
		}
	}

	return mobs
}
