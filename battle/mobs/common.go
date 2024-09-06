package mobs

import (
	"sao/base"
	"sao/battle"
	"sao/types"

	"github.com/google/uuid"
)

type MobEntity struct {
	//ID as mob type
	Id           string
	HP           int
	Effects      []types.ActionEffect
	UUID         uuid.UUID
	Stats        map[types.Stat]int
	Name         string
	Props        map[string]interface{}
	Loot         []types.Loot
	TempSkill    []*types.WithExpire[types.PlayerSkill]
	OnDefeatFunc func(types.PlayerEntity)
	ActionFunc   func(*MobEntity, *battle.Fight) []types.Action
}

func (m *MobEntity) ChangeHP(value int) {
	m.HP += value
}

func (m *MobEntity) GetName() string {
	if m.Name != "" {
		return m.Name
	}

	return "Nieznana istota"
}

func (m *MobEntity) GetCurrentHP() int {
	return m.HP
}

func (m *MobEntity) GetCurrentMana() int {
	return 0
}

func (m *MobEntity) RestoreMana(value int) {}

func (m *MobEntity) GetFlags() types.EntityFlag {
	return types.ENTITY_AUTO
}

func (m *MobEntity) TriggerEvent(trigger types.SkillTrigger, evt types.EventData, target interface{}) []interface{} {
	return []interface{}{}
}

func (m *MobEntity) TakeDMG(dmg types.ActionDamage) []types.Damage {
	return base.TakeDMG(dmg, m)
}

func (m *MobEntity) TakeDMGOrDodge(dmg types.ActionDamage) ([]types.Damage, bool) {
	return base.TakeDMGOrDodge(dmg, m)
}

func (m *MobEntity) DamageShields(dmg int) int {
	_, _, dmgLeft := base.DamageShields(dmg, m)

	return dmgLeft
}

func (m *MobEntity) Action(f types.FightInstance) []types.Action {
	if m.ActionFunc != nil {
		return m.ActionFunc(m, f.(*battle.Fight))
	}

	return m.GetDefaultAction(f.(*battle.Fight))
}

func (m *MobEntity) GetDefaultAction(f types.FightInstance) []types.Action {
	return base.DefaultAction(f, m)
}

func (m *MobEntity) GetUUID() uuid.UUID {
	return m.UUID
}

func (m *MobEntity) GetLoot() []types.Loot {
	return m.Loot
}

func (m *MobEntity) CanDodge() bool {
	return false
}

func (m *MobEntity) ApplyEffect(e types.ActionEffect) {
	m.Effects = append(m.Effects, e)
}

func (m *MobEntity) GetEffectByType(effect types.Effect) *types.ActionEffect {
	return base.GetEffectByType(effect, m)
}

func (m *MobEntity) GetAllEffects() []types.ActionEffect {
	return m.Effects
}

func (m *MobEntity) Heal(value int) {
	base.Heal(m, value)
}

func (m *MobEntity) Cleanse() {
	keepList, _ := base.Cleanse(m)

	m.Effects = keepList
}

func (m *MobEntity) TriggerAllEffects() []types.ActionEffect {
	effects, expiredEffects := base.TriggerAllEffects(m)

	m.Effects = effects

	return expiredEffects
}

func (m *MobEntity) GetSkill(uuid uuid.UUID) types.PlayerSkill {
	for _, skill := range m.TempSkill {
		if skill.Value.GetUUID() == uuid {
			return skill.Value
		}
	}

	return nil
}

func (m *MobEntity) TriggerTempSkills() {
	list := make([]*types.WithExpire[types.PlayerSkill], 0)

	for _, skill := range m.TempSkill {
		if !skill.AfterUsage {
			skill.Expire--

			if skill.Expire > 0 {
				list = append(list, skill)
			} else {
				continue
			}
		}
	}

	m.TempSkill = list
}

func (m *MobEntity) RemoveTempByUUID(uuid uuid.UUID) {
	tempList := make([]*types.WithExpire[types.PlayerSkill], 0)

	for _, skill := range m.TempSkill {
		if skill.Value.GetUUID() != uuid {
			tempList = append(tempList, skill)
		}
	}

	m.TempSkill = tempList
}

func (m *MobEntity) RemoveEffect(uuid uuid.UUID) {
	m.Effects = base.RemoveEffect(uuid, m)
}

func (m *MobEntity) GetEffectByUUID(uuid uuid.UUID) *types.ActionEffect {
	return base.GetEffectByUUID(uuid, m)
}

func (m *MobEntity) GetStat(stat types.Stat) int {
	switch stat {
	case types.STAT_MANA_PLUS:
		return 0
	case types.STAT_HP_PLUS:
		return 0
	}

	statValue := 0
	percentValue := 0

	for _, effect := range m.Effects {
		if effect.Effect == types.EFFECT_STAT_INC {

			if value, ok := effect.Meta.(types.ActionEffectStat); ok {
				if value.Stat != stat {
					continue
				}

				if value.IsPercent {
					percentValue += value.Value
				} else {
					statValue += value.Value
				}
			}
		}

		if effect.Effect == types.EFFECT_STAT_DEC {
			if value, ok := effect.Meta.(types.ActionEffectStat); ok {
				if value.Stat != stat {
					continue
				}

				if value.IsPercent {
					percentValue -= value.Value
				} else {
					statValue -= value.Value
				}
			}
		}
	}

	tempValue := statValue

	if statValue, ok := m.Stats[stat]; ok {
		tempValue += statValue
	}

	return tempValue + (tempValue * percentValue / 100)
}

func (m *MobEntity) AppendTempSkill(skill types.WithExpire[types.PlayerSkill]) {
	m.TempSkill = append(m.TempSkill, &skill)
}

func (m *MobEntity) GetTempSkills() []*types.WithExpire[types.PlayerSkill] {
	return m.TempSkill
}

func (m *MobEntity) HasOnDefeat() bool {
	return m.OnDefeatFunc != nil
}

func (m *MobEntity) OnDefeat(player types.PlayerEntity) {
	if m.OnDefeatFunc != nil {
		m.OnDefeatFunc(player)
	}
}

func Spawn(id string) *MobEntity {
	temp, ok := Mobs[id]

	if !ok {
		return nil
	}

	temp.UUID = uuid.New()

	return &temp
}
