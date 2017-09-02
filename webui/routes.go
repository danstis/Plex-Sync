package webui

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Route defines a route for the WebUI
type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes defines a collection of Routes
type routes []route

// NewRouter creates a new router using all selected Routes
func NewRouter() *mux.Router {

	router := mux.NewRouter()
	for _, route := range uiroutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("./webui/static/")))
	router.PathPrefix("/static/").Handler(static)

	return router
}

var uiroutes = routes{
	route{
		"Home",
		"GET",
		"/",
		rootHandler,
	},
	route{
		"Settings",
		"GET",
		"/settings",
		settingsHandler,
	},
	route{
		"Token",
		"GET",
		"/settings/token",
		tokenHandler,
	},
	route{
		"TokenRequest",
		"GET",
		"/token/request",
		tokenRequestHandler,
	},
}
