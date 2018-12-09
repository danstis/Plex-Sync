package main

import (
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/danstis/Plex-Sync/config"
	"github.com/danstis/Plex-Sync/logger"
	"github.com/danstis/Plex-Sync/plex"
	"github.com/danstis/Plex-Sync/web"
	"github.com/gorilla/handlers"
)

func main() {
	log.Printf("Plex-Sync v%v", plex.Version)

	settings, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	plex.CacheLifetime = settings.CacheLifetime

	r := web.NewRouter()

	loggedRouter := handlers.LoggingHandler(logger.CreateLogger("logs/plex-sync-webserver.log", settings.Logging.MaxLogSize, settings.Logging.MaxLogCount, settings.Logging.MaxLogAge), r)
	go http.ListenAndServe(fmt.Sprintf(":%v", settings.WebServerPort), loggedRouter)
	log.Printf("Started webserver http://localhost:%v", settings.WebServerPort)
	log.SetOutput(logger.CreateLogger("logs/plex-sync.log", settings.Logging.MaxLogSize, settings.Logging.MaxLogCount, settings.Logging.MaxLogAge))

	for {
		token := plex.Token()
		if token != "" {
			if err := settings.LocalServer.GetToken(token); err != nil {
				log.Printf("ERROR: %v", err)
			}
			if err := settings.RemoteServer.GetToken(token); err != nil {
				log.Printf("ERROR: %v", err)
			}
			plex.SyncWatchedTv(settings.LocalServer, settings.RemoteServer)
		}
		log.Printf("Sleeping for %v...", time.Duration(settings.SyncInterval)*time.Second)
		time.Sleep(time.Duration(settings.SyncInterval) * time.Second)
	}
}
