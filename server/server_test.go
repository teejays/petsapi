package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/teejays/clog"

	"../service/pet"
)

func TestRouting(t *testing.T) {

	// Reduce the amount of logs
	clog.LogLevel = 6
	defer func() {
		clog.LogLevel = 0
		pet.ResetData()
	}()

	h := handler()
	srv := httptest.NewServer(h)
	defer srv.Close()

	tests := []struct {
		name            string
		method          string
		route           string
		body            string
		expectedCode    int
		expectedBody    string
		preProcessFunc  func()
		postProcessFunc func()
	}{
		{
			"create pet",
			http.MethodPost,
			"/v1/pets",
			`{"id": 1, "name": "Tommy"}`,
			http.StatusCreated,
			"",
			nil,
			func() { pet.ResetData() },
		},
		{
			"list pets",
			http.MethodGet,
			"/v1/pets",
			``,
			http.StatusOK,
			`[{"id":1,"name":"Tommy"},{"id":2,"name":"Tiger"},{"id":3,"name":"Buddy"},{"id":5,"name":"Kitty"},{"id":8,"name":"Coco"},{"id":13,"name":"Pebbles"}]`,
			func() { pet.PopulateMockPets() },
			func() { pet.ResetData() },
		},
		{
			"get pet by ID",
			http.MethodGet,
			"/v1/pets/1",
			``,
			http.StatusOK,
			`{"id":1,"name":"Tommy"}`,
			func() { pet.PopulateMockPets() },
			func() { pet.ResetData() },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var resp *http.Response
			var err error

			if tt.preProcessFunc != nil {
				tt.preProcessFunc()
			}
			if tt.postProcessFunc != nil {
				defer tt.postProcessFunc()
			}

			url := fmt.Sprintf("%s%s", srv.URL, tt.route)
			switch tt.method {
			case http.MethodGet:
				resp, err = http.Get(url)
			case http.MethodPost:
				buff := bytes.NewBufferString(tt.body)
				resp, err = http.Post(url, "", buff)
			}
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			got, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectedCode, resp.StatusCode)
			assert.Equal(t, tt.expectedBody, string(got))

		})
	}

}
