package main

import (
	"log"

	"github.com/danstis/Plex-Sync/plex"
	"github.com/spf13/viper"
)

func main() {

	cp := plex.CredPrompter{}
	r := plex.TokenRequester{}
	token := plex.Token(cp, r)
	log.Printf("Token = %s", token)

	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No configuration file loaded - using defaults")
	}
	localserver := localserver{
		name:     viper.GetString("localServer.name"),
		hostname: viper.GetString("localServer.hostname"),
		port:     viper.GetInt("localServer.port"),
	}

	log.Println("Local server details:", localserver)
}

type localserver struct {
	name     string
	hostname string
	port     int
}
