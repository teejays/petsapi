package pet

import (
	"fmt"
	"math"
	"sort"
	"sync"
)

var data = []Pet{}
var dataMapID map[int64]int
var dataLock sync.RWMutex

// ResetData is the exported wrapper for resetData()
func ResetData() {
	resetData()
}

func resetData() {
	dataLock.Lock()
	defer dataLock.Unlock()
	data = []Pet{}
	dataMapID = make(map[int64]int)
}

// ErrNotExist represents entity not found in DB error
var ErrNotExist = fmt.Errorf("entity does not exist")

// AddPet adds a new pet
func AddPet(p Pet) error {
	// Validate
	if err := p.Validate(); err != nil {
		return err
	}

	// Apply a mutex to avoid race conditions for IDs
	dataLock.Lock()
	defer dataLock.Unlock()

	if dataMapID == nil {
		dataMapID = make(map[int64]int)
	}
	index, exists := dataMapID[p.ID]
	if exists {
		// replace the item
		data[index] = p
		return nil
	}

	data = append(data, p)
	dataMapID[p.ID] = len(data) - 1
	return nil
}

// GetPetByID gets the Pet with the provided ID
func GetPetByID(id int64) (*Pet, error) {
	// Apply a mutex so we can read safely
	dataLock.RLock()
	defer dataLock.RUnlock()

	index, exists := dataMapID[id]
	if !exists {
		return nil, ErrNotExist
	}
	p := data[index]

	return &p, nil
}

// ListPets gets the Pet with the provided ID
func ListPets() ([]Pet, error) {
	// Apply a mutex so we can read safely
	dataLock.RLock()
	defer dataLock.RUnlock()

	var pets = data
	// sort pets by ID
	sort.Slice(pets, func(i, j int) bool {
		return pets[i].ID < pets[j].ID
	})

	return data, nil
}

// Paginate takes a []Pet and returns only the elements appropriate
// for the given page
func Paginate(pets []Pet, maxPerPage, pageNum int) ([]Pet, int, error) {

	if maxPerPage < 1 {
		return nil, -1, fmt.Errorf("invalid max per page value: should be greater than 0")
	}
	if pageNum < 1 {
		return nil, -1, fmt.Errorf("invalid page number: should be greater than 0")
	}

	var maxPages = int(math.Ceil(float64(len(pets)) / float64(maxPerPage)))
	if maxPages != 0 && pageNum > maxPages {
		return nil, -1, fmt.Errorf("invalid page number: max of %d page(s), got %d", maxPages, pageNum)
	}
	var nextPageNum = pageNum + 1
	var start, end int
	start = (pageNum - 1) * maxPerPage
	end = start + maxPerPage
	if len(pets) <= end {
		end = len(pets)
		nextPageNum = -1 // no next page
	}

	return pets[start:end], nextPageNum, nil
}
