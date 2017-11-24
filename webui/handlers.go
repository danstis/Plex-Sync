package webui

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"

	"github.com/danstis/Plex-Sync/plex"
)

// PageData defines a struct to store the current version information.
type PageData struct {
	Version     string
	FullVersion string
	Shows       []string
}

var ss, _ = plex.SelectedShows()
var v = PageData{
	Version:     plex.ShortVersion,
	FullVersion: plex.Version,
	Shows:       ss,
}

// refreshShows updates the PageData with the latest shows.
func refreshShows() {
	v.Shows, _ = plex.SelectedShows()
}

// RootHandler returns the default page.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(path.Join("webui", "templates", "index.html")))
	refreshShows()
	err := tmpl.Execute(w, v)
	if err != nil {
		log.Println(err)
	}
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(path.Join("webui", "templates", "settings", "settings.html")))
	err := tmpl.Execute(w, v)
	if err != nil {
		log.Println(err)
	}
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(path.Join("webui", "templates", "settings", "promptCredentials.html")))
	err := tmpl.Execute(w, v)
	if err != nil {
		log.Println(err)
	}
}

func logsHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(path.Join("webui", "templates", "logs.html")))
	err := tmpl.Execute(w, v)
	if err != nil {
		log.Println(err)
	}
}

func generalLogHeadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := vars["logfile"]
	fi, err := os.Stat(path.Join("logs", file))
	if err != nil {
		log.Println("Failed to find log file")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%v", fi.Size()))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusOK)
}

func generalLogHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := vars["logfile"]
	http.ServeFile(w, r, path.Join("logs", file))
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
