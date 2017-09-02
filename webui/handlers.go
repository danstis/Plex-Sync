package webui

import (
	"fmt"
	"net/http"
	"path"

	"github.com/danstis/Plex-Sync/plex"
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
	creds := plex.Credentials{
		Username: r.PostFormValue("username"),
		Password: r.PostFormValue("password"),
	}
	err := plex.TokenRequest(creds)
	if err != nil {
		// TODO: Handle errors
		fmt.Fprintf(w, fmt.Sprintf("Error getting new token: %v", err))
		return
	}

	http.Redirect(w, r, "/", http.StatusFound) // redirect user to homepage
}

func tokenRemoveHandler(w http.ResponseWriter, r *http.Request) {
	err := plex.RemoveCachedToken()
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/", http.StatusFound) // redirect user to homepage
}
