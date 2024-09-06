package lua

import (
	"fmt"
	"sao/battle"
	"sao/types"
	"sao/utils"

	"github.com/Shopify/go-lua"
	"github.com/google/uuid"
)

func AddStatTypes(state *lua.State) {
	state.NewTable()

	state.PushInteger(int(types.STAT_NONE))
	state.SetField(-2, "STAT_NONE")

	state.PushInteger(int(types.STAT_HP))
	state.SetField(-2, "STAT_HP")

	state.PushInteger(int(types.STAT_HP_PLUS))
	state.SetField(-2, "STAT_HP_PLUS")

	state.PushInteger(int(types.STAT_SPD))
	state.SetField(-2, "STAT_SPD")

	state.PushInteger(int(types.STAT_AGL))
	state.SetField(-2, "STAT_AGL")

	state.PushInteger(int(types.STAT_AD))
	state.SetField(-2, "STAT_AD")

	state.PushInteger(int(types.STAT_DEF))
	state.SetField(-2, "STAT_DEF")

	state.PushInteger(int(types.STAT_MR))
	state.SetField(-2, "STAT_MR")

	state.PushInteger(int(types.STAT_MANA))
	state.SetField(-2, "STAT_MANA")

	state.PushInteger(int(types.STAT_MANA_PLUS))
	state.SetField(-2, "STAT_MANA_PLUS")

	state.PushInteger(int(types.STAT_AP))
	state.SetField(-2, "STAT_AP")

	state.PushInteger(int(types.STAT_HEAL_SELF))
	state.SetField(-2, "STAT_HEAL_SELF")

	state.PushInteger(int(types.STAT_HEAL_POWER))
	state.SetField(-2, "STAT_HEAL_POWER")

	state.PushInteger(int(types.STAT_LETHAL))
	state.SetField(-2, "STAT_LETHAL")

	state.PushInteger(int(types.STAT_LETHAL_PERCENT))
	state.SetField(-2, "STAT_LETHAL_PERCENT")

	state.PushInteger(int(types.STAT_MAGIC_PEN))
	state.SetField(-2, "STAT_MAGIC_PEN")

	state.PushInteger(int(types.STAT_MAGIC_PEN_PERCENT))
	state.SetField(-2, "STAT_MAGIC_PEN_PERCENT")

	state.PushInteger(int(types.STAT_ADAPTIVE))
	state.SetField(-2, "STAT_ADAPTIVE")

	state.PushInteger(int(types.STAT_ADAPTIVE_PERCENT))
	state.SetField(-2, "STAT_ADAPTIVE_PERCENT")

	state.PushInteger(int(types.STAT_OMNI_VAMP))
	state.SetField(-2, "STAT_OMNI_VAMP")

	state.PushInteger(int(types.STAT_ATK_VAMP))
	state.SetField(-2, "STAT_ATK_VAMP")

	state.SetGlobal("StatsConst")
}

