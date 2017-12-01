package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func Init(db *gorm.DB) {
	db.AutoMigrate(
		&Settings{},
	)
}
