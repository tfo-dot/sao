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
				Unlocked: true,
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
						Unlocked: true,
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
					{
						Name:     "Boss",
						CID:      "1265406962263396463",
						CityPart: false,
						Enemies: []EnemyMeta{
							{
								MinNum: 1,
								MaxNum: 1,
								Enemy:  "LV0_Boss",
							},
						},
						Effects:  []LocationEffect{},
						TP:       false,
						Unlocked: true,
						Flags:    []string{"boss"},
					},
				},
				Effects: []LocationEffect{},
			},
		}
	} else {
		return map[string]Floor{
			"Miasto": {
				Name:             "Miasto",
				CID:              "1160645670798098483",
				Default:          "Główny Plac",
				Unlocked:         true,
				CountsAsUnlocked: false,
				Locations: []Location{
					{
						Name:     "Główny Plac",
						CID:      "1160646008741576794",
						CityPart: true,
						Effects:  []LocationEffect{},
						TP:       true,
						Unlocked: true,
						Flags:    []string{},
					},
					{
						Name:     "Uliczki",
						CID:      "1160646032259027004",
						CityPart: true,
						Effects:  []LocationEffect{},
						TP:       true,
						Unlocked: true,
						Flags:    []string{},
					},
					{
						Name:     "Rynek",
						CID:      "1160646058196619294",
						Effects:  []LocationEffect{},
						TP:       true,
						Unlocked: true,
						Flags:    []string{},
					},
					{
						Name:     "Targ",
						CID:      "1255069551910060124",
						Effects:  []LocationEffect{},
						TP:       true,
						Unlocked: true,
						Flags:    []string{},
					},
					{
						Name:     "Tawerna",
						CID:      "1160646100433240125",
						Effects:  []LocationEffect{},
						TP:       true,
						Unlocked: true,
						Flags:    []string{},
					},
					{
						Name:     "Gildia poszukiwaczy przygód",
						CID:      "1160646142254657576",
						Effects:  []LocationEffect{},
						TP:       true,
						Unlocked: true,
						Flags:    []string{},
					},
					{
						Name:     "Fontanna",
						CID:      "1160646178619273318",
						Effects:  []LocationEffect{},
						TP:       true,
						Unlocked: true,
						Flags:    []string{},
					},
					{
						Name:     "Kuźnia",
						CID:      "1160646203671842937",
						Effects:  []LocationEffect{},
						TP:       true,
						Unlocked: true,
						Flags:    []string{},
					},
					{
						Name:     "Biblioteka",
						CID:      "1160646498032308265",
						Effects:  []LocationEffect{},
						TP:       true,
						Unlocked: true,
						Flags:    []string{},
					},
					{
						Name:     "Arena",
						CID:      "1160646528474566776",
						Effects:  []LocationEffect{},
						TP:       true,
						Unlocked: true,
						Flags:    []string{"arena"},
					},
					{
						Name:     "Trybuny areny",
						CID:      "1255069271239692288",
						Effects:  []LocationEffect{},
						TP:       true,
						Unlocked: true,
						Flags:    []string{},
					},
				},
				Effects: []LocationEffect{},
			},
		}
	}
}
