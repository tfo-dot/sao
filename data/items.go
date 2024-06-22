package data

import (
	"sao/battle"
	"sao/types"
	"sao/utils"

	"github.com/google/uuid"
)

var Items = map[uuid.UUID]types.PlayerItem{
	uuid.MustParse("00000000-0000-0000-0000-000000000000"): {
		UUID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		Name:        "Błogosławieństwo Reimi",
		Description: "Przeleczenie daje tarczę.",
		TakesSlot:   true,
		Stacks:      false,
		Consume:     false,
		Count:       1,
		MaxCount:    1,
		Hidden:      false,
		Stats: map[types.Stat]int{
			types.STAT_AD:        25,
			types.STAT_HP:        100,
			types.STAT_LIFESTEAL: 10,
		},
		Effects: []types.PlayerSkill{ReimiBlessingSkill{}},
	},
}

var ReimiBlessing = Items[uuid.MustParse("00000000-0000-0000-0000-000000000000")]
var ReimiBlessingSkillUUID = uuid.MustParse("00000000-0000-0001-0000-100000000001")
var ReimiBlessingEffectUUID = uuid.MustParse("00000000-0000-0001-0001-100000000001")

type ReimiBlessingSkill struct{}

func (rbs ReimiBlessingSkill) GetName() string {
	return "Uświęcona tarcza"
}

func (rbs ReimiBlessingSkill) GetDescription() string {
	return "Przeleczenie daje tarczę."
}

func (rbs ReimiBlessingSkill) GetTrigger() types.Trigger {
	return types.Trigger{
		Type: types.TRIGGER_PASSIVE,
		Event: &types.EventTriggerDetails{
			TriggerType:   types.TRIGGER_HEAL_SELF,
			TargetType:    []types.TargetTag{types.TARGET_SELF},
			TargetDetails: []types.TargetDetails{types.DETAIL_ALL},
			TargetCount:   1,
		},
	}
}

func (rbs ReimiBlessingSkill) GetCD() int {
	return 0
}

func (rbs ReimiBlessingSkill) GetCost() int {
	return 0
}

func (rbs ReimiBlessingSkill) GetUUID() uuid.UUID {
	return ReimiBlessingSkillUUID
}

func (rbs ReimiBlessingSkill) IsLevelSkill() bool {
	return false
}

func (rbs ReimiBlessingSkill) Execute(owner, target interface{}, fightInstance *interface{}, meta interface{}) {
	ownerEntity := owner.(battle.Entity)

	oldEffect := ownerEntity.GetEffectByUUID(ReimiBlessingEffectUUID)

	maxShield := utils.PercentOf(ownerEntity.GetStat(types.STAT_HP), 25) + utils.PercentOf(ownerEntity.GetStat(types.STAT_AD), 25)

	if oldEffect != nil {
		ownerEntity.RemoveEffect(ReimiBlessingEffectUUID)
	} else {
		oldEffect = &battle.ActionEffect{
			Effect:   battle.EFFECT_SHIELD,
			Value:    0,
			Duration: -1,
			Uuid:     ReimiBlessingEffectUUID,
			Meta:     nil,
			Caster:   ownerEntity.GetUUID(),
			Source:   battle.SOURCE_ITEM,
		}
	}

	if oldEffect.Value < 0 {
		oldEffect.Value = 0
	}

	oldEffect.Value += meta.(battle.ActionEffectHeal).Value

	if oldEffect.Value > maxShield {
		oldEffect.Value = maxShield
	}

	ownerEntity.ApplyEffect(*oldEffect)
}

func (rbs ReimiBlessingSkill) GetEvents() map[types.CustomTrigger]func(owner *interface{}) {
	return map[types.CustomTrigger]func(owner *interface{}){
		types.CUSTOM_TRIGGER_UNLOCK: func(owner *interface{}) {
			(*owner).(battle.Entity).ApplyEffect(battle.ActionEffect{
				Effect:   battle.EFFECT_SHIELD,
				Value:    0,
				Duration: -1,
				Uuid:     ReimiBlessingEffectUUID,
				Meta:     nil,
				Caster:   (*owner).(battle.Entity).GetUUID(),
				Source:   battle.SOURCE_ITEM,
			})
		},
	}
}
