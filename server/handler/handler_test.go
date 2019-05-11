package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/teejays/clog"

	"../../service/pet"
)

func TestHandleCreatePet(t *testing.T) {

	clog.LogLevel = 6
	defer func() {
		clog.LogLevel = 0
		pet.ResetData()
	}()

	tt := []struct {
		name               string
		content            string
		expectedCode       int
		isError            bool
		expectedErrMessage string
		expectedResponse   string
	}{
		{
			name:               "passing an empty body should return 400",
			content:            "",
			expectedCode:       http.StatusBadRequest,
			isError:            true,
			expectedErrMessage: "unexpected end of JSON input",
		},
		{
			name:               "passing an invalid JSON in body should return 400",
			content:            "{...}",
			expectedCode:       http.StatusBadRequest,
			isError:            true,
			expectedErrMessage: "invalid character '.' looking for beginning of object key string",
		},
		{
			name:               "passing a valid JSON but non-pet should return 400",
			content:            "[]",
			expectedCode:       http.StatusBadRequest,
			isError:            true,
			expectedErrMessage: "json: cannot unmarshal array into Go value of type pet.Pet",
		},
		{
			name:               "passing an empty JSON object should return 400",
			content:            "{}",
			expectedCode:       http.StatusBadRequest,
			isError:            true,
			expectedErrMessage: "invalid id: cannot be less than 1",
		},
		{
			name:               "passing a JSON Pet object without an id should return 400",
			content:            `{"name":"Tommy"}`,
			expectedCode:       http.StatusBadRequest,
			isError:            true,
			expectedErrMessage: "invalid id: cannot be less than 1",
		},
		{
			name:               "passing a JSON Pet object with a zero id should return 400",
			content:            `{"id": 0, "name":"Tommy"}`,
			expectedCode:       http.StatusBadRequest,
			isError:            true,
			expectedErrMessage: "invalid id: cannot be less than 1",
		},
		{
			name:               "passing a JSON Pet object with a negative id should return 400",
			content:            `{"id": -1, "name":"Tommy"}`,
			expectedCode:       http.StatusBadRequest,
			isError:            true,
			expectedErrMessage: "invalid id: cannot be less than 1",
		},
		{
			name:               "passing a JSON Pet object with a string as id should return 400",
			content:            `{"id": "str", "name":"Tommy"}`,
			expectedCode:       http.StatusBadRequest,
			isError:            true,
			expectedErrMessage: "json: cannot unmarshal string into Go struct field Pet.id of type int64",
		},
		{
			name:               "passing a JSON Pet object with a bool as id should return 400",
			content:            `{"id": true, "name":"Tommy"}`,
			expectedCode:       http.StatusBadRequest,
			isError:            true,
			expectedErrMessage: "json: cannot unmarshal bool into Go struct field Pet.id of type int64",
		},
		{
			name:               "passing a JSON Pet object with a missing name should return 400",
			content:            `{"id": 1}`,
			expectedCode:       http.StatusBadRequest,
			isError:            true,
			expectedErrMessage: "invalid name: cannot be empty",
		},
		{
			name:               "passing a JSON Pet object with an empty name should return 400",
			content:            `{"id": 1, "name": ""}`,
			expectedCode:       http.StatusBadRequest,
			isError:            true,
			expectedErrMessage: "invalid name: cannot be empty",
		},
		{
			name:               "passing a JSON Pet object with whitespace as name should return 400",
			content:            `{"id": 1, "name": "   "}`,
			expectedCode:       http.StatusBadRequest,
			isError:            true,
			expectedErrMessage: "invalid name: cannot be empty",
		},
		{
			name:               "passing a JSON Pet object with a number as name should return 400",
			content:            `{"id": 1, "name": 123}`,
			expectedCode:       http.StatusBadRequest,
			isError:            true,
			expectedErrMessage: "json: cannot unmarshal number into Go struct field Pet.name of type string",
		},
		{
			name:               "passing a JSON Pet object with a bool as name should return 400",
			content:            `{"id": 1, "name": false}`,
			expectedCode:       http.StatusBadRequest,
			isError:            true,
			expectedErrMessage: "json: cannot unmarshal bool into Go struct field Pet.name of type string",
		},
		{
			name:               "passing a JSON Pet object with a number as tag should return 400",
			content:            `{"id": 1, "name": "Tommy", "tag": 123}`,
			expectedCode:       http.StatusBadRequest,
			isError:            true,
			expectedErrMessage: "json: cannot unmarshal number into Go struct field Pet.tag of type string",
		},
		{
			name:               "passing a JSON Pet object with a bool as tag should return 400",
			content:            `{"id": 1, "name": "Tommy", "tag": true}`,
			expectedCode:       http.StatusBadRequest,
			isError:            true,
			expectedErrMessage: "json: cannot unmarshal bool into Go struct field Pet.tag of type string",
		},
		{
			name:               "passing a valid JSON Pet object without tag field should return 201",
			content:            `{"id": 1, "name": "Tommy"}`,
			expectedCode:       http.StatusCreated,
			isError:            false,
			expectedErrMessage: "",
		},
		{
			name:               "passing a valid JSON Pet object with a tag field should return 201",
			content:            `{"id": 1, "name": "Tommy", "tag": "pets"}`,
			expectedCode:       http.StatusCreated,
			isError:            false,
			expectedErrMessage: "",
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {

			// Create the fake HTTP request
			var buff = bytes.NewBufferString(test.content)
			var req = httptest.NewRequest(http.MethodPost, "/v1/pets", buff)
			var w = httptest.NewRecorder()

			// Call the handler
			HandleCreatePet(w, req)

			// Verify the status code
			assert.Equal(t, test.expectedCode, w.Code)

			// Verify the response
			resp := w.Result()
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			if test.isError {
				var errH Error
				err = json.Unmarshal(body, &errH)
				if err != nil {
					t.Error(err)
				}
				assert.Equal(t, test.expectedCode, int(errH.Code))
				assert.Equal(t, cleanErrMessage(test.expectedErrMessage), errH.Message)
			} else {
				assert.Equal(t, test.expectedResponse, string(body))
			}

		})
	}

}

