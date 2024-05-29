package mobs

import (
	"sao/battle"
	"sao/types"
	"sao/utils"

	"github.com/google/uuid"
)

type MobEntity struct {
	//ID as mob type
	Id      string
	MaxHP   int
	HP      int
	SPD     int
	ATK     int
	Effects EffectList
	UUID    uuid.UUID
	Props   map[string]interface{}
	Loot    []battle.Loot
}

type EffectList []battle.ActionEffect

func (e EffectList) GetEffect(effect battle.Effect) *battle.ActionEffect {
	for _, eff := range e {
		if eff.Effect == effect {
			return &eff
		}
	}

	return nil
}

func (e EffectList) TriggerAllEffects(en battle.Entity) (EffectList, EffectList) {
	effects := make([]battle.ActionEffect, 0)
	expiredEffects := make([]battle.ActionEffect, 0)

	for _, effect := range e {
		if effect.Duration > 0 {
			effect.Duration--
		}

		switch effect.Effect {
		case battle.EFFECT_POISON:
			en.TakeDMG(battle.ActionDamage{
				Damage:   []battle.Damage{{Value: effect.Value, Type: battle.DMG_TRUE, CanDodge: false}},
				CanDodge: false,
			})
		case battle.EFFECT_HEAL:
			if en.GetStat(types.STAT_HEAL_SELF) != 0 {
				en.Heal(utils.PercentOf(en.GetStat(types.STAT_HEAL_SELF), 100+effect.Value))
			} else {
				en.Heal(effect.Value)
			}
		case battle.EFFECT_MANA:
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

func (e EffectList) Cleanse() {
	tempList := make([]battle.ActionEffect, 0)

	for _, effect := range e {
		switch effect.Effect {
		case battle.EFFECT_POISON:
			continue
		case battle.EFFECT_BLIND:
			continue
		case battle.EFFECT_DISARM:
			continue
		case battle.EFFECT_FEAR:
			continue
		case battle.EFFECT_GROUND:
			continue
		case battle.EFFECT_ROOT:
			continue
		case battle.EFFECT_SILENCE:
			continue
		case battle.EFFECT_STUN:
			continue
		}

		tempList = append(tempList, effect)
	}

	e = tempList
}

func (m *MobEntity) GetName() string {
	switch m.Id {
	case "LV0_Rycerz":
		return "Rycerz"
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

func (m *MobEntity) RestoreMana(value int) {
	//No mana XD
}

func (m *MobEntity) GetAP() int {
	return 0
}

func (m *MobEntity) IsAuto() bool {
	return true
}

func (m *MobEntity) TakeDMG(dmg battle.ActionDamage) int {
	startingHP := m.HP

	for _, dmg := range dmg.Damage {
		//Skip shield and such
		if dmg.Type == battle.DMG_TRUE {
			m.HP -= dmg.Value
			continue
		}

		rawDmg := dmg.Value

		switch dmg.Type {
		case battle.DMG_PHYSICAL:
			rawDmg = utils.CalcReducedDamage(dmg.Value, m.GetDEF())
		case battle.DMG_MAGICAL:
			rawDmg = utils.CalcReducedDamage(dmg.Value, m.GetMR())
		}

		m.HP -= m.DamageShields(rawDmg)
	}

	return startingHP - m.HP
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

	tauntEffect := m.GetEffect(battle.EFFECT_TAUNTED)

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

func (m *MobEntity) GetEffect(effect battle.Effect) *battle.ActionEffect {
	return m.Effects.GetEffect(effect)
}

func (m *MobEntity) GetAllEffects() []battle.ActionEffect {
	return m.Effects
}

func (m *MobEntity) Heal(value int) {
	m.HP += value
}

func (m *MobEntity) Cleanse() {
	m.Effects.Cleanse()
}

func (m *MobEntity) TriggerAllEffects() []battle.ActionEffect {
	effects, expiredEffects := m.Effects.TriggerAllEffects(m)

	m.Effects = effects

	return expiredEffects
}

func (m *MobEntity) GetStat(stat types.Stat) int {
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
	}

	return tempValue + (tempValue * percentValue / 100)
}

func Spawn(id string) *MobEntity {

	switch id {
	case "LV0_Rycerz":
		return &MobEntity{
			Id:      id,
			MaxHP:   90,
			HP:      90,
			SPD:     40,
			ATK:     25,
			Effects: make(EffectList, 0),
			UUID:    uuid.New(),
			Props:   make(map[string]interface{}, 0),
			Loot: []battle.Loot{{
				Type: battle.LOOT_EXP, Meta: &map[string]interface{}{"value": 55},
			}},
		}
	}

	return nil
}
