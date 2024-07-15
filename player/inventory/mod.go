package inventory

import (
	"errors"
	"sao/battle"
	"sao/data"
	"sao/types"

	"github.com/google/uuid"
)

type PlayerInventory struct {
	Gold                int
	TempSkills          []*types.WithExpire[types.PlayerSkill]
	Items               []*types.PlayerItem
	Ingredients         map[uuid.UUID]*types.Ingredient
	Cooldowns           map[uuid.UUID]int
	LevelSkillsCDS      map[int]int
	LevelSkills         map[int]PlayerSkillLevel
	LevelSkillsUpgrades map[int]int
}

func (inv *PlayerInventory) AddTempSkill(skill types.WithExpire[types.PlayerSkill]) {
	inv.TempSkills = append(inv.TempSkills, &skill)
}

func (inv *PlayerInventory) Serialize() map[string]interface{} {
	//TODO finish this
	lvlSkills := make([]int, 0)

	for key := range inv.LevelSkills {
		lvlSkills = append(lvlSkills, key)
	}

	return map[string]interface{}{
		"gold":        inv.Gold,
		"items":       []*types.PlayerItem{},
		"ingredients": inv.Ingredients,
		"skills": map[string]interface{}{
			"skills":   lvlSkills,
			"upgrades": inv.LevelSkillsUpgrades,
			"cds":      inv.LevelSkillsCDS,
		},
	}
}

func DeserializeInventory(rawData map[string]interface{}) PlayerInventory {
	inv := PlayerInventory{
		Gold:                int(rawData["gold"].(float64)),
		TempSkills:          []*types.WithExpire[types.PlayerSkill]{},
		Items:               []*types.PlayerItem{},
		Ingredients:         map[uuid.UUID]*types.Ingredient{},
		Cooldowns:           map[uuid.UUID]int{},
		LevelSkillsCDS:      map[int]int{},
		LevelSkills:         map[int]PlayerSkillLevel{},
		LevelSkillsUpgrades: map[int]int{},
	}

	//TODO finish this

	return inv
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

	newItem := data.Items[recipe.Product.UUID]
	newItem.Count = recipe.Product.Count

	inv.Items = append(inv.Items, &newItem)

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
		val, exists := item.Stats[stat]

		if exists {
			value += val
		}
	}

	for _, ingredient := range inv.Ingredients {
		val, exists := ingredient.Stats[stat]

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

	if len(skill) == 1 {
		inv.LevelSkills[lvl] = skill[0]
	} else {
		panic("Multiple options")
		//TODO: Implement skill selection
	}

	skillEvents := inv.LevelSkills[lvl].GetEvents()

	if effect, effectExists := skillEvents[types.CUSTOM_TRIGGER_UNLOCK]; effectExists {
		effect(&player)
	}

	return nil
}

func (inv *PlayerInventory) UpgradeSkill(lvl int, upgradeID string) error {
	if _, exists := inv.LevelSkills[lvl]; !exists {
		return errors.New("SKILL_NOT_UNLOCKED")
	}

	skill, skillExists := inv.LevelSkills[lvl]

	if !skillExists {
		return errors.New("SKILL_NOT_FOUND")
	}

	upgradeIdx := -1

	for idx, upgrade := range skill.GetUpgrades() {
		if upgrade.Id == upgradeID {
			upgradeIdx = idx
			break
		}
	}

	if upgradeIdx == -1 {
		return errors.New("UPGRADE_NOT_FOUND")
	}

	inv.LevelSkillsUpgrades[lvl] = inv.LevelSkillsUpgrades[lvl] & (1 << upgradeIdx)

	return nil
}

func GetDefaultInventory() PlayerInventory {
	return PlayerInventory{
		Gold:                0,
		Ingredients:         map[uuid.UUID]*types.Ingredient{},
		Items:               []*types.PlayerItem{},
		LevelSkills:         map[int]PlayerSkillLevel{},
		LevelSkillsCDS:      map[int]int{},
		LevelSkillsUpgrades: map[int]int{},
	}
}
