package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/danstis/Plex-Sync/models"
	"github.com/danstis/Plex-Sync/plex"
	"github.com/gorilla/mux"
)

const (
	// JSONContentType is a helper to easily set the content type
	JSONContentType = "application/json; charset=utf-8"
)

func apiLogHead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := vars["logfile"]
	fi, err := os.Stat(path.Join("logs", file))
	if err != nil {
		log.Println("Failed to find log file")
		w.WriteHeader(http.StatusNotFound)
		os.Exit(1)
		return
	}

	w.Header().Set("Content-Length", fmt.Sprintf("%v", fi.Size()))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Accept-Ranges", "bytes")
	w.WriteHeader(http.StatusNoContent)
}

func apiLogGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	file := vars["logfile"]
	http.ServeFile(w, r, path.Join("logs", file)) //TODO: Setting could overwrite this path
}

func apiTokenDelete(w http.ResponseWriter, r *http.Request) {
	err := plex.RemoveCachedToken()
	if err != nil {
		log.Printf("unable to remove token file: %v\n", err)
		w.WriteHeader(http.StatusFailedDependency)
	}
	w.WriteHeader(http.StatusNoContent)
}

func apiSettingsGet(w http.ResponseWriter, r *http.Request) {
	var s models.Settings
	var err error

	if err = s.Load(); err != nil {
		log.Printf("Error getting settings from DB: %v\n", err)
		w.WriteHeader(http.StatusFailedDependency)
		return
	}
	jv, err := json.Marshal(s)
	if err != nil {
		log.Printf("Error converting settings to JSON: %v\n", err)
		w.WriteHeader(http.StatusFailedDependency)
		return
	}
	w.Header().Set("Content-Type", JSONContentType)
	fmt.Fprintf(w, string(jv))
}

func apiSettingsCreate(w http.ResponseWriter, r *http.Request) {
	var s models.Settings
	if err := s.Load(); err != nil {
		log.Printf("Error reading settings from DB: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewDecoder(r.Body).Decode(&s)

	if err := s.Save(); err != nil {
		log.Printf("Error writing settings to DB: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", JSONContentType)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&s)
}

func apiVersionGet(w http.ResponseWriter, r *http.Request) {
	type version struct {
		ShortVersion string `json:"shortVersion"`
		Fullversion  string `json:"fullVersion"`
	}
	v := version{
		ShortVersion: plex.ShortVersion,
		Fullversion:  plex.Version,
	}
	w.Header().Set("Content-Type", JSONContentType)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&v)
}

func apiShowsGet(w http.ResponseWriter, r *http.Request) {
	ss, err := plex.SelectedShows()
	if err != nil {
		log.Printf("Error getting selected shows: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", JSONContentType)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&ss)
}
