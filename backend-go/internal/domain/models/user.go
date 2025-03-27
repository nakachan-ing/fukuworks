package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Email    string    `gorm:"unique"`
	Projects []Project `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	// for auth test
	Password string
}
