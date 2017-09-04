package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"net/http"

	"github.com/danstis/Plex-Sync/plex"
	"github.com/danstis/Plex-Sync/webui"
	"github.com/gorilla/handlers"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No configuration file loaded - using defaults")
	}
	sleepInterval := viper.GetDuration("general.interval")
	listeningPort := viper.GetInt("general.webserverport")
	localServer := plex.Host{
		Name:     viper.GetString("localServer.name"),
		Hostname: viper.GetString("localServer.hostname"),
		Port:     viper.GetInt("localServer.port"),
		Ssl:      viper.GetBool("localServer.usessl"),
	}
	remoteServer := plex.Host{
		Name:     viper.GetString("remoteServer.name"),
		Hostname: viper.GetString("remoteServer.hostname"),
		Port:     viper.GetInt("remoteServer.port"),
		Ssl:      viper.GetBool("remoteServer.usessl"),
	}

	r := webui.NewRouter()

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	go http.ListenAndServe(fmt.Sprintf(":%v", listeningPort), loggedRouter)
	log.Printf("Started webserver http://localhost:%v", listeningPort)

	for {
		token := plex.Token()
		if token != "" {
			localServer.GetToken(token)
			remoteServer.GetToken(token)

			plex.SyncWatchedTv(localServer, remoteServer)
		}
		log.Printf("Sleeping for %v...", (sleepInterval * time.Second))
		time.Sleep(sleepInterval * time.Second)
	}
}
