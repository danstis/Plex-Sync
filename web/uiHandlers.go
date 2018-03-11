package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/danstis/Plex-Sync/plex"
)

// PageData defines a struct to store the current version information.
type PageData struct {
	Shows []string
}

var ss, _ = plex.SelectedShows()
var v = PageData{
	Shows: ss,
}

// refreshShows updates the PageData with the latest shows.
func refreshShows() {
	v.Shows, _ = plex.SelectedShows()
}

// RootHandler returns the default page.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(path.Join("web", "templates", "index.html")))
	refreshShows()
	err := tmpl.Execute(w, v)
	if err != nil {
		log.Println(err)
	}
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join("web", "templates", "settings", "settings.html"))
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join("web", "templates", "settings", "promptCredentials.html"))
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join("web", "templates", "logs.html"))
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
