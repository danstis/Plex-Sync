package database

import (
	"github.com/jinzhu/gorm"
)

// Conn stores the database connection for reference by the other packages
var Conn *gorm.DB
