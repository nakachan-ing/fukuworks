package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	ProjectID   uint
	Number      uint
	Title       string
	Description string
	Status      string `gorm:"type:varchar(20)"`
	Priority    string `gorm:"type:varchar(20)"`
	DueDate     time.Time
}

var allowedTaskStatuses = map[string]bool{
	"Todo":  true,
	"Doing": true,
	"Done":  true,
}

var allowedPriorities = map[string]bool{
	"Low":    true,
	"Medium": true,
	"High":   true,
}

func (t *Task) SetStatus(status string) error {
	if !allowedTaskStatuses[status] {
		return errors.New("invalid status" + status)
	}
	t.Status = status
	return nil
}

func (t *Task) SetPriority(priority string) error {
	if !allowedPriorities[priority] {
		return errors.New("invalid priority" + priority)
	}
	t.Priority = priority
	return nil
}
