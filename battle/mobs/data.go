package mobs

import (
	"os"
	"sao/data"
	"sao/types"
	"strings"

	saoParts "sao/parts"

	"github.com/google/uuid"
	"github.com/tfo-dot/parts"
)

var Mobs map[string]MobEntity = GetMobs()

func GetMobs() map[string]MobEntity {
	dirData, err := os.ReadDir(data.Config.GameDataLocation + "/mobs")

	if err != nil {
		panic(err)
	}

	mobs := map[string]MobEntity{}

	for _, file := range dirData {
		if file.IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name(), ".pts") {
			continue
		}

		println("Loading mob: " + file.Name())

		code, err := os.ReadFile(data.Config.GameDataLocation + "/mobs/" + file.Name())

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

		mobEntity := MobEntity{
			Effects:   make([]types.ActionEffect, 0),
			UUID:      uuid.New(),
			Props:     make(map[string]any),
			TempSkill: make([]*types.WithExpire[types.PlayerSkill], 0),
			Stats:     make(map[types.Stat]int),
			Loot:      make([]types.Loot, 0),
		}

		parts.ReadFromParts(vm, &mobEntity)

		if val, err := saoParts.FetchVal(vm, "SPD"); err != nil {
			panic(err)
		} else {
			mobEntity.Stats[types.STAT_SPD] = val.(int)
		}

		if val, err := saoParts.FetchVal(vm, "ATK"); err != nil {
			panic(err)
		} else {
			mobEntity.Stats[types.STAT_AD] = val.(int)
		}

		if val, err := saoParts.FetchVal(vm, "HP"); err != nil {
			panic(err)
		} else {
			mobEntity.Stats[types.STAT_HP] = val.(int)
			mobEntity.HP = val.(int)
		}

		mobs[mobEntity.Id] = mobEntity
	}

	return mobs
}