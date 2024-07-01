package battle

import "sao/types"

type EntitySort struct {
	Entities []Entity
	Order    []types.TargetDetails
	Meta     *map[string]interface{}
}

func (e EntitySort) Len() int {
	return len(e.Entities)
}

func (e EntitySort) Less(i, j int) bool {
	for _, order := range e.Order {
		result := TargetDetailsCheck(e.Entities[i], e.Entities[j], order, e.Meta)

		if result == 0 {
			continue
		}

		return result > 0
	}

	return false
}

func TargetDetailsCheck(left, right interface{}, order types.TargetDetails, meta *map[string]interface{}) int {
	if order == types.DETAIL_ALL {
		return 0
	}

	switch order {
	case types.DETAIL_MAX_HP:
		return left.(Entity).GetStat(types.STAT_HP) - right.(Entity).GetStat(types.STAT_HP)
	case types.DETAIL_LOW_HP:
		return right.(Entity).GetStat(types.STAT_HP) - left.(Entity).GetStat(types.STAT_HP)
	case types.DETAIL_MAX_MP:
		return left.(Entity).GetStat(types.STAT_MANA) - right.(Entity).GetStat(types.STAT_MANA)
	case types.DETAIL_LOW_MP:
		return right.(Entity).GetStat(types.STAT_MANA) - left.(Entity).GetStat(types.STAT_MANA)
	case types.DETAIL_MAX_ATK:
		return left.(Entity).GetStat(types.STAT_AD) - right.(Entity).GetStat(types.STAT_AD)
	case types.DETAIL_LOW_ATK:
		return right.(Entity).GetStat(types.STAT_AD) - left.(Entity).GetStat(types.STAT_AD)
	case types.DETAIL_MAX_DEF:
		return left.(Entity).GetStat(types.STAT_DEF) - right.(Entity).GetStat(types.STAT_DEF)
	case types.DETAIL_LOW_DEF:
		return right.(Entity).GetStat(types.STAT_DEF) - left.(Entity).GetStat(types.STAT_DEF)
	case types.DETAIL_MAX_SPD:
		return left.(Entity).GetStat(types.STAT_SPD) - right.(Entity).GetStat(types.STAT_SPD)
	case types.DETAIL_LOW_SPD:
		return right.(Entity).GetStat(types.STAT_SPD) - left.(Entity).GetStat(types.STAT_SPD)
	case types.DETAIL_MAX_AP:
		return left.(Entity).GetStat(types.STAT_AP) - right.(Entity).GetStat(types.STAT_AP)
	case types.DETAIL_LOW_AP:
		return right.(Entity).GetStat(types.STAT_AP) - left.(Entity).GetStat(types.STAT_AP)
	case types.DETAIL_MAX_RES:
		return left.(Entity).GetStat(types.STAT_MR) - right.(Entity).GetStat(types.STAT_MR)
	case types.DETAIL_LOW_RES:
		return right.(Entity).GetStat(types.STAT_MR) - left.(Entity).GetStat(types.STAT_MR)
	case types.DETAIL_HAS_EFFECT:
		if meta == nil {
			return 0
		}
		leftHas := left.(Entity).GetEffectByType((*meta)["effect"].(Effect))
		rightHas := right.(Entity).GetEffectByType((*meta)["effect"].(Effect))

		if leftHas != nil && rightHas == nil {
			return -1
		}

		if leftHas == nil && rightHas != nil {
			return 1
		}

		return 0
	case types.DETAIL_NO_EFFECT:
		if meta == nil {
			return 0
		}

		leftHas := left.(Entity).GetEffectByType((*meta)["effect"].(Effect))
		rightHas := right.(Entity).GetEffectByType((*meta)["effect"].(Effect))

		if leftHas != nil && rightHas == nil {
			return 1
		}

		if leftHas == nil && rightHas != nil {
			return -1
		}

		return 0
	}

	return 0
}

func (e EntitySort) Swap(i, j int) {
	e.Entities[i], e.Entities[j] = e.Entities[j], e.Entities[i]
}
