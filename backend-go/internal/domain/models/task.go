package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ProjectID   uint
	Title       string
	Description string
	Status      string `gorm:"type:varchar(20)"`
	Priority    string `gorm:"type:varchar(20)"`
	DueDate     time.Time
}
