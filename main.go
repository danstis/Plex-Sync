package main

import (
	"log"

	"github.com/danstis/Plex-Sync/plex"
	"github.com/spf13/viper"
)

func main() {
	cp := plex.CredPrompter{}
	r := plex.TokenRequester{}
	token, err := plex.Token(cp, r)
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
	localServer := plexServer{
		name:     viper.GetString("localServer.name"),
		hostname: viper.GetString("localServer.hostname"),
		port:     viper.GetInt("localServer.port"),
	}
	remoteServer := plexServer{
		name:     viper.GetString("remoteServer.name"),
		hostname: viper.GetString("remoteServer.hostname"),
		port:     viper.GetInt("remoteServer.port"),
	}

	log.Println("Local server details:", localServer)
	log.Println("Remote server details:", remoteServer)

	at, err := plex.ServerAccessToken(token, remoteServer.name)
	if err != nil {
		log.Printf("error getting access token %v", err)
	}
	log.Printf("Remote Access Token %q", at)
}

type plexServer struct {
	name     string
	hostname string
	port     int
}
