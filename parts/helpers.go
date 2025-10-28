package parts

import (
	"sao/types"
	"sao/utils"

	"github.com/google/uuid"
	"github.com/tfo-dot/parts"
)

func AddConsts(vm *parts.VM) {
	vm.Enviroment.AppendValues(map[string]any{
		"STAT_NONE":              int(types.STAT_NONE),
		"STAT_HP":                int(types.STAT_HP),
		"STAT_HP_PLUS":           int(types.STAT_HP_PLUS),
		"STAT_SPD":               int(types.STAT_SPD),
		"STAT_AGL":               int(types.STAT_AGL),
		"STAT_AD":                int(types.STAT_AD),
		"STAT_DEF":               int(types.STAT_DEF),
		"STAT_MR":                int(types.STAT_MR),
		"STAT_MANA":              int(types.STAT_MANA),
		"STAT_MANA_PLUS":         int(types.STAT_MANA_PLUS),
		"STAT_AP":                int(types.STAT_AP),
		"STAT_HEAL_SELF":         int(types.STAT_HEAL_SELF),
		"STAT_HEAL_POWER":        int(types.STAT_HEAL_POWER),
		"STAT_LETHAL":            int(types.STAT_LETHAL),
		"STAT_LETHAL_PERCENT":    int(types.STAT_LETHAL_PERCENT),
		"STAT_MAGIC_PEN":         int(types.STAT_MAGIC_PEN),
		"STAT_MAGIC_PEN_PERCENT": int(types.STAT_MAGIC_PEN_PERCENT),
		"STAT_OMNI_VAMP":         int(types.STAT_OMNI_VAMP),
		"STAT_ATK_VAMP":          int(types.STAT_ATK_VAMP),

		"ENTITY_AUTO":   int(types.ENTITY_AUTO),
		"ENTITY_SUMMON": int(types.ENTITY_SUMMON),

		"DMG_PHYSICAL": int(types.DMG_PHYSICAL),
		"DMG_MAGICAL":  int(types.DMG_MAGICAL),
		"DMG_TRUE":     int(types.DMG_TRUE),

		"EFFECT_DOT":          int(types.EFFECT_DOT),
		"EFFECT_HEAL":         int(types.EFFECT_HEAL),
		"EFFECT_MANA_RESTORE": int(types.EFFECT_MANA_RESTORE),
		"EFFECT_SHIELD":       int(types.EFFECT_SHIELD),
		"EFFECT_STUN":         int(types.EFFECT_STUN),
		"EFFECT_STAT_INC":     int(types.EFFECT_STAT_INC),
		"EFFECT_STAT_DEC":     int(types.EFFECT_STAT_DEC),
		"EFFECT_RESIST":       int(types.EFFECT_RESIST),
		"EFFECT_TAUNT":        int(types.EFFECT_TAUNT),
		"EFFECT_TAUNTED":      int(types.EFFECT_TAUNTED),

		"LOOT_EXP":  int(types.LOOT_EXP),
		"LOOT_GOLD": int(types.LOOT_GOLD),

		"ACTION_ATTACK":  int(types.ACTION_ATTACK),
		"ACTION_DEFEND":  int(types.ACTION_DEFEND),
		"ACTION_SKILL":   int(types.ACTION_SKILL),
		"ACTION_ITEM":    int(types.ACTION_ITEM),
		"ACTION_RUN":     int(types.ACTION_RUN),
		"ACTION_COUNTER": int(types.ACTION_COUNTER),
		"ACTION_EFFECT":  int(types.ACTION_EFFECT),
		"ACTION_DMG":     int(types.ACTION_DMG),
		"ACTION_SUMMON":  int(types.ACTION_SUMMON),

		"SUMMON_FLAG_NONE":   int(types.SUMMON_FLAG_NONE),
		"SUMMON_FLAG_ATTACK": int(types.SUMMON_FLAG_ATTACK),
		"SUMMON_FLAG_EXPIRE": int(types.SUMMON_FLAG_EXPIRE),

		"TRIGGER_NONE":                int(types.TRIGGER_NONE),
		"TRIGGER_ATTACK_BEFORE":       int(types.TRIGGER_ATTACK_BEFORE),
		"TRIGGER_ATTACK_HIT":          int(types.TRIGGER_ATTACK_HIT),
		"TRIGGER_ATTACK_MISS":         int(types.TRIGGER_ATTACK_MISS),
		"TRIGGER_ATTACK_GOT_HIT":      int(types.TRIGGER_ATTACK_GOT_HIT),
		"TRIGGER_EXECUTE":             int(types.TRIGGER_EXECUTE),
		"TRIGGER_TURN":                int(types.TRIGGER_TURN),
		"TRIGGER_CAST_ULT":            int(types.TRIGGER_CAST_ULT),
		"TRIGGER_DAMAGE_BEFORE":       int(types.TRIGGER_DAMAGE_BEFORE),
		"TRIGGER_DAMAGE":              int(types.TRIGGER_DAMAGE),
		"TRIGGER_DAMAGE_GOT_HIT":      int(types.TRIGGER_DAMAGE_GOT_HIT),
		"TRIGGER_HEAL_SELF":           int(types.TRIGGER_HEAL_SELF),
		"TRIGGER_HEAL_OTHER":          int(types.TRIGGER_HEAL_OTHER),
		"TRIGGER_APPLY_CROWD_CONTROL": int(types.TRIGGER_APPLY_CROWD_CONTROL),

		"TRIGGER_PASSIVE":   int(types.TRIGGER_PASSIVE),
		"TRIGGER_ACTIVE":    int(types.TRIGGER_ACTIVE),
		"TRIGGER_TYPE_NONE": int(types.TRIGGER_TYPE_NONE),
	})
}

