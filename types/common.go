package types

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/google/uuid"
)

type DiscordMessageStruct struct {
	ChannelID      string
	MessageContent discord.MessageCreate
	DM             bool
}

type DiscordEvent interface {
	GetEvent() DiscordMessage
	GetData() any
}

type DiscordMessage int

const (
	MSG_SEND DiscordMessage = iota
	MSG_CHOICE
)

type DiscordSendMsg struct {
	Data DiscordMessageStruct
}

func (fsm DiscordSendMsg) GetEvent() DiscordMessage {
	return MSG_SEND
}

func (fsm DiscordSendMsg) GetData() any {
	return fsm.Data
}

type DiscordChoiceMsg struct {
	Data DiscordChoice
}

func (fsm DiscordChoiceMsg) GetEvent() DiscordMessage {
	return MSG_CHOICE
}

func (fsm DiscordChoiceMsg) GetData() any {
	return fsm.Data
}

type Stat int

const (
	STAT_NONE Stat = iota
	STAT_HP
	STAT_HP_PLUS
	STAT_SPD
	STAT_AGL
	STAT_AD
	STAT_DEF
	STAT_MR
	STAT_MANA
	STAT_MANA_PLUS
	STAT_AP
	STAT_HEAL_SELF
	STAT_HEAL_POWER
	STAT_LETHAL
	STAT_LETHAL_PERCENT
	STAT_MAGIC_PEN
	STAT_MAGIC_PEN_PERCENT
	STAT_ADAPTIVE
	STAT_ADAPTIVE_PERCENT
	STAT_OMNI_VAMP
	STAT_ATK_VAMP
)

var StatToString = map[Stat]string{
	STAT_HP:                "HP",
	STAT_SPD:               "SPD",
	STAT_AGL:               "AGL",
	STAT_AD:                "AD",
	STAT_DEF:               "DEF",
	STAT_MR:                "RES",
	STAT_MANA:              "MANA",
	STAT_AP:                "AP",
	STAT_HEAL_SELF:         "Otrzymane leczenie",
	STAT_HEAL_POWER:        "Moc leczenia",
	STAT_LETHAL:            "Przebicie pancerza",
	STAT_LETHAL_PERCENT:    "Przebicie procentowe pancerza",
	STAT_MAGIC_PEN:         "Przebicie odporności na magię",
	STAT_MAGIC_PEN_PERCENT: "Przebicie procentowe odporności na magię",
	STAT_ADAPTIVE:          "Siła adaptacyjna",
}

type AdaptiveAttackType int

const (
	ADAPTIVE_ATK AdaptiveAttackType = iota
	ADAPTIVE_AP
)

type AdaptiveDefenseType int

const (
	ADAPTIVE_DEF AdaptiveDefenseType = iota
	ADAPTIVE_RES
)

type EffectSource int

const (
	SOURCE_ND EffectSource = iota
	SOURCE_PARTY
	SOURCE_LOCATION
	SOURCE_ITEM
)

type EntityFlag int

const (
	ENTITY_AUTO EntityFlag = 1 << iota
	ENTITY_SUMMON
)

var PathToString = map[SkillPath]string{
	PathControl:   "Kontrola",
	PathDamage:    "Obrażenia",
	PathEndurance: "Wytrzymałość",
	PathSpecial:   "Specjalista",
}

type DamageType int

const (
	DMG_PHYSICAL DamageType = iota
	DMG_MAGICAL
	DMG_TRUE
)

type EntityLocation struct {
	Floor    string
	Location string
}

type ActionDamage struct {
	Damage   []Damage
	CanDodge bool
}

type Damage struct {
	Value int
	Type  DamageType
	//Its ignored when []Damage is of 1
	IsPercent bool
	CanDodge  bool
}

type ActionEffect struct {
	Effect   Effect
	Value    int
	Duration int
	Uuid     uuid.UUID
	Meta     any
	Caster   uuid.UUID
	Target   uuid.UUID
	Source   EffectSource
	OnExpire func(owner Entity, fightInstance FightInstance, meta ActionEffect)
}

