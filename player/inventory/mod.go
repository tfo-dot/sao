package inventory

import (
	"errors"
	"sao/data"
	"sao/types"
	"strconv"

	"github.com/google/uuid"
)

type PlayerInventory struct {
	Gold                int
	TempSkills          []*types.WithExpire[types.PlayerSkill]
	Items               []*types.PlayerItem
	ItemSkillCD         map[uuid.UUID]int
	Ingredients         map[uuid.UUID]*types.Ingredient
	LevelSkillsCDS      map[int]int
	LevelSkills         map[int]types.PlayerSkillUpgradable
	LevelChoices        map[int]int
	LevelSkillsUpgrades map[int]int
	FurySkillsCD        map[uuid.UUID]int
}

func (inv *PlayerInventory) AddTempSkill(skill types.WithExpire[types.PlayerSkill]) {
	inv.TempSkills = append(inv.TempSkills, &skill)
}

func (inv *PlayerInventory) Serialize() map[string]interface{} {
	items := make([]map[string]interface{}, 0)

	for _, item := range inv.Items {
		items = append(items, map[string]interface{}{
			"uuid":  item.UUID.String(),
			"count": item.Count,
		})
	}

	itemsCD := make(map[string]interface{})

	for key, value := range inv.ItemSkillCD {
		itemsCD[key.String()] = value
	}

	ingredients := make([]map[string]interface{}, 0)

	for _, ingredient := range inv.Ingredients {
		ingredients = append(ingredients, map[string]interface{}{
			"uuid":  ingredient.UUID.String(),
			"count": ingredient.Count,
		})
	}

	lvlSkills := make(map[int]interface{})

	for key := range inv.LevelSkills {

		choice, notDefaultChoice := inv.LevelChoices[key]

		if !notDefaultChoice {
			choice = 0
		}

		lvlSkills[key] = map[string]interface{}{
			"path":     inv.LevelSkills[key].GetPath(),
			"choice":   choice,
			"upgrades": inv.LevelSkillsUpgrades[key],
		}
	}

	furySkillsCD := make(map[string]interface{})

	for key, value := range inv.FurySkillsCD {
		furySkillsCD[key.String()] = value
	}

	return map[string]interface{}{
		"gold":           inv.Gold,
		"items":          items,
		"itemSkillCD":    itemsCD,
		"ingredients":    ingredients,
		"levelSkillsCDS": inv.LevelSkillsCDS,
		"levelSkills":    lvlSkills,
		"furySkillsCD":   furySkillsCD,
	}
}

func DeserializeInventory(rawData map[string]interface{}) PlayerInventory {
	inv := GetDefaultInventory()

	inv.Gold = int(rawData["gold"].(float64))

	if rawItemData, okay := rawData["items"].([]map[string]interface{}); okay {
		for _, item := range rawItemData {
			uuid, _ := uuid.Parse(item["uuid"].(string))
			copy := data.Items[uuid]

			copy.Count = item["count"].(int)

			inv.Items = append(inv.Items, &copy)
		}
	}

	if rawItemCD, okay := rawData["itemSkillCD"].(map[string]interface{}); okay {
		for key, value := range rawItemCD {
			uuid, _ := uuid.Parse(key)
			inv.ItemSkillCD[uuid] = value.(int)
		}
	}

	if rawIngredientData, okay := rawData["ingredients"].([]map[string]interface{}); okay {
		for _, ingredient := range rawIngredientData {
			uuid, _ := uuid.Parse(ingredient["uuid"].(string))
			copy := data.Ingredients[uuid]

			copy.Count = ingredient["count"].(int)

			inv.Ingredients[uuid] = &copy
		}
	}

	if levelCDData, okay := rawData["levelSkillsCDS"].(map[string]interface{}); okay {
		for key, value := range levelCDData {
			parsed, _ := strconv.Atoi(key)
			inv.LevelSkillsCDS[parsed] = int(value.(float64))
		}
	}

	if _, okay := rawData["levelSkills"].(map[string]interface{}); okay {
		for key, lvlData := range rawData["levelSkills"].(map[string]interface{}) {
			parsed, _ := strconv.Atoi(key)
			data := lvlData.(map[string]interface{})

			inv.LevelSkills[parsed] = AVAILABLE_SKILLS[types.SkillPath(data["path"].(float64))][parsed][int(data["choice"].(float64))]

			inv.LevelSkillsUpgrades[parsed] = int(data["upgrades"].(float64))
		}
	}

	if _, okay := rawData["furySkillsCD"].(map[string]interface{}); okay {
		for key, value := range rawData["furySkillsCD"].(map[string]interface{}) {
			uuid, _ := uuid.Parse(key)
			inv.FurySkillsCD[uuid] = value.(int)
		}
	}

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

	for _, skills := range inv.LevelSkills {
		skillStats := skills.GetStats(inv.LevelSkillsUpgrades[skills.GetLevel()])

		val, exists := skillStats[stat]

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
		TempSkills:          make([]*types.WithExpire[types.PlayerSkill], 0),
		ItemSkillCD:         make(map[uuid.UUID]int),
		Ingredients:         make(map[uuid.UUID]*types.Ingredient),
		Items:               make([]*types.PlayerItem, 0),
		LevelSkills:         make(map[int]types.PlayerSkillUpgradable),
		LevelSkillsCDS:      make(map[int]int),
		LevelSkillsUpgrades: make(map[int]int),
		LevelChoices:        make(map[int]int),
		FurySkillsCD:        make(map[uuid.UUID]int),
	}
}
