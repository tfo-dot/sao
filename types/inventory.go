package types

import (
	"github.com/disgoorg/disgo/discord"
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

	Execute(owner PlayerEntity, target Entity, fightInstance FightInstance, meta any) any
	GetEvents() map[CustomTrigger]func(owner PlayerEntity)
}

type PlayerSkillUpgradable interface {
	PlayerSkill

	GetUpgradableDescription(upgrades int) string

	CanUse(owner PlayerEntity, fightInstance FightInstance) bool
	GetLevel() int
	GetPath() SkillPath
	GetUpgrades() []PlayerSkillUpgrade
	GetCooldown(upgrades int) int
	UpgradableExecute(owner PlayerEntity, target Entity, fightInstance FightInstance, meta any) any
	GetUpgradableTrigger(upgrades int) Trigger
	GetStats(upgrades int) map[Stat]int
	GetDerivedStats(upgrades int) []DerivedStat
	GetUpgradableCost(upgrades int) int
}

type PlayerSkillUpgrade struct {
	Description string
	Id          string
	Events      map[CustomTrigger]func(owner PlayerEntity)
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
	TRIGGER_DAMAGE_GOT_HIT
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

type SkillPath int

const (
	PathControl SkillPath = iota
	PathEndurance
	PathDamage
	PathSpecial
)

type PlayerItem struct {
	UUID        uuid.UUID
	Name        string
	Description string
	TakesSlot   bool `parts:"TakesSlot,ignoreEmpty"`
	Stacks      bool `parts:"Stacks,ignoreEmpty"`
	Consume     bool `parts:"Consume,ignoreEmpty"`
	Count       int
	MaxCount    int
	Hidden      bool `parts:"Hidden,ignoreEmpty"`
	Stats       map[Stat]int `parts:"PartsStats,ignoreEmpty"`
	DerivedStats []DerivedStat `parts:"PartsDerivedStats,ignoreEmpty"`
	Effects     []PlayerSkill `parts:"EffectsList,ignoreEmpty"`
}

func (item *PlayerItem) UseItem(owner PlayerEntity, target Entity, fight FightInstance) {
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

type WithCount[T any] struct {
	Item  T
	Count int
}

type WithExpire[v any] struct {
	Value      v
	AfterUsage bool
	Expire     int
	//After usage or after turns
	Either   bool
	OnExpire func(owner PlayerEntity, fight FightInstance)
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

type DiscordChoice struct {
	Id     string
	Select func(*events.ComponentInteractionCreate)
}

type EventData struct {
	Source Entity
	Target Entity
	Fight  FightInstance
}

type FightInstance interface {
	GetEnemiesFor(uuid.UUID) []Entity
	GetAlliesFor(uuid.UUID) []Entity
	GetEntity(uuid.UUID) Entity

	GetChannelId() string

	AppendEventHandler(uuid.UUID, SkillTrigger, func(owner, target Entity, fi FightInstance, meta any) any) uuid.UUID
	RemoveEventHandler(uuid.UUID)

	HandleAction(Action)

	SendMessage(string, discord.MessageCreate, bool)

	CanSummon(uuid.UUID, int) bool
	GetTurnFor(uuid.UUID) int
}
