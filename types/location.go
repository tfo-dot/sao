package types

type PlayerLocation struct {
	FloorName    string
	LocationName string
}

func DefaultPlayerLocation() PlayerLocation {
	// return PlayerLocation{FloorName: "dev", LocationName: "Las"}
	return PlayerLocation{FloorName: "beta-miasto", LocationName: "Rynek"}
}
