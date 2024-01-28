package location

import (
	"sao/battle"
)

type Location struct {
	Name     string
	CID      string
	CityPart bool
	Enemies  []string
	Effects  []LocationEffect
}

type Floor struct {
	Name      string
	CID       string
	Locations []Location
	Effects   []LocationEffect
}

func (f Floor) FindLocation(str string) *Location {
	for _, loc := range f.Locations {
		if loc.CID == str || loc.Name == str {
			return &loc
		}
	}

	return nil
}

type LocationEffect struct {
	Effect battle.Effect
	Value  int
	Meta   *map[string]interface{}
}
