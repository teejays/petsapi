package pet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {

	tests := []struct {
		name    string
		input   Pet
		isError bool
	}{
		{
			"zero ID should be invalid",
			Pet{ID: 0, Name: "Tommy"},
			true,
		},
		{
			"negative ID should be invalid",
			Pet{ID: -1, Name: "Tommy"},
			true,
		},
		{
			"empty name ID should be invalid",
			Pet{ID: 1, Name: ""},
			true,
		},
		{
			"whitespace name should be invalid",
			Pet{ID: 0, Name: "  "},
			true,
		},
		{
			"no tags should be OK",
			Pet{ID: 1, Name: "Tommy"},
			false,
		},
		{
			"tags should be OK",
			Pet{ID: 1, Name: "Tommy", Tag: "some tag"},
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			err := test.input.Validate()
			assert.Equal(t, test.isError, err != nil)

		})
	}
}

func TestNewPet(t *testing.T) {

	tests := []struct {
		name      string
		inputID   int64
		inputName string
		inputTag  string
		isError   bool
	}{
		{
			"invalid ID should error",
			-1,
			"Tommy",
			"some tag",
			true,
		},
		{
			"invalid name should error",
			1,
			"",
			"some tag",
			true,
		},
		{
			"empty tag should be OK",
			1,
			"Tommy",
			"",
			false,
		},
		{
			"valid params should be OK",
			1,
			"Tommy",
			"some tag",
			false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			p, err := NewPet(test.inputID, test.inputName, test.inputTag)
			assert.Equal(t, test.isError, err != nil)
			if !test.isError {
				assert.Equal(t, test.inputID, p.ID)
				assert.Equal(t, test.inputName, p.Name)
				assert.Equal(t, test.inputTag, p.Tag)
			}

		})
	}
}
