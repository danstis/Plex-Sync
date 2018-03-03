package models

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Settings struct {
	gorm.Model
	General      General
	LocalServer  Server
	RemoteServer Server
	Webui        Webui
}

type General struct {
	gorm.Model
	SyncInterval     time.Duration
	WebserverPort    int
	Logfile          string
	WebserverLogfile string
	MaxLogSize       int
	MaxLogCount      int
	MaxLogAge        int
}

type Webui struct {
	gorm.Model
	CacheLifetime int
}

type Server struct {
	gorm.Model
	Name     string
	Hostname string
	Port     int
	Ssl      bool
	Token    string
}

//GetSettings returns the settings record from the DB
func GetSettings(db *gorm.DB) (Settings, error) {
	var s Settings
	if db.First(&s, 1).RecordNotFound() {
		return defaults(), nil
	}
	return s, db.First(&s, 1).Error
}

func defaults() Settings {
	s := Settings{
		General: General{
			SyncInterval:     600,
			WebserverPort:    8085,
			Logfile:          "logs/plex-sync.log",
			WebserverLogfile: "logs/plex-sync-webserver.log",
			MaxLogSize:       20,
			MaxLogCount:      5,
			MaxLogAge:        1,
		},
		LocalServer:  Server{},
		RemoteServer: Server{},
		Webui: Webui{
			CacheLifetime: 5,
		},
	}
	return s
}
