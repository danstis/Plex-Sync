package main

import (
	"log"

	"github.com/spf13/viper"
)

// Version contains the version of the app.
const Version = "0.0.1"

func main() {

	cp := credPrompter{}
	r := tokenRequester{}
	token := token(cp, r)
	log.Printf("Token = %s", token)

	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No configuration file loaded - using defaults")
	}
	viper.GetString("myplex.token")
}