func AddPlayerFunctions(state *lua.State) {
	state.PushGoFunction(func(state *lua.State) int {
		player := state.ToUserData(1).(types.PlayerEntity)

		rawMap, err := utils.GetTableAsMap(state)

		if err != nil {
			panic(err)
		}

		player.AppendDerivedStat(types.DerivedStat{
			Base:    types.Stat(rawMap["Base"].(float64)),
			Derived: types.Stat(rawMap["Derived"].(float64)),
			Percent: int(rawMap["Percent"].(float64)),
			Source:  uuid.MustParse(rawMap["Source"].(string)),
		})

		return 0
	})

	state.SetGlobal("AppendDerivedStat")

	state.PushGoFunction(func(state *lua.State) int {
		player := state.ToUserData(1).(types.PlayerEntity)
		stat, _ := state.ToInteger(2)

		state.PushInteger(player.GetStat(types.Stat(stat)))

		return 1
	})

	state.SetGlobal("GetStat")

	state.PushGoFunction(func(state *lua.State) int {
		player := state.ToUserData(1).(types.PlayerEntity)

		dataMap, err := utils.GetTableAsMap(state)

		if err != nil {
			panic(err)
		}

		var exec func(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{}
		var trigger types.Trigger

		if value, existsValue := dataMap["Value"]; existsValue {
			value := value.(map[string]interface{})

			trigger = ReadMapAsTrigger(value["Trigger"].(map[string]interface{}))

			if execFuncRef, execExists := value["Execute"]; execExists {
				if execFunc, ok := execFuncRef.(utils.LuaFunctionRef); ok {
					exec = func(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
						state.Global(execFunc.FunctionName)

						state.PushUserData(owner)
						state.PushUserData(target)
						state.PushUserData(fightInstance)
						state.PushUserData(meta)

						state.Call(4, 1)

						if !state.IsNil(-1) {
							rValue, err := utils.GetTableAsMap(state)

							if err != nil {
								panic(err)
							}

							fmt.Println(ParseReturnMeta(rValue, trigger))

							return ParseReturnMeta(rValue, trigger)
						}

						return nil
					}
				}
			} else {
				panic("Execute function wasn't found")
			}
		}

		player.AppendTempSkill(types.WithExpire[types.PlayerSkill]{
			AfterUsage: dataMap["AfterUsage"].(bool),
			Expire:     int(dataMap["Expire"].(float64)),
			Either:     dataMap["Either"].(bool),
			Value:      SimplePlayerSkill{Trigger: trigger, Exec: exec},
		})

		return 0
	})

	state.SetGlobal("AppendTempSkill")

	state.PushGoFunction(func(state *lua.State) int {
		player, _ := state.ToUserData(1).(types.PlayerEntity)
		floor, _ := state.ToString(2)

		player.UnlockFloor(floor)

		return 0
	})

	state.SetGlobal("UnlockFloor")
}

func AddFightFunctions(state *lua.State) {
	state.PushGoFunction(func(state *lua.State) int {
		fight := state.ToUserData(1).(*battle.Fight)
		dataMap, err := utils.GetTableAsMap(state)

		if err != nil {
			panic(err)
		}

		var consumeTurn bool

		if value, existsValue := dataMap["ConsumeTurn"]; existsValue {
			consumeTurn = value.(bool)
		}

		act := types.Action{
			Event:       StringToActionEvent[dataMap["Event"].(string)],
			Source:      uuid.MustParse(dataMap["Source"].(string)),
			Target:      uuid.MustParse(dataMap["Target"].(string)),
			ConsumeTurn: &consumeTurn,
		}

		switch act.Event {
		case types.ACTION_SUMMON:
			panic("Summon not implemented")
		case types.ACTION_DMG:
			localMeta := dataMap["Meta"].(map[string]interface{})

			dmg := types.ActionDamage{
				Damage:   make([]types.Damage, 0),
				CanDodge: localMeta["CanDodge"].(bool),
			}

			for _, value := range localMeta["Damage"].([]interface{}) {
				value := value.(map[string]interface{})

				dmg.Damage = append(dmg.Damage, types.Damage{
					Value: value["Value"].(int),
					Type:  types.DamageType(value["Type"].(int)),
				})
			}

			act.Meta = dmg
		case types.ACTION_EFFECT:
			localMeta := dataMap["Meta"].(map[string]interface{})

			effectType := StringToEffectType[localMeta["Effect"].(string)]

			actionEffect := types.ActionEffect{
				Effect:   effectType,
				Value:    int(localMeta["Value"].(float64)),
				Uuid:     uuid.MustParse(localMeta["Uuid"].(string)),
				Duration: int(localMeta["Duration"].(float64)),
				Caster:   uuid.MustParse(localMeta["Caster"].(string)),
				Target:   uuid.MustParse(localMeta["Target"].(string)),
				Source:   types.EffectSource(0),
				OnExpire: func(owner types.Entity, fightInstance types.FightInstance, meta types.ActionEffect) {
					value := localMeta["OnExpire"].(utils.LuaFunctionRef)

					state.Global(value.FunctionName)

					state.PushUserData(owner)
					state.PushUserData(fightInstance)
					state.PushUserData(meta)

					state.Call(3, 0)
				},
			}

			effectMetaDetails := localMeta["Meta"].(map[string]interface{})

			switch effectType {
			case types.EFFECT_HEAL:
				actionEffect.Meta = types.ActionEffectHeal{
					Value: int(effectMetaDetails["Value"].(float64)),
				}
			case types.EFFECT_STAT_DEC:
				actionEffect.Meta = types.ActionEffectStat{
					Stat:      types.Stat(effectMetaDetails["Stat"].(float64)),
					Value:     int(effectMetaDetails["Value"].(float64)),
					IsPercent: effectMetaDetails["IsPercent"].(bool),
				}
			case types.EFFECT_STAT_INC:
				actionEffect.Meta = types.ActionEffectStat{
					Stat:      types.Stat(effectMetaDetails["Stat"].(float64)),
					Value:     int(effectMetaDetails["Value"].(float64)),
					IsPercent: effectMetaDetails["IsPercent"].(bool),
				}
			case types.EFFECT_RESIST:
				actionEffect.Meta = types.ActionEffectResist{
					Value:     int(effectMetaDetails["Value"].(float64)),
					IsPercent: effectMetaDetails["IsPercent"].(bool),
					DmgType:   int(effectMetaDetails["DmgType"].(float64)),
				}
			}

			act.Meta = actionEffect
		}

		fight.HandleAction(act)

		return 0
	})

	state.SetGlobal("HandleAction")

	state.PushGoFunction(func(state *lua.State) int {
		fight := state.ToUserData(1).(*battle.Fight)
		playerUuid, _ := state.ToString(2)

		ent := fight.GetAlliesFor(uuid.MustParse(playerUuid))

		state.NewTable()

		for idx, entity := range ent {
			state.PushInteger(idx + 1)
			state.PushUserData(entity)
			state.SetTable(-3)
		}

		return 1
	})

	state.SetGlobal("GetAlliesFor")

	state.PushGoFunction(func(state *lua.State) int {
		fight := state.ToUserData(1).(*battle.Fight)
		playerUuid, _ := state.ToString(2)

		ent := fight.GetEnemiesFor(uuid.MustParse(playerUuid))

		state.NewTable()

		for idx, entity := range ent {
			state.PushInteger(idx + 1)
			state.PushUserData(entity)
			state.SetTable(-3)
		}

		return 1
	})

	state.SetGlobal("GetEnemiesFor")

	state.PushGoFunction(func(state *lua.State) int {
		fight := state.ToUserData(1).(*battle.Fight)

		entityUuid, _ := state.ToString(2)

		turn := fight.GetTurnFor(uuid.MustParse(entityUuid))

		state.PushInteger(turn)

		return 1
	})

	state.SetGlobal("GetTurnFor")
}

func AddEntityFunctions(state *lua.State) {
	state.PushGoFunction(func(state *lua.State) int {
		entity := state.ToUserData(1).(types.Entity)

		state.PushString(entity.GetUUID().String())

		return 1
	})

	state.SetGlobal("GetUUID")

	state.PushGoFunction(func(state *lua.State) int {
		entity := state.ToUserData(1).(types.Entity)

		state.PushInteger(entity.GetCurrentHP())

		return 1
	})

	state.SetGlobal("GetCurrentHP")

	state.PushGoFunction(func(state *lua.State) int {
		entity := state.ToUserData(1).(types.Entity)

		rawUuid, _ := state.ToString(2)

		effectUuid := uuid.MustParse(rawUuid)

		effect := entity.GetEffectByUUID(effectUuid)

		if effect == nil {
			state.PushNil()
			return 1
		}

		state.NewTable()

		state.PushInteger(int(effect.Effect))
		state.SetField(-2, "Effect")

		state.PushInteger(effect.Value)
		state.SetField(-2, "Value")

		state.PushInteger(effect.Duration)
		state.SetField(-2, "Duration")

		state.PushString(effect.Uuid.String())
		state.SetField(-2, "Uuid")

		state.PushString(effect.Caster.String())
		state.SetField(-2, "Caster")

		state.PushString(effect.Target.String())
		state.SetField(-2, "Target")

		state.PushInteger(int(effect.Source))
		state.SetField(-2, "Source")

		state.PushNil()
		state.SetField(-2, "Meta")

		return 1
	})

	state.SetGlobal("GetEffectByUUID")

	state.PushGoFunction(func(state *lua.State) int {
		entity := state.ToUserData(1).(types.Entity)

		rawUuid, _ := state.ToString(2)

		effectUuid := uuid.MustParse(rawUuid)

		entity.RemoveEffect(effectUuid)

		return 0

	})

	state.SetGlobal("RemoveEffect")

	state.PushGoFunction(func(state *lua.State) int {
		entity := state.ToUserData(1).(types.Entity)

		dataMap, err := utils.GetTableAsMap(state)

		if err != nil {
			panic(err)
		}

		localMeta := dataMap["Meta"].(map[string]interface{})

		effectType := StringToEffectType[localMeta["Effect"].(string)]

		actionEffect := types.ActionEffect{
			Effect:   effectType,
			Value:    int(localMeta["Value"].(float64)),
			Uuid:     uuid.MustParse(localMeta["Uuid"].(string)),
			Duration: int(localMeta["Duration"].(float64)),
			Caster:   uuid.MustParse(localMeta["Caster"].(string)),
			Target:   uuid.MustParse(localMeta["Target"].(string)),
			Source:   types.EffectSource(0),
			OnExpire: func(owner types.Entity, fightInstance types.FightInstance, meta types.ActionEffect) {
				value := localMeta["OnExpire"].(utils.LuaFunctionRef)

				state.Global(value.FunctionName)

				state.PushUserData(owner)
				state.PushUserData(fightInstance)
				state.PushUserData(meta)

				state.Call(3, 0)
			},
		}

		effectMetaDetails := localMeta["Meta"].(map[string]interface{})

		switch effectType {
		case types.EFFECT_DOT:
			panic("DOT event not implemented")
		case types.EFFECT_HEAL:
			actionEffect.Meta = types.ActionEffectHeal{
				Value: int(effectMetaDetails["Value"].(float64)),
			}
		case types.EFFECT_STAT_DEC:
			actionEffect.Meta = types.ActionEffectStat{
				Stat:      types.Stat(effectMetaDetails["Stat"].(float64)),
				Value:     int(effectMetaDetails["Value"].(float64)),
				IsPercent: effectMetaDetails["IsPercent"].(bool),
			}
		case types.EFFECT_STAT_INC:
			actionEffect.Meta = types.ActionEffectStat{
				Stat:      types.Stat(effectMetaDetails["Stat"].(float64)),
				Value:     int(effectMetaDetails["Value"].(float64)),
				IsPercent: effectMetaDetails["IsPercent"].(bool),
			}
		case types.EFFECT_RESIST:
			actionEffect.Meta = types.ActionEffectResist{
				Value:     int(effectMetaDetails["Value"].(float64)),
				IsPercent: effectMetaDetails["IsPercent"].(bool),
				DmgType:   int(effectMetaDetails["DmgType"].(float64)),
			}
		}

		entity.ApplyEffect(actionEffect)

		return 0
	})

	state.SetGlobal("ApplyEffect")

	state.PushGoFunction(func(state *lua.State) int {
		entity := state.ToUserData(1).(types.Entity)

		rawName, _ := state.ToString(2)

		effect := entity.GetEffectByType(StringToEffectType[rawName])

		if effect == nil {
			state.PushNil()
			return 1
		}

		state.NewTable()

		state.PushInteger(int(effect.Effect))
		state.SetField(-2, "Effect")

		state.PushInteger(effect.Value)
		state.SetField(-2, "Value")

		state.PushInteger(effect.Duration)
		state.SetField(-2, "Duration")

		state.PushString(effect.Uuid.String())
		state.SetField(-2, "Uuid")

		state.PushString(effect.Caster.String())
		state.SetField(-2, "Caster")

		state.PushString(effect.Target.String())
		state.SetField(-2, "Target")

		state.PushInteger(int(effect.Source))
		state.SetField(-2, "Source")

		state.PushNil()
		state.SetField(-2, "Meta")

		return 1
	})

	state.SetGlobal("GetEffectByType")

	state.PushGoFunction(func(state *lua.State) int {
		entity := state.ToUserData(1).(types.MobEntity)

		fight := state.ToUserData(2).(*battle.Fight)

		act := entity.GetDefaultAction(fight)

		state.NewTable()

		for idx, action := range act {
			state.PushInteger(idx + 1)

			SerializeAction(action, state)

			state.SetTable(-3)
		}

		return 1
	})

	state.SetGlobal("DefaultAction")
}

func ReadMapAsTrigger(dataMap map[string]interface{}) types.Trigger {
	trigger := types.Trigger{}

	for key, value := range dataMap {
		switch key {
		case "Type":
			switch value {
			case "PASSIVE":
				trigger.Type = types.TRIGGER_PASSIVE
			case "ACTIVE":
				trigger.Type = types.TRIGGER_ACTIVE
			case "TYPE_NONE":
				trigger.Type = types.TRIGGER_TYPE_NONE
			}
		case "Event":
			trigger.Event = StringToTriggerType[value.(string)]
		case "Cooldown":
			cooldown := types.CooldownMeta{
				PassEvent: types.TRIGGER_NONE,
			}

			for key, value := range value.(map[string]interface{}) {
				switch key {
				case "PassEvent":
					cooldown.PassEvent = StringToTriggerType[value.(string)]
				}
			}

			if cooldown.PassEvent != types.TRIGGER_NONE {
				trigger.Cooldown = &cooldown
			}
		case "Flags":
			trigger.Flags = types.SkillFlag(value.(int))
		}
	}

	return trigger
}

var StringToActionEvent map[string]types.ActionEnum = map[string]types.ActionEnum{
	"ACTION_ATTACK":  types.ACTION_ATTACK,
	"ACTION_DEFEND":  types.ACTION_DEFEND,
	"ACTION_SKILL":   types.ACTION_SKILL,
	"ACTION_ITEM":    types.ACTION_ITEM,
	"ACTION_RUN":     types.ACTION_RUN,
	"ACTION_COUNTER": types.ACTION_COUNTER,
	"ACTION_EFFECT":  types.ACTION_EFFECT,
	"ACTION_DMG":     types.ACTION_DMG,
	"ACTION_SUMMON":  types.ACTION_SUMMON,
}

var StringToTriggerType map[string]types.SkillTrigger = map[string]types.SkillTrigger{
	"NONE":                types.TRIGGER_NONE,
	"ATTACK_BEFORE":       types.TRIGGER_ATTACK_BEFORE,
	"ATTACK_HIT":          types.TRIGGER_ATTACK_HIT,
	"ATTACK_MISS":         types.TRIGGER_ATTACK_MISS,
	"ATTACK_GOT_HIT":      types.TRIGGER_ATTACK_GOT_HIT,
	"EXECUTE":             types.TRIGGER_EXECUTE,
	"TURN":                types.TRIGGER_TURN,
	"CAST_ULT":            types.TRIGGER_CAST_ULT,
	"DAMAGE_BEFORE":       types.TRIGGER_DAMAGE_BEFORE,
	"DAMAGE":              types.TRIGGER_DAMAGE,
	"HEAL_SELF":           types.TRIGGER_HEAL_SELF,
	"HEAL_OTHER":          types.TRIGGER_HEAL_OTHER,
	"APPLY_CROWD_CONTROL": types.TRIGGER_APPLY_CROWD_CONTROL,
}

var StringToEffectType map[string]types.Effect = map[string]types.Effect{
	"EFFECT_DOT":          types.EFFECT_DOT,
	"EFFECT_HEAL":         types.EFFECT_HEAL,
	"EFFECT_MANA_RESTORE": types.EFFECT_MANA_RESTORE,
	"EFFECT_SHIELD":       types.EFFECT_SHIELD,
	"EFFECT_STUN":         types.EFFECT_STUN,
	"EFFECT_STAT_INC":     types.EFFECT_STAT_INC,
	"EFFECT_STAT_DEC":     types.EFFECT_STAT_DEC,
	"EFFECT_RESIST":       types.EFFECT_RESIST,
	"EFFECT_TAUNT":        types.EFFECT_TAUNT,
	"EFFECT_TAUNTED":      types.EFFECT_TAUNTED,
}

var StringToDamageType = map[string]types.DamageType{
	"DMG_TRUE":     types.DMG_TRUE,
	"DMG_PHYSICAL": types.DMG_PHYSICAL,
	"DMG_MAGICAL":  types.DMG_MAGICAL,
}

func ParseReturnMeta(dataMap map[string]interface{}, trigger types.Trigger) interface{} {
	switch trigger.Event {
	case types.TRIGGER_ATTACK_BEFORE:
		var shouldMiss bool
		if value, exists := dataMap["ShouldMiss"]; exists {
			shouldMiss = value.(bool)
		}

		var shouldHit bool
		if value, exists := dataMap["ShouldHit"]; exists {
			shouldHit = value.(bool)
		}

		effects := make([]types.DamagePartial, 0)

		for _, effect := range dataMap["Effects"].([]interface{}) {
			effect := effect.(map[string]interface{})

			effects = append(effects, types.DamagePartial{
				Value:   int(effect["Value"].(float64)),
				Type:    types.DamageType(effect["Type"].(float64)),
				Percent: effect["Percent"].(bool),
			})
		}

		return types.AttackTriggerMeta{
			ShouldMiss: shouldMiss,
			ShouldHit:  shouldHit,
			Effects:    effects,
		}
	case types.TRIGGER_ATTACK_HIT:
		var shouldMiss bool
		if value, exists := dataMap["ShouldMiss"]; exists {
			shouldMiss = value.(bool)
		}

		var shouldHit bool
		if value, exists := dataMap["ShouldHit"]; exists {
			shouldHit = value.(bool)
		}

		effects := make([]types.DamagePartial, 0)

		for _, effect := range dataMap["Effects"].([]interface{}) {
			effect := effect.(map[string]interface{})

			effects = append(effects, types.DamagePartial{
				Value:   int(effect["Value"].(float64)),
				Type:    types.DamageType(effect["Type"].(float64)),
				Percent: effect["Percent"].(bool),
			})
		}

		return types.AttackTriggerMeta{
			ShouldMiss: shouldMiss,
			ShouldHit:  shouldHit,
			Effects:    effects,
		}
	case types.TRIGGER_ATTACK_GOT_HIT:
		var shouldMiss bool
		if value, exists := dataMap["ShouldMiss"]; exists {
			shouldMiss = value.(bool)
		}

		var shouldHit bool
		if value, exists := dataMap["ShouldHit"]; exists {
			shouldHit = value.(bool)
		}

		effects := make([]types.DamagePartial, 0)

		for _, effect := range dataMap["Effects"].([]interface{}) {
			effect := effect.(map[string]interface{})

			effects = append(effects, types.DamagePartial{
				Value:   int(effect["Value"].(float64)),
				Type:    types.DamageType(effect["Type"].(float64)),
				Percent: effect["Percent"].(bool),
			})
		}

		return types.AttackTriggerMeta{
			ShouldMiss: shouldMiss,
			ShouldHit:  shouldHit,
			Effects:    effects,
		}
	case types.TRIGGER_DAMAGE_BEFORE:
		effects := make([]types.DamagePartial, 0)

		for _, effect := range dataMap["Effects"].([]interface{}) {
			effect := effect.(map[string]interface{})

			effects = append(effects, types.DamagePartial{
				Value:   int(effect["Value"].(float64)),
				Type:    types.DamageType(effect["Type"].(float64)),
				Percent: effect["Percent"].(bool),
			})
		}

		return types.DamageTriggerMeta{
			Effects: effects,
		}
	case types.TRIGGER_DAMAGE:
		effects := make([]types.DamagePartial, 0)

		for _, effect := range dataMap["Effects"].([]interface{}) {
			effect := effect.(map[string]interface{})

			effects = append(effects, types.DamagePartial{
				Value:   int(effect["Value"].(float64)),
				Type:    types.DamageType(effect["Type"].(float64)),
				Percent: effect["Percent"].(bool),
			})
		}

		return types.DamageTriggerMeta{
			Effects: effects,
		}
	case types.TRIGGER_HEAL_SELF:
		effects := make([]types.IncreasePartial, 0)

		for _, effect := range dataMap["Effects"].([]interface{}) {
			effect := effect.(map[string]interface{})

			effects = append(effects, types.IncreasePartial{
				Value:   int(effect["Value"].(float64)),
				Percent: effect["Percent"].(bool),
			})
		}

		return types.EffectTriggerMeta{
			Effects: effects,
		}
	case types.TRIGGER_HEAL_OTHER:
		effects := make([]types.IncreasePartial, 0)

		for _, effect := range dataMap["Effects"].([]interface{}) {
			effect := effect.(map[string]interface{})

			effects = append(effects, types.IncreasePartial{
				Value:   int(effect["Value"].(float64)),
				Percent: effect["Percent"].(bool),
			})
		}

		return types.EffectTriggerMeta{
			Effects: effects,
		}
	case types.TRIGGER_APPLY_CROWD_CONTROL:
		effects := make([]types.IncreasePartial, 0)

		for _, effect := range dataMap["Effects"].([]interface{}) {
			effect := effect.(map[string]interface{})

			effects = append(effects, types.IncreasePartial{
				Value:   int(effect["Value"].(float64)),
				Percent: effect["Percent"].(bool),
			})
		}

		return types.EffectTriggerMeta{
			Effects: effects,
		}
	}

	return nil
}

func ParseActionReturn(dataMap map[string]interface{}, state *lua.State) types.Action {
	var consumeTurn bool

	if value, existsValue := dataMap["ConsumeTurn"]; existsValue {
		consumeTurn = value.(bool)
	}

	act := types.Action{
		Event:       StringToActionEvent[dataMap["Event"].(string)],
		Source:      uuid.MustParse(dataMap["Source"].(string)),
		Target:      uuid.MustParse(dataMap["Target"].(string)),
		ConsumeTurn: &consumeTurn,
	}

	switch act.Event {
	case types.ACTION_SUMMON:
		panic("Summon not implemented")
	case types.ACTION_DMG:
		localMeta := dataMap["Meta"].(map[string]interface{})

		dmg := types.ActionDamage{
			Damage:   make([]types.Damage, 0),
			CanDodge: localMeta["CanDodge"].(bool),
		}

		for _, value := range localMeta["Damage"].([]interface{}) {
			value := value.(map[string]interface{})

			dmg.Damage = append(dmg.Damage, types.Damage{
				Value: value["Value"].(int),
				Type:  types.DamageType(value["Type"].(int)),
			})
		}

		act.Meta = dmg
	case types.ACTION_EFFECT:
		localMeta := dataMap["Meta"].(map[string]interface{})

		effectType := StringToEffectType[localMeta["Effect"].(string)]

		actionEffect := types.ActionEffect{
			Effect:   effectType,
			Value:    int(localMeta["Value"].(float64)),
			Uuid:     uuid.MustParse(localMeta["Uuid"].(string)),
			Duration: int(localMeta["Duration"].(float64)),
			Caster:   uuid.MustParse(localMeta["Caster"].(string)),
			Target:   uuid.MustParse(localMeta["Target"].(string)),
			Source:   types.EffectSource(0),
			OnExpire: func(owner types.Entity, fightInstance types.FightInstance, meta types.ActionEffect) {
				value := localMeta["OnExpire"].(utils.LuaFunctionRef)

				state.Global(value.FunctionName)

				state.PushUserData(owner)
				state.PushUserData(fightInstance)
				state.PushUserData(meta)

				state.Call(3, 0)
			},
		}

		effectMetaDetails := localMeta["Meta"].(map[string]interface{})

		switch effectType {
		case types.EFFECT_HEAL:
			actionEffect.Meta = types.ActionEffectHeal{
				Value: int(effectMetaDetails["Value"].(float64)),
			}
		case types.EFFECT_STAT_DEC:
			actionEffect.Meta = types.ActionEffectStat{
				Stat:      types.Stat(effectMetaDetails["Stat"].(float64)),
				Value:     int(effectMetaDetails["Value"].(float64)),
				IsPercent: effectMetaDetails["IsPercent"].(bool),
			}
		case types.EFFECT_STAT_INC:
			actionEffect.Meta = types.ActionEffectStat{
				Stat:      types.Stat(effectMetaDetails["Stat"].(float64)),
				Value:     int(effectMetaDetails["Value"].(float64)),
				IsPercent: effectMetaDetails["IsPercent"].(bool),
			}
		case types.EFFECT_RESIST:
			actionEffect.Meta = types.ActionEffectResist{
				Value:     int(effectMetaDetails["Value"].(float64)),
				IsPercent: effectMetaDetails["IsPercent"].(bool),
				DmgType:   int(effectMetaDetails["DmgType"].(float64)),
			}
		}

		act.Meta = actionEffect
	}

	return act
}

func SerializeAction(act types.Action, state *lua.State) {
	state.NewTable()

	for key, value := range StringToActionEvent {
		if value == act.Event {
			state.PushString(key)
			state.SetField(-2, "Event")
		}
	}

	state.PushString(act.Source.String())
	state.SetField(-2, "Source")

	state.PushString(act.Target.String())
	state.SetField(-2, "Target")

	if act.ConsumeTurn != nil {
		state.PushBoolean(*act.ConsumeTurn)
		state.SetField(-2, "ConsumeTurn")
	}

	switch act.Event {
	case types.ACTION_SUMMON:
		panic("Summon not implemented")
	case types.ACTION_DMG:
		dmg := act.Meta.(types.ActionDamage)

		state.NewTable()

		state.PushBoolean(dmg.CanDodge)
		state.SetField(-2, "CanDodge")

		state.NewTable()

		for idx, value := range dmg.Damage {
			state.PushInteger(idx + 1)

			state.NewTable()

			for key, valueNew := range StringToDamageType {
				if value.Type == valueNew {
					state.PushString(key)
					state.SetField(-2, "Type")
				}
			}

			state.PushInteger(value.Value)
			state.SetField(-2, "Value")

			state.SetTable(-3)
		}

		state.SetField(-2, "Damage")

		state.SetField(-2, "Meta")
	case types.ACTION_EFFECT:
		effect := act.Meta.(types.ActionEffect)

		state.NewTable()

		for key, value := range StringToEffectType {
			if value == effect.Effect {
				state.PushString(key)
				state.SetField(-2, "Effect")
			}
		}

		state.PushInteger(effect.Value)
		state.SetField(-2, "Value")

		state.PushString(effect.Uuid.String())
		state.SetField(-2, "Uuid")

		state.PushString(effect.Caster.String())
		state.SetField(-2, "Caster")

		state.PushString(effect.Target.String())
		state.SetField(-2, "Target")

		state.PushInteger(effect.Duration)
		state.SetField(-2, "Duration")

		if effect.Meta != nil {
			state.NewTable()

			switch effect.Effect {
			case types.EFFECT_HEAL:
				heal := effect.Meta.(types.ActionEffectHeal)

				state.PushInteger(heal.Value)
				state.SetField(-2, "Value")

				state.SetField(-2, "Meta")
			case types.EFFECT_STAT_DEC:
				stat := effect.Meta.(types.ActionEffectStat)

				state.PushInteger(int(stat.Stat))
				state.SetField(-2, "Stat")

				state.PushInteger(stat.Value)
				state.SetField(-2, "Value")

				state.PushBoolean(stat.IsPercent)
				state.SetField(-2, "IsPercent")

				state.SetField(-2, "Meta")
			case types.EFFECT_STAT_INC:
				stat := effect.Meta.(types.ActionEffectStat)

				state.PushInteger(int(stat.Stat))
				state.SetField(-2, "Stat")

				state.PushInteger(stat.Value)

				state.SetField(-2, "Value")
				state.PushBoolean(stat.IsPercent)

				state.SetField(-2, "IsPercent")
				state.SetField(-2, "Meta")

			case types.EFFECT_RESIST:
				resist := effect.Meta.(types.ActionEffectResist)

				state.PushInteger(resist.Value)
				state.SetField(-2, "Value")

				state.PushBoolean(resist.IsPercent)
				state.SetField(-2, "IsPercent")

				state.PushInteger(resist.DmgType)
				state.SetField(-2, "DmgType")

				state.SetField(-2, "Meta")
			}

		}

		state.SetField(-2, "Meta")
	}
}

type SimplePlayerSkill struct {
	Trigger types.Trigger
	Exec    func(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{}
}

func (s SimplePlayerSkill) Execute(owner types.PlayerEntity, target types.Entity, fightInstance types.FightInstance, meta interface{}) interface{} {
	return s.Exec(owner, target, fightInstance, meta)
}

func (s SimplePlayerSkill) GetEvents() map[types.CustomTrigger]func(owner types.PlayerEntity) {
	return nil
}

func (s SimplePlayerSkill) GetUUID() uuid.UUID {
	return uuid.New()
}

func (s SimplePlayerSkill) GetName() string {
	return ""
}

func (s SimplePlayerSkill) GetDescription() string {
	return ""
}

func (s SimplePlayerSkill) GetCD() int {
	return 0
}

func (s SimplePlayerSkill) GetCost() int {
	return 0
}

func (s SimplePlayerSkill) GetTrigger() types.Trigger {
	return s.Trigger
}

func (s SimplePlayerSkill) IsLevelSkill() bool {
	return false
}
