package battle

import (
	"sao/types"

	"github.com/google/uuid"
)

type FightMessage byte

const (
	MSG_ACTION_NEEDED FightMessage = iota
	MSG_FIGHT_END
	MSG_ENTITY_DIED
)

type Stat int

const (
	STAT_HP Stat = iota
	STAT_SPD
	STAT_AGL
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
	ACTION_DEFEND
	ACTION_SKILL
	ACTION_ITEM
	//Helper events
	ACTION_EFFECT
	ACTION_DMG
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
	EFFECT_RESIST
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
	GetAGL() int
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

	//TODO Effect stat dec/inc working correctly
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
	SetDefendingState(bool)
	GetDefendingState() bool
}

type DodgeEntity interface {
	TakeDMGOrDodge(ActionDamage) (int, bool)
}

type EntitySort struct {
	Entities []Entity
	Order    []types.TargetDetails
	Meta     *map[string]interface{}
}

func (e EntitySort) Len() int {
	return len(e.Entities)
}

func (e EntitySort) Less(i, j int) bool {
	for _, order := range e.Order {
		result := TargetDetailsCheck(e.Entities[i], e.Entities[j], order, e.Meta)

		if result == 0 {
			continue
		}

		return result > 0
	}

	return false
}

func TargetDetailsCheck(left, right interface{}, order types.TargetDetails, meta *map[string]interface{}) int {
	if order == types.DETAIL_ALL {
		return 0
	}

	switch order {
	case types.DETAIL_MAX_HP:
		return left.(Entity).GetMaxHP() - right.(Entity).GetMaxHP()
	case types.DETAIL_LOW_HP:
		return right.(Entity).GetMaxHP() - left.(Entity).GetMaxHP()
	case types.DETAIL_MAX_MP:
		return left.(Entity).GetMaxMana() - right.(Entity).GetMaxMana()
	case types.DETAIL_LOW_MP:
		return right.(Entity).GetMaxMana() - left.(Entity).GetMaxMana()
	case types.DETAIL_MAX_ATK:
		return left.(Entity).GetATK() - right.(Entity).GetATK()
	case types.DETAIL_LOW_ATK:
		return right.(Entity).GetATK() - left.(Entity).GetATK()
	case types.DETAIL_MAX_DEF:
		return left.(Entity).GetDEF() - right.(Entity).GetDEF()
	case types.DETAIL_LOW_DEF:
		return right.(Entity).GetDEF() - left.(Entity).GetDEF()
	case types.DETAIL_MAX_SPD:
		return left.(Entity).GetSPD() - right.(Entity).GetSPD()
	case types.DETAIL_LOW_SPD:
		return right.(Entity).GetSPD() - left.(Entity).GetSPD()
	case types.DETAIL_MAX_AP:
		return left.(Entity).GetAP() - right.(Entity).GetAP()
	case types.DETAIL_LOW_AP:
		return right.(Entity).GetAP() - left.(Entity).GetAP()
	case types.DETAIL_MAX_RES:
		return left.(Entity).GetMR() - right.(Entity).GetMR()
	case types.DETAIL_LOW_RES:
		return right.(Entity).GetMR() - left.(Entity).GetMR()
	case types.DETAIL_HAS_EFFECT:
		if meta == nil {
			return 0
		}
		leftHas := left.(Entity).HasEffect((*meta)["effect"].(Effect))
		rightHas := right.(Entity).HasEffect((*meta)["effect"].(Effect))

		if leftHas && !rightHas {
			return -1
		}

		if !leftHas && rightHas {
			return 1
		}

		return 0
	case types.DETAIL_NO_EFFECT:
		if meta == nil {
			return 0
		}

		leftHas := left.(Entity).HasEffect((*meta)["effect"].(Effect))
		rightHas := right.(Entity).HasEffect((*meta)["effect"].(Effect))

		if leftHas && !rightHas {
			return 1
		}

		if !leftHas && rightHas {
			return -1
		}

		return 0
	}

	return 0
}

func (e EntitySort) Swap(i, j int) {
	e.Entities[i], e.Entities[j] = e.Entities[j], e.Entities[i]
}
