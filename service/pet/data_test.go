package pet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPet(t *testing.T) {

	// Clean the data set once test is done
	defer resetData()

	tests := []struct {
		name    string
		input   Pet
		isError bool
	}{
		{
			"passing empty Pet should return a validation err",
			Pet{},
			true,
		},
		{
			"passing a valid Pet should not return an error",
			Pet{ID: 1, Name: "Tommy"},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			err := AddPet(test.input)
			assert.Equal(t, test.isError, err != nil)
			// if we saved it, let's make sure it's saved right
			if !test.isError {
				p, err := GetPetByID(test.input.ID)
				if err != nil {
					t.Error(err)
				}
				assert.Equal(t, &test.input, p)
			}
		})
	}
}

func TestGetPetByID(t *testing.T) {

	// Clean the data set once test is done
	defer resetData()

	// Popuate data with mock
	mockPets := getMockPets()
	err := populateMockPets(mockPets)
	if err != nil {
		t.Fatalf("Could not populate mock data: %v", err)
	}

	tests := []struct {
		name    string
		input   int64
		output  *Pet
		isError bool
	}{
		{
			"should be able to get a mock Pet",
			mockPets[0].ID,
			&mockPets[0],
			false,
		},
		{
			"non-existant ID should be return an error",
			42,
			nil,
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			p, err := GetPetByID(test.input)
			assert.Equal(t, test.isError, err != nil)
			assert.Equal(t, test.output, p)

		})
	}
}

func TestListPets(t *testing.T) {

	// Test 1: With no pets, ListPets should return empty slice
	t.Run(
		"clean slate should return empty slice",
		func(t *testing.T) {
			pets, err := ListPets()
			assert.Nil(t, err)
			assert.Equal(t, []Pet{}, pets)
		},
	)

	// Test 2: With mock data, ListPets should mockData
	t.Run(
		"with data populated, shoudl return all Pets",
		func(t *testing.T) {
			// Clean the data set once test is done
			defer resetData()
			// Popuate data with mock
			mockPets := getMockPets()
			err := populateMockPets(mockPets)
			if err != nil {
				t.Fatalf("Could not populate mock data: %v", err)
			}

			pets, err := ListPets()
			assert.Nil(t, err)
			assert.Equal(t, mockPets, pets)
		},
	)
}

func TestPaginate(t *testing.T) {

	// get mock pets
	mockPets := getMockPets()

	tests := []struct {
		name              string
		inputPets         []Pet
		inputMaxPerPage   int
		inputPageNum      int
		outputPets        []Pet
		outputNextPageNum int
		isError           bool
	}{
		{
			"max per page of less than 1 should error",
			[]Pet{},
			0,
			1,
			nil,
			-1,
			true,
		},
		{
			"input page num of less than 1 should error",
			mockPets,
			1,
			0,
			nil,
			-1,
			true,
		},
		{
			"input page num of more than available pages should error",
			mockPets,
			len(mockPets),
			2,
			nil,
			-1,
			true,
		},
		{
			"getting first of more than page should return 2 as next page",
			mockPets,
			len(mockPets) / 2,
			1,
			mockPets[:len(mockPets)/2],
			2,
			false,
		},
		{
			"getting second of two pages should return invalid next page",
			mockPets,
			len(mockPets) / 2,
			2,
			mockPets[len(mockPets)/2:],
			-1,
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			pets, nextPage, err := Paginate(test.inputPets, test.inputMaxPerPage, test.inputPageNum)
			assert.Equal(t, test.isError, err != nil)
			assert.Equal(t, test.outputPets, pets)
			assert.Equal(t, test.outputNextPageNum, nextPage)

		})
	}

}
