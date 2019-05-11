package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/teejays/clog"

	"./route"
)

// StartServer initializes and runs the HTTP server
func StartServer(addr string, port int) error {

	h := handler()
	http.Handle("/", h)

	// Start the server
	clog.Infof("Listenining on: %s:%d", addr, port)
	return http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), nil)

}

func handler() http.Handler {
	// Get all the routes
	routes := route.GetRoutes()

	// Start the router
	m := mux.NewRouter()

	// Set up middlewares
	m.Use(loggerMiddleware)
	m.Use(setHeaderMiddleware)

	// Range over routes and set them up
	for _, r := range routes {
		m.HandleFunc(r.GetPattern(), r.HandlerFunc).
			Methods(r.Method)
	}

	return m
}

// loggerMiddleware is a http.Handler middleware function that logs any request received
func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		clog.Debugf("Server: HTTP request received for %s %s", r.Method, r.URL.Path)
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

// setHeaderMiddleware sets the header for the response
func setHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the header
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
