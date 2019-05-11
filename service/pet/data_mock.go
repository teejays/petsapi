package pet

// PopulateMockPets populates the data with mock pets
func PopulateMockPets() {
	mockPets := getMockPets()
	populateMockPets(mockPets)
}

func getMockPets() []Pet {
	// Popuate data with mock
	var mockPets = []Pet{
		{
			ID:   1,
			Name: "Tommy",
		},
		{
			ID:   2,
			Name: "Tiger",
		},
		{
			ID:   3,
			Name: "Buddy",
		},
		{
			ID:   5,
			Name: "Kitty",
		},
		{
			ID:   8,
			Name: "Coco",
		},
		{
			ID:   13,
			Name: "Pebbles",
		},
	}
	return mockPets
}

func populateMockPets(mockPets []Pet) error {
	// Popuate data with mock
	for _, p := range mockPets {
		err := AddPet(p)
		if err != nil {
			return err
		}
	}

	return nil
}
