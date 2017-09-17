package webui

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"

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

func logsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join("webui", "templates", "logs.html"))
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
