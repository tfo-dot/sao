package inventory

import (
	"sao/battle"
	"sao/types"

	"github.com/google/uuid"
)

var AVAILABLE_SKILLS = []PSkill{
	{
		Path:        PathControl,
		Name:        "Ogłuszenie",
		Description: "Ogłusza przeciwnika na 1 turę",
		CD:          3,
		Cost:        1,
		ForLevel:    1,
		Trigger: types.Trigger{
			Type: types.TRIGGER_ACTIVE,
			Event: &types.EventTriggerDetails{
				TriggerType:   types.TRIGGER_NONE,
				TargetType:    []types.TargetTag{types.TARGET_ENEMY},
				TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
				TargetCount:   1,
				Meta:          map[string]interface{}{},
			},
		},
		Upgrades: []PSkillUpgrade{
			{
				Name:        "Ochłodzenie",
				Description: "Zmniejsza czas odnowienia o 1 turę",
			},
		},
		Execute: func(owner battle.Entity, target battle.Entity, fightInstance *battle.Fight, skill *PSkill) {
			fightInstance.HandleAction(
				battle.Action{
					Event:  battle.ACTION_EFFECT,
					Source: owner.GetUUID(),
					Target: target.GetUUID(),
					Meta: battle.ActionEffect{
						Effect:   battle.EFFECT_STUN,
						Duration: 1,
						Value:    0,
					},
				},
			)
		},
		OnEvent: nil,
	},
	{
		Path:        PathControl,
		Name:        "Mana++",
		Description: "Zwiększa maksymalną ilość many o 5",
		CD:          -1,
		Cost:        -1,
		ForLevel:    2,
		Trigger: types.Trigger{
			Type: types.TRIGGER_PASSIVE,
			Event: &types.EventTriggerDetails{
				TriggerType:   types.TRIGGER_NONE,
				TargetType:    []types.TargetTag{},
				TargetDetails: []types.TargetDetails{},
				TargetCount:   -2,
				Meta:          map[string]interface{}{},
			},
		},
		Upgrades: []PSkillUpgrade{},
		Execute:  func(owner battle.Entity, target battle.Entity, fightInstance *battle.Fight, skill *PSkill) {},
		OnEvent: &map[types.CustomTrigger]func(owner *battle.PlayerEntity){
			types.CUSTOM_TRIGGER_UNLOCK: func(owner *battle.PlayerEntity) {
				(*owner).AddItem(PlayerItem{
					UUID:        uuid.New(),
					Name:        "Mana++",
					Description: "Zwiększa maksymalną ilość many o 5",
					//It's item from skill since I have no actual way of representing it
					TakesSlot: false,
					Stacks:    false,
					Consume:   false,
					Count:     1,
					MaxCount:  1,
					Stats: map[battle.Stat]int{
						battle.STAT_MANA: 5,
					},
					Effects: []types.Skill{},
					Hidden:  true,
				})

				(*owner).RestoreMana(5)
			},
		},
	},
}

type PSkill struct {
	Name        string
	Description string
	Type        string
	Path        SkillPath
	ForLevel    int
	Trigger     types.Trigger
	CD          int
	Cost        int
	Upgrades    []PSkillUpgrade
	Execute     func(owner battle.Entity, target battle.Entity, fightInstance *battle.Fight, skill *PSkill)
	OnEvent     *map[types.CustomTrigger]func(owner *battle.PlayerEntity)
}

type SkillPath int

const (
	PathControl SkillPath = iota
	PathEndurance
	PathDamage
	PathMobility
)

type PSkillUpgrade struct {
	Name        string
	Description string
	//TODO add upgrade effects
}
