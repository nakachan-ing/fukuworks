package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	UserID       uint
	Title        string
	Description  string
	Platform     string
	Client       string
	EstimatedFee float64 `gorm:"type:REAL"`
	Status       string  `gorm:"type:varchar(20)"`
	Deadline     time.Time
	Tasks        []Task `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}

const (
	StatusOpen       = "Open"
	StatusInProgress = "In Progress"
	StatusCompleted  = "Completed"
	StatusCanceled   = "Canceled"
)

func (p *Project) SetStatus(status string) {
	if status == StatusOpen || status == StatusInProgress || status == StatusCompleted || status == StatusCanceled {
		p.Status = status
	}
}
