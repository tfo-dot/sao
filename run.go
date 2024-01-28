package main

import (
	"fmt"
	"os"
	"sao/player"
	"sao/world"
)

func main() {
	world := world.CreateWorld(os.Args[1] == "test")

	go world.StartClock()

	world.RegisterNewPlayer(player.MALE, "tfo", "344048874656366592")

	for _, pl := range world.Players {
		fmt.Printf("Player `%s` registered \n", pl.GetName())
	}
}
