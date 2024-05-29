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
	Name      string
	CID       string
	Default   string
	Locations []Location
	Effects   []LocationEffect
	Unlocked  bool
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
	Effect int
	Value  int
	Meta   *map[string]interface{}
}

func GetFloors(test bool) map[string]Floor {
	if test {
		return map[string]Floor{
			"dev": {
				Name:     "dev",
				CID:      "1162450076438900958",
				Default:  "Las",
				Unlocked: false,
				Locations: []Location{
					{
						Name:     "Rynek",
						CID:      "1162450122249076756",
						CityPart: true,
						Effects:  []LocationEffect{},
						TP:       true,
						Unlocked: true,
						Flags:    []string{},
					},
					{
						Name:     "Las",
						CID:      "1162450159251234876",
						CityPart: false,
						Enemies: []EnemyMeta{
							{
								MinNum: 1,
								MaxNum: 3,
								Enemy:  "LV0_Rycerz",
							},
						},
						Effects:  []LocationEffect{},
						TP:       false,
						Unlocked: false,
						Flags:    []string{},
					},
					{
						Name:     "Arena",
						CID:      "1241069460073222225",
						CityPart: true,
						Effects:  []LocationEffect{},
						TP:       true,
						Unlocked: true,
						Flags:    []string{"arena"},
					},
				},
				Effects: []LocationEffect{},
			},
		}
	} else {
		//TODO fill this
		return map[string]Floor{}
	}
}
