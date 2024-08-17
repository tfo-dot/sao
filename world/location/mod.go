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
			"beta-miasto": {
				Name:     "beta-miasto",
				CID:      "1272232915308249119",
				Default:  "Rynek",
				Unlocked: true,
				Locations: []Location{
					{
						Name:     "Rynek",
						CID:      "1272233088608637050",
						CityPart: true,
						TP:       true,
						Unlocked: true,
					},
					{
						Name:     "Tawerna",
						CID:      "1272233328925343895",
						CityPart: true,
						TP:       true,
						Unlocked: true,
					},
					{
						Name:     "Arena",
						CID:      "1272233289159278622",
						CityPart: true,
						Unlocked: true,
						Flags:    []string{"arena"},
					},
					{
						Name:     "Brama główna",
						CID:      "1272233404900839586",
						CityPart: true,
						Unlocked: true,
						TP:       true,
					},
					{
						Name:     "Kuźnia",
						CID:      "1272233428649250928",
						CityPart: true,
						Unlocked: true,
						TP:       true,
					},
					{
						Name:     "Trybuny areny",
						CID:      "1272233365935886347",
						CityPart: true,
						Unlocked: true,
						TP:       true,
					},
				},
			},
			"beta-poza-miastem": {
				Name:     "beta-poza-miastem",
				CID:      "1272232994010169426",
				Default:  "Las",
				Unlocked: true,
				Locations: []Location{
					{
						Name:     "Las",
						CID:      "1272233789200007269",
						TP:       true,
						CityPart: false,
						Unlocked: true,

						Enemies: []EnemyMeta{
							{
								MinNum: 1,
								MaxNum: 2,
								Enemy:  "LV0_Wilk",
							},
						},
					},
					{
						Name:     "Polana",
						CID:      "1272233801837183036",
						TP:       true,
						CityPart: false,
						Unlocked: true,

						Enemies: []EnemyMeta{
							{
								MinNum: 1,
								MaxNum: 2,
								Enemy:  "LV0_Wilk",
							},
						},
					},
					{
						Name:     "Podnóże wulkanu",
						CID:      "1272233986617376829",
						TP:       true,
						CityPart: false,
						Unlocked: true,

						Enemies: []EnemyMeta{
							{
								MinNum: 1,
								MaxNum: 2,
								Enemy:  "LV0_Rycerz",
							},
						},
					},
					{
						Name:     "Wulkan",
						CID:      "1272234002174050394",
						TP:       true,
						CityPart: false,
						Unlocked: true,

						Flags: []string{"boss"},
						Enemies: []EnemyMeta{
							{
								MinNum: 1,
								MaxNum: 1,
								Enemy:  "LV0_Boss",
							},
						},
					},
				},
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
