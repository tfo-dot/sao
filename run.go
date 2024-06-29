package main

import (
	"embed"
	"fmt"
	"os"
	"sao/discord"
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

	isTest := false

	if len(os.Args) > 1 && os.Args[1] == "test" {
		isTest = true
	}

	world := world.CreateWorld(string(botKey), isTest)

	world.LoadBackup()

	go world.StartClock()

	discord.World = &world

	go discord.StartClient(string(botKey))

	for {
		continue
	}
}
