package main

import (
	"log"

	"github.com/spf13/viper"
)

func main() {

	cp := credPrompter{}
	token := token(cp)
	log.Printf("Token = %s", token)

	viper.SetConfigName("config")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Println("No configuration file loaded - using defaults")
	}
	viper.GetString("myplex.token")
}