type Effect int

const (
	EFFECT_DOT Effect = iota
	EFFECT_HEAL
	EFFECT_MANA_RESTORE
	EFFECT_SHIELD
	EFFECT_STUN
	EFFECT_STAT_INC
	EFFECT_STAT_DEC
	EFFECT_RESIST
	EFFECT_TAUNT
	EFFECT_TAUNTED
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
	Type ItemType
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

type Action struct {
	Event       ActionEnum
	Target      uuid.UUID
	Source      uuid.UUID
	ConsumeTurn *bool
	Meta        any
}

type ActionSummon struct {
	Flags       SummonFlags
	ExpireTimer int
	//For max count
	EntityType uuid.UUID
	Entity     Entity
}

type SummonFlags int

const (
	SUMMON_FLAG_NONE SummonFlags = 1 << iota
	SUMMON_FLAG_ATTACK
	SUMMON_FLAG_EXPIRE
)

func HasFlag[T ~int](flags T, flag T) bool {
	return flags&flag != 0
}

type ActionEffectHeal struct {
	Value int
}

type ActionEffectStat struct {
	Stat      Stat
	Value     int
	IsPercent bool
}

type ActionEffectResist struct {
	All       bool
	Value     int
	IsPercent bool
	//4 stands for all
	DmgType int
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

type Entity interface {
	GetCurrentHP() int
	GetCurrentMana() int

	GetStat(Stat) int

	Action(FightInstance) []Action
	TakeDMG(ActionDamage) []Damage
	DamageShields(int) int

	Heal(int)
	RestoreMana(int)
	Cleanse()

	GetLoot() []Loot
	CanDodge() bool

	GetFlags() EntityFlag

	GetName() string
	GetUUID() uuid.UUID

	ApplyEffect(ActionEffect)
	GetEffectByType(Effect) *ActionEffect
	GetEffectByUUID(uuid.UUID) *ActionEffect
	GetSkill(uuid.UUID) PlayerSkill
	GetAllEffects() []ActionEffect
	RemoveEffect(uuid.UUID)
	TriggerAllEffects() []ActionEffect

	AppendTempSkill(WithExpire[PlayerSkill])
	GetTempSkills() []*WithExpire[PlayerSkill]
	RemoveTempByUUID(uuid.UUID)
	TriggerTempSkills()
	TriggerEvent(SkillTrigger, EventData, interface{}) []interface{}

	HasOnDefeat() bool

	ChangeHP(int)
	TakeDMGOrDodge(ActionDamage) ([]Damage, bool)
}

type MobEntity interface {
	Entity

	GetDefaultAction(FightInstance) []Action
}

type DefeatableEntity interface {
	Entity

	OnDefeat(PlayerEntity)
}

type PlayerEntity interface {
	Entity

	ClearFight()

	GetUpgrades(int) int
	GetLvlSkill(int) PlayerSkill

	SetLvlCD(int, int)
	GetLvlCD(int) int

	SetDefendingState(bool)
	GetDefendingState() bool

	GetAllItems() []*PlayerItem
	AddItem(*PlayerItem)
	RemoveItem(int)

	GetLvl() int
	GetSkills() []PlayerSkill

	AppendDerivedStat(DerivedStat)
	SetLevelStat(Stat, int)
	GetLevelStat(Stat) int
	GetDefaultStat(Stat) int
	ReduceCooldowns(SkillTrigger)

	SetLevelSkillMeta(int, interface{})
	GetLevelSkillMeta(int) interface{}

	UnlockFloor(string)
	UseItem(uuid.UUID, Entity, FightInstance)
}

type NPCStore struct {
	Uuid     uuid.UUID
	Name     string
	Location EntityLocation
	Stock    []*Stock
}

type Stock struct {
	ItemType ItemType
	ItemUUID uuid.UUID
	Price    int
}
