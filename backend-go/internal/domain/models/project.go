package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	UserID       uint
	Number       uint
	Title        string
	Description  string
	Platform     string
	Client       string
	EstimatedFee float64 `gorm:"type:REAL"`
	Status       string  `gorm:"type:varchar(20)"`
	Deadline     time.Time
	Tasks        []Task `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}

var allowedStatuses = map[string]bool{
	"NotStarted": true,
	"InProgress": true,
	"Completed":  true,
	"Canceled":   true,
}

func (p *Project) SetStatus(status string) error {
	if !allowedStatuses[status] {
		return errors.New("invalid status" + status)
	}
	p.Status = status
	return nil
}
