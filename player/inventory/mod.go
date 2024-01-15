package inventory

import (
	"sao/battle"
	"sao/types"

	"github.com/google/uuid"
)

type PlayerInventory struct {
	Gold     int
	Capacity int
	Items    []PlayerItem
	CDS      map[uuid.UUID]int
	Skills   []types.PlayerSkill
}

func (inv PlayerInventory) GetStat(stat battle.Stat) int {
	value := 0

	for _, item := range inv.Items {
		val, exists := item.Stats[stat]

		if exists {
			value += val
		}
	}

	return value
}

func (inv PlayerInventory) UseItem(itemUuid uuid.UUID, owner interface{}) {
	for i, item := range inv.Items {
		if item.UUID == itemUuid {
			if item.Consume && item.Count <= 0 {
				return
			}

			if item.Consume {
				item.Count--
			}

			inv.Items[i].UseItem(owner)

			if item.Count == 0 && item.Consume {
				inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			}
		}
	}
}

func GetDefaultInventory() PlayerInventory {
	return PlayerInventory{
		Gold:     0,
		Capacity: 10,
		Items:    []PlayerItem{},
		CDS:      map[uuid.UUID]int{},
		Skills:   []types.PlayerSkill{},
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

type PlayerItem struct {
	UUID        uuid.UUID
	Name        string
	Description string
	TakesSlot   bool
	Stacks      bool
	Consume     bool
	Count       int
	MaxCount    int
	Stats       map[battle.Stat]int
	Effects     []types.Skill
}

func (item *PlayerItem) UseItem(owner interface{}) {

	if item.Count < 0 {
		return
	}

	if item.Consume {
		item.Count--
	}

	for _, effect := range item.Effects {
		if effect.Trigger.Type == types.TRIGGER_PASSIVE {
			continue
		}

		effect.Execute(owner, nil, nil)
	}
}
