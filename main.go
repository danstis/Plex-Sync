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
	localServer := plex.Host{
		Name:     viper.GetString("localServer.name"),
		Hostname: viper.GetString("localServer.hostname"),
		Port:     viper.GetInt("localServer.port"),
		Ssl:      viper.GetBool("usessl"),
	}
	remoteServer := plex.Host{
		Name:     viper.GetString("remoteServer.name"),
		Hostname: viper.GetString("remoteServer.hostname"),
		Port:     viper.GetInt("remoteServer.port"),
		Ssl:      viper.GetBool("usessl"),
	}

	log.Println("Local server details:", localServer)
	log.Println("Remote server details:", remoteServer)

	at, err := plex.ServerAccessToken(token, remoteServer.Name)
	if err != nil {
		log.Printf("error getting access token %v", err)
	}
	log.Printf("Remote Access Token %q", at)
}
