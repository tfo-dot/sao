package battle

import (
	"sao/types"

	"github.com/google/uuid"
)

type FightMessage byte

const (
	MSG_ACTION_NEEDED FightMessage = iota
	MSG_FIGHT_END
)

type Stat int

const (
	STAT_HP Stat = iota
	STAT_SPD
	STAT_DGD
	STAT_AD
	STAT_DEF
	STAT_MR
	STAT_MANA
	STAT_AP
	STAT_HEAL_POWER
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
	ACTION_EFFECT
	ACTION_DODGE
	ACTION_DEFEND
	ACTION_SKILL
)

type Effect int

const (
	EFFECT_POISON Effect = iota
	EFFECT_FEAR
	EFFECT_VAMP
	EFFECT_HEAL
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
	Meta     *map[string]interface{}
}

type ActionEffectHeal struct {
	Value int
}

type ActionEffectStat struct {
	Stat  Stat
	Value int
}

func (d Damage) ToActionMeta() ActionDamage {
	return ActionDamage{
		Damage:   []Damage{d},
		CanDodge: d.CanDodge,
	}
}

func ActionMetaFromList(list []Damage, dodge bool) ActionDamage {
	return ActionDamage{
		Damage:   list,
		CanDodge: dodge,
	}
}

type Damage struct {
	Value    int
	Type     DamageType
	CanDodge bool
}

type SkillTrigger int

const (
	ATTACK SkillTrigger = iota
	DEFEND
	DODGE
	HIT
	FIGHT_START
	FIGHT_END
	EXECUTE
	TURN
	HEALTH
	MANA
)

const SPEED_GAUGE = 100

type Entity interface {
	GetCurrentHP() int
	GetMaxHP() int

	GetSPD() int
	GetATK() int
	GetDEF() int
	GetMR() int
	GetDGD() int
	GetMaxMana() int
	GetCurrentMana() int
	GetAP() int

	Action(*Fight) int
	TakeDMG(ActionDamage) int
	Heal(int)

	GetLoot() []Loot
	CanDodge() bool

	IsAuto() bool
	GetName() string
	GetUUID() uuid.UUID

	ApplyEffect(ActionEffect)
	HasEffect(Effect) bool
	GetEffect(Effect) *ActionEffect
	GetAllEffects() []ActionEffect
	TriggerAllEffects()
}

type PlayerEntity interface {
	ReceiveLoot(Loot)
	ReceiveMultipleLoot([]Loot)

	GetAllSkills() []types.PlayerSkill
}

type DodgeEntity interface {
	TakeDMGOrDodge(ActionDamage) int
}