func TestHandleListPets(t *testing.T) {

	// Reduce the amount of logs
	clog.LogLevel = 6
	defer func() {
		clog.LogLevel = 0
		pet.ResetData()
	}()

	type request struct {
		query string
		body  string
	}

	type response struct {
		statusCode int
		isError    bool
		errMessage string
		body       string
		headers    map[string]string
	}

	tests := []struct {
		name           string
		input          request
		preProcessFunc func(*http.Request)
		expected       response
	}{
		{
			name:           "a clean slate should return empty array",
			input:          request{},
			preProcessFunc: nil,
			expected: response{
				statusCode: http.StatusOK,
				isError:    false,
				errMessage: "",
				body:       `[]`,
			},
		},
		{
			name:  "a mock data state should return all mock pets",
			input: request{},
			preProcessFunc: func(r *http.Request) {
				pet.PopulateMockPets()
			},
			expected: response{
				statusCode: http.StatusOK,
				isError:    false,
				errMessage: "",
				body:       `[{"id":1,"name":"Tommy"},{"id":2,"name":"Tiger"},{"id":3,"name":"Buddy"},{"id":5,"name":"Kitty"},{"id":8,"name":"Coco"},{"id":13,"name":"Pebbles"}]`,
			},
		},
		{
			name: "limit of more than 100 should error",
			input: request{
				query: "?limit=300",
			},
			preProcessFunc: nil,
			expected: response{
				statusCode: http.StatusBadRequest,
				isError:    true,
				errMessage: "max limit allowed is 100",
			},
		},
		{
			name: "limit of less then num elements, and no page, should include next page in header",
			input: request{
				query: "?limit=1",
			},
			preProcessFunc: func(r *http.Request) {
				pet.PopulateMockPets()
			},
			expected: response{
				statusCode: http.StatusOK,
				isError:    false,
				body:       `[{"id":1,"name":"Tommy"}]`,
				headers:    map[string]string{"x-next": "/v1/pets?limit=1&page=2"},
			},
		},
		{
			name: "limit of less then num elements, and explicit not last page, should include next page in header",
			input: request{
				query: "?limit=1&page=2",
			},
			preProcessFunc: func(r *http.Request) {
				pet.PopulateMockPets()
			},
			expected: response{
				statusCode: http.StatusOK,
				isError:    false,
				body:       `[{"id":2,"name":"Tiger"}]`,
				headers:    map[string]string{"x-next": "/v1/pets?limit=1&page=3"},
			},
		},
		{
			name: "passing a string as limit should error",
			input: request{
				query: "?limit=abc&page=2",
			},
			preProcessFunc: nil,
			expected: response{
				statusCode: http.StatusBadRequest,
				isError:    true,
				errMessage: `error parsing limit value to an int: strconv.Atoi: parsing "abc": invalid syntax`,
			},
		},
		{
			name: "passing a string as page should error",
			input: request{
				query: "?limit=1&page=abc",
			},
			preProcessFunc: nil,
			expected: response{
				statusCode: http.StatusBadRequest,
				isError:    true,
				errMessage: `error parsing page value to an int: strconv.Atoi: parsing "abc": invalid syntax`,
			},
		},
		{
			name: "passing multiple values for limit should error",
			input: request{
				query: "?limit=1&limit=2&page=2",
			},
			preProcessFunc: nil,
			expected: response{
				statusCode: http.StatusBadRequest,
				isError:    true,
				errMessage: `multiple URL form values found for limit`,
			},
		},
		{
			name: "page num of more than available pages should error",
			input: request{
				query: "?limit=50&page=3",
			},
			preProcessFunc: func(r *http.Request) {
				pet.PopulateMockPets()
			},
			expected: response{
				statusCode: http.StatusBadRequest,
				isError:    true,
				errMessage: "invalid page number: max of 1 page(s), got 3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Create the fake HTTP request
			var buff = bytes.NewBufferString(tt.input.body)
			path := fmt.Sprintf("%s%s", "/v1/pets", tt.input.query)
			var r = httptest.NewRequest(http.MethodGet, path, buff)
			var w = httptest.NewRecorder()

			if tt.preProcessFunc != nil {
				tt.preProcessFunc(r)
			}

			// Call the handler
			HandleListPets(w, r)

			// Verify the status code
			assert.Equal(t, tt.expected.statusCode, w.Code)

			// Verify the response
			resp := w.Result()
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
				t.Skip()
			}

			if tt.expected.isError {
				var errH Error
				err = json.Unmarshal(body, &errH)
				if err != nil {
					t.Error(err)
				}
				assert.Equal(t, tt.expected.statusCode, int(errH.Code))
				assert.Equal(t, cleanErrMessage(tt.expected.errMessage), errH.Message)
			} else {
				assert.Equal(t, tt.expected.body, string(body))
			}

			if tt.expected.headers != nil {
				for k, v := range tt.expected.headers {
					assert.Equal(t, v, w.Header().Get(k))
				}
			}

		})
	}
}

