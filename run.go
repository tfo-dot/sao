package main

import (
	"fmt"
	"net"
	"net/http"
	"sao/world"
)

func main() {
	world := world.CreateWorld()

	go world.StartClock()

	ln, errSoc := net.Listen("tcp", "127.0.0.1:8080")

	http.HandleFunc("/entity/", world.HTTPGetEntity)
	http.HandleFunc("/player/details/", world.HTTPGetPlayer)
	http.HandleFunc("/player/actions/", world.HTTPGetPlayerActions)
	http.HandleFunc("/player/store/", world.HTTPGetPlayerStore)
	http.HandleFunc("/fight/", world.HTTPGetFight)
	http.HandleFunc("/time/", world.HTTPGetTime)
	http.HandleFunc("/store/", world.HTTPGetStore)

	go http.ListenAndServe("127.0.0.1:8124", nil)

	if errSoc != nil {
		panic(errSoc)
	}

	//Register character speedrun
	world.HandlePacket([]byte{0, 0, 18, 51, 52, 52, 48, 52, 56, 56, 55, 52, 54, 53, 54, 51, 54, 54, 53, 57, 50, 3, 116, 102, 111})

	for i, pl := range world.Players {
		// world.PlayerEncounter(i)

		fmt.Printf("Player %s: %s\n", i, pl.GetName())

		/* potionUUID := uuid.New()

		//CONSUMABLE EXAMPLE

			pl.Inventory.Items = append(pl.Inventory.Items, inventory.PlayerItem{
				UUID:     potionUUID,
				Name:     "Health Potion",
				Stats:    map[battle.Stat]int{},
				Consume:  true,
				Count:    1,
				MaxCount: 5,
				Stacks:   true,
				Effects: []inventory.Skill{
					{
						Name: "Heal",
						Trigger: inventory.Trigger{
							Type: inventory.Active,
						},
						Cost: nil,
						Execute: func(owner interface{}) {
							owner.(*player.Player).Heal(100)
						},
					}},
			})

			pl.Inventory.UseItem(potionUUID, &pl) */

		//Passive item example
		/*
			itemUUID := uuid.New()
			pl.Inventory.Items = append(pl.Inventory.Items, inventory.PlayerItem{
				UUID: itemUUID,
				Name: "Bami's Cinder",
				Stats: map[battle.Stat]int{
					battle.STAT_HP: 100,
				},
				Consume: false,
				Count:   1,
				Stacks:  false,
				Effects: []types.Skill{
					{
						Name: "Burning",
						Trigger: types.Trigger{
							Type: types.TRIGGER_PASSIVE,
							Event: &types.EventTriggerDetails{
								TriggerType:   types.TRIGGER_TURN,
								TargetType:    []types.TargetTag{types.TARGET_ENEMY},
								TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
								TargetCount:   -1,
							},
						},
						Cost: nil,
						Execute: func(source, target, fight interface{}) {
							//Fight is here so use that to add dmg event
							target.(battle.Entity).TakeDMG(battle.ActionDamage{
								Damage: []battle.Damage{
									{
										Type:  battle.DMG_TRUE,
										Value: utils.PercentOf(target.(battle.Entity).GetMaxHP(), 1),
									},
								},
								CanDodge: false,
							})
						},
					}},
			})*/
	}

	for {
		conn, err := ln.Accept()

		if err != nil {
			fmt.Printf("Error (1): %s", err.Error())
			continue
		}

		if world.Conn == nil {
			world.SetConnection(&conn)
		} else {
			conn.Close()
			fmt.Println("Connection sus")
		}
	}
}
