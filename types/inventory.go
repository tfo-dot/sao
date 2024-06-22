package types

import (
	"github.com/google/uuid"
)

type PlayerSkill interface {
	GetName() string
	GetDescription() string
	GetUUID() uuid.UUID

	GetCD() int
	GetCost() int

	GetTrigger() Trigger

	IsLevelSkill() bool

	//Meta is used for passive events mostly ig
	Execute(owner, target interface{}, fightInstance *interface{}, meta interface{})
	GetEvents() map[CustomTrigger]func(owner *interface{})
}

type PlayerSkillLevel interface {
	PlayerSkill

	GetLevel() int
	GetPath() SkillPath
}

type SkillTriggerType int

const (
	TRIGGER_PASSIVE SkillTriggerType = iota
	TRIGGER_ACTIVE
)

type Trigger struct {
	Type  SkillTriggerType
	Event *EventTriggerDetails
}

type EventTriggerDetails struct {
	TriggerType   SkillTrigger
	TargetType    []TargetTag
	TargetDetails []TargetDetails
	//-1 for no limit
	TargetCount int
	Meta        map[string]interface{}
}

type SkillTrigger int

const (
	TRIGGER_ATTACK SkillTrigger = iota
	TRIGGER_DEFEND
	TRIGGER_DODGE
	TRIGGER_HIT
	TRIGGER_FIGHT_START
	TRIGGER_FIGHT_END
	TRIGGER_EXECUTE
	TRIGGER_TURN
	TRIGGER_HEALTH
	TRIGGER_MANA
	TRIGGER_EFFECT
	TRIGGER_COUNTER
	TRIGGER_CAST
	TRIGGER_DAMAGE
	TRIGGER_NONE
	TRIGGER_HEAL_SELF
	TRIGGER_HEAL_OTHER
)

type TargetTag int

const (
	TARGET_SELF TargetTag = iota
	TARGET_ENEMY
	TARGET_ALLY
)

type TargetDetails int

const (
	DETAIL_LOW_HP TargetDetails = iota
	DETAIL_MAX_HP
	DETAIL_LOW_MP
	DETAIL_MAX_MP
	DETAIL_LOW_ATK
	DETAIL_MAX_ATK
	DETAIL_LOW_DEF
	DETAIL_MAX_DEF
	DETAIL_LOW_SPD
	DETAIL_MAX_SPD
	DETAIL_LOW_AP
	DETAIL_MAX_AP
	DETAIL_LOW_RES
	DETAIL_MAX_RES
	DETAIL_HAS_EFFECT
	DETAIL_NO_EFFECT
	DETAIL_ALL
)

type CustomTrigger int

const (
	CUSTOM_TRIGGER_UNLOCK CustomTrigger = iota
	CUSTOM_TRIGGER_AFTER_EXECUTE
	CUSTOM_TRIGGER_BEFORE_EXECUTE
)

type PlayerItem struct {
	UUID        uuid.UUID
	Name        string
	Description string
	TakesSlot   bool
	Stacks      bool
	Consume     bool
	Count       int
	MaxCount    int
	Hidden      bool
	Stats       map[Stat]int
	Effects     []PlayerSkill
}

type SkillPath int

const (
	PathControl SkillPath = iota
	PathEndurance
	PathDamage
	PathMobility
)

func (item *PlayerItem) UseItem(owner interface{}, target interface{}, fight *interface{}) {

	if item.Count < 0 {
		return
	}

	if item.Consume {
		item.Count--
	}

	for _, effect := range item.Effects {
		if effect.GetTrigger().Type == TRIGGER_PASSIVE {
			continue
		}

		effect.Execute(owner, target, fight, nil)
	}
}

type ItemType int

const (
	ITEM_OTHER ItemType = iota
	ITEM_MATERIAL
)

type Ingredient struct {
	UUID  uuid.UUID
	Name  string
	Stats map[Stat]int
	Count int
}

type Recipe struct {
	UUID        uuid.UUID
	Name        string
	Ingredients []WithCount[uuid.UUID]
	Cost        int
	Product     ResultItem
}

type ResultItem struct {
	UUID  uuid.UUID
	Type  ItemType
	Count int
}

type WithCount[T any] struct {
	Item  T
	Count int
}