func TestHandleGetByPetID(t *testing.T) {

	// Reduce the amount of logs
	clog.LogLevel = 6
	defer func() {
		clog.LogLevel = 0
		pet.ResetData()
	}()

	type request struct {
		pathAppend string
		body       string
	}

	type response struct {
		statusCode int
		isError    bool
		errMessage string
		body       string
		headers    map[string]string
	}

	tests := []struct {
		name           string
		input          request
		preProcessFunc func(*http.Request)
		expected       response
	}{
		{
			name:           "passing empty id should error",
			input:          request{},
			preProcessFunc: nil,
			expected: response{
				statusCode: http.StatusBadRequest,
				isError:    true,
				errMessage: "could not find var id in the route",
			},
		},
		{
			name: "passing negative id should error",
			input: request{
				pathAppend: "-1",
			},
			preProcessFunc: nil,
			expected: response{
				statusCode: http.StatusNotFound,
				isError:    true,
				errMessage: "entity does not exist",
			},
		},
		{
			name: "passing a string id should error",
			input: request{
				pathAppend: "abc",
			},
			preProcessFunc: nil,
			expected: response{
				statusCode: http.StatusBadRequest,
				isError:    true,
				errMessage: "could not convert var id to an int64: strconv.Atoi: parsing \"abc\": invalid syntax",
			},
		},
		{
			name: "passing a int ID but with no data in system should give NotFound error",
			input: request{
				pathAppend: "1",
			},
			preProcessFunc: nil,
			expected: response{
				statusCode: http.StatusNotFound,
				isError:    true,
				errMessage: "entity does not exist",
			},
		},
		{
			name: "passing a valid ID but should return the pet",
			input: request{
				pathAppend: "3",
			},
			preProcessFunc: func(r *http.Request) {
				pet.PopulateMockPets()
			},
			expected: response{
				statusCode: http.StatusOK,
				isError:    false,
				body:       `{"id":3,"name":"Buddy"}`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Create the fake HTTP request
			var buff = bytes.NewBufferString(tt.input.body)
			path := fmt.Sprintf("%s%s", "/v1/pets/", tt.input.pathAppend)
			var r = httptest.NewRequest(http.MethodGet, path, buff)
			r = mux.SetURLVars(r, map[string]string{"id": tt.input.pathAppend})
			var w = httptest.NewRecorder()

			if tt.preProcessFunc != nil {
				tt.preProcessFunc(r)
			}

			// Call the handler
			HandleGetPetByID(w, r)

			// Verify the status code
			assert.Equal(t, tt.expected.statusCode, w.Code)

			// Verify the response
			resp := w.Result()
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
				t.Skip()
			}

			if w.Code != http.StatusOK && tt.expected.isError {
				var errH Error
				err = json.Unmarshal(body, &errH)
				if err != nil {
					t.Error(err)
				}
				assert.Equal(t, tt.expected.statusCode, int(errH.Code))
				assert.Equal(t, cleanErrMessage(tt.expected.errMessage), errH.Message)
			} else {
				assert.Equal(t, tt.expected.body, string(body))
			}

			if tt.expected.headers != nil {
				for k, v := range tt.expected.headers {
					assert.Equal(t, v, w.Header().Get(k))
				}
			}

		})
	}
}
