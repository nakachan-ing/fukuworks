package persistence

import (
	"log"

	"github.com/nakachan-ing/fukuworks/backend-go/internal/domain/models"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Task{},
	)
	if err != nil {
		return err
	}
	log.Printf("Database migration completed successfully.")
	return nil
}
