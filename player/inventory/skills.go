package inventory

import (
	"sao/battle"
	"sao/types"

	"github.com/google/uuid"
)

var AVAILABLE_SKILLS = []PlayerSkill{
	{
		Path:        PathControl,
		Name:        "Ogłuszenie",
		Description: "Ogłusza przeciwnika na 1 turę",
		CD: &CDMeta{
			Calc: func(meta PlayerSkill, upgrades []string) int {
				for _, upgrade := range upgrades {
					if upgrade == "CD_DEC" {
						return 2
					}
				}

				return 3
			},
		},
		Cost:     1,
		ForLevel: 1,
		Trigger: &types.Trigger{
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
				Id:          "CD_DEC",
				Description: "Zmniejsza czas odnowienia o 1 turę",
				OnEvent:     nil,
			},
			{
				Name:        "Spowolnienie",
				Id:          "SLOW",
				Description: "Po wyjściu z ogłuszenia spowolnia o 10 turę",
				OnEvent:     nil,
			},
			{
				Name:        "Przyśpieszenie",
				Id:          "SPD_INC",
				Description: "Przyśpiesza użytkownika o 10 na 1 turę",
				OnEvent:     nil,
			},
		},
		Execute: func(owner battle.PlayerEntity, target battle.Entity, fightInstance *battle.Fight) {
			effectUuid := uuid.New()

			fightInstance.HandleAction(
				battle.Action{
					Event:  battle.ACTION_EFFECT,
					Source: owner.GetUUID(),
					Target: target.GetUUID(),
					Meta: battle.ActionEffect{
						Effect:   battle.EFFECT_STUN,
						Duration: 1,
						Value:    0,
						Uuid:     effectUuid,
					},
				},
			)

			unlockedUpgrades := owner.GetUpgrades(1)

			for _, upgrade := range unlockedUpgrades {
				if upgrade == "SPD_INC" {
					(owner).ApplyEffect(battle.ActionEffect{
						Effect:   battle.EFFECT_STAT_INC,
						Duration: 1,
						Value:    10,
						Meta: &map[string]interface{}{
							"stat": types.STAT_SPD,
						},
					})
				}

				if upgrade == "SLOW" {
					fightInstance.EventHandlers[effectUuid] = types.Skill{
						Name: "Spowolnienie",
						Trigger: types.Trigger{
							Event: &types.EventTriggerDetails{
								Meta: map[string]interface{}{
									"caster": owner.GetUUID(),
								},
							},
						},
						Cost: 0,
						Execute: func(owner, target interface{}, fightInstance *interface{}) {
							(owner).(battle.Entity).ApplyEffect(battle.ActionEffect{
								Effect:   battle.EFFECT_STAT_DEC,
								Duration: 1,
								Value:    10,
								Meta: &map[string]interface{}{
									"stat": types.STAT_SPD,
								},
							})
						},
					}
				}
			}
		},
		OnEvent: nil,
	},
	{
		Path:        PathEndurance,
		Name:        "Tarcza",
		Description: "Daje tarcze pochłaniającą 25 DMG na 1 turę",
		CD: &CDMeta{
			Calc: func(meta PlayerSkill, upgrades []string) int {
				for _, upgrade := range upgrades {
					if upgrade == "CD_DEC" {
						return 2
					}
				}

				return 3
			},
		},
		Cost:     1,
		ForLevel: 1,
		Trigger: &types.Trigger{
			Type: types.TRIGGER_ACTIVE,
			Event: &types.EventTriggerDetails{
				TriggerType:   types.TRIGGER_NONE,
				TargetType:    []types.TargetTag{types.TARGET_SELF},
				TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
				TargetCount:   1,
				Meta:          map[string]interface{}{},
			},
		},
		Upgrades: []PSkillUpgrade{
			{
				Name:        "Ochłodzenie",
				Id:          "CD_DEC",
				Description: "Zmniejsza czas odnowienia o 1 turę",
				OnEvent:     nil,
			},
			{
				Name:        "Wartość++",
				Id:          "EFFECT_INC",
				Description: "Zwiększa tarczę o X",
				OnEvent:     nil,
			},
			{
				Name:        "Czas++",
				Id:          "DURATION_INC",
				Description: "Zwiększenie czasu trwania tarczy o 1 turę",
				OnEvent:     nil,
			},
		},
		Execute: func(owner battle.PlayerEntity, target battle.Entity, fightInstance *battle.Fight) {
			value := 25
			duration := 1

			unlockedUpgrades := owner.GetUpgrades(1)

			for _, upgrade := range unlockedUpgrades {
				if upgrade == "EFFECT_INC" {
					value += 0
				}

				if upgrade == "DURATION_INC" {
					duration++
				}
			}

			owner.ApplyEffect(battle.ActionEffect{
				Effect:   battle.EFFECT_SHIELD,
				Duration: duration,
				Value:    value,
				Meta:     nil,
			})
		},
	},
	{
		Path:        PathDamage,
		Name:        "Dmg++",
		Description: "Dodatkowe 25 dmg do kolejnego ataku",
		CD: &CDMeta{
			Calc: func(meta PlayerSkill, upgrades []string) int {

				for _, upgrade := range upgrades {
					if upgrade == "CD_DEC" {
						return 2
					}
				}

				return 3
			},
		},
		Cost:     1,
		ForLevel: 1,
		Trigger: &types.Trigger{
			Type: types.TRIGGER_ACTIVE,
			Event: &types.EventTriggerDetails{
				TriggerType:   types.TRIGGER_NONE,
				TargetType:    []types.TargetTag{types.TARGET_SELF},
				TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
				TargetCount:   1,
				Meta:          map[string]interface{}{},
			},
		},
		Upgrades: []PSkillUpgrade{
			{
				Name:        "Ochłodzenie",
				Id:          "CD_DEC",
				Description: "Zmniejsza czas odnowienia o 1 turę",
				OnEvent:     nil,
			},
			{
				Name:        "Wartość++",
				Id:          "EFFECT_INC",
				Description: "Zwiększa obrażenia o X",
				OnEvent:     nil,
			},
			{
				Name:        "Czas++",
				Id:          "DURATION_INC",
				Description: "Zwiększenie czasu trwania efektu o 1 turę",
				OnEvent:     nil,
			},
		},
		Execute: func(owner battle.PlayerEntity, target battle.Entity, fightInstance *battle.Fight) {
			value := 25
			duration := 1

			unlockedUpgrades := owner.GetUpgrades(1)

			for _, upgrade := range unlockedUpgrades {
				if upgrade == "EFFECT_INC" {
					value += 0
				}

				if upgrade == "DURATION_INC" {
					duration++
				}
			}

			owner.ApplyEffect(battle.ActionEffect{
				Effect:   battle.EFFECT_ON_HIT,
				Duration: duration,
				Value:    value,
				Meta: battle.ActionEffectOnHit{
					Skill:     false,
					Attack:    true,
					IsPercent: false,
				},
			})
		},
	},
	{
		Path:        PathMobility,
		Name:        "Sprint",
		Description: "Daje 10 SPD i AGL na 1 turę",
		CD: &CDMeta{
			Calc: func(meta PlayerSkill, upgrades []string) int {
				for _, upgrade := range upgrades {
					if upgrade == "CD_DEC" {
						return 2
					}
				}

				return 3
			},
		},
		Cost:     1,
		ForLevel: 1,
		Trigger: &types.Trigger{
			Type: types.TRIGGER_ACTIVE,
			Event: &types.EventTriggerDetails{
				TriggerType:   types.TRIGGER_NONE,
				TargetType:    []types.TargetTag{types.TARGET_SELF},
				TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
				TargetCount:   1,
				Meta:          map[string]interface{}{},
			},
		},
		Upgrades: []PSkillUpgrade{
			{
				Name:        "Ochłodzenie",
				Id:          "CD_DEC",
				Description: "Zmniejsza czas odnowienia o 1 turę",
				OnEvent:     nil,
			},
			{
				Name:        "Wartość++",
				Id:          "EFFECT_INC",
				Description: "Zwiększa SPD i AGL o X",
				OnEvent:     nil,
			},
			{
				Name:        "Czas++",
				Id:          "DURATION_INC",
				Description: "Zwiększenie czasu trwania efektu o 1 turę",
				OnEvent:     nil,
			},
		},
		Execute: func(owner battle.PlayerEntity, target battle.Entity, fightInstance *battle.Fight) {
			value := 10
			duration := 1

			unlockedUpgrades := owner.GetUpgrades(1)

			for _, upgrade := range unlockedUpgrades {
				if upgrade == "EFFECT_INC" {
					value += 0
				}

				if upgrade == "DURATION_INC" {
					duration++
				}
			}

			owner.ApplyEffect(battle.ActionEffect{
				Effect:   battle.EFFECT_STAT_INC,
				Duration: duration,
				Value:    0,
				Meta: battle.ActionEffectStat{
					Stat:      types.STAT_SPD,
					Value:     value,
					IsPercent: false,
				},
			})

			owner.ApplyEffect(battle.ActionEffect{
				Effect:   battle.EFFECT_STAT_INC,
				Duration: duration,
				Value:    value,
				Meta: battle.ActionEffectStat{
					Stat:      types.STAT_AGL,
					Value:     value,
					IsPercent: false,
				},
			})
		},
	},
	{
		Path:        PathControl,
		Name:        "Doładowanie",
		Description: "Zwiększa maksymalną ilość many o 5",
		CD:          nil,
		Upgrades: []PSkillUpgrade{
			{
				Name:        "Wsparcie!",
				Description: "Zwiększa moc leczenia i tarczy o 10%",
				Id:          "HEAL_INC",
				OnEvent: &map[types.CustomTrigger]func(owner *battle.PlayerEntity){
					types.CUSTOM_TRIGGER_UNLOCK: func(owner *battle.PlayerEntity) {
						(*owner).AddItem(&types.PlayerItem{
							UUID:        uuid.New(),
							Name:        "Księga Wsparcia",
							TakesSlot:   false,
							Description: "Zwiększa moc leczenia i tarczy o 10%",
							Stacks:      false,
							Consume:     false,
							Count:       1,
							MaxCount:    1,
							Hidden:      true,
							Stats: map[int]int{
								int(types.STAT_HEAL_POWER): 10,
							},
						})
					},
				},
			},
		},
		ForLevel: 2,
		Trigger:  nil,
		Cost:     0,
		OnEvent: &map[types.CustomTrigger]func(owner *battle.PlayerEntity){
			types.CUSTOM_TRIGGER_UNLOCK: func(owner *battle.PlayerEntity) {
				(*owner).AddItem(&types.PlayerItem{
					UUID:        uuid.New(),
					Name:        "Księga Doładowania",
					TakesSlot:   false,
					Description: "Zwiększa maksymalną ilość many o 5",
					Stacks:      false,
					Consume:     false,
					Count:       1,
					MaxCount:    1,
					Hidden:      true,
					Stats: map[int]int{
						int(types.STAT_MANA): 5,
					},
				})
			},
		},
	},
	{
		Path:        PathMobility,
		Name:        "Speed!",
		Description: "Zwiększa SPD i DGD o 5",
		CD:          nil,
		ForLevel:    2,
		Trigger: &types.Trigger{
			Type: types.TRIGGER_PASSIVE,
		},
		Cost:     0,
		Upgrades: []PSkillUpgrade{},
		Execute:  nil,
	},
}

type PlayerSkill struct {
	Name        string
	Description string
	Path        SkillPath
	ForLevel    int
	Trigger     *types.Trigger
	CD          *CDMeta
	Cost        int
	Upgrades    []PSkillUpgrade
	Execute     func(owner battle.PlayerEntity, target battle.Entity, fightInstance *battle.Fight)
	OnEvent     *map[types.CustomTrigger]func(owner *battle.PlayerEntity)
}

type CDMeta struct {
	Calc func(PlayerSkill, []string) int
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
	Id          string
	Execute     func(owner battle.PlayerEntity, target battle.Entity, fightInstance *battle.Fight)
	OnEvent     *map[types.CustomTrigger]func(owner *battle.PlayerEntity)
}
