package models

import "time"

type Notification struct {
	ID        uint `gorm:"primaryKey"`
	Message   string
	Status    string
	CreatedAt time.Time
}