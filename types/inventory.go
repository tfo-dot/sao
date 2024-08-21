package types

import (
	"github.com/disgoorg/disgo/events"
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

	Execute(owner, target, fightInstance, meta interface{}) interface{}
	GetEvents() map[CustomTrigger]func(owner interface{})
}

type PlayerSkillLevel interface {
	PlayerSkill

	GetLevel() int
	GetPath() SkillPath
}

type PlayerSkillUpgradable interface {
	PlayerSkill

	GetUpgradableDescription(upgrades int) string

	CanUse(owner interface{}, fightInstance interface{}, upgrades int) bool
	GetLevel() int
	GetPath() SkillPath
	GetUpgrades() []PlayerSkillUpgrade
	GetCooldown(upgrades int) int
	UpgradableExecute(owner, target, fightInstance, meta interface{}, upgrades int) interface{}
	GetUpgradableTrigger(upgrades int) Trigger
	GetStats(upgrades int) map[Stat]int
	GetUpgradableCost(upgrades int) int
}

type PlayerSkillUpgrade struct {
	Description string
	Id          string
	Events      *map[CustomTrigger]func(owner interface{})
}

type SkillTriggerType int

const (
	TRIGGER_PASSIVE SkillTriggerType = iota
	TRIGGER_ACTIVE
	TRIGGER_TYPE_NONE
)

type Trigger struct {
	Type     SkillTriggerType
	Event    SkillTrigger
	Cooldown *CooldownMeta
	Flags    SkillFlag
	Target   *TargetTrigger
}

type TargetTrigger struct {
	Target     TargetTag
	MaxTargets int
}

type CooldownMeta struct {
	//Default TRIGGER_TURN
	PassEvent SkillTrigger
}

type SkillTrigger int

const (
	TRIGGER_NONE SkillTrigger = iota
	TRIGGER_ATTACK_BEFORE
	TRIGGER_ATTACK_HIT
	TRIGGER_ATTACK_MISS
	TRIGGER_ATTACK_GOT_HIT
	TRIGGER_EXECUTE
	TRIGGER_TURN
	TRIGGER_CAST_ULT
	TRIGGER_DAMAGE_BEFORE
	TRIGGER_DAMAGE
	TRIGGER_HEAL_SELF
	TRIGGER_HEAL_OTHER
	TRIGGER_APPLY_CROWD_CONTROL
)

type IncreasePartial struct {
	Value   int
	Percent bool
}

type DamagePartial struct {
	Value   int
	Percent bool
	Type    DamageType
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
	TARGET_SELF TargetTag = iota << 1
	TARGET_ENEMY
	TARGET_ALLY
)

type SkillFlag int

const (
	FLAG_IGNORE_CC SkillFlag = 1 << iota
	FLAG_INSTANT_SKILL
)

type CustomTrigger int

const (
	CUSTOM_TRIGGER_UNLOCK CustomTrigger = iota
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

type DelayedAction struct {
	Target    uuid.UUID
	Execute   func(owner, fightInstance interface{})
	Turns     int
	PassEvent SkillTrigger
}

type DiscordChoice struct {
	Id     string
	Select func(*events.ComponentInteractionCreate)
}

type EventData struct {
	Source interface{}
	Target interface{}
	Fight  interface{}
}

type BaseAttackIncreaseSkill struct {
	Calculate func(meta AttackTriggerMeta) AttackTriggerMeta
}

func (base BaseAttackIncreaseSkill) Execute(owner, target, fightInstance, meta interface{}) interface{} {
	return base.Calculate(meta.(AttackTriggerMeta))
}

func (base BaseAttackIncreaseSkill) GetCD() int {
	return 0
}

func (base BaseAttackIncreaseSkill) GetCost() int {
	return 0
}

func (base BaseAttackIncreaseSkill) GetDescription() string {
	return "Zwiększa obrażenia ataku"
}

func (base BaseAttackIncreaseSkill) GetEvents() map[CustomTrigger]func(interface{}) {
	return map[CustomTrigger]func(interface{}){}
}

func (base BaseAttackIncreaseSkill) GetName() string {
	return "Template Skill"
}

func (base BaseAttackIncreaseSkill) GetTrigger() Trigger {
	return Trigger{
		Type:  TRIGGER_PASSIVE,
		Event: TRIGGER_ATTACK_BEFORE,
	}
}

func (base BaseAttackIncreaseSkill) GetUUID() uuid.UUID {
	return uuid.New()
}

func (base BaseAttackIncreaseSkill) IsLevelSkill() bool {
	return false
}
