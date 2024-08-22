package mobs

import (
	"sao/battle"
	"sao/types"
	"sao/utils"

	"github.com/google/uuid"
)

type MobEntity struct {
	//ID as mob type
	Id        string
	MaxHP     int
	HP        int
	SPD       int
	ATK       int
	Effects   EffectList
	UUID      uuid.UUID
	Props     map[string]interface{}
	Loot      []battle.Loot
	TempSkill []*types.WithExpire[types.PlayerSkill]
}

type EffectList []battle.ActionEffect

func (e EffectList) GetEffectByType(effect battle.Effect) *battle.ActionEffect {
	for _, eff := range e {
		if eff.Effect == effect {
			return &eff
		}
	}

	return nil
}

func (e EffectList) GetEffectByUUID(eUuid uuid.UUID) *battle.ActionEffect {
	for _, eff := range e {
		if eff.Uuid == eUuid {
			return &eff
		}
	}

	return nil
}

func (e EffectList) RemoveEffect(eUuid uuid.UUID) EffectList {
	tempList := make([]battle.ActionEffect, 0)

	for _, eff := range e {
		if eff.Uuid != eUuid {
			tempList = append(tempList, eff)
		}
	}

	return tempList
}

func (e EffectList) TriggerAllEffects(en battle.Entity) (EffectList, EffectList) {
	effects := make([]battle.ActionEffect, 0)
	expiredEffects := make([]battle.ActionEffect, 0)

	for _, effect := range e {
		if effect.Duration > 0 {
			effect.Duration--
		}

		switch effect.Effect {
		case battle.EFFECT_DOT:
			en.TakeDMG(battle.ActionDamage{
				Damage:   []battle.Damage{{Value: effect.Value, Type: types.DMG_TRUE, CanDodge: false}},
				CanDodge: false,
			})
		case battle.EFFECT_HEAL_SELF:
			en.Heal(effect.Value)
		case battle.EFFECT_MANA_RESTORE:
			en.RestoreMana(effect.Value)
		}

		if effect.Duration == 0 {
			expiredEffects = append(expiredEffects, effect)
		} else {
			effects = append(effects, effect)
		}
	}

	return effects, expiredEffects
}

func (e EffectList) Cleanse() EffectList {
	tempList := make([]battle.ActionEffect, 0)

	for _, effect := range e {
		switch effect.Effect {
		case battle.EFFECT_DOT:
			continue
		case battle.EFFECT_STUN:
			continue
		case battle.EFFECT_TAUNTED:
			continue
		case battle.EFFECT_STAT_DEC:
			continue
		}

		tempList = append(tempList, effect)
	}

	return tempList
}

func (m *MobEntity) GetName() string {
	switch m.Id {
	case "LV0_Rycerz":
		return "Rycerz"
	case "LV0_Boss":
		return "Boss"
	case "LV0_Wilk":
		return "Wilk"
	}

	return "Nieznana istota"
}

func (m *MobEntity) GetCurrentHP() int {
	return m.HP
}

func (m *MobEntity) GetMaxHP() int {
	return m.MaxHP
}

func (m *MobEntity) GetATK() int {
	return m.ATK
}

func (m *MobEntity) GetSPD() int {
	return m.SPD
}

func (m *MobEntity) GetDEF() int {
	return 0
}

func (m *MobEntity) GetMR() int {
	return 0
}

func (m *MobEntity) GetAGL() int {
	return 0
}

func (m *MobEntity) GetMaxMana() int {
	return 0
}

func (m *MobEntity) GetCurrentMana() int {
	return 0
}

func (m *MobEntity) RestoreMana(value int) {}

func (m *MobEntity) GetAP() int {
	return 0
}

func (m *MobEntity) GetFlags() types.EntityFlag {
	return types.ENTITY_AUTO
}

func (m *MobEntity) TriggerEvent(trigger types.SkillTrigger, evt types.EventData, target interface{}) []interface{} {
	return []interface{}{}
}

func (m *MobEntity) TakeDMG(dmg battle.ActionDamage) []battle.Damage {
	dmgStats := []battle.Damage{
		{Value: 0, Type: types.DMG_PHYSICAL},
		{Value: 0, Type: types.DMG_MAGICAL},
		{Value: 0, Type: types.DMG_TRUE},
	}

	for _, dmg := range dmg.Damage {
		//Skip shield and such
		if dmg.Type == types.DMG_TRUE {
			m.HP -= dmg.Value
			dmgStats[2].Value += dmg.Value
			continue
		}

		rawDmg := dmg.Value

		switch dmg.Type {
		case types.DMG_PHYSICAL:
			rawDmg = utils.CalcReducedDamage(dmg.Value, m.GetDEF())
		case types.DMG_MAGICAL:
			rawDmg = utils.CalcReducedDamage(dmg.Value, m.GetMR())
		}

		value := m.DamageShields(rawDmg)

		dmgStats[dmg.Type].Value += value

		m.HP -= value
	}

	return dmgStats
}

