package mobs

import (
	"sao/battle"
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
	Effects      EffectList
	CustomAction func(self *SummonEntity, f *battle.Fight) []battle.Action
	OnSummon     func(f *battle.Fight, s *SummonEntity)
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

func (s *SummonEntity) TakeDMG(dmg battle.ActionDamage) []battle.Damage {
	dmgStats := []battle.Damage{
		{Value: 0, Type: types.DMG_PHYSICAL},
		{Value: 0, Type: types.DMG_MAGICAL},
		{Value: 0, Type: types.DMG_TRUE},
	}

	for _, dmg := range dmg.Damage {
		//Skip shield and such
		if dmg.Type == types.DMG_TRUE {
			s.CurrentHP -= dmg.Value
			dmgStats[2].Value += dmg.Value
			continue
		}

		rawDmg := dmg.Value

		switch dmg.Type {
		case types.DMG_PHYSICAL:
			rawDmg = utils.CalcReducedDamage(dmg.Value, s.GetStat(types.STAT_DEF))
		case types.DMG_MAGICAL:
			rawDmg = utils.CalcReducedDamage(dmg.Value, s.GetStat(types.STAT_MR))
		}

		value := s.DamageShields(rawDmg)

		dmgStats[dmg.Type].Value += value

		s.CurrentHP -= value
	}

	return dmgStats
}

func (s *SummonEntity) DamageShields(dmg int) int {
	leftOverDmg := dmg
	idxToRemove := make([]int, 0)

	for idx, effect := range s.Effects {
		if effect.Effect == battle.EFFECT_SHIELD {
			newShieldValue := effect.Value - leftOverDmg

			if newShieldValue <= 0 {
				leftOverDmg = newShieldValue * -1

				idxToRemove = append(idxToRemove, idx)
			} else {
				effect.Value = newShieldValue
				leftOverDmg = 0
			}
		}
	}

	for _, idx := range idxToRemove {
		s.Effects = append(s.Effects[:idx], s.Effects[idx+1:]...)
	}

	return leftOverDmg
}

func (s *SummonEntity) Action(f *battle.Fight) []battle.Action {
	if s.CustomAction != nil {
		return s.CustomAction(s, f)
	}

	enemies := f.GetEnemiesFor(s.UUID)

	if len(enemies) == 0 {
		return []battle.Action{}
	}

	tauntEffect := s.GetEffectByType(battle.EFFECT_TAUNTED)

	if tauntEffect != nil {
		return []battle.Action{
			{
				Event:  battle.ACTION_ATTACK,
				Source: s.UUID,
				Target: tauntEffect.Meta.(uuid.UUID),
			},
		}
	}

	return []battle.Action{
		{
			Event:  battle.ACTION_ATTACK,
			Source: s.UUID,
			Target: utils.RandomElement(enemies).GetUUID(),
		},
	}
}

func (s *SummonEntity) GetUUID() uuid.UUID {
	return s.UUID
}

func (s *SummonEntity) GetLoot() []battle.Loot {
	return []battle.Loot{}
}

func (s *SummonEntity) CanDodge() bool {
	return false
}

func (s *SummonEntity) ApplyEffect(e battle.ActionEffect) {
	s.Effects = append(s.Effects, e)
}

func (s *SummonEntity) GetEffectByType(effect battle.Effect) *battle.ActionEffect {
	return s.Effects.GetEffectByType(effect)
}

func (s *SummonEntity) GetAllEffects() []battle.ActionEffect {
	return s.Effects
}

func (s *SummonEntity) Heal(value int) {
	if s.GetStat(types.STAT_HEAL_POWER) != 0 {
		value = utils.PercentOf(value, 100+s.GetStat(types.STAT_HEAL_POWER))
	}

	s.CurrentHP += value
}

func (s *SummonEntity) Cleanse() {
	s.Effects = s.Effects.Cleanse()
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

func (s *SummonEntity) TriggerAllEffects() []battle.ActionEffect {
	effects, expiredEffects := s.Effects.TriggerAllEffects(s)

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
	s.Effects = s.Effects.RemoveEffect(uuid)
}

func (s *SummonEntity) GetEffectByUUID(uuid uuid.UUID) *battle.ActionEffect {
	return s.Effects.GetEffectByUUID(uuid)
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
		if effect.Effect == battle.EFFECT_STAT_INC {

			if value, ok := effect.Meta.(battle.ActionEffectStat); ok {
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

		if effect.Effect == battle.EFFECT_STAT_DEC {

			if value, ok := effect.Meta.(battle.ActionEffectStat); ok {
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

	switch stat {
	case types.STAT_SPD:
		if stat, exists := s.Stats[types.STAT_SPD]; exists {
			tempValue += stat
		}
	case types.STAT_AGL:
		if stat, exists := s.Stats[types.STAT_AGL]; exists {
			tempValue += stat
		}
	case types.STAT_AD:
		if stat, exists := s.Stats[types.STAT_AD]; exists {
			tempValue += stat
		}
	case types.STAT_DEF:
		if stat, exists := s.Stats[types.STAT_DEF]; exists {
			tempValue += stat
		}
	case types.STAT_MR:
		if stat, exists := s.Stats[types.STAT_MR]; exists {
			tempValue += stat
		}
	case types.STAT_MANA:
		tempValue += 0
	case types.STAT_AP:
		if stat, exists := s.Stats[types.STAT_AP]; exists {
			tempValue += stat
		}
	case types.STAT_HEAL_POWER:
		if stat, exists := s.Stats[types.STAT_HEAL_POWER]; exists {
			tempValue += stat
		}
	case types.STAT_HEAL_SELF:
		if stat, exists := s.Stats[types.STAT_HEAL_SELF]; exists {
			tempValue += stat
		}
	case types.STAT_HP:
		if stat, exists := s.Stats[types.STAT_HP]; exists {
			tempValue += stat
		}
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
