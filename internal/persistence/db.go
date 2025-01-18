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
		panic(fmt.Errorf("failed to connect database: %w", err))
	}
	if err := db.AutoMigrate(
		&model.OnedriveAccount{},
		&model.TokenData{},
		&model.DriveInfo{},
	); err != nil {
		panic(fmt.Errorf("failed to migrate database: %w", err))
	}
	return db
}

// GetDB returns a DB pool instance
func GetDB(file string) *gorm.DB {
	if db == nil {
		db = newDB(file)
	}

	return db
}
