package mobs

import (
	"fmt"
	"sao/base"
	"sao/battle"
	"sao/types"
	"sao/utils"
	"slices"

	"github.com/google/uuid"
)

type MobEntity struct {
	Id              string
	HP              int
	Effects         []types.ActionEffect
	UUID            uuid.UUID
	Stats           map[types.Stat]int
	Name            string
	Props           map[string]any
	Loot            []types.Loot
	TempSkill       []*types.WithExpire[types.PlayerSkill]
	PartsActionFunc func(...any) (any, error) `parts:"Action,ignoreEmpty"`
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

func (m *MobEntity) TriggerEvent(trigger types.SkillTrigger, evt types.EventData, target any) []any {
	return []any{}
}

func (m *MobEntity) TakeDMG(dmg types.ActionDamage) map[types.DamageType]int {
	return base.TakeDMG(dmg, m)
}

func (m *MobEntity) TakeDMGOrDodge(dmg types.ActionDamage) (map[types.DamageType]int, bool) {
	return base.TakeDMGOrDodge(dmg, m)
}

func (m *MobEntity) DamageShields(dmg int) int {
	leftOverEffects, dmgLeft := base.DamageShields(dmg, m)

	m.Effects = leftOverEffects

	return dmgLeft
}

func (m *MobEntity) Action(f types.FightInstance) []types.Action {
	if m.PartsActionFunc != nil {
		ret, err := m.PartsActionFunc(m, f.(*battle.Fight))

		if err != nil {
			panic(err)
		}

		temp := make([]types.Action, len(ret.([]any)))

		for idx, val := range ret.([]any) {
			if casted, ok := val.(types.Action); !ok {
				mapData := val.(map[string]any)

				act := types.Action{
					Event:  types.ActionEnum(mapData["RTEvent"].(int)),
					Target: *mapData["RTTarget"].(*uuid.UUID),
					Source: *mapData["RTSource"].(*uuid.UUID),
				}

				switch act.Event {
				case types.ACTION_EFFECT:
					metaData := mapData["RTMeta"].(map[string]any)

					effectType := types.Effect(metaData["RTEffect"].(int))

					if effectType != types.EFFECT_DOT {
						panic(fmt.Errorf("unsuported event type %d", effectType))
					}

					act.Meta = types.ActionEffect{
						Effect: effectType, Value: metaData["RTValue"].(int), Duration: metaData["RTDuration"].(int),
					}
				}

				temp[idx] = act
			} else {
				temp[idx] = casted
			}

		}

		return temp
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
	if effect := m.GetEffectByType(types.EFFECT_LOOT_INCREASE); effect != nil {
		for i := range m.Loot {
			m.Loot[i].Count = utils.PercentOf(m.Loot[i].Count, 100+effect.Value)
		}

		return m.Loot
	}

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
	m.Effects = base.Cleanse(m)
}

func (m *MobEntity) TriggerAllEffects() {
	effects := base.TriggerAllEffects(m)

	m.Effects = effects
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
	for idx, skill := range m.TempSkill {
		if skill.Value.GetUUID() != uuid {
			slices.Delete(m.TempSkill, idx, idx+1)
		}
	}
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
					percentValue += effect.Value
				} else {
					statValue += effect.Value
				}
			}
		}

		if effect.Effect == types.EFFECT_STAT_DEC {
			if value, ok := effect.Meta.(types.ActionEffectStat); ok {
				if value.Stat != stat {
					continue
				}

				if value.IsPercent {
					percentValue -= effect.Value
				} else {
					statValue -= effect.Value
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

func Spawn(id string) *MobEntity {
	temp, ok := Mobs[id]

	if !ok {
		return nil
	}

	temp.UUID = uuid.New()

	return &temp
}
