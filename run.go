package main

import (
	"fmt"
	"sao/player"
	"sao/world"
)

func main() {
	world := world.CreateWorld()

	go world.StartClock()

	world.RegisterNewPlayer(player.MALE, "tfo", "<id here>")

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
}
