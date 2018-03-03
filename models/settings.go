package models

import (
	"time"

	"github.com/danstis/Plex-Sync/plex"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Settings struct {
	gorm.Model
	CacheLifetime  int
	MaxLogAge      int
	MaxLogCount    int
	MaxLogSize     int
	SyncInterval   time.Duration
	WebserverPort  int
	LocalServer    plex.Host
	LocalServerID  uint
	RemoteServer   plex.Host
	RemoteServerID uint
}

//GetSettings returns the settings record from the DB
func GetSettings(db *gorm.DB) (Settings, error) {
	var s Settings
	if db.Set("gorm:auto_preload", true).First(&s, 1).RecordNotFound() {
		return defaults(), nil
	}
	return s, db.Set("gorm:auto_preload", true).First(&s, 1).Error
}

func defaults() Settings {
	s := Settings{
		SyncInterval:  600 * time.Second,
		WebserverPort: 8085,
		MaxLogSize:    20,
		MaxLogCount:   5,
		MaxLogAge:     1,
		CacheLifetime: 5,
		LocalServer:   plex.Host{},
		RemoteServer:  plex.Host{},
	}
	return s
}
