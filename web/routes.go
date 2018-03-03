package web

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
	for _, route := range uiRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	for _, route := range apiRoutes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	static := http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/")))
	cache := http.StripPrefix("/cache/", http.FileServer(http.Dir("./.cache/")))
	router.PathPrefix("/static/").Handler(static)
	router.PathPrefix("/cache/").Handler(cache)

	return router
}

var uiRoutes = routes{
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
		"POST",
		"/token/request",
		tokenRequestHandler,
	},
	route{
		"TokenRemove",
		"GET",
		"/token/remove",
		tokenRemoveHandler,
	},
	route{
		"Logs",
		"GET",
		"/logs",
		logsHandler,
	},
}

var apiRoutes = routes{
	route{
		"ApiLogHead",
		"HEAD",
		"/api/log/{logfile}",
		apiLogHead,
	},
	route{
		"ApiLogGet",
		"GET",
		"/api/log/{logfile}",
		apiLogGet,
	},
	route{
		"ApiTokenDelete",
		"DELETE",
		"/api/token",
		apiTokenDelete,
	},
	route{
		"ApiSettingsGet",
		"GET",
		"/api/settings",
		apiSettingsGet,
	},
	route{
		"ApiSettingsPost",
		"POST",
		"/api/settings",
		apiSettingsCreate,
	},
}
