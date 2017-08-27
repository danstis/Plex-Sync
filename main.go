package main

import (
	"log"
	"time"

	"github.com/danstis/Plex-Sync/plex"
	"github.com/spf13/viper"
)

func main() {
	cp := plex.CredPrompter{}
	r := plex.TokenRequester{}
	token, err := plex.Token(cp, r)
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
