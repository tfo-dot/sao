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
	Type LootType
	Meta *map[string]interface{}
}

type DamageType int

const (
	DMG_PHYSICAL DamageType = iota
	DMG_MAGICAL
	DMG_TRUE
)

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
)

type Effect int

const (
	EFFECT_POISON Effect = iota
	EFFECT_FEAR
	EFFECT_HEAL_SELF
	EFFECT_HEAL_OTHER
	EFFECT_HEAL_REDUCE
	EFFECT_MANA_RESTORE
	EFFECT_VAMP
	EFFECT_LIFESTEAL
	EFFECT_SHIELD
	EFFECT_BLIND
	EFFECT_DISARM
	EFFECT_GROUND
	EFFECT_ROOT
	EFFECT_SILENCE
	EFFECT_STUN
	EFFECT_IMMUNE
	EFFECT_MARK
	EFFECT_STAT_INC
	EFFECT_STAT_DEC
	EFFECT_RESIST
	EFFECT_FASTEN
	EFFECT_TAUNT
	EFFECT_TAUNTED
	EFFECT_ON_HIT
)

type Action struct {
	Event  ActionEnum
	Target uuid.UUID
	Source uuid.UUID
	Meta   any
}

type ActionPartial struct {
	Event ActionEnum
	Meta  *uuid.UUID
}

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
	Source   EffectSource
}

type EffectSource int

const (
	SOURCE_ND EffectSource = iota
	SOURCE_PARTY
	SOURCE_LOCATION
	SOURCE_ITEM
)

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

type ActionEffectOnHit struct {
	Skill     bool
	Attack    bool
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
	Type  DamageType
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

	IsAuto() bool
	GetName() string
	GetUUID() uuid.UUID

	ApplyEffect(ActionEffect)
	GetEffectByType(Effect) *ActionEffect
	GetEffectByUUID(uuid.UUID) *ActionEffect
	GetAllEffects() []ActionEffect
	RemoveEffect(uuid.UUID)
	TriggerAllEffects() []ActionEffect
}

type DodgeEntity interface {
	Entity

	TakeDMGOrDodge(ActionDamage) ([]Damage, bool)
}

type PlayerEntity interface {
	DodgeEntity

	GetUID() string

	GetAllSkills() []types.PlayerSkill
	GetUpgrades(int) []string
	GetLvlSkill(int) types.PlayerSkill
	GetSkill(uuid.UUID) types.PlayerSkill

	SetLvlCD(int, int)
	GetLvlCD(int) int

	SetCD(uuid.UUID, int)
	GetCD(uuid.UUID) int

	GetLevelSkillsCD() map[int]int
	GetSkillsCD() map[uuid.UUID]int

	SetDefendingState(bool)
	GetDefendingState() bool

	GetAllItems() []*types.PlayerItem
	AddItem(*types.PlayerItem)
	RemoveItem(int)

	GetParty() *uuid.UUID

	AppendDerivedStat(types.DerivedStat)
}

type FightEvent interface {
	GetEvent() FightMessage
	GetData() any
}

type FightStartMsg struct{}

func (fsm FightStartMsg) GetEvent() FightMessage {
	return MSG_FIGHT_START
}

func (fsm FightStartMsg) GetData() any {
	return nil
}

type FightEndMsg struct{}

func (fsm FightEndMsg) GetEvent() FightMessage {
	return MSG_FIGHT_END
}

func (fsm FightEndMsg) GetData() any {
	return nil
}

type FightActionNeededMsg struct {
	Entity uuid.UUID
}

func (fsm FightActionNeededMsg) GetEvent() FightMessage {
	return MSG_ACTION_NEEDED
}

func (fsm FightActionNeededMsg) GetData() any {
	return fsm.Entity
}