func (m *MobEntity) DamageShields(dmg int) int {
	leftOverDmg := dmg
	idxToRemove := make([]int, 0)

	for idx, effect := range m.Effects {
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
		m.Effects = append(m.Effects[:idx], m.Effects[idx+1:]...)
	}

	return leftOverDmg
}

func (m *MobEntity) Action(f *battle.Fight) []battle.Action {
	enemies := f.GetEnemiesFor(m.UUID)

	if len(enemies) == 0 {
		return []battle.Action{}
	}

	tauntEffect := m.GetEffectByType(battle.EFFECT_TAUNTED)

	if tauntEffect != nil {
		return []battle.Action{
			{
				Event:  battle.ACTION_ATTACK,
				Source: m.UUID,
				Target: tauntEffect.Meta.(uuid.UUID),
			},
		}
	}

	return []battle.Action{
		{
			Event:  battle.ACTION_ATTACK,
			Source: m.UUID,
			Target: utils.RandomElement(enemies).GetUUID(),
		},
	}
}

func (m *MobEntity) GetUUID() uuid.UUID {
	return m.UUID
}

func (m *MobEntity) GetLoot() []battle.Loot {
	return m.Loot
}

func (m *MobEntity) CanDodge() bool {
	return false
}

func (m *MobEntity) ApplyEffect(e battle.ActionEffect) {
	m.Effects = append(m.Effects, e)
}

func (m *MobEntity) GetEffectByType(effect battle.Effect) *battle.ActionEffect {
	return m.Effects.GetEffectByType(effect)
}

func (m *MobEntity) GetAllEffects() []battle.ActionEffect {
	return m.Effects
}

func (m *MobEntity) Heal(value int) {
	if m.GetStat(types.STAT_HEAL_POWER) != 0 {
		value = utils.PercentOf(value, 100+m.GetStat(types.STAT_HEAL_POWER))
	}

	m.HP += value
}

func (m *MobEntity) Cleanse() {
	m.Effects = m.Effects.Cleanse()
}

func (m *MobEntity) TriggerAllEffects() []battle.ActionEffect {
	effects, expiredEffects := m.Effects.TriggerAllEffects(m)

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
	m.Effects = m.Effects.RemoveEffect(uuid)
}

func (m *MobEntity) GetEffectByUUID(uuid uuid.UUID) *battle.ActionEffect {
	return m.Effects.GetEffectByUUID(uuid)
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
		tempValue += m.SPD
	case types.STAT_AGL:
		tempValue += m.GetAGL()
	case types.STAT_AD:
		tempValue += m.GetATK()
	case types.STAT_DEF:
		tempValue += m.GetDEF()
	case types.STAT_MR:
		tempValue += m.GetMR()
	case types.STAT_MANA:
		tempValue += m.GetCurrentMana()
	case types.STAT_AP:
		tempValue += m.GetAP()
	case types.STAT_HEAL_POWER:
		tempValue += 0
	case types.STAT_HP:
		tempValue += m.GetMaxHP()
	}

	return tempValue + (tempValue * percentValue / 100)
}

func (m *MobEntity) AppendTempSkill(skill types.WithExpire[types.PlayerSkill]) {
	m.TempSkill = append(m.TempSkill, &skill)
}

func (m *MobEntity) GetTempSkills() []*types.WithExpire[types.PlayerSkill] {
	return m.TempSkill
}

func Spawn(id string) *MobEntity {

	switch id {
	case "LV0_Rycerz":
		return &MobEntity{
			Id:        id,
			MaxHP:     90,
			HP:        90,
			SPD:       40,
			ATK:       25,
			Effects:   make(EffectList, 0),
			UUID:      uuid.New(),
			Props:     make(map[string]interface{}, 0),
			Loot:      []battle.Loot{{Type: battle.LOOT_EXP, Count: 55}},
			TempSkill: make([]*types.WithExpire[types.PlayerSkill], 0),
		}
	case "LV0_Wilk":
		return &MobEntity{
			Id:        id,
			MaxHP:     90,
			HP:        90,
			SPD:       40,
			ATK:       40,
			Effects:   make(EffectList, 0),
			UUID:      uuid.New(),
			Props:     make(map[string]interface{}, 0),
			Loot:      []battle.Loot{{Type: battle.LOOT_EXP, Count: 55}},
			TempSkill: make([]*types.WithExpire[types.PlayerSkill], 0),
		}
	case "LV0_Boss":
		return &MobEntity{
			Id:        id,
			MaxHP:     500,
			HP:        500,
			SPD:       30,
			ATK:       100,
			Effects:   make(EffectList, 0),
			UUID:      uuid.New(),
			Props:     make(map[string]interface{}, 0),
			Loot:      []battle.Loot{{Type: battle.LOOT_EXP, Count: 55}},
			TempSkill: make([]*types.WithExpire[types.PlayerSkill], 0),
		}
	}

	return nil
}
