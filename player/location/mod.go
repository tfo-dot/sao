package location

type PlayerLocation struct {
	FloorName    string
	LocationName string
}

func DefaultLocation() PlayerLocation {
	return PlayerLocation{FloorName: "dev", LocationName: "Las"}
	// return PlayerLocation{FloorName: "dev", LocationName: "Rynek"}
}
