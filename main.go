package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"net/http"

	"github.com/danstis/Plex-Sync/config"
	"github.com/danstis/Plex-Sync/models"
	"github.com/danstis/Plex-Sync/plex"
	"github.com/danstis/Plex-Sync/webui"
	"github.com/gorilla/handlers"
	"github.com/jinzhu/gorm"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	log.Printf("Plex-Sync v%v", plex.Version)

	db, err := gorm.Open("sqlite3", "config/Plex-Sync.db")
	if err != nil {
		log.Printf("ERROR: %v", err)
	}

	defer db.Close()

	models.Init(db)
	var settings models.Settings
	log.Printf("%v", db.First(&settings, 1))

	plex.CacheLifetime = config.Settings.Webui.CacheLifetime

	r := webui.NewRouter()

	loggedRouter := handlers.LoggingHandler(createLogger(config.Settings.General.WebserverLogfile), r)
	go http.ListenAndServe(fmt.Sprintf(":%v", config.Settings.General.WebserverPort), loggedRouter)
	log.Printf("Started webserver http://localhost:%v", config.Settings.General.WebserverPort)
	log.SetOutput(createLogger(config.Settings.General.Logfile))

	for {
		token := plex.Token()
		if token != "" {
			if err := config.Settings.LocalServer.GetToken(token); err != nil {
				log.Printf("ERROR: %v", err)
			}
			if err := config.Settings.RemoteServer.GetToken(token); err != nil {
				log.Printf("ERROR: %v", err)
			}
			plex.SyncWatchedTv(config.Settings.LocalServer, config.Settings.RemoteServer)
		}
		log.Printf("Sleeping for %v...", config.Settings.General.SyncInterval)
		time.Sleep(config.Settings.General.SyncInterval)
	}
}

func createLogger(filename string) io.Writer {
	return io.MultiWriter(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    config.Settings.General.MaxLogSize, // megabytes
		MaxBackups: config.Settings.General.MaxLogCount,
		MaxAge:     config.Settings.General.MaxLogAge, //days
	}, os.Stdout)
}
