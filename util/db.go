package util

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// CoreData/Download/downloadTask.db
func OpenDb(db string) (*gorm.DB, error) {
	return gorm.Open("sqlite3", db)
}
