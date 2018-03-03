package web

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/danstis/Plex-Sync/plex"
	"github.com/gorilla/mux"
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
