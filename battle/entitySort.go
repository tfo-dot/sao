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
		return left.(Entity).GetMaxHP() - right.(Entity).GetMaxHP()
	case types.DETAIL_LOW_HP:
		return right.(Entity).GetMaxHP() - left.(Entity).GetMaxHP()
	case types.DETAIL_MAX_MP:
		return left.(Entity).GetMaxMana() - right.(Entity).GetMaxMana()
	case types.DETAIL_LOW_MP:
		return right.(Entity).GetMaxMana() - left.(Entity).GetMaxMana()
	case types.DETAIL_MAX_ATK:
		return left.(Entity).GetATK() - right.(Entity).GetATK()
	case types.DETAIL_LOW_ATK:
		return right.(Entity).GetATK() - left.(Entity).GetATK()
	case types.DETAIL_MAX_DEF:
		return left.(Entity).GetDEF() - right.(Entity).GetDEF()
	case types.DETAIL_LOW_DEF:
		return right.(Entity).GetDEF() - left.(Entity).GetDEF()
	case types.DETAIL_MAX_SPD:
		return left.(Entity).GetSPD() - right.(Entity).GetSPD()
	case types.DETAIL_LOW_SPD:
		return right.(Entity).GetSPD() - left.(Entity).GetSPD()
	case types.DETAIL_MAX_AP:
		return left.(Entity).GetAP() - right.(Entity).GetAP()
	case types.DETAIL_LOW_AP:
		return right.(Entity).GetAP() - left.(Entity).GetAP()
	case types.DETAIL_MAX_RES:
		return left.(Entity).GetMR() - right.(Entity).GetMR()
	case types.DETAIL_LOW_RES:
		return right.(Entity).GetMR() - left.(Entity).GetMR()
	case types.DETAIL_HAS_EFFECT:
		if meta == nil {
			return 0
		}
		leftHas := left.(Entity).GetEffect((*meta)["effect"].(Effect))
		rightHas := right.(Entity).GetEffect((*meta)["effect"].(Effect))

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

		leftHas := left.(Entity).GetEffect((*meta)["effect"].(Effect))
		rightHas := right.(Entity).GetEffect((*meta)["effect"].(Effect))

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
