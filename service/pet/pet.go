package pet

import (
	"fmt"
	"strings"
)

// Pet represents the model for pet entity
type Pet struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Tag  string `json:"tag,omitempty"`
}

var ErrInvalidID = fmt.Errorf("invalid id: cannot be less than 1")
var ErrInvalidName = fmt.Errorf("invalid name: cannot be empty")

// Validate returns an error if any of the fields in Pet is not valid
func (p Pet) Validate() error {
	// Validate params
	if p.ID < 1 {
		return ErrInvalidID
	}
	if strings.TrimSpace(p.Name) == "" {
		return ErrInvalidName
	}
	return nil
}

// NewPet create a new instance of a Pet
func NewPet(id int64, name, tag string) (Pet, error) {
	// Validate again
	// Create a new instance
	p := Pet{
		ID:   id,
		Name: name,
		Tag:  tag,
	}

	if err := p.Validate(); err != nil {
		return Pet{}, err
	}

	return p, nil
}
