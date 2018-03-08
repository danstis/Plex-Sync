package models

import (
	"time"

	"github.com/danstis/Plex-Sync/database"
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

func (s *Settings) Save() error {
	return database.Conn.Save(&s).Error
}

func (s *Settings) Load() error {
	if database.Conn.Set("gorm:auto_preload", true).First(&s, 1).RecordNotFound() {
		s.SyncInterval = 600
		s.WebserverPort = 8085
		s.MaxLogSize = 20
		s.MaxLogCount = 5
		s.MaxLogAge = 1
		s.CacheLifetime = 5
		s.LocalServer = plex.Host{}
		s.RemoteServer = plex.Host{}
		return nil
	}
	return database.Conn.Set("gorm:auto_preload", true).First(&s, 1).Error
}
