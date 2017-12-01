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
