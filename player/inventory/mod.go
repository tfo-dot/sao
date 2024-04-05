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
	Skills              []types.PlayerSkill
	CDS                 map[uuid.UUID]int
	LevelSkillsCDS      map[int]int
	LevelSkills         map[int]PlayerSkillLevel
	LevelSkillsUpgrades map[int][]string
}

func (inv *PlayerInventory) AddIngredient(ingredient *types.Ingredient) {
	if _, exists := inv.Ingredients[ingredient.UUID]; exists {
		inv.Ingredients[ingredient.UUID].Count += ingredient.Count
		return
	}

	inv.Ingredients[ingredient.UUID] = ingredient
}

func (inv *PlayerInventory) Craft(recipe types.Recipe) error {
	for _, ingredient := range recipe.Ingredients {
		if inv.Ingredients[ingredient.Item].Count < ingredient.Count {
			return errors.New("MISSING_INGREDIENT")
		}
	}

	for _, ingredient := range recipe.Ingredients {
		inv.Ingredients[ingredient.Item].Count -= ingredient.Count
	}

	//TODO add result item to inventory

	return nil
}

func (inv *PlayerInventory) AddItem(item *types.PlayerItem) {
	for _, invItem := range inv.Items {
		if invItem.UUID == item.UUID && invItem.Stacks && invItem.Count < invItem.MaxCount {
			if invItem.Count+item.Count > invItem.MaxCount {
				item.Count -= invItem.MaxCount - invItem.Count
				invItem.Count = invItem.MaxCount
			} else {
				invItem.Count += item.Count
				return
			}
		}
	}

	//This will bite me in the ass later
	inv.Items = append(inv.Items, item)
}

func (inv PlayerInventory) HasIngredients(ingredients []types.Ingredient) bool {
	for _, ingredient := range ingredients {
		entry, exists := inv.Ingredients[ingredient.UUID]

		if !exists || ingredient.Count == 0 {
			return false
		}

		if entry.Count < ingredient.Count {
			return false
		}
	}

	return true
}

func (inv *PlayerInventory) RemoveIngredients(ingredients []types.Ingredient) {
	for _, ingredient := range ingredients {
		entry, exists := inv.Ingredients[ingredient.UUID]

		if !exists {
			continue
		}

		entry.Count -= ingredient.Count

		if entry.Count <= 0 {
			delete(inv.Ingredients, ingredient.UUID)
		}
	}
}

func (inv PlayerInventory) GetStat(stat types.Stat) int {
	value := 0

	for _, item := range inv.Items {
		val, exists := item.Stats[int(stat)]

		if exists {
			value += val
		}
	}

	return value
}

func (inv *PlayerInventory) UseItem(itemUuid uuid.UUID, owner interface{}, target interface{}, fightInstance *interface{}) {
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

func (inv *PlayerInventory) UnlockSkill(path types.SkillPath, lvl int, playerLvl int, player *battle.PlayerEntity) error {
	if lvl > playerLvl {
		return errors.New("PLAYER_LVL_TOO_LOW")
	}

	if _, exists := inv.LevelSkills[lvl]; exists {
		return errors.New("SKILL_ALREADY_UNLOCKED")
	}

	skill, skillExists := AVAILABLE_SKILLS[path][lvl]

	if !skillExists {
		return errors.New("SKILL_NOT_FOUND")
	}

	inv.LevelSkills[lvl] = skill

	skillEvents := skill.GetEvents()

	if effect, effectExists := skillEvents[types.CUSTOM_TRIGGER_UNLOCK]; effectExists {
		var tempPlayer interface{} = player

		effect(&tempPlayer)
	}

	return nil
}

func (inv *PlayerInventory) UpgradeSkill(lvl int, upgradeName string) error {
	if _, exists := inv.LevelSkills[lvl]; !exists {
		return errors.New("SKILL_NOT_UNLOCKED")
	}

	skill, skillExists := inv.LevelSkills[lvl]

	if !skillExists {
		return errors.New("SKILL_NOT_FOUND")
	}

	if exists := skill.HasUpgrade(upgradeName); !exists {
		return errors.New("UPGRADE_NOT_FOUND")
	}

	skill.UnlockUpgrade(upgradeName)

	inv.LevelSkillsUpgrades[lvl] = append(inv.LevelSkillsUpgrades[lvl], upgradeName)

	return nil
}

func GetDefaultInventory() PlayerInventory {
	return PlayerInventory{
		Gold:                0,
		Capacity:            10,
		Ingredients:         map[uuid.UUID]*types.Ingredient{},
		Items:               []*types.PlayerItem{},
		CDS:                 map[uuid.UUID]int{},
		Skills:              []types.PlayerSkill{},
		LevelSkills:         map[int]PlayerSkillLevel{},
		LevelSkillsCDS:      map[int]int{},
		LevelSkillsUpgrades: map[int][]string{},
	}
}

func (inv *PlayerInventory) UseSkill(skillUuid uuid.UUID, owner, target interface{}, fightInstance *interface{}) {
	for _, skill := range inv.Skills {
		if skill.GetTrigger().Type != types.TRIGGER_ACTIVE {
			continue
		}

		if skill.GetUUID() == skillUuid {
			skill.Execute(owner, target, fightInstance)

			inv.CDS[skill.GetUUID()] = skill.GetCD()

			return
		}
	}
}
