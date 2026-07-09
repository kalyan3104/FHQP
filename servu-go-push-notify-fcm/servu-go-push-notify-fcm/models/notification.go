package models 

import "time"

type Notification struct {
	ID        string    `gorm:"primaryKey"`
	UserID    string
	Title     string
	Message   string
	IsRead    bool
	CreatedAt time.Time
}