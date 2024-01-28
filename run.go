package main

import (
	"embed"
	"fmt"
	"os"
	"sao/discord"
	"sao/player"
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

	go discord.StartClient(string(botKey))

	world := world.CreateWorld(os.Args[1] == "test")

	go world.StartClock()

	world.RegisterNewPlayer(player.MALE, "tfo", "344048874656366592")

	for _, pl := range world.Players {
		fmt.Printf("Player `%s` registered \n", pl.GetName())
	}

	for {
		continue
	}
}
