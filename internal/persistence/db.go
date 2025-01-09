package persistence

import (
	"fmt"
	"github.com/eldius/onedrive-client/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func newDB(file string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(file), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("failed to connect database: %ww", err))
	}
	if err := db.AutoMigrate(
		&model.OnedriveAccount{},
		&model.TokenData{},
		&model.DriveInfo{},
	); err != nil {
		panic(fmt.Sprintf("failed to migrate database: %ww", err))
	}
	return db
}

// GetDB returns a DB pool instance
func GetDB() *gorm.DB {
	if db == nil {
		db = newDB(".db")
	}

	return db
}
