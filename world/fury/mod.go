package fury

import (
	"sao/types"

	"github.com/google/uuid"
)

type Fury struct {
	Name        string
	Master      *uuid.UUID
	Tiers       []FuryTier
	CurrentTier int
	XP          FuryXP
	LvlStats    func(lvl int, tier int) map[types.Stat]int
}

type FuryTier struct {
	Stats       map[types.Stat]int
	Skills      []types.PlayerSkill
	Ingredients []types.Ingredient
}

type FuryXP struct {
	XP  int
	LVL int
}

func (f *Fury) NextLvlXPGauge() int {
	return f.CurrentTier*1000 + f.XP.LVL*100
}

func (f *Fury) AddXP(xp int) {
	if f.XP.LVL == 10 {
		f.XP.XP = 0
	} else {
		f.XP.XP += xp
	}

	for f.XP.XP >= f.NextLvlXPGauge() && f.XP.LVL < 10 {
		f.XP.LVL++
		f.XP.XP -= f.NextLvlXPGauge()
	}

	if f.XP.LVL == 10 {
		f.XP.XP = 0
	}
}

func (f *Fury) GetStats() map[types.Stat]int {
	baseStats := f.LvlStats(f.XP.LVL, f.CurrentTier)

	for i := 0; i < f.CurrentTier; i++ {
		for k, v := range f.Tiers[i].Stats {

			if _, ok := baseStats[k]; !ok {
				baseStats[k] = 0
			}

			baseStats[k] += v
		}
	}

	return baseStats
}

func (f *Fury) GetStat(stat types.Stat) int {
	if value, ok := f.GetStats()[stat]; ok {
		return value
	}

	return 0
}

func (f *Fury) GetSkills() []types.PlayerSkill {
	skills := []types.PlayerSkill{}

	for i := 0; i < f.CurrentTier; i++ {
		skills = append(skills, f.Tiers[i].Skills...)
	}

	return skills
}

func (f *Fury) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"name":        f.Name,
		"master":      f.Master,
		"tiers":       f.Tiers,
		"currentTier": f.CurrentTier,
		"xp":          f.XP,
	}
}

func Deserialize(data map[string]interface{}) *Fury {
	return &Fury{
		Name:        data["name"].(string),
		Master:      data["master"].(*uuid.UUID),
		Tiers:       data["tiers"].([]FuryTier),
		CurrentTier: data["currentTier"].(int),
		XP:          data["xp"].(FuryXP),
	}
}
