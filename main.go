package main

import (
	"log"
	"os"

	"net/http"

	"time"

	"path"

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
	r.HandleFunc("/settings/token", tokenHandler)
	http.Handle("/settings/token", r)

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	go http.ListenAndServe(":80", loggedRouter)
	log.Println("Started webserver http://localhost")

	cp := plex.CredPrompter{}
	tr := plex.TokenRequester{}
	token, err := plex.Token(cp, tr)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	log.Printf("Token = %s", token)

	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	err = viper.ReadInConfig()
	if err != nil {
		log.Println("No configuration file loaded - using defaults")
	}
	localserver := localserver{
		name:     viper.GetString("localServer.name"),
		hostname: viper.GetString("localServer.hostname"),
		port:     viper.GetInt("localServer.port"),
	}

	log.Println("Local server details:", localserver)

	for {
		log.Println("Waiting")
		time.Sleep(60 * time.Second)
	}
}

type localserver struct {
	name     string
	hostname string
	port     int
}

// RootHandler returns the default page.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	// w.WriteHeader(http.StatusOK)
	http.ServeFile(w, r, path.Join("templates", "index.html"))
}

func tokenHandler(w http.ResponseWriter, r *http.Request) {
	// tmplt := template.New("promptCredentials.html")
	// tmplt, err := tmplt.ParseFiles(path.Join("templates", "settings", "promptCredentials.html"))
	// if err != nil {
	// 	w.WriteHeader(http.StatusOK)
	// 	return
	// }
	//w.WriteHeader(http.StatusOK)
	// tmplt.Execute(w, "something")
	http.ServeFile(w, r, path.Join("templates", "settings", "promptCredentials.html"))
}
