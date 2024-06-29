package main

import (
	"embed"
	"fmt"
	"os"
	"sao/data"
	"sao/discord"
	"sao/world"

	"github.com/google/uuid"
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

	isTest := false

	if len(os.Args) > 1 && os.Args[1] == "test" {
		isTest = true
	}

	world := world.CreateWorld(string(botKey), isTest)

	world.LoadBackup()

	tp := world.Players[uuid.MustParse("e2eb0ff0-f98e-47e4-b90c-86a21dc58259")]

	itemTemp := data.Items[data.AttackVisageUUID]

	tp.AddItem(&itemTemp)

	go world.StartClock()

	discord.World = &world

	go discord.StartClient(string(botKey))

	for {
		continue
	}
}
