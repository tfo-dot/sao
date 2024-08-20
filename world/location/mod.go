package location

type Location struct {
	Name     string
	CID      string
	CityPart bool
	Effects  []LocationEffect
	TP       bool
	Enemies  []EnemyMeta
	Unlocked bool
	Flags    []string
}

type EnemyMeta struct {
	MinNum int
	MaxNum int
	Enemy  string
}

type Floor struct {
	Name             string
	CID              string
	Default          string
	Locations        []Location
	Effects          []LocationEffect
	Unlocked         bool
	CountsAsUnlocked bool
}

type LocationEffect struct {
	Effect int
	Value  int
	Meta   *map[string]interface{}
}

func (f Floor) FindLocation(str string) *Location {
	for _, loc := range f.Locations {
		if loc.CID == str || loc.Name == str {
			return &loc
		}
	}

	return nil
}

func GetFloors() map[string]Floor {
	return map[string]Floor{}
}
