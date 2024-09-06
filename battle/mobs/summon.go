package mobs

import (
	"sao/base"
	"sao/types"
	"sao/utils"

	"github.com/google/uuid"
)

type SummonEntity struct {
	Owner        uuid.UUID
	UUID         uuid.UUID
	Name         string
	Stats        map[types.Stat]int
	CurrentHP    int
	TempSkill    []*types.WithExpire[types.PlayerSkill]
	Effects      []types.ActionEffect
	CustomAction func(self *SummonEntity, f types.FightInstance) []types.Action
	OnSummon     func(f types.FightInstance, s *SummonEntity)
}

func (s *SummonEntity) HasOnDefeat() bool {
	return false
}

func (s *SummonEntity) GetName() string {
	return s.Name
}

func (s *SummonEntity) RestoreMana(value int) {}

func (s *SummonEntity) GetFlags() types.EntityFlag {
	return types.ENTITY_AUTO | types.ENTITY_SUMMON
}

func (s *SummonEntity) TriggerEvent(trigger types.SkillTrigger, evt types.EventData, target interface{}) []interface{} {
	return []interface{}{}
}

func (s *SummonEntity) TakeDMGOrDodge(dmg types.ActionDamage) ([]types.Damage, bool) {
	return base.TakeDMGOrDodge(dmg, s)
}

func (s *SummonEntity) TakeDMG(dmg types.ActionDamage) []types.Damage {
	return base.TakeDMG(dmg, s)
}

func (s *SummonEntity) DamageShields(dmg int) int {
	_, _, damageLeft := base.DamageShields(dmg, s)

	return damageLeft
}

func (s *SummonEntity) Action(f types.FightInstance) []types.Action {
	if s.CustomAction != nil {
		return s.CustomAction(s, f)
	}

	return base.DefaultAction(f, s)
}

func (s *SummonEntity) GetUUID() uuid.UUID {
	return s.UUID
}

func (s *SummonEntity) GetLoot() []types.Loot {
	return []types.Loot{}
}

func (s *SummonEntity) CanDodge() bool {
	return false
}

func (s *SummonEntity) ApplyEffect(e types.ActionEffect) {
	s.Effects = append(s.Effects, e)
}

func (s *SummonEntity) GetEffectByType(effect types.Effect) *types.ActionEffect {
	return base.GetEffectByType(effect, s)
}

func (s *SummonEntity) GetAllEffects() []types.ActionEffect {
	return s.Effects
}

func (s *SummonEntity) Heal(value int) {
	if s.GetStat(types.STAT_HEAL_POWER) != 0 {
		value = utils.PercentOf(value, 100+s.GetStat(types.STAT_HEAL_POWER))
	}

	s.CurrentHP += value
}

func (s *SummonEntity) Cleanse() {
	keepList, _ := base.Cleanse(s)

	s.Effects = keepList
}

func (s *SummonEntity) GetTempSkills() []*types.WithExpire[types.PlayerSkill] {
	return s.TempSkill
}

func (s *SummonEntity) TriggerTempSkills() {
	list := make([]*types.WithExpire[types.PlayerSkill], 0)

	for _, skill := range s.TempSkill {
		if !skill.AfterUsage {
			skill.Expire--

			if skill.Expire > 0 {
				list = append(list, skill)
			} else {
				continue
			}
		}
	}

	s.TempSkill = list
}

func (s *SummonEntity) RemoveTempByUUID(uuid uuid.UUID) {
	tempList := make([]*types.WithExpire[types.PlayerSkill], 0)

	for _, skill := range s.TempSkill {
		if skill.Value.GetUUID() != uuid {
			tempList = append(tempList, skill)
		}
	}

	s.TempSkill = tempList
}

func (s *SummonEntity) TriggerAllEffects() []types.ActionEffect {
	effects, expiredEffects := base.TriggerAllEffects(s)

	s.Effects = effects

	return expiredEffects
}

func (m *SummonEntity) GetSkill(uuid uuid.UUID) types.PlayerSkill {
	for _, skill := range m.TempSkill {
		if skill.Value.GetUUID() == uuid {
			return skill.Value
		}
	}

	return nil
}

func (s *SummonEntity) RemoveEffect(uuid uuid.UUID) {
	s.Effects = base.RemoveEffect(uuid, s)
}

func (s *SummonEntity) GetEffectByUUID(uuid uuid.UUID) *types.ActionEffect {
	return base.GetEffectByUUID(uuid, s)
}

func (s *SummonEntity) GetStat(stat types.Stat) int {
	switch stat {
	case types.STAT_MANA_PLUS:
		return 0
	case types.STAT_HP_PLUS:
		return 0
	}

	statValue := 0
	percentValue := 0

	for _, effect := range s.Effects {
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

	if statValue, exists := s.Stats[stat]; exists {
		tempValue += statValue
	}

	return tempValue + (tempValue * percentValue / 100)
}

func (s *SummonEntity) AppendTempSkill(skill types.WithExpire[types.PlayerSkill]) {
	s.TempSkill = append(s.TempSkill, &skill)
}

func (s *SummonEntity) GetCurrentHP() int {
	return s.CurrentHP
}

func (s *SummonEntity) GetCurrentMana() int {
	return 0
}

func (s *SummonEntity) ChangeHP(value int) {
	s.CurrentHP += value
}
