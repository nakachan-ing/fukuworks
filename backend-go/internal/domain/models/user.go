package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Email    string    `gorm:"unique"`
	Projects []Project `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Password string
	// for role test
	Role string `gorm:"default:user"` // user or admin

}
