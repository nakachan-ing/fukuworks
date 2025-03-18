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

const (
	TaskStatusTodo  = "Todo"
	TaskStatusDoing = "Doing"
	TaskStatusDone  = "Done"

	PriorityLow    = "Low"
	PriorityMedium = "Medium"
	PriorityHigh   = "High"
)

func (t *Task) SetStatus(status string) {
	if status == TaskStatusTodo || status == TaskStatusDoing || status == TaskStatusDone {
		t.Status = status
	}
}

func (t *Task) SetPriority(priority string) {
	if priority == PriorityLow || priority == PriorityMedium || priority == PriorityHigh {
		t.Priority = priority
	}
}
