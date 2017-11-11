package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"net/http"

	"github.com/danstis/Plex-Sync/plex"
	"github.com/danstis/Plex-Sync/webui"
	"github.com/gorilla/handlers"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	log.Printf("Plex-Sync v%v", plex.Version)
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

	plex.CacheLifetime = viper.GetInt("webui.cacheLifetime") * -1
	r := webui.NewRouter()

	loggedRouter := handlers.LoggingHandler(createLogger(viper.GetString("general.webserverlogfile")), r)
	go http.ListenAndServe(fmt.Sprintf(":%v", listeningPort), loggedRouter)
	log.Printf("Started webserver http://localhost:%v", listeningPort)
	log.SetOutput(createLogger(viper.GetString("general.logfile")))

	for {
		token := plex.Token()
		if token != "" {
			if err = localServer.GetToken(token); err != nil {
				log.Printf("ERROR: %v", err)
			}
			if err = remoteServer.GetToken(token); err != nil {
				log.Printf("ERROR: %v", err)
			}
			plex.SyncWatchedTv(localServer, remoteServer)
		}
		log.Printf("Sleeping for %v...", (sleepInterval * time.Second))
		time.Sleep(sleepInterval * time.Second)
	}
}

func createLogger(filename string) io.Writer {
	return io.MultiWriter(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    viper.GetInt("general.maxlogsize"), // megabytes
		MaxBackups: viper.GetInt("general.maxlogcount"),
		MaxAge:     viper.GetInt("general.maxlogage"), //days
	}, os.Stdout)
}
