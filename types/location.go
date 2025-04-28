package types

type Location struct {
	Name     string
	CID      string
	CityPart bool
	TP       bool
	Enemies  []EnemyMeta `parts:"Enemies,ignoreEmpty"`
	Flags    []string `parts:"Flags,ignoreEmpty"`
	Unlocked bool
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
	Locations        []Location `parts:"Locations,ignoreEmpty"`
	Flags            []string `parts:"Flags,ignoreEmpty"`
	Unlocked         bool
	CountsAsUnlocked bool
}

func (f Floor) FindLocation(str string) *Location {
	for _, loc := range f.Locations {
		if loc.CID == str || loc.Name == str {
			return &loc
		}
	}

	return nil
}