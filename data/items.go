package data

import (
	"fmt"
	"os"
	"sao/config"
	saoLua "sao/lua"
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

		state.PushGoFunction(func(state *lua.State) int {
			value := lua.CheckInteger(state, 1)
			percent := lua.CheckInteger(state, 2)

			state.PushInteger(utils.PercentOf(value, percent))

			return 1
		})

		state.SetField(-2, "PercentOf")

		state.PushGoFunction(func(state *lua.State) int {
			state.PushString(uuid.New().String())

			return 1
		})

		state.SetField(-2, "GenerateUUID")

		state.SetGlobal("utils")

		saoLua.AddStatTypes(state)
		saoLua.AddPlayerFunctions(state)
		saoLua.AddEntityFunctions(state)
		saoLua.AddFightFunctions(state)

		println("Loading item: " + file.Name())

		err := lua.DoFile(state, config.Config.GameDataLocation+"/items/"+file.Name())

		if err != nil {
			panic(err)
		}

		item := types.PlayerItem{
			UUID:        uuid.MustParse(utils.GetLuaString(state, "UUID")),
			Name:        utils.GetLuaString(state, "Name"),
			Description: utils.GetLuaString(state, "Description"),
			TakesSlot:   utils.GetLuaBool(state, "TakesSlot"),
			Stacks:      utils.GetLuaBool(state, "Stacks"),
			Consume:     utils.GetLuaBool(state, "Consume"),
			Count:       utils.GetLuaInt(state, "Count"),
			MaxCount:    utils.GetLuaInt(state, "MaxCount"),
			Hidden:      utils.GetLuaBool(state, "Hidden"),
			Stats:       map[types.Stat]int{},
			Effects:     []types.PlayerSkill{},
		}

		state.Global("Stats")

		tempStats, err := utils.GetTableAsMap(state)

		if err != nil {
			panic(err)
		}

		for key, value := range tempStats {
			item.Stats[utils.StringToStat[key]] = int(value.(float64))
		}

		state.Global("Effects")

		if state.IsNil(-1) {
			state.Pop(1)
			items[item.UUID] = item
			continue
		}

		tempEffects, err := utils.GetTableAsArray(state)

		if err != nil {
			panic(err)
		}

		for idx, effect := range tempEffects {
			item.Effects = append(item.Effects, ItemEffect{
				State:      state,
				Idx:        idx,
				EffectData: effect.(map[string]interface{}),
			})
		}

		items[item.UUID] = item
	}

	return items
}

type ItemEffect struct {
	State      *lua.State
	Idx        int
	EffectData map[string]interface{}
}

func (ie ItemEffect) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	if execute, exists := ie.EffectData["Execute"]; exists {
		if val, ok := execute.(utils.LuaFunctionRef); ok {
			ie.State.Global(val.FunctionName)

			ie.State.PushUserData(owner)
			ie.State.PushUserData(target)
			ie.State.PushUserData(fightInstance)

			//TODO push as table
			ie.State.PushUserData(meta)

			ie.State.Call(4, 1)

			if !ie.State.IsNil(-1) {
				rValue, err := utils.GetTableAsMap(ie.State)

				if err != nil {
					panic(err)
				}

				trigger := ie.GetTrigger()

				return saoLua.ParseReturnMeta(rValue, trigger)
			}

			return nil
		}

		panic(fmt.Sprintf("Execute (#%d) is not a function", ie.Idx))
	}

	return nil
}

func (ie ItemEffect) GetEvents() map[types.CustomTrigger]func(owner types.PlayerEntity) {
	eventData, exists := ie.EffectData["Events"]

	if !exists {
		return nil
	}

	events := map[types.CustomTrigger]func(owner types.PlayerEntity){}

	for key, value := range eventData.(map[string]interface{}) {
		if key == "TRIGGER_UNLOCK" {
			if val, ok := value.(utils.LuaFunctionRef); ok {
				events[types.CUSTOM_TRIGGER_UNLOCK] = func(owner types.PlayerEntity) {
					ie.State.Global(val.FunctionName)

					ie.State.PushUserData(owner)

					ie.State.Call(1, 0)
				}
			}
		}
	}

	return events
}

func (ie ItemEffect) GetUUID() uuid.UUID {
	return uuid.New()
}

func (ie ItemEffect) GetName() string {
	return ""
}

func (ie ItemEffect) GetDescription() string {
	return ""
}

func (ie ItemEffect) GetCD() int {
	if cd, exists := ie.EffectData["CD"]; exists {
		return int(cd.(float64))
	}

	return 0
}

func (ie ItemEffect) GetCost() int {
	return 0
}

func (ie ItemEffect) GetTrigger() types.Trigger {
	if trigger, exists := ie.EffectData["Trigger"]; exists {
		return saoLua.ReadMapAsTrigger(trigger.(map[string]interface{}))
	} else {
		panic(fmt.Sprintf("Trigger (#%d) is not defined", ie.Idx))
	}
}

func (ie ItemEffect) IsLevelSkill() bool {
	return false
}

type SimplePlayerSkill struct {
	Trigger types.Trigger
	Exec    func(owner, target, fightInstance, meta interface{}) interface{}
}

func (s SimplePlayerSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	return s.Exec(owner, target, fightInstance, meta)
}

func (s SimplePlayerSkill) GetEvents() map[types.CustomTrigger]func(owner interface{}) {
	return nil
}

func (s SimplePlayerSkill) GetUUID() uuid.UUID {
	return uuid.New()
}

func (s SimplePlayerSkill) GetName() string {
	return ""
}

func (s SimplePlayerSkill) GetDescription() string {
	return ""
}

func (s SimplePlayerSkill) GetCD() int {
	return 0
}

func (s SimplePlayerSkill) GetCost() int {
	return 0
}

func (s SimplePlayerSkill) GetTrigger() types.Trigger {
	return s.Trigger
}

func (s SimplePlayerSkill) IsLevelSkill() bool {
	return false
}
