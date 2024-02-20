package main

import (
	"embed"
	"fmt"
	"os"
	"sao/discord"
	"sao/player"
	"sao/player/inventory"
	"sao/types"
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

	world := world.CreateWorld(os.Args[1] == "test")

	go world.StartClock()

	discord.World = &world

	testPlayer := world.RegisterNewPlayer("tfo", "344048874656366592")

	testPlayer.Inventory.Items = append(testPlayer.Inventory.Items, &types.PlayerItem{
		UUID:        uuid.New(),
		Name:        "Mikstura",
		Description: "Przywraca 25 punktów życia",
		Count:       1,
		MaxCount:    5,
		TakesSlot:   true,
		Stacks:      true,
		Consume:     true,
		Hidden:      false,
		Stats:       map[int]int{},
		Effects: []types.Skill{
			{
				Name: "Uzdrowienie",
				Trigger: types.Trigger{
					Type: types.TRIGGER_ACTIVE,
					Event: &types.EventTriggerDetails{
						TriggerType:   types.TRIGGER_NONE,
						TargetType:    []types.TargetTag{types.TARGET_SELF},
						TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
						TargetCount:   1,
						Meta:          map[string]interface{}{},
					},
				},
				Cost: 0,
				Execute: func(source, target interface{}, fight *interface{}) {
					source.(*player.Player).Heal(25)
				},
			},
		},
	})

	err = testPlayer.Inventory.UnlockSkill(inventory.PathControl, 1, 1, testPlayer)

	if err != nil {
		fmt.Println(err)
	}

	go discord.StartClient(string(botKey))

	// world.PlayerSearch(testPlayer.GetUUID())

	for {
		continue
	}
}
