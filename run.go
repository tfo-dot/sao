package main

import (
	"embed"
	"fmt"
	"os"
	"sao/discord"
	"sao/player/inventory"
	"sao/world"
)

//go:embed botkey.txt
var botkey embed.FS

func main() {
	botkeyRaw, err := botkey.Open("botkey.txt")

	if err != nil {
		fmt.Println("Cannot read config file")
		return
	}

	botKeyFileStats, err := botkeyRaw.Stat()

	if err != nil {
		fmt.Println("Cannot read config file")
		return
	}

	botKey := make([]byte, botKeyFileStats.Size())

	_, err = botkeyRaw.Read(botKey)

	if err != nil {
		fmt.Println("Cannot read config file")
		return
	}

	world := world.CreateWorld(os.Args[1] == "test")

	go world.StartClock()

	discord.World = &world

	world.RegisterNewPlayer("tfo", "344048874656366592")

	for _, pl := range world.Players {
		if pl.Name == "tfo" {

			pl.AddEXP(1000)

			fmt.Println(pl.Inventory.UnlockSkill(inventory.PathControl, 2, pl.XP.Level, pl))
		}
	}

	go discord.StartClient(string(botKey))

	for {
		continue
	}
}
