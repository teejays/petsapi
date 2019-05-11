package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"../../service/pet"
	"github.com/gorilla/mux"
	"github.com/teejays/clog"
)

// HandleListPets returns all the pets
func HandleListPets(w http.ResponseWriter, r *http.Request) {
	clog.Debugf("Request Path: %+v", r.URL)
	// Get the query params
	defaultLimit := 100
	limit, err := getQueryParamInt(r, "limit", defaultLimit)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, false)
		return
	}
	clog.Debugf("limit = %d", limit)
	if limit > defaultLimit {
		writeError(w, http.StatusBadRequest, fmt.Errorf("max limit allowed is %d", defaultLimit), false)
		return
	}

	defaultPage := 1
	page, err := getQueryParamInt(r, "page", defaultPage)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, false)
		return
	}

	// Get the pets
	pets, err := pet.ListPets()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err, true)
		return
	}

	// Handle pagination
	var nextPage int
	pets, nextPage, err = pet.Paginate(pets, limit, page)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, false)
		return
	}

	// Add the header for next page url
	if nextPage > 0 {
		w.Header().Set("x-next", fmt.Sprintf("%s?limit=%d&page=%d", r.URL.Path, limit, nextPage))
	}

	// Set the response
	writeResponse(w, http.StatusOK, pets)
	return

}

// HandleCreatePet creates a new pet and stores it
func HandleCreatePet(w http.ResponseWriter, r *http.Request) {

	// Read the HTTP request body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, false)
		return
	}
	defer r.Body.Close()

	// Unmarshal JSON into Go type
	var p pet.Pet
	err = json.Unmarshal(body, &p)
	if err != nil {
		writeError(w, http.StatusBadRequest, err, false)
		return
	}

	// Validate that it is good to save
	err = p.Validate()
	if err != nil {
		writeError(w, http.StatusBadRequest, err, false)
		return
	}

	// Save the new pet
	err = pet.AddPet(p)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err, true)
		return
	}

	writeResponse(w, http.StatusCreated, nil)

}

// HandleGetPetByID fetches the pet that has the provided ID
func HandleGetPetByID(w http.ResponseWriter, r *http.Request) {
	clog.Debugf("Request Path: %+v", r.URL)

	// Get the Pet ID
	id, err := getMuxParamrInt(r, "id")
	if err != nil {
		writeError(w, http.StatusBadRequest, err, false)
		return
	}

	// Get the pet
	p, err := pet.GetPetByID(id)
	if err == pet.ErrNotExist {
		writeError(w, http.StatusNotFound, err, false)
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, err, true)
		return
	}

	// Write the response
	writeResponse(w, http.StatusOK, p)
}

func getQueryParamInt(r *http.Request, name string, defaultVal int) (int, error) {
	err := r.ParseForm()
	if err != nil {
		return defaultVal, err
	}
	values, exist := r.Form[name]
	clog.Debugf("URL values for %s: %+v", name, values)
	if !exist {
		return defaultVal, nil
	}
	if len(values) > 1 {
		return defaultVal, fmt.Errorf("multiple URL form values found for %s", name)
	}

	val, err := strconv.Atoi(values[0])
	if err != nil {
		return defaultVal, fmt.Errorf("error parsing %s value to an int: %v", name, err)
	}
	return val, nil
}

// getMuxParamrInt extracts the userid param out of the request route
func getMuxParamrInt(r *http.Request, name string) (int64, error) {

	var vars = mux.Vars(r)
	clog.Debugf("MUX vars are: %+v", vars)
	valStr := vars[name]
	if strings.TrimSpace(valStr) == "" {
		return -1, fmt.Errorf("could not find var %s in the route", name)
	}

	val, err := strconv.Atoi(valStr)
	if err != nil {
		return -1, fmt.Errorf("could not convert var %s to an int64: %v", name, err)
	}

	return int64(val), nil
}

func writeResponse(w http.ResponseWriter, code int, v interface{}) {
	w.WriteHeader(code)

	if v == nil {
		return
	}

	// Json marshal the resp
	data, err := json.Marshal(v)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err, true)
		return
	}
	// Write the response
	_, err = w.Write(data)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err, true)
		return
	}
}

func writeError(w http.ResponseWriter, code int, err error, hide bool) {
	errMessage := cleanErrMessage(err.Error())
	clog.Error(errMessage)

	if hide {
		errMessage = apiErrMessageClean
	}

	errE := NewError(code, errMessage)

	w.WriteHeader(code)
	data, err := json.Marshal(errE)
	if err != nil {
		panic(fmt.Sprintf("Failed to json.Unmarshal an error for http response: %v", err))
	}
	_, err = w.Write(data)
	if err != nil {
		panic(fmt.Sprintf("Failed to write error to the http response: %v", err))
	}
}

func cleanErrMessage(msg string) string {
	return fmt.Sprintf("There was an error processing the request: %v", msg)
}
