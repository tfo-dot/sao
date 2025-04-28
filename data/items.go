package data

import (
	"fmt"
	"os"
	saoParts "sao/parts"
	"sao/types"
	"strings"

	"github.com/google/uuid"
	"github.com/tfo-dot/parts"
)

var Items = GetItems()

func GetItems() map[uuid.UUID]types.PlayerItem {
	dirData, err := os.ReadDir(Config.GameDataLocation + "/items")

	if err != nil {
		panic(err)
	}

	items := map[uuid.UUID]types.PlayerItem{}

	for _, file := range dirData {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".pts") {
			continue
		}

		println("Loading item: " + file.Name())

		code, err := os.ReadFile(Config.GameDataLocation + "/items/" + file.Name())

		if err != nil {
			panic(err)
		}

		vm, err := parts.GetVMWithSource(string(code))

		if err != nil {
			panic(err)
		}

		saoParts.AddConsts(vm)
		saoParts.AddFunctions(vm)

		item := types.PlayerItem{
			TakesSlot: true,
			Stacks:    false,
			Consume:   false,
			Count:     1,
			MaxCount:  1,
			Hidden:    false,
			Effects:   []types.PlayerSkill{},
			Stats:     make(map[types.Stat]int),
		}

		err = vm.Run()

		parts.ReadFromParts(vm, &item)

		if err != nil {
			panic(err)
		}

		statMap, err := saoParts.FetchVal(vm, "Stats")

		if err != nil {
			panic(err)
		}

		for key, value := range statMap.(map[string]any) {
			if key == "RTDerived" {
				for _, stat := range value.([]any) {
					item.DerivedStats = append(item.DerivedStats, types.DerivedStat{
						Base:    types.Stat(stat.(map[string]any)["RTBase"].(int)),
						Derived: types.Stat(stat.(map[string]any)["RTDerived"].(int)),
						Percent: stat.(map[string]any)["RTPercent"].(int),
					})
				}

				continue
			}

			keyRaw, err := vm.Enviroment.Resolve(fmt.Sprintf("STAT_%s", strings.TrimPrefix(key, "RT")))

			if err != nil {
				panic(err)
			}

			keyVal := types.Stat(keyRaw.Value.(int))

			item.Stats[keyVal] = value.(int)
		}

		if vm.Enviroment.Has("Effects") {
			effectList, err := saoParts.FetchVal(vm, "Effects")

			if err != nil {
				panic(err)
			}

			for _, val := range effectList.([]any) {
				item.Effects = append(item.Effects, ItemEffect{val.(map[string]any)})
			}
		}

		rawUUID, err := saoParts.FetchVal(vm, "UUID")

		if err != nil {
			panic(err)
		}

		item.UUID = uuid.MustParse(rawUUID.(string))

		items[item.UUID] = item
	}

	return items
}

type ItemEffect struct {
	EffectData map[string]any
}

func (ie ItemEffect) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	if execute, exists := ie.EffectData["Execute"]; exists {
		res, err := execute.(func(...any) (any, error))(owner, target, fightInstance, meta)

		if err != nil {
			panic(err)
		}

		return res
	}

	return nil
}

func (ie ItemEffect) GetEvents() map[types.CustomTrigger]func(owner types.PlayerEntity) {
	eventData, exists := ie.EffectData["Events"]

	if !exists {
		return nil
	}

	events := map[types.CustomTrigger]func(owner types.PlayerEntity){}

	for key, value := range eventData.(map[int]any) {
		switch types.CustomTrigger(key) {
		case types.CUSTOM_TRIGGER_UNLOCK:
			events[types.CUSTOM_TRIGGER_UNLOCK] = func(owner types.PlayerEntity) {
				value.(func(...any) (any, error))(owner)
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
		return cd.(int)
	}

	return 0
}

func (ie ItemEffect) GetCost() int {
	return 0
}

func (ie ItemEffect) GetTrigger() types.Trigger {
	if _, exists := ie.EffectData["Trigger"]; exists {
		//TODO parse this
		return types.Trigger{}
	} else {
		panic(fmt.Sprintf("Trigger is not defined"))
	}
}

func (ie ItemEffect) IsLevelSkill() bool {
	return false
}