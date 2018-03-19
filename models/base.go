package models

import (
	"github.com/danstis/Plex-Sync/plex"
	"github.com/jinzhu/gorm"
	// The sqlite DB driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Init starts the DB migration
func Init(db *gorm.DB) {
	db.AutoMigrate(
		&Settings{},
		&plex.Host{},
	)
}
