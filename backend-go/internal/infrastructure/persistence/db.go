package persistence

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDB(databasePath string) (*gorm.DB, error) {
	// Connect database
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	err = AutoMigrate(db)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to SQLite:", databasePath)
	return db, nil
}
