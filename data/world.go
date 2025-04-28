package data

import (
	"os"
	saoParts "sao/parts"

	"github.com/tfo-dot/parts"
)

var WorldConfig WorldConfigStruct = GetWorldConfig()

type WorldConfigStruct struct {
	Hardcore   bool `parts:"HARDCORE_MODE"`
	SpeedGauge int  `parts:"SPEED_GAUGE"`
}

func GetWorldConfig() WorldConfigStruct {
	var tempConfig WorldConfigStruct

	println("Loading world config:", Config.GameDataLocation+"/world/config.pts")

	code, err := os.ReadFile(Config.GameDataLocation + "/world/config.pts")

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

	parts.ReadFromParts(vm, &tempConfig)

	return tempConfig
}
