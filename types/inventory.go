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
	Execute(owner, target, fightInstance, meta interface{}) interface{}
	GetEvents() map[CustomTrigger]func(owner interface{})
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
	TRIGGER_TYPE_NONE
)

type Trigger struct {
	Type     SkillTriggerType
	Event    *EventTriggerDetails
	Cooldown *CooldownMeta
	Flags    int
}

type CooldownMeta struct {
	//Default TRIGGER_TURN
	PassEvent SkillTrigger
}

type EventTriggerDetails struct {
	TriggerType   SkillTrigger
	TargetType    []TargetTag
	TargetDetails []TargetDetails
	//If og trigger won't trigger it can work as a fallback i. e. attack miss / attack hit
	OptionalEvent SkillTrigger
	Meta          map[string]interface{}
}

type SkillTrigger int

const (
	TRIGGER_NONE SkillTrigger = iota
	TRIGGER_ATTACK_BEFORE
	TRIGGER_ATTACK_HIT
	TRIGGER_ATTACK_MISS
	TRIGGER_ATTACK_GOT_HIT
	TRIGGER_DEFEND_START
	TRIGGER_DEFEND_END
	TRIGGER_DODGE
	TRIGGER_FIGHT_START
	TRIGGER_FIGHT_END
	TRIGGER_EXECUTE
	TRIGGER_TURN
	TRIGGER_HEALTH
	TRIGGER_MANA
	TRIGGER_COUNTER_ATTEMPT
	TRIGGER_COUNTER_HIT
	TRIGGER_COUNTER_MISS
	TRIGGER_CAST_LVL
	TRIGGER_CAST_ULT
	TRIGGER_DAMAGE_BEFORE
	TRIGGER_DAMAGE_AFTER
	TRIGGER_DAMAGE
	TRIGGER_HEAL_SELF
	TRIGGER_HEAL_OTHER
	TRIGGER_SHIELD_SELF
	TRIGGER_SHIELD_OTHER
	TRIGGER_APPLY_BUFF
	TRIGGER_APPLY_DEBUFF
	TRIGGER_REMOVE_BUFF
	TRIGGER_REMOVE_DEBUFF
	TRIGGER_APPLY_CROWD_CONTROL
	TRIGGER_REMOVE_CROWD_CONTROL
	TRIGGER_APPLY_EFFECT
	TRIGGER_REMOVE_EFFECT
)

type IncreasePartial struct {
	Value   int
	Percent bool
}

type DamagePartial struct {
	Value   int
	Percent bool
	//0 for physical, 1 for magical, 2 for true
	Type int
}

type AttackTriggerMeta struct {
	Effects    []DamagePartial
	ShouldMiss bool
	ShouldHit  bool
}

type DamageTriggerMeta struct {
	Effects []DamagePartial
}

type EffectTriggerMeta struct {
	Effects []IncreasePartial
}

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
	PathSpecial
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

type WithExpire[v any] struct {
	Value      v
	AfterUsage bool
	Expire     int
	//After usage or after turns
	Either bool
}

type WithTarget[v any] struct {
	Value  v
	Target uuid.UUID
}

type DerivedStat struct {
	Base    Stat
	Derived Stat
	Percent int
	Source  uuid.UUID
}

type SkillChoice struct {
	Choice int
}
