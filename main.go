package main

import (
	"fmt"
	"log"
	"time"

	"net/http"

	"github.com/danstis/Plex-Sync/logger"
	"github.com/danstis/Plex-Sync/models"
	"github.com/danstis/Plex-Sync/plex"
	"github.com/danstis/Plex-Sync/webui"
	"github.com/gorilla/handlers"
	"github.com/jinzhu/gorm"
)

func main() {
	log.Printf("Plex-Sync v%v", plex.Version)

	db, err := gorm.Open("sqlite3", "config/Plex-Sync.db")
	if err != nil {
		log.Printf("ERROR: %v", err)
	}

	defer db.Close()

	models.Init(db)
	settings, err := models.GetSettings(db)
	if err != nil {
		log.Fatal(err)
	}

	plex.CacheLifetime = settings.Webui.CacheLifetime

	r := webui.NewRouter()

	loggedRouter := handlers.LoggingHandler(logger.CreateLogger(settings.General.WebserverLogfile), r)
	go http.ListenAndServe(fmt.Sprintf(":%v", settings.General.WebserverPort), loggedRouter)
	log.Printf("Started webserver http://localhost:%v", settings.General.WebserverPort)
	log.SetOutput(logger.CreateLogger(settings.General.Logfile))

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
		log.Printf("Sleeping for %v...", settings.General.SyncInterval)
		time.Sleep(settings.General.SyncInterval)
	}
}
