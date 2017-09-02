package webui

import (
	"net/http"
	"path"
)

// RootHandler returns the default page.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join("webui", "templates", "index.html"))
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join("webui", "templates", "settings", "settings.html"))
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join("webui", "templates", "settings", "promptCredentials.html"))
}

func tokenRequestHandler(w http.ResponseWriter, r *http.Request) {
	// request new token using credentials passed in form
	http.ServeFile(w, r, path.Join("webui", "templates", "settings", "promptCredentials.html"))
}
