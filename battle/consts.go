package battle

import (
	"sao/types"

	"github.com/google/uuid"
)

type FightMessage byte

const (
	MSG_ACTION_NEEDED FightMessage = iota
	MSG_FIGHT_START
	MSG_FIGHT_END
	MSG_ENTITY_RESCUE
	MSG_ENTITY_DIED
)

type LootType int

const (
	LOOT_ITEM LootType = iota
	LOOT_EXP
	LOOT_GOLD
)

type Loot struct {
	Type  LootType
	Count int
	Meta  *LootMeta
}

// Only for items
type LootMeta struct {
	Type types.ItemType
	Uuid uuid.UUID
}

type ActionEnum int

const (
	ACTION_ATTACK ActionEnum = iota
	ACTION_DEFEND
	ACTION_SKILL
	ACTION_ITEM
	ACTION_RUN
	//Helper events
	ACTION_COUNTER
	ACTION_EFFECT
	ACTION_DMG
	ACTION_SUMMON
)

type Effect int

const (
	//TODO Rewrite it to damage over time
	EFFECT_POISON Effect = iota
	EFFECT_HEAL_SELF
	EFFECT_HEAL_OTHER
	//TODO rewrite it to HEAL_SHIELD_POWER, STAT_DEC
	EFFECT_HEAL_REDUCE
	EFFECT_MANA_RESTORE
	EFFECT_SHIELD
	EFFECT_STUN
	EFFECT_STAT_INC
	EFFECT_STAT_DEC
	EFFECT_RESIST
	EFFECT_TAUNT
	EFFECT_TAUNTED
)

type Action struct {
	Event  ActionEnum
	Target uuid.UUID
	Source uuid.UUID
	Meta   any
}

type ActionSummon struct {
	Flags       SummonFlags
	ExpireTimer int
	Entity      Entity
}

type SummonFlags int

const (
	SUMMON_FLAG_NONE SummonFlags = 1 << iota
	SUMMON_FLAG_ATTACK
	SUMMON_FLAG_EXPIRE
)

type ActionDamage struct {
	Damage   []Damage
	CanDodge bool
}

type ActionEffect struct {
	Effect   Effect
	Value    int
	Duration int
	Uuid     uuid.UUID
	Meta     any
	Caster   uuid.UUID
	Target   uuid.UUID
	Source   types.EffectSource
	OnExpire func(owner, fightInstance interface{}, meta ActionEffect)
}

type ActionEffectHeal struct {
	Value int
}

type ActionEffectStat struct {
	Stat      types.Stat
	Value     int
	IsPercent bool
}

type ActionEffectResist struct {
	Value     int
	IsPercent bool
}

type ActionSkillMeta struct {
	Lvl        int
	IsForLevel bool
	SkillUuid  uuid.UUID
	Targets    []uuid.UUID
}

type ActionItemMeta struct {
	Item    uuid.UUID
	Targets []uuid.UUID
}

type Damage struct {
	Value int
	Type  types.DamageType
	//Its ignored when []Damage is of 1
	IsPercent bool
	CanDodge  bool
}

const SPEED_GAUGE = 100

type Entity interface {
	GetCurrentHP() int
	GetCurrentMana() int

	GetStat(types.Stat) int

	Action(*Fight) []Action
	TakeDMG(ActionDamage) []Damage
	DamageShields(int) int

	Heal(int)
	RestoreMana(int)
	Cleanse()

	GetLoot() []Loot
	CanDodge() bool

	GetFlags() types.EntityFlag

	GetName() string
	GetUUID() uuid.UUID

	ApplyEffect(ActionEffect)
	GetEffectByType(Effect) *ActionEffect
	GetEffectByUUID(uuid.UUID) *ActionEffect
	GetAllEffects() []ActionEffect
	RemoveEffect(uuid.UUID)
	TriggerAllEffects() []ActionEffect

	AppendTempSkill(types.WithExpire[types.PlayerSkill])
	TriggerEvent(types.SkillTrigger, interface{}) []interface{}
}

type DodgeEntity interface {
	Entity

	TakeDMGOrDodge(ActionDamage) ([]Damage, bool)
}

type PlayerEntity interface {
	DodgeEntity

	GetUpgrades(int) int
	GetLvlSkill(int) types.PlayerSkill

	SetLvlCD(int, int)
	GetLvlCD(int) int

	SetDefendingState(bool)
	GetDefendingState() bool

	GetAllItems() []*types.PlayerItem
	AddItem(*types.PlayerItem)
	RemoveItem(int)

	AppendDerivedStat(types.DerivedStat)
	SetLevelStat(types.Stat, int)
	GetDefaultStat(types.Stat) int
	ReduceCooldowns(types.SkillTrigger)
}