func AddFunctions(vm *parts.VM) {
	env := parts.VMEnviroment{
		Enclosing: nil,
		Values:    make(map[string]*parts.Literal),
	}

	env.DefineFunction("GetTurnFor", func(f types.FightInstance, uid *uuid.UUID) int {
		return f.GetTurnFor(*uid)
	})

	env.DefineFunction("GetRandomEnemy", func(fight types.FightInstance, mobUuid *uuid.UUID) *uuid.UUID {
		randomElt := utils.RandomElement(fight.GetEnemiesFor(*mobUuid)).GetUUID()

		return &randomElt
	})

	env.DefineFunction("GetRandomAlly", func(fight types.FightInstance, mobUuid *uuid.UUID) *uuid.UUID {
		randomElt := utils.RandomElement(fight.GetAlliesFor(*mobUuid)).GetUUID()

		return &randomElt
	})

	env.DefineFunction("DefaultAction", func(f types.FightInstance, m types.MobEntity) []any {
		actions := m.GetDefaultAction(f)

		temp := make([]any, len(actions))

		for idx, val := range actions {
			temp[idx] = val
		}

		return temp
	})

	env.DefineFunction("GetUUID", func(ent types.Entity) *uuid.UUID {
		uid := ent.GetUUID()
		return &uid
	})

	env.DefineFunction("GenerateUUID", func() *uuid.UUID {
		val := uuid.New()
		return &val
	})

	env.DefineFunction("PercentOf", func(value, percent int) int {
		return utils.PercentOf(value, percent)
	})

	env.DefineFunction("RandomInt", func(min, max int) int {
		return utils.RandomNumber(min, max)
	})

	env.DefineFunction("ForEach", func(elements []any, execute func(idx int, elt any)) {
		for idx, elt := range elements {
			execute(idx, elt)
		}
	})

	env.DefineFunction("GetEffectByType", func(ent types.Entity, effect types.Effect) any {
		efc := ent.GetEffectByType(effect)

		if efc == nil {
			return -1
		}
		
		return efc
	})

	env.DefineFunction("GetEffectByUUID", func(ent types.Entity, uid string) any {
		efc := ent.GetEffectByUUID(uuid.MustParse(uid))

		if efc == nil {
			return -1
		}
		
		return efc
	})

	env.DefineFunction("RemoveEffect", func(ent types.Entity, uid string) {
		ent.RemoveEffect(uuid.MustParse(uid))
	})

	env.DefineFunction("ApplyEffect", func(ent types.Entity, effect any) {
		dataMap := effect.(map[string]any)

		actEffect := types.ActionEffect{
			Effect: types.Effect(dataMap["RTEffect"].(int)),
			Value: dataMap["RTValue"].(int),
			Duration: dataMap["RTDuration"].(int),
			Uuid: uuid.MustParse(dataMap["RTUuid"].(string)),
		}

		//TODO parse meta

		ent.ApplyEffect(actEffect)
	})

	vm.Enviroment.Append(&env)
}

func FetchVal(vm *parts.VM, key string) (any, error) {
	rawVal, err := vm.Enviroment.Resolve(key)

	if err != nil {
		return nil, err
	}

	val, err := rawVal.ToGoTypes(vm)

	if err != nil {
		return nil, err
	}

	return val, nil
}
