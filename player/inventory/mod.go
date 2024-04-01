package inventory

import (
	"errors"
	"sao/battle"
	"sao/types"

	"github.com/google/uuid"
)

type PlayerInventory struct {
	Gold                int
	Capacity            int
	Items               []*types.PlayerItem
	Ingredients         map[uuid.UUID]*types.Ingredient
	Skills              []*types.PlayerSkill
	CDS                 map[uuid.UUID]int
	LevelSkillsCDS      map[int]int
	LevelSkills         map[int]PlayerSkill
	LevelSkillsUpgrades map[int][]string
}

func (inv PlayerInventory) AddIngredient(ingredient *types.Ingredient) {
	if _, exists := inv.Ingredients[ingredient.UUID]; exists {
		inv.Ingredients[ingredient.UUID].Count += ingredient.Count
		return
	}

	inv.Ingredients[ingredient.UUID] = ingredient
}

func (inv PlayerInventory) GetStat(stat battle.Stat) int {
	value := 0

	for _, item := range inv.Items {
		val, exists := item.Stats[int(stat)]

		if exists {
			value += val
		}
	}

	return value
}

func (inv PlayerInventory) UseItem(itemUuid uuid.UUID, owner interface{}, target interface{}, fightInstance *interface{}) {
	for i, item := range inv.Items {
		if item.UUID == itemUuid {
			if item.Consume && item.Count <= 0 {
				return
			}

			if item.Consume {
				item.Count--
			}

			inv.Items[i].UseItem(owner, target, fightInstance)

			if item.Count == 0 && item.Consume {
				inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			}
		}
	}
}

func (inv PlayerInventory) UnlockSkill(path SkillPath, lvl int, playerLvl int, player battle.PlayerEntity) error {
	if lvl > playerLvl {
		return errors.New("PLAYER_LVL_TOO_LOW")
	}

	if _, exists := inv.LevelSkills[lvl]; exists {
		return errors.New("SKILL_ALREADY_UNLOCKED")
	}

	for _, skill := range AVAILABLE_SKILLS {
		if skill.Path == path && skill.ForLevel == lvl {
			inv.LevelSkills[lvl] = skill

			if skill.OnEvent != nil {
				if fn, exists := (*skill.OnEvent)[types.CUSTOM_TRIGGER_UNLOCK]; exists {
					fn(&player)
				}
			}

			return nil
		}
	}

	return errors.New("SKILL_NOT_FOUND")
}

func (inv PlayerInventory) UpgradeSkill(path SkillPath, lvl int, upgradeName string) error {
	if _, exists := inv.LevelSkills[lvl]; !exists {
		return errors.New("SKILL_NOT_UNLOCKED")
	}

	for _, skill := range AVAILABLE_SKILLS {
		if skill.Path == path && skill.ForLevel == lvl {
			for _, upgrade := range skill.Upgrades {
				if upgrade.Name == upgradeName {
					if _, exists := inv.LevelSkillsUpgrades[lvl]; !exists {
						inv.LevelSkillsUpgrades[lvl] = []string{}
					}

					inv.LevelSkillsUpgrades[lvl] = append(inv.LevelSkillsUpgrades[lvl], upgrade.Name)
					return nil
				}
			}
		}
	}

	return errors.New("SKILL_NOT_FOUND")
}

func GetDefaultInventory() PlayerInventory {
	return PlayerInventory{
		Gold:        0,
		Capacity:    10,
		Items:       []*types.PlayerItem{},
		CDS:         map[uuid.UUID]int{},
		Skills:      []*types.PlayerSkill{},
		LevelSkills: map[int]PlayerSkill{},
	}
}

func (inv *PlayerInventory) UseSkill(skillUuid uuid.UUID, owner, target, fightInstance interface{}) {
	for _, skill := range inv.Skills {
		if skill.Trigger.Type != types.TRIGGER_ACTIVE {
			continue
		}

		if skill.UUID == skillUuid {
			skill.Action(owner, target, fightInstance)

			inv.CDS[skill.UUID] = skill.CD

			return
		}
	}
}
