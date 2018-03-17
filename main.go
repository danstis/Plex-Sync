package main

import (
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/danstis/Plex-Sync/database"
	"github.com/danstis/Plex-Sync/logger"
	"github.com/danstis/Plex-Sync/models"
	"github.com/danstis/Plex-Sync/plex"
	"github.com/danstis/Plex-Sync/web"
	"github.com/gorilla/handlers"
	"github.com/jinzhu/gorm"
)

func main() {
	log.Printf("Plex-Sync v%v", plex.Version)

	var err error
	database.Conn, err = gorm.Open("sqlite3", "config/Plex-Sync.db")
	if err != nil {
		log.Printf("ERROR: %v", err)
	}
	defer database.Conn.Close()

	models.Init(database.Conn)
	var settings models.Settings
	if err := settings.Load(); err != nil {
		log.Fatal(err)
	}

	plex.CacheLifetime = settings.CacheLifetime

	r := web.NewRouter()

	loggedRouter := handlers.LoggingHandler(logger.CreateLogger("logs/plex-sync-webserver.log", settings.MaxLogSize, settings.MaxLogCount, settings.MaxLogAge), r)
	go http.ListenAndServe(fmt.Sprintf(":%v", settings.WebserverPort), loggedRouter)
	log.Printf("Started webserver http://localhost:%v", settings.WebserverPort)
	log.SetOutput(logger.CreateLogger("logs/plex-sync.log", settings.MaxLogSize, settings.MaxLogCount, settings.MaxLogAge))

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
		log.Printf("Sleeping for %v...", settings.SyncInterval*time.Second)
		time.Sleep(settings.SyncInterval * time.Second)
	}
}
