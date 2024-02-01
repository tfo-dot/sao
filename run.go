package main

import (
	"embed"
	"fmt"
	"os"
	"sao/discord"
	"sao/world"
	"sao/world/party"

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

	world := world.CreateWorld(os.Args[1] == "test")

	go world.StartClock()

	discord.World = &world

	world.RegisterNewPlayer("tfo", "344048874656366592")

	for _, pl := range world.Players {
		if pl.Name == "tfo" {

			pUuid := uuid.New()

			world.Parties[pUuid] = &party.Party{
				Players: []uuid.UUID{pl.GetUUID()},
				Roles: &map[party.PartyRole][]uuid.UUID{
					party.Leader: {pl.GetUUID()},
				},
			}

			pl.Meta.Party = &pUuid
		}
	}

	go discord.StartClient(string(botKey))

	for {
		continue
	}
}
