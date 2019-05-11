package route

import (
	"fmt"
	"net/http"

	"../handler"
)

// Route represents a standard route object
type Route struct {
	Method      string
	Version     int
	Path        string
	HandlerFunc http.HandlerFunc
}

// GetPattern returns the url match pattern for the route
func (r Route) GetPattern() string {
	return fmt.Sprintf("/v%d/%s", r.Version, r.Path)
}

var routes = []Route{
	{
		Method:      http.MethodGet,
		Version:     1,
		Path:        "pets",
		HandlerFunc: handler.HandleListPets,
	},
	{
		Method:      http.MethodPost,
		Version:     1,
		Path:        "pets",
		HandlerFunc: handler.HandleCreatePet,
	},
	{
		Method:      http.MethodGet,
		Version:     1,
		Path:        "pets/{id:[0-9]+}",
		HandlerFunc: handler.HandleGetPetByID,
	},
}

// GetRoutes provides all the routes for this server
func GetRoutes() []Route {
	return routes
}
