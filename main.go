package main

import (
	"log"
	"os"
	"path"
	"time"

	"net/http"

	"github.com/danstis/Plex-Sync/plex"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/", rootHandler)
	s := http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))
	r.PathPrefix("/static/").Handler(s)
	http.Handle("/", r)
	r.HandleFunc("/settings", settingsHandler)
	http.Handle("/settings", r)
	r.HandleFunc("/settings/token", tokenHandler)
	http.Handle("/settings/token", r)
	r.HandleFunc("/token/request", tokenRequestHandler)
	http.Handle("/token/request", r)

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	go http.ListenAndServe(":80", loggedRouter)
	log.Println("Started webserver http://localhost")

	cp := plex.CredPrompter{}
	tr := plex.TokenRequester{}
	token, err := plex.Token(cp, tr)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	err = viper.ReadInConfig()
	if err != nil {
		log.Println("No configuration file loaded - using defaults")
	}
	sleepInterval := viper.GetDuration("general.interval")
	localServer := plex.Host{
		Name:      viper.GetString("localServer.name"),
		Hostname:  viper.GetString("localServer.hostname"),
		Port:      viper.GetInt("localServer.port"),
		Ssl:       viper.GetBool("localServer.usessl"),
		TvSection: viper.GetInt("localServer.tvsection"),
	}
	remoteServer := plex.Host{
		Name:      viper.GetString("remoteServer.name"),
		Hostname:  viper.GetString("remoteServer.hostname"),
		Port:      viper.GetInt("remoteServer.port"),
		Ssl:       viper.GetBool("remoteServer.usessl"),
		TvSection: viper.GetInt("remoteServer.tvsection"),
	}

	localServer.GetToken(token)
	remoteServer.GetToken(token)

	for {
		plex.SyncWatchedTv(localServer, remoteServer)
		log.Printf("Sleeping for %v...", (sleepInterval * time.Second))
		time.Sleep(sleepInterval * time.Second)
	}
}

// RootHandler returns the default page.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join("templates", "index.html"))
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join("templates", "settings", "settings.html"))
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, path.Join("templates", "settings", "promptCredentials.html"))
}

func tokenRequestHandler(w http.ResponseWriter, r *http.Request) {
	// request new token using credentials passed in form
	http.ServeFile(w, r, path.Join("templates", "settings", "promptCredentials.html"))
}
