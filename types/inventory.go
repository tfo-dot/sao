package types

import "github.com/google/uuid"

type Skill struct {
	Name    string
	Trigger Trigger
	Cost    *SkillCost
	Execute func(source, target interface{}, fight interface{})
}

type PlayerSkill struct {
	Name    string
	Trigger Trigger
	Cost    SkillCost
	UUID    uuid.UUID
	Grade   SkillGrade
	//TODO conditional cd, like start counting after the shield breaks
	//TODO passive skills have cd too
	CD          int
	SkillSource SkillSource
	//TODO remember fight events? Ye
	/* Source, target - entities, fight is fight instance */
	Action func(source, target interface{}, fight interface{})
}

type Resource int

const (
	ManaResource Resource = iota
)

type SkillCost struct {
	Cost     int
	Resource Resource
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

type SkillGrade int

const (
	GradeCommon SkillGrade = iota
	GradeUltimate
)

type SkillSource int

const (
	SourceLVL SkillSource = iota
	SourceItem
	SourceFurry
)

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
	TRIGGER_NONE
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
)
